package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Unset any env vars that might interfere
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("PORT")
	os.Unsetenv("SECRET_KEY")

	cfg := Load()

	if cfg.DBUser != "sebasdb" {
		t.Errorf("expected default DB_USER 'sebasdb', got %q", cfg.DBUser)
	}
	if cfg.DBPassword != "challenge_app_07032026" {
		t.Errorf("expected default DB_PASSWORD 'challenge_app_07032026', got %q", cfg.DBPassword)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("expected default DB_HOST 'localhost', got %q", cfg.DBHost)
	}
	if cfg.DBPort != "5432" {
		t.Errorf("expected default DB_PORT '5432', got %q", cfg.DBPort)
	}
	if cfg.DBName != "portal_plan_public_app" {
		t.Errorf("expected default DB_NAME 'portal_plan_public_app', got %q", cfg.DBName)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected default PORT '8080', got %q", cfg.Port)
	}
	if cfg.SecretKey != "challenge-secret-key-2026" {
		t.Errorf("expected default SECRET_KEY 'challenge-secret-key-2026', got %q", cfg.SecretKey)
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("DB_USER", "custom_user")
	os.Setenv("DB_PASSWORD", "custom_pass")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "9999")
	os.Setenv("DB_NAME", "custom_db")
	os.Setenv("PORT", "3000")
	os.Setenv("SECRET_KEY", "custom-secret")
	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("PORT")
		os.Unsetenv("SECRET_KEY")
	}()

	cfg := Load()

	if cfg.DBUser != "custom_user" {
		t.Errorf("expected DB_USER 'custom_user', got %q", cfg.DBUser)
	}
	if cfg.DBPassword != "custom_pass" {
		t.Errorf("expected DB_PASSWORD 'custom_pass', got %q", cfg.DBPassword)
	}
	if cfg.Port != "3000" {
		t.Errorf("expected PORT '3000', got %q", cfg.Port)
	}
}
