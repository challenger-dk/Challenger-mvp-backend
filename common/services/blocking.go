package services

import (
	"server/common/appError"
	"server/common/config"
	"server/common/models"

	"gorm.io/gorm"
)

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
