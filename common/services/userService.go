package services

import (
	"server/common/appError"
	"server/common/config"
	"server/common/dto"
	"server/common/models"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserCursor struct {
	LastName  string
	FirstName string
	ID        uint
}

// GetUsers returns a paginated list of users, excluding those blocked by or blocking the requester.
// Cursor pagination: results are ordered by last_name, first_name, id.
// Pass cursor=nil for first page. Use returned nextCursor for subsequent pages.
func GetUsers(requestingUserID uint, searchQuery string, limit int, cursor *UserCursor) ([]models.User, *UserCursor, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	// Users that I have blocked
	iBlockedSubQuery := config.DB.Table("user_blocked_users").
		Select("blocked_user_id").
		Where("user_id = ?", requestingUserID)

	// Users that have blocked me
	blockedMeSubQuery := config.DB.Table("user_blocked_users").
		Select("user_id").
		Where("blocked_user_id = ?", requestingUserID)

	q := config.DB.
		Model(&models.User{}).
		Preload("FavoriteSports").
		Where("id != ?", requestingUserID).
		Where("id NOT IN (?)", iBlockedSubQuery).
		Where("id NOT IN (?)", blockedMeSubQuery)

	// Search (FB-like: match first, last, full name)
	searchQuery = strings.TrimSpace(searchQuery)
	if searchQuery != "" {
		like := "%" + searchQuery + "%"
		q = q.Where(`(
			first_name ILIKE ? OR
			last_name ILIKE ? OR
			(first_name || ' ' || last_name) ILIKE ?
		)`, like, like, like)
	}

	// Cursor pagination: fetch only rows "after" the cursor in the sort order
	// (last_name, first_name, id) is a stable ordering.
	if cursor != nil {
		q = q.Where(`
			(last_name, first_name, id) > (?, ?, ?)
		`, cursor.LastName, cursor.FirstName, cursor.ID)
	}

	// IMPORTANT: consistent ordering (cursor relies on this)
	q = q.Order("last_name ASC").Order("first_name ASC").Order("id ASC").Limit(limit + 1)

	var users []models.User
	if err := q.Find(&users).Error; err != nil {
		return nil, nil, err
	}

	// Determine next cursor (limit+1 trick)
	var nextCursor *UserCursor
	if len(users) > limit {
		last := users[limit-1]
		nextCursor = &UserCursor{
			LastName:  last.LastName,
			FirstName: last.FirstName,
			ID:        last.ID,
		}
		users = users[:limit]
	}

	return users, nextCursor, nil
}

// GetUserByID fetches a user by ID directly.
// IMPORTANT: Use GetVisibleUser for controller logic to ensure blocking rules are applied.
func GetUserByID(userID uint) (*models.User, error) {
	var user models.User

	err := config.DB.Preload("FavoriteSports").
		Preload("Friends").
		Preload("Teams").
		Preload("JoinedChallenges").
		Preload("CreatedChallenges").
		Preload("EmergencyContacts").
		First(&user, userID).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetVisibleUser fetches a user only if they are not blocked by the requester.
func GetVisibleUser(requestingUserID, targetUserID uint) (*models.User, error) {
	// If asking for self, just return ID
	if requestingUserID == targetUserID {
		return GetUserByID(targetUserID)
	}

	// Check if a block exists between the two users
	if IsBlocked(requestingUserID, targetUserID) {
		// Return UserNotFound to avoid leaking that the user exists but is blocked
		return nil, appError.ErrUserNotFound
	}

	return GetUserByID(targetUserID)
}

func GetUserByIDWithSettings(userID uint) (*models.User, error) {
	var user models.User

	err := config.DB.
		Preload("FavoriteSports").
		Preload("Friends").
		Preload("Teams").
		Preload("Settings").
		Preload("JoinedChallenges", func(db *gorm.DB) *gorm.DB {
			return db.Order("date ASC").Order("start_time ASC")
		}).
		Preload("JoinedChallenges.Location").
		Preload("JoinedChallenges.Creator").
		Preload("JoinedChallenges.Users").
		Preload("JoinedChallenges.Teams").
		Preload("EmergencyContacts").
		First(&user, userID).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserSettings(userID uint) (*models.UserSettings, error) {
	var settings models.UserSettings
	err := config.DB.First(&settings, userID).Error

	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func GetInCommonStats(currentUserID, targetUserID uint) (dto.CommonStatsDto, error) {
	var stats dto.CommonStatsDto

	// Check blocking
	if IsBlocked(currentUserID, targetUserID) {
		return stats, appError.ErrUserNotFound
	}

	db := config.DB

	// 1. Common Teams Count
	var count int64
	err := db.Table("user_teams as t1").
		Joins("JOIN user_teams as t2 ON t1.team_id = t2.team_id").
		Where("t1.user_id = ? AND t2.user_id = ?", currentUserID, targetUserID).
		Count(&count).Error

	if err != nil {
		return stats, err
	}
	stats.CommonTeamsCount = count

	// 2. Common Friends Count
	err = db.Table("user_friends as f1").
		Joins("JOIN user_friends as f2 ON f1.friend_id = f2.friend_id").
		Where("f1.user_id = ? AND f2.user_id = ?", currentUserID, targetUserID).
		Count(&count).Error

	if err != nil {
		return stats, err
	}
	stats.CommonFriendsCount = count

	// 3. Common Sports (Favorites)
	var commonSports []models.Sport
	err = db.Table("sports").
		Joins("JOIN user_favorite_sports as us1 ON us1.sport_id = sports.id").
		Joins("JOIN user_favorite_sports as us2 ON us2.sport_id = sports.id").
		Where("us1.user_id = ? AND us2.user_id = ?", currentUserID, targetUserID).
		Find(&commonSports).Error

	if err != nil {
		return stats, err
	}

	// Convert models to DTOs
	stats.CommonSports = make([]dto.SportDto, len(commonSports))
	for i, s := range commonSports {
		stats.CommonSports[i] = dto.ToSportDto(s)
	}

	return stats, nil
}

func CreateUser(newUser models.User, password string) (*models.User, error) {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var existingUser models.User

		err := tx.Where("email = ?", newUser.Email).
			First(&existingUser).
			Error

		if err == nil {
			return appError.ErrUserExists
		}

		if password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			hashedPasswordStr := string(hashedPassword)
			newUser.Password = &hashedPasswordStr
		}

		err = tx.Omit("FavoriteSports").Create(&newUser).Error
		if err != nil {
			return err
		}

		// Associate favorite sports if provided
		if len(newUser.FavoriteSports) > 0 {
			sports := make([]string, len(newUser.FavoriteSports))
			for i, sport := range newUser.FavoriteSports {
				sports[i] = sport.Name
			}
			if err := associateFavoriteSports(tx, newUser.ID, sports); err != nil {
				return err
			}
			// Reload user with favorite sports
			if err := tx.Preload("FavoriteSports").First(&newUser, newUser.ID).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func UpdateUser(userID uint, user dto.UserUpdateDto) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var existingUser models.User
		if err := tx.Preload("FavoriteSports").First(&existingUser, userID).Error; err != nil {
			return err
		}

		if user.FirstName != "" {
			existingUser.FirstName = user.FirstName
		}

		if user.LastName != "" {
			existingUser.LastName = user.LastName
		}

		if user.ProfilePicture != "" {
			existingUser.ProfilePicture = user.ProfilePicture
		}

		if user.Bio != "" {
			existingUser.Bio = user.Bio
		}

		if user.BirthDate != (time.Time{}) {
			existingUser.BirthDate = user.BirthDate
		}

		if user.City != "" {
			existingUser.City = user.City
		}

		if err := tx.Save(&existingUser).Error; err != nil {
			return err
		}

		// Update favorite sports if provided
		if user.FavoriteSports != nil {
			if err := associateFavoriteSports(tx, userID, user.FavoriteSports); err != nil {
				return err
			}
		}

		return nil
	})
}

func UpdateUserSettings(userID uint, settingsDto dto.UserSettingsUpdateDto) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var settings models.UserSettings

		if err := tx.First(&settings, userID).Error; err != nil {
			return err
		}

		if settingsDto.NotifyTeamInvites != nil {
			settings.NotifyTeamInvites = *settingsDto.NotifyTeamInvites
		}
		if settingsDto.NotifyTeamMembership != nil {
			settings.NotifyTeamMembership = *settingsDto.NotifyTeamMembership
		}

		if settingsDto.NotifyFriendRequests != nil {
			settings.NotifyFriendRequests = *settingsDto.NotifyFriendRequests
		}
		if settingsDto.NotifyFriendUpdates != nil {
			settings.NotifyFriendUpdates = *settingsDto.NotifyFriendUpdates
		}

		if settingsDto.NotifyChallengeInvites != nil {
			settings.NotifyChallengeInvites = *settingsDto.NotifyChallengeInvites
		}
		if settingsDto.NotifyChallengeUpdates != nil {
			settings.NotifyChallengeUpdates = *settingsDto.NotifyChallengeUpdates
		}
		if settingsDto.NotifyChallengeReminders != nil {
			settings.NotifyChallengeReminders = *settingsDto.NotifyChallengeReminders
		}

		return tx.Save(&settings).Error
	})
}

