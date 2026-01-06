package services

import (
	"errors"
	"strings"

	"ieq/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	DB  *gorm.DB
	JWT *JWTService
}

func NewAuthService(db *gorm.DB, jwtSvc *JWTService) *AuthService {
	return &AuthService{DB: db, JWT: jwtSvc}
}

func (s *AuthService) Login(email, password string) (string, *models.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return "", nil, ErrInvalidCredentials
	}

	var user models.User
	err := s.DB.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	if !user.Active {
		return "", nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.JWT.Sign(user.ID, user.Role.Slug)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}
