package integration

import (
	"server/common/config"
	"server/common/models"
	"server/common/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestEmergencyInfoService_CreateEmergencyContact(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// 1. Create a user
	user, err := services.CreateUser(models.User{
		Email:     "emergency@test.com",
		FirstName: "Emergency",
		LastName:  "User",
	}, "password123")
	assert.NoError(t, err)

	// 2. Create emergency contact successfully
	emergencyInfo := models.EmergencyInfo{
		Name:         "John Doe",
		PhoneNumber:  "+1234567890",
		Relationship: "Father",
	}
	err = services.CreateEmergencyContact(*user, emergencyInfo)
	assert.NoError(t, err)

	// 3. Verify emergency contact was created
	var createdContact models.EmergencyInfo
	err = config.DB.Where("user_id = ? AND name = ?", user.ID, "John Doe").First(&createdContact).Error
	assert.NoError(t, err)
	assert.Equal(t, user.ID, createdContact.UserID)
	assert.Equal(t, "John Doe", createdContact.Name)
	assert.Equal(t, "+1234567890", createdContact.PhoneNumber)
	assert.Equal(t, "Father", createdContact.Relationship)
	assert.NotZero(t, createdContact.ID)

	// 4. Create multiple emergency contacts for same user
	emergencyInfo2 := models.EmergencyInfo{
		Name:         "Jane Doe",
		PhoneNumber:  "+0987654321",
		Relationship: "Mother",
	}
	err = services.CreateEmergencyContact(*user, emergencyInfo2)
	assert.NoError(t, err)

	// 5. Verify both contacts exist
	var contacts []models.EmergencyInfo
	err = config.DB.Where("user_id = ?", user.ID).Find(&contacts).Error
	assert.NoError(t, err)
	assert.Len(t, contacts, 2)

	// 6. Try to create contact for non-existent user
	nonExistentUser := models.User{ID: 99999}
	err = services.CreateEmergencyContact(nonExistentUser, emergencyInfo)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestEmergencyInfoService_DeleteEmergencyContact(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// 1. Create a user
	user, err := services.CreateUser(models.User{
		Email:     "delete@test.com",
		FirstName: "Delete",
		LastName:  "User",
	}, "password123")
	assert.NoError(t, err)

	// 2. Create emergency contact
	emergencyInfo := models.EmergencyInfo{
		Name:         "Contact To Delete",
		PhoneNumber:  "+1111111111",
		Relationship: "Friend",
	}
	err = services.CreateEmergencyContact(*user, emergencyInfo)
	assert.NoError(t, err)

	// 3. Get the created contact ID
	var createdContact models.EmergencyInfo
	err = config.DB.Where("user_id = ? AND name = ?", user.ID, "Contact To Delete").First(&createdContact).Error
	assert.NoError(t, err)
	contactID := createdContact.ID

	// 4. Delete emergency contact successfully
	err = services.DeleteEmergencyContact(*user, contactID)
	assert.NoError(t, err)

	// 5. Verify contact was deleted
	var deletedContact models.EmergencyInfo
	err = config.DB.First(&deletedContact, contactID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// 6. Create another user and contact
	user2, err := services.CreateUser(models.User{
		Email:     "user2@test.com",
		FirstName: "User",
		LastName:  "Two",
	}, "password123")
	assert.NoError(t, err)

	emergencyInfo2 := models.EmergencyInfo{
		Name:         "User2 Contact",
		PhoneNumber:  "+2222222222",
		Relationship: "Sibling",
	}
	err = services.CreateEmergencyContact(*user2, emergencyInfo2)
	assert.NoError(t, err)

	var user2Contact models.EmergencyInfo
	err = config.DB.Where("user_id = ?", user2.ID).First(&user2Contact).Error
	assert.NoError(t, err)

	// 7. Try to delete another user's contact (should fail)
	err = services.DeleteEmergencyContact(*user, user2Contact.ID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// 8. Verify user2's contact still exists
	var stillExists models.EmergencyInfo
	err = config.DB.First(&stillExists, user2Contact.ID).Error
	assert.NoError(t, err)

	// 9. Try to delete non-existent contact
	err = services.DeleteEmergencyContact(*user, 99999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestEmergencyInfoService_UpdateEmergencyContact(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// 1. Create a user
	user, err := services.CreateUser(models.User{
		Email:     "update@test.com",
		FirstName: "Update",
		LastName:  "User",
	}, "password123")
	assert.NoError(t, err)

	// 2. Create emergency contact
	originalContact := models.EmergencyInfo{
		Name:         "Original Name",
		PhoneNumber:  "+3333333333",
		Relationship: "Original Relationship",
	}
	err = services.CreateEmergencyContact(*user, originalContact)
	assert.NoError(t, err)

	// 3. Get the created contact ID
	var createdContact models.EmergencyInfo
	err = config.DB.Where("user_id = ? AND name = ?", user.ID, "Original Name").First(&createdContact).Error
	assert.NoError(t, err)
	contactID := createdContact.ID

	// 4. Update emergency contact successfully
	updatedContact := models.EmergencyInfo{
		Name:         "Updated Name",
		PhoneNumber:  "+4444444444",
		Relationship: "Updated Relationship",
	}
	err = services.UpdateEmergencyContact(*user, updatedContact, contactID)
	assert.NoError(t, err)

	// 5. Verify contact was updated
	var updated models.EmergencyInfo
	err = config.DB.First(&updated, contactID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, "+4444444444", updated.PhoneNumber)
	assert.Equal(t, "Updated Relationship", updated.Relationship)
	assert.Equal(t, user.ID, updated.UserID)
	assert.Equal(t, contactID, updated.ID) // ID should not change

	// 6. Create another user and contact
	user2, err := services.CreateUser(models.User{
		Email:     "user2update@test.com",
		FirstName: "User",
		LastName:  "Two",
	}, "password123")
	assert.NoError(t, err)

	user2Contact := models.EmergencyInfo{
		Name:         "User2 Contact",
		PhoneNumber:  "+5555555555",
		Relationship: "Colleague",
	}
	err = services.CreateEmergencyContact(*user2, user2Contact)
	assert.NoError(t, err)

	var user2Created models.EmergencyInfo
	err = config.DB.Where("user_id = ?", user2.ID).First(&user2Created).Error
	assert.NoError(t, err)

	// 7. Try to update another user's contact (should fail)
	unauthorizedUpdate := models.EmergencyInfo{
		Name:         "Hacked Name",
		PhoneNumber:  "+6666666666",
		Relationship: "Hacked",
	}
	err = services.UpdateEmergencyContact(*user, unauthorizedUpdate, user2Created.ID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// 8. Verify user2's contact was not modified
	var unchanged models.EmergencyInfo
	err = config.DB.First(&unchanged, user2Created.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "User2 Contact", unchanged.Name)
	assert.Equal(t, "+5555555555", unchanged.PhoneNumber)
	assert.Equal(t, "Colleague", unchanged.Relationship)

	// 9. Try to update non-existent contact
	nonExistentUpdate := models.EmergencyInfo{
		Name:         "Non Existent",
		PhoneNumber:  "+7777777777",
		Relationship: "None",
	}
	err = services.UpdateEmergencyContact(*user, nonExistentUpdate, 99999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// 10. Update with partial changes (all fields should be updated)
	partialUpdate := models.EmergencyInfo{
		Name:         "Only Name Changed",
		PhoneNumber:  "+8888888888",
		Relationship: "Only Relationship Changed",
	}
	err = services.UpdateEmergencyContact(*user, partialUpdate, contactID)
	assert.NoError(t, err)

	var finalContact models.EmergencyInfo
	err = config.DB.First(&finalContact, contactID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Only Name Changed", finalContact.Name)
	assert.Equal(t, "+8888888888", finalContact.PhoneNumber)
	assert.Equal(t, "Only Relationship Changed", finalContact.Relationship)
}

func TestEmergencyInfoService_MultipleContactsPerUser(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// 1. Create a user
	user, err := services.CreateUser(models.User{
		Email:     "multiple@test.com",
		FirstName: "Multiple",
		LastName:  "Contacts",
	}, "password123")
	assert.NoError(t, err)

	// 2. Create multiple emergency contacts
	contacts := []models.EmergencyInfo{
		{Name: "Contact 1", PhoneNumber: "+1111111111", Relationship: "Father"},
		{Name: "Contact 2", PhoneNumber: "+2222222222", Relationship: "Mother"},
		{Name: "Contact 3", PhoneNumber: "+3333333333", Relationship: "Sibling"},
	}

	for _, contact := range contacts {
		err = services.CreateEmergencyContact(*user, contact)
		assert.NoError(t, err)
	}

	// 3. Verify all contacts exist
	var allContacts []models.EmergencyInfo
	err = config.DB.Where("user_id = ?", user.ID).Find(&allContacts).Error
	assert.NoError(t, err)
	assert.Len(t, allContacts, 3)

	// 4. Update one contact
	var contactToUpdate models.EmergencyInfo
	err = config.DB.Where("user_id = ? AND name = ?", user.ID, "Contact 1").First(&contactToUpdate).Error
	assert.NoError(t, err)

	updated := models.EmergencyInfo{
		Name:         "Updated Contact 1",
		PhoneNumber:  "+9999999999",
		Relationship: "Updated Father",
	}
	err = services.UpdateEmergencyContact(*user, updated, contactToUpdate.ID)
	assert.NoError(t, err)

	// 5. Delete one contact
	var contactToDelete models.EmergencyInfo
	err = config.DB.Where("user_id = ? AND name = ?", user.ID, "Contact 2").First(&contactToDelete).Error
	assert.NoError(t, err)

	err = services.DeleteEmergencyContact(*user, contactToDelete.ID)
	assert.NoError(t, err)

	// 6. Verify remaining contacts
	var remainingContacts []models.EmergencyInfo
	err = config.DB.Where("user_id = ?", user.ID).Find(&remainingContacts).Error
	assert.NoError(t, err)
	assert.Len(t, remainingContacts, 2)

	// Verify updated contact exists
	var updatedContact models.EmergencyInfo
	err = config.DB.First(&updatedContact, contactToUpdate.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Updated Contact 1", updatedContact.Name)

	// Verify deleted contact is gone
	var deletedContact models.EmergencyInfo
	err = config.DB.First(&deletedContact, contactToDelete.ID).Error
	assert.Error(t, err)
}
