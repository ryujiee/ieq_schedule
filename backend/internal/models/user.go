package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"size:120;uniqueIndex;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	Name         string    `gorm:"size:120"`
	Active       bool      `gorm:"not null;default:true"`

	RoleID uint `gorm:"not null;index"`
	Role   Role `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
