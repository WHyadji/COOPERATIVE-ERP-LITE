package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config menyimpan konfigurasi aplikasi
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

// DatabaseConfig berisi konfigurasi database
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ServerConfig berisi konfigurasi server
type ServerConfig struct {
	Port    string
	GinMode string
}

// JWTConfig berisi konfigurasi JWT
type JWTConfig struct {
	Secret           string
	ExpirationHours  int
}

// CORSConfig berisi konfigurasi CORS
type CORSConfig struct {
	AllowedOrigins []string
}

// LoadConfig memuat konfigurasi dari environment variables
func LoadConfig() (*Config, error) {
	// Load .env file jika ada (untuk development)
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	expirationHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		expirationHours = 24
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "koperasi_erp"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
			ExpirationHours: expirationHours,
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{
				getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
			},
		},
	}

	return config, nil
}

// GetDSN mengembalikan connection string untuk database
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// getEnv mengambil nilai environment variable atau menggunakan default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
