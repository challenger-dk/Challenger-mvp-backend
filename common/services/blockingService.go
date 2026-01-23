package services

import (
	"server/common/appError"
	"server/common/config"
	"server/common/models"

	"gorm.io/gorm"
)

func GetBlockedUsers(userID uint) ([]models.User, error) {
	var user models.User
	// Only get basic info about blocked users
	err := config.DB.Preload("BlockedUsers", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "first_name", "last_name", "email", "created_at", "updated_at").Order("created_at DESC")
	}).First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	return user.BlockedUsers, nil
}

func BlockUser(userIdA uint, userIdB uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Ids must be different
		if userIdA == userIdB {
			return appError.ErrSameUser
		}

		var userA, userB models.User

		err := tx.First(&userA, userIdA).Error
		if err != nil {
			return err
		}

		err = tx.First(&userB, userIdB).Error
		if err != nil {
			return err
		}

		// 1. Remove Friendship (if exists)
		// Since blocking breaks the relationship, they should no longer be friends.
		// We try to delete the association in both directions.
		err = tx.Model(&userA).Association("Friends").Delete(&userB)
		if err != nil {
			return err
		}
		err = tx.Model(&userB).Association("Friends").Delete(&userA)
		if err != nil {
			return err
		}

		// 2. Add to BlockedUsers (Symmetric blocking)
		err = tx.Model(&userA).
			Association("BlockedUsers").
			Append(&userB)

		if err != nil {
			return err
		}

		err = tx.Model(&userB).
			Association("BlockedUsers").
			Append(&userA)

		if err != nil {
			return err
		}

		return nil
	})
}

func UnblockUser(userIdA uint, userIdB uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Ids must be different
		if userIdA == userIdB {
			return appError.ErrSameUser
		}

		var userA, userB models.User

		err := tx.First(&userA, userIdA).Error
		if err != nil {
			return err
		}

		err = tx.First(&userB, userIdB).Error
		if err != nil {
			return err
		}

		err = tx.Model(&userA).
			Association("BlockedUsers").
			Delete(&userB)

		if err != nil {
			return err
		}

		err = tx.Model(&userB).
			Association("BlockedUsers").
			Delete(&userA)

		if err != nil {
			return err
		}

		return nil
	})
}

// IsBlocked checks if 'blockerID' has blocked 'targetID'.
// Returns true if blocked.
func IsBlocked(blockerID, targetID uint) bool {
	var count int64
	// Check if blockerID has blocked targetID
	// Based on the User model: BlockedUsers []User `gorm:"many2many:user_blocked_users;joinForeignKey:UserID;JoinReferences:BlockedUserID"`
	// user_id is the blocker, blocked_user_id is the target.
	config.DB.Table("user_blocked_users").
		Where("user_id = ? AND blocked_user_id = ?", blockerID, targetID).
		Count(&count)

	return count > 0
}

// GetBlockedUserIDs returns all user IDs that are blocked by or have blocked the given user.
// This includes bidirectional blocking (if A blocks B, both A and B are in each other's blocked list).
func GetBlockedUserIDs(userID uint) []uint {
	var blockedIDs []uint

	// Get all users that are blocked in either direction
	// If userID blocked someone OR someone blocked userID
	config.DB.Table("user_blocked_users").
		Where("user_id = ? OR blocked_user_id = ?", userID, userID).
		Select("CASE WHEN user_id = ? THEN blocked_user_id ELSE user_id END as blocked_id", userID).
		Pluck("blocked_id", &blockedIDs)

	return blockedIDs
}

// ExcludeBlockedUsers returns a GORM scope that filters out blocked users.
// When used on the users table directly (e.g., in Preload), it filters by "id".
// When used on other tables, it filters by "user_id".
// Usage: db.Scopes(services.ExcludeBlockedUsers(currentUserID)).Find(&users)
func ExcludeBlockedUsers(userID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		blockedIDs := GetBlockedUserIDs(userID)
		if len(blockedIDs) == 0 {
			return db
		}

		// Check if we're querying the users table directly
		// In that case, filter by "id" instead of "user_id"
		stmt := db.Statement
		if stmt.Table == "users" || stmt.Table == "" {
			// When preloading users, filter by id
			return db.Where("id NOT IN ?", blockedIDs)
		}

		// For other tables, filter by user_id
		return db.Where("user_id NOT IN ?", blockedIDs)
	}
}

// ExcludeBlockedUsersOn returns a GORM scope that filters out content from blocked users
// based on a specific field name (e.g., "creator_id", "sender_id", etc.).
// Usage: db.Scopes(services.ExcludeBlockedUsersOn(currentUserID, "creator_id")).Find(&challenges)
func ExcludeBlockedUsersOn(userID uint, fieldName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		blockedIDs := GetBlockedUserIDs(userID)
		if len(blockedIDs) > 0 {
			return db.Where(fieldName+" NOT IN ?", blockedIDs)
		}
		return db
	}
}

// ExcludeBlockedUsersMultipleFields returns a GORM scope that filters out content from blocked users
// based on multiple field names. Content is excluded if ANY of the fields contain a blocked user.
// Usage: db.Scopes(services.ExcludeBlockedUsersMultipleFields(currentUserID, "creator_id", "owner_id")).Find(&items)
func ExcludeBlockedUsersMultipleFields(userID uint, fieldNames ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		blockedIDs := GetBlockedUserIDs(userID)
		if len(blockedIDs) == 0 {
			return db
		}

		// Build OR conditions for each field
		for i, fieldName := range fieldNames {
			if i == 0 {
				db = db.Where(fieldName+" NOT IN ?", blockedIDs)
			} else {
				db = db.Where(fieldName+" NOT IN ?", blockedIDs)
			}
		}
		return db
	}
}
