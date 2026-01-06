package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string

	AppPort string
	AppEnv  string

	JWTSecret string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPass:     os.Getenv("DB_PASS"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		AppPort:    getenvDefault("APP_PORT", "8080"),
		AppEnv:     getenvDefault("APP_ENV", "dev"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}
}

func getenvDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
