package models

import "time"

// tabela pivot (member_functions)
type MemberFunction struct {
	MemberID       uint      `gorm:"primaryKey;autoIncrement:false"`
	TeamFunctionID uint      `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}
