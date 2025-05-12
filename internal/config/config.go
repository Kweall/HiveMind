package config

import "os"

type Config struct {
	StorageType string
	DatabaseURL string
}

func Load() *Config {
	return &Config{
		StorageType: getEnv("STORAGE_TYPE", "memory"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
