package models

import (
	"time"

	"gorm.io/gorm"
)

type ScheduleAssignment struct {
	ID uint `gorm:"primaryKey"`

	// Data SEM hora (date)
	Date time.Time `gorm:"type:date;not null;index;uniqueIndex:idx_date_function"`

	TeamFunctionID uint         `gorm:"not null;uniqueIndex:idx_date_function"`
	TeamFunction   TeamFunction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	MemberID uint   `gorm:"not null"`
	Member   Member `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
