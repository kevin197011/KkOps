package service

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
	"github.com/kkops/backend/internal/utils"
)

type SettingsService interface {
	GetSetting(key string) (*models.SystemSettings, error)
	GetSettingsByCategory(category string) (map[string]string, error)
	GetAllSettings() (map[string]string, error)
	UpdateSetting(key, value, category, description string, updatedBy uint64) error
	GetSaltConfig() (*config.SaltConfig, error)
	UpdateSaltConfig(cfg *config.SaltConfig, updatedBy uint64) error
	GetAuditLogSettings() (*AuditLogSettings, error)
	UpdateAuditLogSettings(settings *AuditLogSettings, updatedBy uint64) error
	GetAuditLogStats() (*AuditLogStats, error)
	CleanupAuditLogs() (int64, error)
}

type AuditLogSettings struct {
	RetentionDays int `json:"retention_days"`
}

type AuditLogStats struct {
	Total        int64 `json:"total"`
	OldRecords   int64 `json:"oldRecords"`
	RetentionDays int   `json:"retentionDays"`
}

type settingsService struct {
	settingsRepo repository.SettingsRepository
	auditRepo    repository.AuditRepository
}

func NewSettingsService(settingsRepo repository.SettingsRepository, auditRepo repository.AuditRepository) SettingsService {
	return &settingsService{
		settingsRepo: settingsRepo,
		auditRepo:    auditRepo,
	}
}

func (s *settingsService) GetSetting(key string) (*models.SystemSettings, error) {
	return s.settingsRepo.GetByKey(key)
}

func (s *settingsService) GetSettingsByCategory(category string) (map[string]string, error) {
	settings, err := s.settingsRepo.GetByCategory(category)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		// 如果是敏感值（密码），进行掩码处理
		if strings.Contains(setting.Key, "password") {
			result[setting.Key] = "***"
		} else {
			result[setting.Key] = setting.Value
		}
	}

	return result, nil
}

func (s *settingsService) GetAllSettings() (map[string]string, error) {
	settings, err := s.settingsRepo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		// 如果是敏感值（密码），进行掩码处理
		if strings.Contains(setting.Key, "password") {
			result[setting.Key] = "***"
		} else {
			result[setting.Key] = setting.Value
		}
	}

	return result, nil
}

func (s *settingsService) UpdateSetting(key, value, category, description string, updatedBy uint64) error {
	// 验证设置值（根据key进行特定验证）
	if err := s.validateSetting(key, value); err != nil {
		return err
	}

	// 如果是敏感值（密码），进行加密
	if strings.Contains(key, "password") {
		encrypted, err := utils.Encrypt(value)
		if err != nil {
			return fmt.Errorf("failed to encrypt password: %w", err)
		}
		value = encrypted
	}

	setting := &models.SystemSettings{
		Key:         key,
		Value:       value,
		Category:    category,
		Description: description,
		UpdatedBy:   updatedBy,
		UpdatedAt:   time.Now(),
	}

	return s.settingsRepo.CreateOrUpdate(setting)
}

func (s *settingsService) validateSetting(key, value string) error {
	// 根据key进行特定验证
	if strings.Contains(key, "url") {
		_, err := url.Parse(value)
		if err != nil {
			return fmt.Errorf("invalid URL format: %w", err)
		}
	} else if strings.Contains(key, "timeout") {
		timeout, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("timeout must be a number: %w", err)
		}
		if timeout < 1 || timeout > 300 {
			return errors.New("timeout must be between 1 and 300 seconds")
		}
	} else if strings.Contains(key, "verify_ssl") {
		// 布尔值验证
		if value != "true" && value != "false" {
			return errors.New("verify_ssl must be 'true' or 'false'")
		}
	}

	return nil
}

func (s *settingsService) GetSaltConfig() (*config.SaltConfig, error) {
	settings, err := s.GetSettingsByCategory("salt")
	if err != nil {
		return nil, err
	}

	cfg := &config.SaltConfig{
		APIURL:    config.AppConfig.Salt.APIURL,    // 默认值
		Username:  config.AppConfig.Salt.Username,  // 默认值
		Password:  config.AppConfig.Salt.Password,  // 默认值
		EAuth:     config.AppConfig.Salt.EAuth,     // 默认值
		Timeout:   config.AppConfig.Salt.Timeout,   // 默认值
		VerifySSL: config.AppConfig.Salt.VerifySSL, // 默认值
	}

	// 从数据库设置中覆盖
	if apiURL, ok := settings["salt.api_url"]; ok && apiURL != "***" {
		cfg.APIURL = apiURL
	}
	if username, ok := settings["salt.username"]; ok && username != "***" {
		cfg.Username = username
	}
	if eauth, ok := settings["salt.eauth"]; ok && eauth != "***" {
		cfg.EAuth = eauth
	}
	if timeoutStr, ok := settings["salt.timeout"]; ok && timeoutStr != "***" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil {
			cfg.Timeout = timeout
		}
	}
	if verifySSLStr, ok := settings["salt.verify_ssl"]; ok && verifySSLStr != "***" {
		cfg.VerifySSL = verifySSLStr == "true"
	}

	// 密码需要特殊处理（需要解密）
	if passwordSetting, err := s.settingsRepo.GetByKey("salt.password"); err == nil {
		decrypted, err := utils.Decrypt(passwordSetting.Value)
		if err == nil {
			cfg.Password = decrypted
		}
	}

	return cfg, nil
}

