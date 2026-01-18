// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package deployment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/service/deployment"
)

// Handler handles deployment management HTTP requests
type Handler struct {
	service *deployment.Service
}

// NewHandler creates a new deployment handler
func NewHandler(service *deployment.Service) *Handler {
	return &Handler{service: service}
}

// CreateModule handles deployment module creation
// @Summary Create deployment module
// @Description Create a new deployment module
// @Tags deployment
// @Accept json
// @Produce json
// @Param request body deployment.CreateModuleRequest true "Create module request"
// @Success 201 {object} deployment.ModuleResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/deployment-modules [post]
func (h *Handler) CreateModule(c *gin.Context) {
	var req deployment.CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	resp, err := h.service.CreateModule(&req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetModule handles deployment module retrieval
// @Summary Get deployment module
// @Description Get deployment module by ID
// @Tags deployment
// @Produce json
// @Param id path int true "Module ID"
// @Success 200 {object} deployment.ModuleResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/deployment-modules/{id} [get]
func (h *Handler) GetModule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid module ID"})
		return
	}

	resp, err := h.service.GetModule(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "module not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListModules handles deployment module list retrieval
// @Summary List deployment modules
// @Description Get list of deployment modules
// @Tags deployment
// @Produce json
// @Param project_id query int false "Filter by project ID"
// @Success 200 {array} deployment.ModuleResponse
// @Router /api/v1/deployment-modules [get]
func (h *Handler) ListModules(c *gin.Context) {
	var projectID *uint
	if pidStr := c.Query("project_id"); pidStr != "" {
		pid, err := strconv.ParseUint(pidStr, 10, 32)
		if err == nil {
			p := uint(pid)
			projectID = &p
		}
	}

	modules, err := h.service.ListModules(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, modules)
}

// UpdateModule handles deployment module update
// @Summary Update deployment module
// @Description Update deployment module by ID
// @Tags deployment
// @Accept json
// @Produce json
// @Param id path int true "Module ID"
// @Param request body deployment.UpdateModuleRequest true "Update module request"
// @Success 200 {object} deployment.ModuleResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/deployment-modules/{id} [put]
func (h *Handler) UpdateModule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid module ID"})
		return
	}

	var req deployment.UpdateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.UpdateModule(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteModule handles deployment module deletion
// @Summary Delete deployment module
// @Description Delete deployment module by ID
// @Tags deployment
// @Param id path int true "Module ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /api/v1/deployment-modules/{id} [delete]
func (h *Handler) DeleteModule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid module ID"})
		return
	}

	if err := h.service.DeleteModule(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "module not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetVersions handles version list retrieval from version source
// @Summary Get versions
// @Description Get available versions from version source URL
// @Tags deployment
// @Produce json
// @Param id path int true "Module ID"
// @Success 200 {object} deployment.VersionSourceResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/deployment-modules/{id}/versions [get]
func (h *Handler) GetVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid module ID"})
		return
	}

	versions, err := h.service.GetVersions(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// Deploy handles deployment execution
// @Summary Execute deployment
// @Description Execute deployment for a module with selected version
// @Tags deployment
// @Accept json
// @Produce json
// @Param id path int true "Module ID"
// @Param request body deployment.DeployRequest true "Deploy request"
// @Success 201 {object} deployment.DeploymentResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/deployment-modules/{id}/deploy [post]
func (h *Handler) Deploy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid module ID"})
		return
	}

	var req deployment.DeployRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	resp, err := h.service.Deploy(uint(id), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetDeployment handles deployment record retrieval
// @Summary Get deployment
// @Description Get deployment record by ID
// @Tags deployment
// @Produce json
// @Param id path int true "Deployment ID"
// @Success 200 {object} deployment.DeploymentResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/deployments/{id} [get]
func (h *Handler) GetDeployment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deployment ID"})
		return
	}

	resp, err := h.service.GetDeployment(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "deployment not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListDeployments handles deployment list retrieval
// @Summary List deployments
// @Description Get list of deployment records
// @Tags deployment
// @Produce json
// @Param module_id query int false "Filter by module ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/deployments [get]
func (h *Handler) ListDeployments(c *gin.Context) {
	var moduleID *uint
	if midStr := c.Query("module_id"); midStr != "" {
		mid, err := strconv.ParseUint(midStr, 10, 32)
		if err == nil {
			m := uint(mid)
			moduleID = &m
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	deployments, total, err := h.service.ListDeployments(moduleID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  deployments,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// CancelDeployment handles deployment cancellation
// @Summary Cancel deployment
// @Description Cancel a running deployment
// @Tags deployment
// @Param id path int true "Deployment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/deployments/{id}/cancel [post]
func (h *Handler) CancelDeployment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deployment ID"})
		return
	}

	if err := h.service.CancelDeployment(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deployment cancelled"})
}

// ExportModuleConfig 导出配置结构
type ExportModuleConfig struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	ProjectName      string   `json:"project_name"`
	EnvironmentName  string   `json:"environment_name,omitempty"`
	TemplateName     string   `json:"template_name,omitempty"`
	VersionSourceURL string   `json:"version_source_url,omitempty"`
	DeployScript     string   `json:"deploy_script"`
	ScriptType       string   `json:"script_type"`
	Timeout          int      `json:"timeout"`
	TargetHosts      []string `json:"target_hosts,omitempty"` // 目标主机列表（主机名或IP）
}

// ExportConfig 导出配置
type ExportConfig struct {
	Version   string               `json:"version"`
	ExportAt  string               `json:"export_at"`
	Modules   []ExportModuleConfig `json:"modules"`
}

// ExportModules handles deployment modules export
// @Summary Export deployment modules
// @Description Export all deployment modules to JSON
// @Tags deployment
// @Produce json
// @Success 200 {object} ExportConfig
// @Router /api/v1/deployment-modules/export [get]
func (h *Handler) ExportModules(c *gin.Context) {
	modules, err := h.service.ListModules(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	exportConfig := ExportConfig{
		Version:  "1.0",
		ExportAt: time.Now().Format(time.RFC3339),
		Modules:  make([]ExportModuleConfig, len(modules)),
	}

	for i, m := range modules {
		// 获取目标主机名称
		targetHosts := h.service.GetAssetHostnames(m.AssetIDs)

		exportConfig.Modules[i] = ExportModuleConfig{
			Name:             m.Name,
			Description:      m.Description,
			ProjectName:      m.ProjectName,
			EnvironmentName:  m.EnvironmentName,
			TemplateName:     m.TemplateName,
			VersionSourceURL: m.VersionSourceURL,
			DeployScript:     m.DeployScript,
			ScriptType:       m.ScriptType,
			Timeout:          m.Timeout,
			TargetHosts:      targetHosts,
		}
	}

	// 设置响应头，以便浏览器下载
	filename := fmt.Sprintf("deployment-modules-%s.json", time.Now().Format("20060102-150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, exportConfig)
}

// ImportModuleConfig 导入配置结构
type ImportModuleConfig struct {
	Name             string   `json:"name" binding:"required"`
	Description      string   `json:"description"`
	ProjectName      string   `json:"project_name" binding:"required"`
	EnvironmentName  string   `json:"environment_name"`
	TemplateName     string   `json:"template_name"`
	VersionSourceURL string   `json:"version_source_url"`
	DeployScript     string   `json:"deploy_script" binding:"required"`
	ScriptType       string   `json:"script_type"`
	Timeout          int      `json:"timeout"`
	TargetHosts      []string `json:"target_hosts"` // 目标主机列表（主机名或IP）
}

// ImportConfig 导入配置
type ImportConfig struct {
	Version string               `json:"version"`
	Modules []ImportModuleConfig `json:"modules" binding:"required"`
}

// ImportResult 导入结果
type ImportResult struct {
	Total     int      `json:"total"`
	Success   int      `json:"success"`
	Failed    int      `json:"failed"`
	Errors    []string `json:"errors,omitempty"`
	Skipped   []string `json:"skipped,omitempty"`
}

// ImportModules handles deployment modules import
// @Summary Import deployment modules
// @Description Import deployment modules from JSON
// @Tags deployment
// @Accept json
// @Produce json
// @Param request body ImportConfig true "Import config"
// @Success 200 {object} ImportResult
// @Failure 400 {object} map[string]string
// @Router /api/v1/deployment-modules/import [post]
func (h *Handler) ImportModules(c *gin.Context) {
	var importConfig ImportConfig
	if err := c.ShouldBindJSON(&importConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的导入配置: " + err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	result := ImportResult{
		Total:   len(importConfig.Modules),
		Errors:  []string{},
		Skipped: []string{},
	}

	for _, m := range importConfig.Modules {
		// 查找项目
		projectID, err := h.service.FindProjectByName(m.ProjectName)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("模块 %s: 项目 %s 不存在", m.Name, m.ProjectName))
			continue
		}

		// 查找环境（可选）
		var environmentID *uint
		if m.EnvironmentName != "" {
			envID, err := h.service.FindEnvironmentByName(m.EnvironmentName)
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("模块 %s: 环境 %s 不存在", m.Name, m.EnvironmentName))
				continue
			}
			environmentID = &envID
		}

		// 查找模板（可选）
		var templateID *uint
		if m.TemplateName != "" {
			tplID, err := h.service.FindTemplateByName(m.TemplateName)
			if err != nil {
				// 模板不存在只是警告，不阻止导入
				result.Skipped = append(result.Skipped, fmt.Sprintf("模块 %s: 模板 %s 不存在，已跳过模板关联", m.Name, m.TemplateName))
			} else {
				templateID = &tplID
			}
		}

		// 解析目标主机
		var assetIDs []uint
		var missingHosts []string
		if len(m.TargetHosts) > 0 {
			assetIDs, missingHosts = h.service.FindAssetsByHostnames(m.TargetHosts)
			if len(missingHosts) > 0 {
				result.Skipped = append(result.Skipped, fmt.Sprintf("模块 %s: 主机 %v 不存在，已跳过", m.Name, missingHosts))
			}
		}

		// 检查是否已存在同名模块
		exists, existingID := h.service.ModuleExistsByName(m.Name, projectID)
		
		scriptType := m.ScriptType
		if scriptType == "" {
			scriptType = "shell"
		}
		timeout := m.Timeout
		if timeout == 0 {
			timeout = 600
		}

		if exists {
			// 更新现有模块
			updateReq := &deployment.UpdateModuleRequest{
				ProjectID:        &projectID,
				EnvironmentID:    environmentID,
				TemplateID:       templateID,
				Name:             m.Name,
				Description:      m.Description,
				VersionSourceURL: m.VersionSourceURL,
				DeployScript:     m.DeployScript,
				ScriptType:       scriptType,
				Timeout:          timeout,
				AssetIDs:         assetIDs,
			}
			_, err := h.service.UpdateModule(existingID, updateReq)
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("模块 %s: 更新失败 - %s", m.Name, err.Error()))
				continue
			}
			result.Skipped = append(result.Skipped, fmt.Sprintf("模块 %s: 已存在，已更新配置", m.Name))
		} else {
			// 创建新模块
			createReq := &deployment.CreateModuleRequest{
				ProjectID:        projectID,
				EnvironmentID:    environmentID,
				TemplateID:       templateID,
				Name:             m.Name,
				Description:      m.Description,
				VersionSourceURL: m.VersionSourceURL,
				DeployScript:     m.DeployScript,
				ScriptType:       scriptType,
				Timeout:          timeout,
				AssetIDs:         assetIDs,
			}
			_, err := h.service.CreateModule(createReq, userID)
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("模块 %s: 创建失败 - %s", m.Name, err.Error()))
				continue
			}
		}
		result.Success++
	}

	c.JSON(http.StatusOK, result)
}

