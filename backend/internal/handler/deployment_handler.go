package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kronos/backend/internal/models"
	"github.com/kronos/backend/internal/service"
)

type DeploymentConfigHandler struct {
	configService service.DeploymentConfigService
}

func NewDeploymentConfigHandler(configService service.DeploymentConfigService) *DeploymentConfigHandler {
	return &DeploymentConfigHandler{configService: configService}
}

func (h *DeploymentConfigHandler) CreateConfig(c *gin.Context) {
	var config models.DeploymentConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	config.CreatedBy = userID.(uint64)

	if err := h.configService.CreateConfig(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"config": config})
}

func (h *DeploymentConfigHandler) GetConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config id"})
		return
	}

	config, err := h.configService.GetConfig(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "config not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"config": config})
}

func (h *DeploymentConfigHandler) ListConfigs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if applicationName := c.Query("application_name"); applicationName != "" {
		filters["application_name"] = applicationName
	}
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if environment := c.Query("environment"); environment != "" {
		filters["environment"] = environment
	}

	configs, total, err := h.configService.ListConfigs(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"configs":   configs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *DeploymentConfigHandler) UpdateConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config id"})
		return
	}

	var config models.DeploymentConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.configService.UpdateConfig(id, &config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "config updated successfully"})
}

func (h *DeploymentConfigHandler) DeleteConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config id"})
		return
	}

	if err := h.configService.DeleteConfig(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "config deleted successfully"})
}

type DeploymentHandler struct {
	deploymentService service.DeploymentService
}

func NewDeploymentHandler(deploymentService service.DeploymentService) *DeploymentHandler {
	return &DeploymentHandler{deploymentService: deploymentService}
}

func (h *DeploymentHandler) CreateDeployment(c *gin.Context) {
	var req struct {
		ConfigID    uint64   `json:"config_id" binding:"required"`
		Version     string   `json:"version" binding:"required"`
		TargetHosts []uint64 `json:"target_hosts" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	deployment, err := h.deploymentService.StartDeployment(
		req.ConfigID,
		req.Version,
		req.TargetHosts,
		userID.(uint64),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"deployment": deployment})
}

func (h *DeploymentHandler) GetDeployment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deployment id"})
		return
	}

	deployment, err := h.deploymentService.GetDeployment(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "deployment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deployment": deployment})
}

func (h *DeploymentHandler) ListDeployments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if configID := c.Query("config_id"); configID != "" {
		if id, err := strconv.ParseUint(configID, 10, 64); err == nil {
			filters["config_id"] = id
		}
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if version := c.Query("version"); version != "" {
		filters["version"] = version
	}

	deployments, total, err := h.deploymentService.ListDeployments(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deployments": deployments,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
	})
}

func (h *DeploymentHandler) GetDeploymentStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deployment id"})
		return
	}

	deployment, err := h.deploymentService.GetDeploymentStatus(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "deployment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deployment": deployment})
}

func (h *DeploymentHandler) RollbackDeployment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deployment id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	rollbackDeployment, err := h.deploymentService.RollbackDeployment(id, userID.(uint64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"deployment": rollbackDeployment})
}

type DeploymentVersionHandler struct {
	versionService service.DeploymentVersionService
}

func NewDeploymentVersionHandler(versionService service.DeploymentVersionService) *DeploymentVersionHandler {
	return &DeploymentVersionHandler{versionService: versionService}
}

func (h *DeploymentVersionHandler) CreateVersion(c *gin.Context) {
	var version models.DeploymentVersion
	if err := c.ShouldBindJSON(&version); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	version.CreatedBy = userID.(uint64)

	if err := h.versionService.CreateVersion(&version); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"version": version})
}

func (h *DeploymentVersionHandler) GetVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version id"})
		return
	}

	version, err := h.versionService.GetVersion(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"version": version})
}

func (h *DeploymentVersionHandler) ListVersions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if applicationName := c.Query("application_name"); applicationName != "" {
		filters["application_name"] = applicationName
	}

	versions, total, err := h.versionService.ListVersions(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"versions":  versions,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *DeploymentVersionHandler) GetVersionsByApplication(c *gin.Context) {
	applicationName := c.Param("application_name")
	if applicationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "application_name is required"})
		return
	}

	versions, err := h.versionService.GetVersionsByApplication(applicationName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

