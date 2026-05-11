// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	SSO        SSOConfig        `mapstructure:"sso"`
	IdP        IdPConfig        `mapstructure:"idp"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
	Log        LogConfig        `mapstructure:"log"`
	Alerts     AlertsConfig     `mapstructure:"alerts"`
}

// AlertsConfig configures alert ingestion (webhook).
type AlertsConfig struct {
	WebhookSecret string `mapstructure:"webhook_secret"` // header X-KkOps-Alert-Secret
}

// IdPConfig holds KkOps-as-IdP (OIDC provider) settings
type IdPConfig struct {
	Enabled       bool          `mapstructure:"enabled"`
	Issuer        string        `mapstructure:"issuer"`         // e.g. http://localhost:3000/oidc (public URL of this IdP)
	SessionSecret string        `mapstructure:"session_secret"` // for cookie session (recommend 32+ bytes)
	SAML          SAMLIdPConfig `mapstructure:"saml"`
	LDAP          LDAPIdPConfig `mapstructure:"ldap"`
}

// SAMLIdPConfig gates SAML 2.0 scaffold endpoints.
type SAMLIdPConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// LDAPIdPConfig gates an experimental LDAP read/bind listener scaffold.
type LDAPIdPConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	ListenAddr string `mapstructure:"listen_addr"`   // e.g. :389 or :1389
	TLSCert    string `mapstructure:"tls_cert_file"` // optional PEM
	TLSKey     string `mapstructure:"tls_key_file"`  // optional PEM
}

// SSOConfig holds SSO (OIDC) configuration for unified auth with ops IdPs
type SSOConfig struct {
	Enabled bool       `mapstructure:"enabled"`
	OIDC    OIDCConfig `mapstructure:"oidc"`
}

// OIDCConfig holds OpenID Connect provider settings
type OIDCConfig struct {
	IssuerURL       string   `mapstructure:"issuer_url"`
	ClientID        string   `mapstructure:"client_id"`
	ClientSecret    string   `mapstructure:"client_secret"`
	RedirectURL     string   `mapstructure:"redirect_url"`      // e.g. https://kkops.example.com/api/v1/auth/sso/callback
	FrontendBaseURL string   `mapstructure:"frontend_base_url"` // e.g. https://app.kkops.com for redirect after login; if empty, use relative /auth/callback
	Scopes          []string `mapstructure:"scopes"`            // default: openid, profile, email
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	SSLMode      string `mapstructure:"sslmode"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string `mapstructure:"secret"`
	ExpiresIn        int    `mapstructure:"expires_in"`         // seconds
	RefreshExpiresIn int    `mapstructure:"refresh_expires_in"` // seconds
}

// EncryptionConfig holds encryption configuration
type EncryptionConfig struct {
	Key string `mapstructure:"key"` // Encryption key for SSH keys and sensitive data
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `mapstructure:"level"`  // debug, info, warn, error
	Format string `mapstructure:"format"` // json, text
	Output string `mapstructure:"output"` // stdout, file path
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	viper.SetConfigType("yaml")

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("./configs")
		viper.AddConfigPath(".")
	}

	// Set defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("jwt.expires_in", 7200)
	viper.SetDefault("jwt.refresh_expires_in", 604800)
	viper.SetDefault("sso.enabled", false)
	viper.SetDefault("sso.oidc.scopes", []string{"openid", "profile", "email"})
	viper.SetDefault("idp.enabled", true)
	viper.SetDefault("idp.session_secret", "change-idp-session-secret-in-production")
	viper.SetDefault("idp.saml.enabled", false)
	viper.SetDefault("idp.ldap.enabled", false)
	viper.SetDefault("idp.ldap.listen_addr", ":1389")
	viper.SetDefault("encryption.key", "change-this-encryption-key-in-production")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text") // text for dev, json for production
	viper.SetDefault("log.output", "stdout")

	// Read from environment variables
	viper.AutomaticEnv()

	// Read from config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Override with environment variables if set
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &config.Database.Port)
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		config.Database.DBName = dbname
	}
	if v := os.Getenv("ALERTS_WEBHOOK_SECRET"); v != "" {
		config.Alerts.WebhookSecret = v
	}

	return &config, nil
}

// DSN returns the PostgreSQL connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
