package models

import "gorm.io/gorm"

type Settings struct {
	gorm.Model

	// lembrete diário
	ReminderEnabled  bool   `gorm:"not null;default:false"`
	ReminderTimeHHMM string `gorm:"not null;default:'08:00'"` // "HH:MM"
	ReminderTemplate string `gorm:"type:text;not null;default:'Olá {Nome}; gostaria de lembrar que hoje é o seu dia de trabalhar e servir na igreja na função {Funcao}.'"`

	// qual canal usar
	ReminderChannelID *uint
}

func DefaultSettings() *Settings {
	return &Settings{
		ReminderEnabled:  false,
		ReminderTimeHHMM: "08:00",
		ReminderTemplate: "Olá {Nome}; gostaria de lembrar que hoje é o seu dia de trabalhar e servir na igreja na função {Funcao}.",
	}
}
