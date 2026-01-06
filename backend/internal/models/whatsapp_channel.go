package models

import "gorm.io/gorm"

type WhatsAppChannel struct {
	gorm.Model
	Name   string `gorm:"not null"`
	Active bool   `gorm:"not null;default:true"`

	// sessão/estado (por enquanto só placeholder)
	Status string `gorm:"not null;default:'disconnected'"` // disconnected, connecting, connected
}
