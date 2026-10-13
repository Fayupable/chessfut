package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DatabaseURL       string
	RedisAddr         string
	RedisPassword     string
	AdminAPIKey       string
	ChessComUserAgent string
	ChessComBaseURL   string
	MigrationsPath    string
	TrustedProxies    []string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		AdminAPIKey:       getEnv("ADMIN_API_KEY", ""),
		ChessComUserAgent: mustGetEnv("CHESSCOM_USER_AGENT"),
		ChessComBaseURL:   mustGetEnv("CHESSCOM_BASE_URL"),
		MigrationsPath:    getEnv("MIGRATIONS_PATH", "db/migration"),
		TrustedProxies:    splitCSV(getEnv("TRUSTED_PROXY_IPS", "")),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required environment variable: %s", key)
	}
	return v
}

func splitCSV(v string) []string {
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
