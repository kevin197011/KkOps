package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
	"github.com/kkops/backend/internal/service"
	"github.com/kkops/backend/internal/salt"
	"github.com/kkops/backend/internal/utils"
)

type SettingsHandler struct {
	settingsService service.SettingsService
	saltManager     *salt.Manager
}

// NewSettingsHandler 创建设置处理器
func NewSettingsHandler(settingsService service.SettingsService, saltManager *salt.Manager) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		saltManager:     saltManager,
	}
}

// GetAllSettings 获取所有设置
func (h *SettingsHandler) GetAllSettings(c *gin.Context) {
	settings, err := h.settingsService.GetAllSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// GetSettingsByCategory 按分类获取设置
func (h *SettingsHandler) GetSettingsByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category parameter is required"})
		return
	}

	settings, err := h.settingsService.GetSettingsByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// GetSetting 获取单个设置
func (h *SettingsHandler) GetSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key parameter is required"})
		return
	}

	// key已经是完整的路径参数，不需要额外处理

	setting, err := h.settingsService.GetSetting(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "setting not found"})
		return
	}

	// 掩码敏感值
	if strings.Contains(setting.Key, "password") {
		setting.Value = "***"
	}

	c.JSON(http.StatusOK, gin.H{"setting": setting})
}

// UpdateSetting 更新单个设置
func (h *SettingsHandler) UpdateSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key parameter is required"})
		return
	}

	// key已经是完整的路径参数，不需要额外处理

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		Value       string `json:"value" binding:"required"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果没有提供category，尝试从key推断
	category := req.Category
	if category == "" {
		if parts := strings.Split(key, "."); len(parts) > 0 {
			category = parts[0]
		} else {
			category = "general"
		}
	}

	if err := h.settingsService.UpdateSetting(key, req.Value, category, req.Description, userID.(uint64)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 注意：Salt客户端重新初始化需要在应用层面处理
	// 这里只更新数据库配置，Salt客户端会在下次重启时重新初始化

	c.JSON(http.StatusOK, gin.H{"message": "setting updated successfully"})
}

// UpdateSaltConfig 更新Salt配置（便捷端点）
func (h *SettingsHandler) UpdateSaltConfig(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		APIURL    string `json:"api_url" binding:"required"`
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password"` // 密码可以为空（如果为空则保留现有密码）
		EAuth     string `json:"eauth" binding:"required"`
		Timeout   int    `json:"timeout" binding:"required"`
		VerifySSL bool   `json:"verify_ssl"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果密码为空，从数据库读取现有密码
	password := req.Password
	if password == "" {
		existingConfig, err := h.settingsService.GetSaltConfig()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required for new configuration, or please provide existing configuration"})
			return
		}
		password = existingConfig.Password
		if password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
			return
		}
	}

	cfg := &config.SaltConfig{
		APIURL:    req.APIURL,
		Username:  req.Username,
		Password:  password,
		EAuth:     req.EAuth,
		Timeout:   req.Timeout,
		VerifySSL: req.VerifySSL,
	}

	if err := h.settingsService.UpdateSaltConfig(cfg, userID.(uint64)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 重新加载Salt客户端（热重载）
	if h.saltManager != nil {
		// 先重新加载配置（从数据库）
		if err := reloadSaltConfigFromDatabase(); err != nil {
			log.Printf("Warning: Failed to reload Salt config from database: %v", err)
		}
		// 然后重新加载客户端
		if err := h.saltManager.ReloadClient(); err != nil {
			// 重载失败，记录日志但不返回错误（配置已保存）
			c.JSON(http.StatusOK, gin.H{
				"message": "Salt configuration saved, but client reload failed. Please restart the service.",
				"warning": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Salt configuration updated and client reloaded successfully"})
}

// GetSaltConfig 获取Salt配置
func (h *SettingsHandler) GetSaltConfig(c *gin.Context) {
	cfg, err := h.settingsService.GetSaltConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 检查数据库中是否存在密码配置
	// 通过直接查询数据库来判断是否已配置密码，而不是依赖 cfg.Password 的值
	// 因为 cfg.Password 可能包含从环境变量读取的默认值
	_, err = h.settingsService.GetSetting("salt.password")
	if err == nil {
		// 数据库中存在密码配置，返回掩码值
		cfg.Password = "***"
	} else {
		// 数据库中不存在密码配置，返回空字符串
		cfg.Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"config": cfg})
}

// TestSaltConnection 测试Salt API连接
func (h *SettingsHandler) TestSaltConnection(c *gin.Context) {
	var req struct {
		APIURL    string `json:"api_url"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		EAuth     string `json:"eauth"`
		Timeout   int    `json:"timeout"`
		VerifySSL bool   `json:"verify_ssl"`
	}

	// 如果请求体中有配置，使用请求体中的配置；否则使用当前保存的配置
	var testConfig *config.SaltConfig
	if err := c.ShouldBindJSON(&req); err == nil && req.APIURL != "" {
		// 使用请求体中的配置进行测试
		// 如果密码为空，从数据库读取现有密码（用于测试已保存的配置）
		password := req.Password
		if password == "" {
			existingConfig, err := h.settingsService.GetSaltConfig()
			if err == nil && existingConfig.Password != "" {
				password = existingConfig.Password
			}
		}

		testConfig = &config.SaltConfig{
			APIURL:    req.APIURL,
			Username:  req.Username,
			Password:  password,
			EAuth:     req.EAuth,
			Timeout:   req.Timeout,
			VerifySSL: req.VerifySSL,
		}

		// 验证必填字段
		if testConfig.APIURL == "" || testConfig.Username == "" || testConfig.Password == "" || testConfig.EAuth == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "API URL, Username, Password, and EAuth are required"})
			return
		}
	} else {
		// 使用当前保存的配置进行测试
		cfg, err := h.settingsService.GetSaltConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load current Salt configuration"})
			return
		}
		testConfig = cfg
	}

	// 创建临时Salt客户端进行测试
	testClient := salt.NewClient(*testConfig)
	
	// 测试连接
	if err := testClient.TestConnection(); err != nil {
		// 提供更友好的错误消息
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "401") || strings.Contains(errorMsg, "Unauthorized") {
			errorMsg = "认证失败：请检查用户名、密码和EAuth类型是否正确"
		} else if strings.Contains(errorMsg, "connection refused") || strings.Contains(errorMsg, "no such host") {
			errorMsg = "无法连接到Salt API：请检查API URL是否正确，以及Salt Master服务是否运行"
		} else if strings.Contains(errorMsg, "timeout") {
			errorMsg = "连接超时：请检查网络连接和超时设置"
		}
		
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   errorMsg,
			"message": "连接测试失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "连接测试成功",
	})
}

