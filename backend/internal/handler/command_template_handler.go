package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/service"
	"gorm.io/datatypes"
)

type CommandTemplateHandler struct {
	templateService service.CommandTemplateService
}

func NewCommandTemplateHandler(templateService service.CommandTemplateService) *CommandTemplateHandler {
	return &CommandTemplateHandler{
		templateService: templateService,
	}
}

// CreateTemplate 创建命令模板
func (h *CommandTemplateHandler) CreateTemplate(c *gin.Context) {
	var req struct {
		Name            string        `json:"name" binding:"required"`
		Description     string        `json:"description"`
		Category        string        `json:"category"`
		CommandFunction string        `json:"command_function" binding:"required"`
		CommandArgs     []interface{} `json:"command_args"`
		Icon            string        `json:"icon"`
		IsPublic        bool          `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 转换命令参数为 JSON
	var commandArgsJSON datatypes.JSON
	if len(req.CommandArgs) > 0 {
		argsBytes, err := json.Marshal(req.CommandArgs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid command_args format"})
			return
		}
		commandArgsJSON = datatypes.JSON(argsBytes)
	}

	template := &models.CommandTemplate{
		Name:            req.Name,
		Description:     req.Description,
		Category:        req.Category,
		CommandFunction: req.CommandFunction,
		CommandArgs:     commandArgsJSON,
		Icon:            req.Icon,
		CreatedBy:       userID.(uint64),
		IsPublic:        req.IsPublic,
	}

	if err := h.templateService.CreateTemplate(template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"template": template})
}

// ListTemplates 列出命令模板
func (h *CommandTemplateHandler) ListTemplates(c *gin.Context) {
	filters := make(map[string]interface{})
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if isPublic := c.Query("is_public"); isPublic != "" {
		if isPublic == "true" {
			filters["is_public"] = true
		}
	}

	// 获取当前用户ID（用于过滤可见模板）
	userID, exists := c.Get("user_id")
	if exists {
		filters["created_by"] = userID.(uint64)
		filters["is_public"] = false // 查询公开模板或用户自己的模板
	}

	templates, err := h.templateService.ListTemplates(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// GetTemplate 获取模板详情
func (h *CommandTemplateHandler) GetTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	template, err := h.templateService.GetTemplate(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"template": template})
}

// UpdateTemplate 更新模板
func (h *CommandTemplateHandler) UpdateTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	var req struct {
		Name            string        `json:"name"`
		Description     string        `json:"description"`
		Category        string        `json:"category"`
		CommandFunction string        `json:"command_function"`
		CommandArgs     []interface{} `json:"command_args"`
		Icon            string        `json:"icon"`
		IsPublic        bool          `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 转换命令参数为 JSON
	var commandArgsJSON datatypes.JSON
	if len(req.CommandArgs) > 0 {
		argsBytes, err := json.Marshal(req.CommandArgs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid command_args format"})
			return
		}
		commandArgsJSON = datatypes.JSON(argsBytes)
	}

	template := &models.CommandTemplate{
		Name:            req.Name,
		Description:     req.Description,
		Category:        req.Category,
		CommandFunction: req.CommandFunction,
		CommandArgs:     commandArgsJSON,
		Icon:            req.Icon,
		IsPublic:        req.IsPublic,
	}

	if err := h.templateService.UpdateTemplate(uint(id), template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template updated"})
}

// DeleteTemplate 删除模板
func (h *CommandTemplateHandler) DeleteTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	if err := h.templateService.DeleteTemplate(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted"})
}
