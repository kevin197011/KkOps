package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Salt     SaltConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	Expiration int // hours
}

type SaltConfig struct {
	APIURL      string `json:"api_url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	EAuth       string `json:"eauth"` // eauth type: pam, ldap, etc.
	Timeout     int    `json:"timeout"` // seconds
	VerifySSL   bool   `json:"verify_ssl"`
}

var AppConfig *Config

func Load() error {
	// Load .env file if exists
	_ = godotenv.Load()

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "kkops"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiration: 24, // 24 hours
		},
		Salt: SaltConfig{
			APIURL:    getEnv("SALT_API_URL", "http://localhost:8000"),
			Username:  getEnv("SALT_USERNAME", "salt"),
			Password:  getEnv("SALT_PASSWORD", ""),
			EAuth:     getEnv("SALT_EAUTH", "pam"),
			Timeout:   30, // 30 seconds
			VerifySSL: getEnv("SALT_VERIFY_SSL", "true") == "true",
		},
	}

	return nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