// DeleteUser permanently deletes a user and all associated data.
// This function handles all user relationships in the correct order to maintain referential integrity.
//
// Relationships cleaned up:
// - Many-to-many: favorite sports, team memberships, challenge participations, friendships, blocks
// - Foreign keys: created challenges, created teams, invitations, notifications, messages, conversation participants, reports
// - Cascading deletes: emergency contacts, user settings (handled by database constraints)
//
// Important: This is a hard delete operation that cannot be undone.
func DeleteUser(user models.User, email string) error {
	userID := user.ID

	// Verify mail to make sure its the right user.
	if user.Email != email {
		return appError.ErrInvalidCredentials
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		// 1. Delete conversation participants
		if err := tx.Where("user_id = ?", userID).
			Delete(&models.ConversationParticipant{}).Error; err != nil {
			return err
		}

		// 2. Delete all messages sent by this user
		if err := tx.Where("sender_id = ?", userID).
			Delete(&models.Message{}).Error; err != nil {
			return err
		}

		// 3. Delete messages where user is recipient (legacy field)
		if err := tx.Where("recipient_id = ?", userID).
			Delete(&models.Message{}).Error; err != nil {
			return err
		}

		// 4. Delete reports created by this user
		if err := tx.Where("reporter_id = ?", userID).
			Delete(&models.Report{}).Error; err != nil {
			return err
		}

		// 5. Delete notifications where user is the recipient
		if err := tx.Where("user_id = ?", userID).
			Delete(&models.Notification{}).Error; err != nil {
			return err
		}

		// 6. Delete notifications where user is the actor (triggered by this user)
		if err := tx.Where("actor_id = ?", userID).
			Delete(&models.Notification{}).Error; err != nil {
			return err
		}

		// 7. Delete invitations sent or received by this user
		if err := tx.Where("inviter_id = ? OR invitee_id = ?", userID, userID).
			Delete(&models.Invitation{}).Error; err != nil {
			return err
		}

		// 8. Remove user from many-to-many: user_challenges (challenge participations)
		if err := tx.Exec("DELETE FROM user_challenges WHERE user_id = ?", userID).Error; err != nil {
			return err
		}

		// 9. Delete challenges created by this user
		// Get all challenges created by this user
		var challenges []models.Challenge
		if err := tx.Where("creator_id = ?", userID).Find(&challenges).Error; err != nil {
			return err
		}

		// For each challenge, clean up relationships before hard deleting
		for _, challenge := range challenges {
			// Delete challenge_teams relationships
			if err := tx.Exec("DELETE FROM challenge_teams WHERE challenge_id = ?", challenge.ID).Error; err != nil {
				return err
			}

			// Delete user_challenges relationships (already done above for the user, but clean up for other users)
			if err := tx.Exec("DELETE FROM user_challenges WHERE challenge_id = ?", challenge.ID).Error; err != nil {
				return err
			}

			// Hard delete the challenge itself (use Unscoped to bypass soft delete)
			if err := tx.Unscoped().Delete(&challenge).Error; err != nil {
				return err
			}
		}

		// 10. Remove user from many-to-many: user_teams (team memberships)
		if err := tx.Exec("DELETE FROM user_teams WHERE user_id = ?", userID).Error; err != nil {
			return err
		}

		// 11. Delete teams created by this user
		// First, get all teams created by this user
		var teams []models.Team
		if err := tx.Where("creator_id = ?", userID).Find(&teams).Error; err != nil {
			return err
		}

		// For each team, clean up relationships before deleting
		for _, team := range teams {
			// Delete challenge_teams relationships (teams in challenges)
			if err := tx.Exec("DELETE FROM challenge_teams WHERE team_id = ?", team.ID).Error; err != nil {
				return err
			}

			// Delete team_sports relationships
			if err := tx.Exec("DELETE FROM team_sports WHERE team_id = ?", team.ID).Error; err != nil {
				return err
			}

			// Delete conversations associated with this team
			// First, get conversation IDs for this team
			var conversationIDs []uint
			if err := tx.Model(&models.Conversation{}).
				Where("team_id = ?", team.ID).
				Pluck("id", &conversationIDs).Error; err != nil {
				return err
			}

			// Delete conversation participants for these conversations
			if len(conversationIDs) > 0 {
				if err := tx.Where("conversation_id IN ?", conversationIDs).
					Delete(&models.ConversationParticipant{}).Error; err != nil {
					return err
				}

				// Delete messages in these conversations
				if err := tx.Where("conversation_id IN ?", conversationIDs).
					Delete(&models.Message{}).Error; err != nil {
					return err
				}

				// Now delete the conversations themselves
				if err := tx.Where("id IN ?", conversationIDs).
					Delete(&models.Conversation{}).Error; err != nil {
					return err
				}
			}

			// Delete user_teams relationships (already done above for the user, but clean up for other users)
			if err := tx.Exec("DELETE FROM user_teams WHERE team_id = ?", team.ID).Error; err != nil {
				return err
			}

			// Finally, hard delete the team itself
			if err := tx.Unscoped().Delete(&team).Error; err != nil {
				return err
			}
		}

		// 12. Remove user from many-to-many: user_friends (friendships - bidirectional)
		// Remove where user is the primary user
		if err := tx.Exec("DELETE FROM user_friends WHERE user_id = ?", userID).Error; err != nil {
			return err
		}
		// Remove where user is the friend
		if err := tx.Exec("DELETE FROM user_friends WHERE friend_id = ?", userID).Error; err != nil {
			return err
		}

		// 13. Remove user from many-to-many: user_blocked_users (blocks - bidirectional)
		// Remove where user is the blocker
		if err := tx.Exec("DELETE FROM user_blocked_users WHERE user_id = ?", userID).Error; err != nil {
			return err
		}
		// Remove where user is the blocked user
		if err := tx.Exec("DELETE FROM user_blocked_users WHERE blocked_user_id = ?", userID).Error; err != nil {
			return err
		}

		// 14. Remove user from many-to-many: user_favorite_sports
		if err := tx.Exec("DELETE FROM user_favorite_sports WHERE user_id = ?", userID).Error; err != nil {
			return err
		}

		// 15. Delete emergency contacts (has CASCADE constraint, but explicit delete for clarity)
		if err := tx.Where("user_id = ?", userID).
			Delete(&models.EmergencyInfo{}).Error; err != nil {
			return err
		}

		// 16. Delete user settings (has CASCADE constraint, but explicit delete for clarity)
		if err := tx.Where("user_id = ?", userID).
			Delete(&models.UserSettings{}).Error; err != nil {
			return err
		}

		// 17. Finally, delete the user record itself
		if err := tx.Unscoped().Delete(&user).Error; err != nil {
			return err
		}

		return nil
	})
}