// reloadSaltConfigFromDatabase 从数据库重新加载Salt配置
func reloadSaltConfigFromDatabase() error {
	if models.DB == nil {
		return nil
	}

	settingsRepo := repository.NewSettingsRepository(models.DB)
	saltSettings, err := settingsRepo.GetByCategory("salt")
	if err != nil {
		return err
	}

	if len(saltSettings) > 0 {
		for _, setting := range saltSettings {
			switch setting.Key {
			case "salt.api_url":
				config.AppConfig.Salt.APIURL = setting.Value
			case "salt.username":
				config.AppConfig.Salt.Username = setting.Value
			case "salt.password":
				decrypted, err := utils.Decrypt(setting.Value)
				if err == nil {
					config.AppConfig.Salt.Password = decrypted
				}
			case "salt.eauth":
				config.AppConfig.Salt.EAuth = setting.Value
			case "salt.timeout":
				if timeout, err := strconv.Atoi(setting.Value); err == nil {
					config.AppConfig.Salt.Timeout = timeout
				}
			case "salt.verify_ssl":
				config.AppConfig.Salt.VerifySSL = setting.Value == "true"
			}
		}
	}

	return nil
}

// GetAuditLogSettings 获取审计日志设置
func (h *SettingsHandler) GetAuditLogSettings(c *gin.Context) {
	settings, err := h.settingsService.GetAuditLogSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateAuditLogSettings 更新审计日志设置
func (h *SettingsHandler) UpdateAuditLogSettings(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		RetentionDays int `json:"retention_days" binding:"required,min=1,max=3650"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings := &service.AuditLogSettings{
		RetentionDays: req.RetentionDays,
	}

	if err := h.settingsService.UpdateAuditLogSettings(settings, userID.(uint64)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "audit log settings updated successfully"})
}

// GetAuditLogStats 获取审计日志统计信息
func (h *SettingsHandler) GetAuditLogStats(c *gin.Context) {
	stats, err := h.settingsService.GetAuditLogStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// CleanupAuditLogs 手动清理审计日志
func (h *SettingsHandler) CleanupAuditLogs(c *gin.Context) {
	deletedCount, err := h.settingsService.CleanupAuditLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "audit logs cleaned up successfully",
		"deleted_count": deletedCount,
	})
}

