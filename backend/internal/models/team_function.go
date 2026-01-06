package models

import (
	"time"

	"gorm.io/gorm"
)

type TeamFunction struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"size:120;not null;uniqueIndex"`
	Active    bool           `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
