package models

import "time"

type Role struct {
	ID        uint   `gorm:"primaryKey"`
	Slug      string `gorm:"size:50;uniqueIndex;not null"` // ex: "admin"
	Name      string `gorm:"size:80;not null"`             // ex: "Administrador"
	CreatedAt time.Time
	UpdatedAt time.Time
}
