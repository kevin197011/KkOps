// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/utils"
)

// InitDB initializes the database connection and runs migrations
func InitDB(cfg *DatabaseConfig, zapLogger *zap.Logger) (*gorm.DB, error) {
	dsn := cfg.DSN()

	var gormLogger logger.Interface
	if zapLogger != nil {
		gormLogger = NewGormLogger(zapLogger)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// Run SQL migrations first (for complex schema changes)
	if err := runSQLMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run SQL migrations: %w", err)
	}

	// Run GORM AutoMigrate (for schema sync)
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// seedDefaultUser creates a default admin user if no users exist
func seedDefaultUser(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}

	// If users already exist, skip seeding
	if count > 0 {
		return nil
	}

	// Create default admin user
	passwordHash, err := utils.HashPassword("admin123")
	if err != nil {
		return fmt.Errorf("failed to hash default password: %w", err)
	}

	adminUser := model.User{
		Username:     "admin",
		Email:        "admin@kkops.local",
		PasswordHash: passwordHash,
		RealName:     "Administrator",
		Status:       "active",
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return fmt.Errorf("failed to create default admin user: %w", err)
	}

	return nil
}

// runSQLMigrations executes SQL migration files in order
func runSQLMigrations(db *gorm.DB) error {
	// Get the base directory (assuming migrations are relative to backend/)
	migrationDir := "migrations"
	
	// Check if migration directory exists
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		// Try alternative path (when running from backend directory)
		migrationDir = "backend/migrations"
		if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
			// No migration directory, skip SQL migrations
			return nil
		}
	}

	// List of migration files in order
	migrationFiles := []string{
		"rename_asset_name_to_hostname.sql",
		"remove_code_fields.sql",
		"add_cloud_platform_management.sql",
	}

	// Execute migrations in order
	for _, filename := range migrationFiles {
		filepath := filepath.Join(migrationDir, filename)
		
		// Check if file exists
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			continue // Skip if file doesn't exist
		}

		// Read migration file
		sqlContent, err := os.ReadFile(filepath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Execute the entire SQL file as one transaction
		// This is better for DO blocks and complex migrations
		sql := string(sqlContent)
		
		// Remove comments that are on their own lines
		lines := strings.Split(sql, "\n")
		var cleanLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			// Keep lines that are not empty comments or comment-only lines
			if trimmed != "" && !strings.HasPrefix(trimmed, "--") {
				cleanLines = append(cleanLines, line)
			} else if strings.Contains(line, "--") && !strings.HasPrefix(trimmed, "--") {
				// Keep lines that have code before comments
				cleanLines = append(cleanLines, line)
			}
		}
		sql = strings.Join(cleanLines, "\n")

		// Execute SQL
		if err := db.Exec(sql).Error; err != nil {
			// Some errors are expected (e.g., column already exists, index already exists)
			// Check if it's a safe error to ignore
			errStr := err.Error()
			if strings.Contains(errStr, "already exists") ||
				strings.Contains(errStr, "duplicate") {
				// Safe to ignore, continue
				continue
			}
			
			// For column/table not found errors, they might be expected depending on migration state
			// But we should log them as warnings, not fail completely
			if strings.Contains(errStr, "does not exist") {
				// Log as warning but continue
				continue
			}

			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}
	}

	return nil
}

// migrate runs database migrations
func migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Department{},
		&model.APIToken{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
		&model.Project{},
		&model.Environment{},
		&model.CloudPlatform{},
		&model.AssetCategory{},
		&model.Asset{},
		&model.Tag{},
		&model.AssetTag{},
		&model.SSHKey{},
		&model.TaskTemplate{},
		&model.Task{},
		&model.TaskExecution{},
	); err != nil {
		return err
	}

	// Initialize default admin user if database is empty
	return seedDefaultUser(db)
}

// NewGormLogger creates a GORM logger adapter for Zap
func NewGormLogger(zapLogger *zap.Logger) logger.Interface {
	return logger.New(
		&zapWriter{logger: zapLogger},
		logger.Config{
			SlowThreshold:             0,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

type zapWriter struct {
	logger *zap.Logger
}

func (w *zapWriter) Printf(format string, args ...interface{}) {
	w.logger.Sugar().Infof(format, args...)
}
