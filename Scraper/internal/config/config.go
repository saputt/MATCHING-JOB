package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
	AppPort     string
	Headless    bool
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		DatabaseUrl: getEnv("DATABASE_URL", ""),
		AppPort:     getEnv("PORT", "8000"),
		Headless:    getEnvBool("HEADLESS", false),
	}
}
