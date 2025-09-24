package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	AppEnv      string
	JWTSecret   string
	JWTExpires  string
}

func Load() *Config {
	// Load .env file if present
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found, falling back to system env vars")
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		AppEnv:      os.Getenv("APP_ENV"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTExpires:  os.Getenv("JWT_EXPIRES_HOURS"),
	}

	// Fail fast if DATABASE_URL is missing
	if cfg.DatabaseURL == "" {
		log.Fatal("❌ missing required env var: DATABASE_URL")
	}

	// Default port if empty
	if cfg.Port == "" {
		cfg.Port = "3000"
	}

	return cfg
}
