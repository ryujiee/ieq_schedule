package models

import (
	"time"

	"gorm.io/gorm"
)

type Member struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"size:120;not null;index"`
	Phone  string `gorm:"size:30;index"`
	Active bool   `gorm:"not null;default:true"`

	// N:N com TeamFunction via tabela member_functions
	Functions []TeamFunction `gorm:"many2many:member_functions;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
