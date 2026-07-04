package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DatabaseURL       string
	RedisAddr         string
	RedisPassword     string
	AdminAPIKey       string
	ChessComUserAgent string
	MigrationsPath    string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		AdminAPIKey:       getEnv("ADMIN_API_KEY", ""),
		ChessComUserAgent: getEnv("CHESSCOM_USER_AGENT", "chessfut/1.0"),
		MigrationsPath:    getEnv("MIGRATIONS_PATH", "db/migration"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
