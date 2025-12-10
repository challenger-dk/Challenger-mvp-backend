package models

type EmergencyInfo struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   `gorm:"not null"`
	Name         string `gorm:"not null"`
	PhoneNumber  string `gorm:"not null"`
	Relationship string `gorm:"not null"`
}