func (s *settingsService) UpdateSaltConfig(cfg *config.SaltConfig, updatedBy uint64) error {
	// 验证配置
	if cfg.APIURL == "" {
		return errors.New("API URL is required")
	}
	if _, err := url.Parse(cfg.APIURL); err != nil {
		return fmt.Errorf("invalid API URL format: %w", err)
	}
	if cfg.Username == "" {
		return errors.New("username is required")
	}
	if cfg.EAuth == "" {
		return errors.New("eauth type is required")
	}
	if cfg.Timeout < 1 || cfg.Timeout > 300 {
		return errors.New("timeout must be between 1 and 300 seconds")
	}

	// 规范化 API URL（移除尾部斜杠）
	cfg.APIURL = strings.TrimSuffix(cfg.APIURL, "/")

	// 更新所有Salt相关设置
	saltSettings := []struct {
		key         string
		value       string
		description string
	}{
		{"salt.api_url", cfg.APIURL, "Salt API URL"},
		{"salt.username", cfg.Username, "Salt API username"},
		{"salt.password", cfg.Password, "Salt API password (encrypted)"},
		{"salt.eauth", cfg.EAuth, "Salt API eauth type"},
		{"salt.timeout", strconv.Itoa(cfg.Timeout), "Salt API timeout in seconds"},
		{"salt.verify_ssl", fmt.Sprintf("%v", cfg.VerifySSL), "Verify SSL certificates"},
	}

	for _, setting := range saltSettings {
		// 如果是密码字段且值为空，跳过更新（保留现有密码）
		if setting.key == "salt.password" && setting.value == "" {
			continue
		}
		if err := s.UpdateSetting(setting.key, setting.value, "salt", setting.description, updatedBy); err != nil {
			return fmt.Errorf("failed to update %s: %w", setting.key, err)
		}
	}

	return nil
}

func (s *settingsService) GetAuditLogSettings() (*AuditLogSettings, error) {
	// 获取审计日志保存天数设置，默认30天
	retentionDays := 30
	if setting, err := s.settingsRepo.GetByKey("audit.retention_days"); err == nil {
		if days, err := strconv.Atoi(setting.Value); err == nil && days > 0 {
			retentionDays = days
		}
	}

	return &AuditLogSettings{
		RetentionDays: retentionDays,
	}, nil
}

func (s *settingsService) UpdateAuditLogSettings(settings *AuditLogSettings, updatedBy uint64) error {
	// 验证保存天数
	if settings.RetentionDays < 1 || settings.RetentionDays > 3650 {
		return errors.New("retention days must be between 1 and 3650")
	}

	// 更新设置
	return s.UpdateSetting(
		"audit.retention_days",
		strconv.Itoa(settings.RetentionDays),
		"audit",
		"Audit log retention period in days",
		updatedBy,
	)
}

func (s *settingsService) GetAuditLogStats() (*AuditLogStats, error) {
	// 获取当前保存天数设置
	settings, err := s.GetAuditLogSettings()
	if err != nil {
		return nil, err
	}

	// 获取总记录数
	total, err := s.auditRepo.GetTotalCount()
	if err != nil {
		return nil, err
	}

	// 获取过期记录数
	cutoffTime := time.Now().AddDate(0, 0, -settings.RetentionDays)
	oldRecords, err := s.auditRepo.GetCountBeforeTime(cutoffTime)
	if err != nil {
		return nil, err
	}

	return &AuditLogStats{
		Total:         total,
		OldRecords:    oldRecords,
		RetentionDays: settings.RetentionDays,
	}, nil
}

func (s *settingsService) CleanupAuditLogs() (int64, error) {
	// 获取保存天数设置
	settings, err := s.GetAuditLogSettings()
	if err != nil {
		return 0, err
	}

	// 计算截止时间
	cutoffTime := time.Now().AddDate(0, 0, -settings.RetentionDays)

	// 删除过期记录
	return s.auditRepo.DeleteBeforeTime(cutoffTime)
}

