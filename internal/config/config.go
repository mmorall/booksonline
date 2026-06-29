package config

import (
	"errors"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DatabaseURL string
	AdminUser   string
	AdminPass   string
}

func Load() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, reading configuration from environment")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL environment variable is required")
	}

	user := os.Getenv("ADMIN_USER")
	if user == "" {
		return nil, errors.New("ADMIN_USER environment variable is required")
	}

	pass := os.Getenv("ADMIN_PASS")
	if pass == "" {
		return nil, errors.New("ADMIN_PASS environment variable is required")
	}

	return &AppConfig{
		DatabaseURL: dbURL,
		AdminUser:   user,
		AdminPass:   pass,
	}, nil
}
