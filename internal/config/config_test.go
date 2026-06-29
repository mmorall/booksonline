package config_test

import (
	"testing"

	"github.com/mmorall/booksonline/internal/config"
)

func TestLoad_Success(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/db")
	t.Setenv("ADMIN_USER", "admin")
	t.Setenv("ADMIN_PASS", "secret")

	cfg, err := config.Load()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.AdminUser != "admin" {
		t.Errorf("expected admin user 'admin', got '%s'", cfg.AdminUser)
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	// Intentionally omit DATABASE_URL
	t.Setenv("ADMIN_USER", "admin")
	t.Setenv("ADMIN_PASS", "secret")

	_, err := config.Load()

	if err == nil {
		t.Fatal("expected error due to missing DATABASE_URL, got nil")
	}
}

func TestLoad_MissingCredentials(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost/db")
	// Omit credentials

	_, err := config.Load()

	if err == nil {
		t.Fatal("expected error due to missing credentials, got nil")
	}
}
