package db

import (
	"errors"

	"ieq/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SeedConfig struct {
	AdminEmail    string
	AdminPassword string
	AdminName     string
}

func AutoMigrateAndSeed(database *gorm.DB, cfg SeedConfig) error {
	if database == nil {
		return errors.New("database is nil")
	}

	if err := database.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.TeamFunction{},

		&models.Member{},
		&models.MemberFunction{},

		&models.ScheduleAssignment{},
		&models.Settings{},
		&models.WhatsAppChannel{},
	); err != nil {
		return err
	}

	adminRole, err := ensureAdminRole(database)
	if err != nil {
		return err
	}

	var count int64
	if err := database.Model(&models.User{}).Where("role_id = ?", adminRole.ID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		Email:        cfg.AdminEmail,
		PasswordHash: string(hash),
		Name:         cfg.AdminName,
		Active:       true,
		RoleID:       adminRole.ID,
	}

	return database.Create(&admin).Error
}

func ensureAdminRole(database *gorm.DB) (*models.Role, error) {
	var role models.Role
	err := database.Where("slug = ?", "admin").First(&role).Error
	if err == nil {
		return &role, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role = models.Role{
		Slug: "admin",
		Name: "Administrador",
	}
	if err := database.Create(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
