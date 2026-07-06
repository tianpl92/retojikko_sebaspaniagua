package config

import (
	"os"
)

type Config struct {
	DBUser         string
	DBPassword     string
	DBHost         string
	DBPort         string
	DBName         string
	Port           string
	SecretKey      string
	IntegrationURL string
}

func Load() *Config {
	return &Config{
		DBUser:         getEnv("DB_USER", "sebasdb"),
		DBPassword:     getEnv("DB_PASSWORD", "challenge_app_07032026"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBName:         getEnv("DB_NAME", "portal_plan_public_app"),
		Port:           getEnv("PORT", "8080"),
		SecretKey:      getEnv("SECRET_KEY", "challenge-secret-key-2026"),
		IntegrationURL: getEnv("INTEGRATION_URL", "https://datos.gov.co/resource/p6dx-8zbt.json"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
