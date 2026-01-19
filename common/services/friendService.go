package services

import (
	"server/common/appError"
	"server/common/config"
	"server/common/models"

	"gorm.io/gorm"
)

func GetFriends(userID uint) ([]models.User, error) {
	var user models.User

	err := config.DB.Preload("Friends").
		First(&user, userID).
		Error

	if err != nil {
		return nil, err
	}

	return user.Friends, nil
}

// tries to find suggested friends for a user based on common friends, teams, challenges, and favorite sports
func GetSuggestedFriends(userID uint) ([]models.User, error) {
	// Get the current user with their friends, teams, challenges, and sports
	var user models.User
	err := config.DB.Preload("Friends").
		Preload("Teams").
		Preload("JoinedChallenges").
		Preload("FavoriteSports").
		First(&user, userID).Error

	if err != nil {
		return nil, err
	}

	// Get IDs of current friends to exclude them
	friendIDs := make([]uint, len(user.Friends))
	for i, friend := range user.Friends {
		friendIDs[i] = friend.ID
	}

	// Get IDs of blocked users to exclude them
	var blockedUserIDs []uint
	err = config.DB.Table("user_blocked_users").
		Select("blocked_user_id").
		Where("user_id = ?", userID).
		Pluck("blocked_user_id", &blockedUserIDs).Error

	if err != nil {
		return nil, err
	}

	// Also get users who have blocked the current user
	var blockingUserIDs []uint
	err = config.DB.Table("user_blocked_users").
		Select("user_id").
		Where("blocked_user_id = ?", userID).
		Pluck("user_id", &blockingUserIDs).Error

	if err != nil {
		return nil, err
	}

	// Combine all blocked user IDs
	allBlockedIDs := append(blockedUserIDs, blockingUserIDs...)

	// Get all potential candidates (users who are not friends and not blocked and not invited by or to the user)
	var candidates []models.User
	query := config.DB.Preload("Friends").
		Preload("Teams").
		Preload("JoinedChallenges").
		Preload("FavoriteSports").
		Where("id != ?", userID). // Exclude self
		// Exclude users who have pending friend invitations with the current user
		Where("id NOT IN (SELECT invitee_id FROM invitations WHERE inviter_id = ? AND resource_type = 'friend' AND status = 'pending')", userID). // Exclude users invited by current user
		Where("id NOT IN (SELECT inviter_id FROM invitations WHERE invitee_id = ? AND resource_type = 'friend' AND status = 'pending')", userID)  // Exclude users who invited current user

	// Exclude friends
	if len(friendIDs) > 0 {
		query = query.Where("id NOT IN ?", friendIDs)
	}

	// Exclude blocked users
	if len(allBlockedIDs) > 0 {
		query = query.Where("id NOT IN ?", allBlockedIDs)
	}

	err = query.Find(&candidates).Error
	if err != nil {
		return nil, err
	}

	// Score each candidate
	type scoredUser struct {
		user                  models.User
		commonFriendsCount    int
		commonTeamsCount      int
		commonChallengesCount int
		commonSportsCount     int
		totalScore            float64
	}

	scoredUsers := make([]scoredUser, 0, len(candidates))

	// Create maps for efficient lookup
	userFriendMap := make(map[uint]bool)
	for _, friend := range user.Friends {
		userFriendMap[friend.ID] = true
	}

	userTeamMap := make(map[uint]bool)
	for _, team := range user.Teams {
		userTeamMap[team.ID] = true
	}

	userChallengeMap := make(map[uint]bool)
	for _, challenge := range user.JoinedChallenges {
		userChallengeMap[challenge.ID] = true
	}

	userSportMap := make(map[uint]bool)
	for _, sport := range user.FavoriteSports {
		userSportMap[sport.ID] = true
	}

	// Calculate scores for each candidate
	for _, candidate := range candidates {
		scored := scoredUser{user: candidate}

		// Count common friends
		for _, candidateFriend := range candidate.Friends {
			if userFriendMap[candidateFriend.ID] {
				scored.commonFriendsCount++
			}
		}

		// Count common teams
		for _, candidateTeam := range candidate.Teams {
			if userTeamMap[candidateTeam.ID] {
				scored.commonTeamsCount++
			}
		}

		// Count common challenges
		for _, candidateChallenge := range candidate.JoinedChallenges {
			if userChallengeMap[candidateChallenge.ID] {
				scored.commonChallengesCount++
			}
		}

		// Count common sports
		for _, candidateSport := range candidate.FavoriteSports {
			if userSportMap[candidateSport.ID] {
				scored.commonSportsCount++
			}
		}

		// Calculate weighted score
		// Weights: Common Friends (4.0) > Common Teams (3.0) > Common Challenges (2.0) > Common Sports (1.0)
		scored.totalScore = float64(scored.commonFriendsCount)*4.0 +
			float64(scored.commonTeamsCount)*3.0 +
			float64(scored.commonChallengesCount)*2.0 +
			float64(scored.commonSportsCount)*1.0

		// Only include users with at least some connection
		if scored.totalScore > 0 {
			scoredUsers = append(scoredUsers, scored)
		}
	}

	// Sort by score (descending)
	for i := 0; i < len(scoredUsers); i++ {
		for j := i + 1; j < len(scoredUsers); j++ {
			if scoredUsers[j].totalScore > scoredUsers[i].totalScore {
				scoredUsers[i], scoredUsers[j] = scoredUsers[j], scoredUsers[i]
			}
		}
	}

	// Return top 10 suggestions
	maxSuggestions := 10
	if len(scoredUsers) < maxSuggestions {
		maxSuggestions = len(scoredUsers)
	}

	suggestions := make([]models.User, maxSuggestions)
	for i := 0; i < maxSuggestions; i++ {
		suggestions[i] = scoredUsers[i].user
	}

	return suggestions, nil
}

// DeleteFriendship removes both users from each other's friends list
func RemoveFriend(userIdA uint, userIdB uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {

		// Ids must be different
		if userIdA == userIdB {
			return appError.ErrInvalidFriendship
		}

		var userA, userB models.User

		if err := tx.First(&userA, userIdA).Error; err != nil {
			return err
		}

		if err := tx.First(&userB, userIdB).Error; err != nil {
			return err
		}

		err := tx.Model(&userA).
			Association("Friends").
			Delete(&userB)

		if err != nil {
			return err
		}

		err = tx.Model(&userB).
			Association("Friends").
			Delete(&userA)

		if err != nil {
			return err
		}

		// No notification here

		return nil
	})
}

// Package private
// createFriendship adds both users to each other's friends list
func createFriendship(userIdA uint, userIdB uint, db *gorm.DB) error {
	var userA, userB models.User

	if err := db.First(&userA, userIdA).Error; err != nil {
		return err
	}
	if err := db.First(&userB, userIdB).Error; err != nil {
		return err
	}

	err := db.Model(&userA).
		Association("Friends").
		Append(&userB)

	if err != nil {
		return err
	}

	err = db.Model(&userB).
		Association("Friends").
		Append(&userA)

	if err != nil {
		return err
	}

	return nil
}