// PreviewImport handles deployment modules import preview
// @Summary Preview import deployment modules
// @Description Preview import result without actually importing
// @Tags deployment
// @Accept json
// @Produce json
// @Param request body ImportConfig true "Import config"
// @Success 200 {object} ImportResult
// @Failure 400 {object} map[string]string
// @Router /api/v1/deployment-modules/import/preview [post]
func (h *Handler) PreviewImport(c *gin.Context) {
	var importConfig ImportConfig
	
	// 尝试解析 JSON
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取请求体失败"})
		return
	}

	if err := json.Unmarshal(body, &importConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 JSON 格式: " + err.Error()})
		return
	}

	result := ImportResult{
		Total:   len(importConfig.Modules),
		Errors:  []string{},
		Skipped: []string{},
	}

	for _, m := range importConfig.Modules {
		// 验证项目
		_, err := h.service.FindProjectByName(m.ProjectName)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("模块 %s: 项目 %s 不存在", m.Name, m.ProjectName))
			continue
		}

		// 验证环境（可选）
		if m.EnvironmentName != "" {
			_, err := h.service.FindEnvironmentByName(m.EnvironmentName)
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("模块 %s: 环境 %s 不存在", m.Name, m.EnvironmentName))
				continue
			}
		}

		// 验证模板（可选）
		if m.TemplateName != "" {
			_, err := h.service.FindTemplateByName(m.TemplateName)
			if err != nil {
				result.Skipped = append(result.Skipped, fmt.Sprintf("模块 %s: 模板 %s 不存在，将跳过模板关联", m.Name, m.TemplateName))
			}
		}

		// 验证目标主机
		if len(m.TargetHosts) > 0 {
			_, missingHosts := h.service.FindAssetsByHostnames(m.TargetHosts)
			if len(missingHosts) > 0 {
				result.Skipped = append(result.Skipped, fmt.Sprintf("模块 %s: 主机 %v 不存在，将跳过", m.Name, missingHosts))
			}
		}

		// 检查是否已存在
		exists, _ := h.service.ModuleExistsByName(m.Name, 0)
		if exists {
			result.Skipped = append(result.Skipped, fmt.Sprintf("模块 %s: 已存在，将更新配置", m.Name))
		}

		result.Success++
	}

	c.JSON(http.StatusOK, result)
}
