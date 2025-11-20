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

	// Validate configuration before returning
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
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

// Validate checks that critical configuration values are properly set
// This prevents running in production with insecure default values
func (c *Config) Validate() error {
	// Check if running in production mode
	isProduction := c.Server.GinMode == "release" || os.Getenv("APP_ENV") == "production"

	// In production, enforce strict security requirements
	if isProduction {
		// JWT Secret must not be the default value
		defaultJWTSecret := "your-super-secret-jwt-key-change-this-in-production"
		if c.JWT.Secret == defaultJWTSecret {
			return fmt.Errorf("JWT_SECRET must be changed from default value in production")
		}

		// JWT Secret must be at least 32 characters (256 bits)
		if len(c.JWT.Secret) < 32 {
			return fmt.Errorf("JWT_SECRET must be at least 32 characters long in production (current: %d)", len(c.JWT.Secret))
		}

		// Database password must not be the default value
		if c.Database.Password == "postgres" {
			return fmt.Errorf("DB_PASSWORD must be changed from default value in production")
		}

		// Database password must be strong (at least 12 characters)
		if len(c.Database.Password) < 12 {
			return fmt.Errorf("DB_PASSWORD must be at least 12 characters long in production (current: %d)", len(c.Database.Password))
		}

		// SSL must be enabled in production
		if c.Database.SSLMode == "disable" {
			log.Println("WARNING: Database SSL is disabled in production. This is not recommended.")
		}

		// CORS origins must not include localhost
		for _, origin := range c.CORS.AllowedOrigins {
			if origin == "http://localhost:3000" || origin == "http://localhost:8080" {
				return fmt.Errorf("ALLOWED_ORIGINS must not include localhost URLs in production (found: %s)", origin)
			}
		}
	}

	// Always validate that critical values are not empty
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET cannot be empty")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST cannot be empty")
	}

	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME cannot be empty")
	}

	if c.JWT.ExpirationHours <= 0 {
		return fmt.Errorf("JWT_EXPIRATION_HOURS must be positive (current: %d)", c.JWT.ExpirationHours)
	}

	// Log configuration status
	if isProduction {
		log.Println("✓ Configuration validation passed (PRODUCTION mode)")
	} else {
		log.Println("✓ Configuration validation passed (DEVELOPMENT mode)")
	}

	return nil
}

// getEnv mengambil nilai environment variable atau menggunakan default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
