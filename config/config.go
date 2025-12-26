package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppPort       string
	WorkStartHour int
	WorkEndHour   int
	SafeNetwork   string
	DBName        string
	JWTSecret     string
}

// LoadConfig membaca .env dan memvalidasi isinya
func LoadConfig() *AppConfig {
	// Coba load file .env, abaikan error jika file tidak ada (misal di production pakai environment asli)
	_ = godotenv.Load()

	return &AppConfig{
		AppPort:       getEnv("APP_PORT", ":8080"),
		WorkStartHour: getEnvAsInt("WORK_START_HOUR", 8),
		WorkEndHour:   getEnvAsInt("WORK_END_HOUR", 17),
		SafeNetwork:   getEnv("SAFE_NETWORK_PREFIX", "127.0.0.1"),
		DBName:        getEnv("DB_NAME", "audit.db"),
		JWTSecret:     getEnv("JWT_SECRET", "default_insecure_secret"),
	}
}

// Helper untuk membaca string env dengan fallback default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper untuk membaca int env dengan fallback default
func getEnvAsInt(key string, fallback int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
		log.Printf("⚠️ Warning: Config %s tidak valid, menggunakan default %d", key, fallback)
	}
	return fallback
}