package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config menampung seluruh konfigurasi aplikasi dari environment variable
type Config struct {
	AppPort    string
	AppEnv     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// LoadConfig memuat file .env (jika ada) dan membaca nilai environment variable
func LoadConfig() *Config {
	// Memuat file .env jika tersedia (tidak masalah jika tidak ada di environment server produksi)
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	return &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "pajak_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// GetDSN menyusun Connection String (Data Source Name) untuk PostgreSQL
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// getEnv adalah fungsi pembantu untuk membaca env dengan nilai default (fallback)
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

