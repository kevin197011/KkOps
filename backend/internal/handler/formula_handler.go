package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/service"
)

type FormulaHandler struct {
	formulaService service.FormulaService
}

// 安全地获取用户ID
func getUserID(c *gin.Context) uint64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 1 // 系统默认用户ID
	}

	switch v := userID.(type) {
	case uint64:
		return v
	case uint:
		return uint64(v)
	case int:
		return uint64(v)
	case int64:
		return uint64(v)
	default:
		return 1 // 系统默认用户ID
	}
}

func NewFormulaHandler(formulaService service.FormulaService) *FormulaHandler {
	return &FormulaHandler{
		formulaService: formulaService,
	}
}

// Formula管理API
func (h *FormulaHandler) CreateFormula(c *gin.Context) {

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Path        string `json:"path" binding:"required"`
		Repository  string `json:"repository" binding:"required"`
		Icon        string `json:"icon"`
		Tags        []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	formula := &models.Formula{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Path:        req.Path,
		Repository:  req.Repository,
		Icon:        req.Icon,
		CreatedBy:   getUserID(c),
	}

	if err := h.formulaService.CreateFormula(formula); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"formula": formula})
}

func (h *FormulaHandler) ListFormulas(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := make(map[string]interface{})
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if repository := c.Query("repository"); repository != "" {
		filters["repository"] = repository
	}

	formulas, total, err := h.formulaService.ListFormulas(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"formulas": formulas,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func (h *FormulaHandler) GetFormula(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid formula ID"})
		return
	}

	formula, err := h.formulaService.GetFormula(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Formula not found"})
		return
	}

	// 获取参数信息
	params, err := h.formulaService.GetFormulaParameters(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get parameters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"formula":   formula,
		"parameters": params,
	})
}

// Formula参数管理API
func (h *FormulaHandler) UpdateFormulaParameters(c *gin.Context) {
	formulaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid formula ID"})
		return
	}

	var req struct {
		Parameters []models.FormulaParameter `json:"parameters" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.formulaService.UpdateFormulaParameters(uint(formulaID), req.Parameters); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parameters updated successfully"})
}

// Formula模板管理API
func (h *FormulaHandler) CreateFormulaTemplate(c *gin.Context) {

	var req struct {
		FormulaID   uint                   `json:"formula_id" binding:"required"`
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		PillarData  map[string]interface{} `json:"pillar_data"`
		IsPublic    bool                   `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template := &models.FormulaTemplate{
		FormulaID:   req.FormulaID,
		Name:        req.Name,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		CreatedBy:   getUserID(c),
	}

	// 序列化Pillar数据
	if req.PillarData != nil {
		pillarJSON, err := json.Marshal(req.PillarData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pillar data"})
			return
		}
		template.PillarData = pillarJSON
	}

	if err := h.formulaService.CreateFormulaTemplate(template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"template": template})
}

func (h *FormulaHandler) GetFormulaTemplates(c *gin.Context) {
	formulaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid formula ID"})
		return
	}

	templates, err := h.formulaService.GetFormulaTemplates(uint(formulaID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// Formula部署API
func (h *FormulaHandler) CreateDeployment(c *gin.Context) {

	var req struct {
		FormulaID   uint                   `json:"formula_id" binding:"required"`
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		TargetHosts []string               `json:"target_hosts" binding:"required"`
		PillarData  map[string]interface{} `json:"pillar_data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deployment := &models.FormulaDeployment{
		FormulaID:   req.FormulaID,
		Name:        req.Name,
		Description: req.Description,
		StartedBy:   getUserID(c),
	}

	// 序列化目标主机
	targetHostsJSON, err := json.Marshal(req.TargetHosts)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target hosts"})
		return
	}
	deployment.TargetHosts = targetHostsJSON

	// 序列化Pillar数据
	if req.PillarData != nil {
		pillarJSON, err := json.Marshal(req.PillarData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pillar data"})
			return
		}
		deployment.PillarData = pillarJSON
	}

	if err := h.formulaService.CreateDeployment(deployment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"deployment": deployment})
}

func (h *FormulaHandler) ExecuteDeployment(c *gin.Context) {
	deploymentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deployment ID"})
		return
	}

	if err := h.formulaService.ExecuteDeployment(uint(deploymentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deployment started successfully"})
}

func (h *FormulaHandler) ListDeployments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := make(map[string]interface{})
	if formulaID := c.Query("formula_id"); formulaID != "" {
		if id, err := strconv.ParseUint(formulaID, 10, 64); err == nil {
			filters["formula_id"] = uint(id)
		}
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	deployments, total, err := h.formulaService.ListDeployments(page, pageSize, filters)
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

func (h *FormulaHandler) GetDeployment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deployment ID"})
		return
	}

	deployment, err := h.formulaService.GetDeployment(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deployment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deployment": deployment})
}

func (h *FormulaHandler) CancelDeployment(c *gin.Context) {
	deploymentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deployment ID"})
		return
	}

	if err := h.formulaService.CancelDeployment(uint(deploymentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deployment cancelled successfully"})
}

func (h *FormulaHandler) CleanupOldDeployments(c *gin.Context) {
	// 清理1个月前的部署记录
	oneMonthAgo := time.Now().AddDate(0, -1, 0)
	count, err := h.formulaService.CleanupOldDeployments(oneMonthAgo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cleanup completed successfully",
		"count":   count,
	})
}

// Formula仓库管理API
func (h *FormulaHandler) CreateRepository(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		URL       string `json:"url" binding:"required"`
		Branch    string `json:"branch"`
		LocalPath string `json:"local_path"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo := &models.FormulaRepository{
		Name:      req.Name,
		URL:       req.URL,
		Branch:    req.Branch,
		LocalPath: req.LocalPath,
		CreatedBy: getUserID(c),
	}

	if repo.Branch == "" {
		repo.Branch = "master"
	}

	if err := h.formulaService.CreateRepository(repo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"repository": repo})
}

func (h *FormulaHandler) ListRepositories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := make(map[string]interface{})
	if isActive := c.Query("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			filters["is_active"] = active
		}
	}

	repos, total, err := h.formulaService.ListRepositories(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"repositories": repos,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
	})
}

func (h *FormulaHandler) SyncRepository(c *gin.Context) {
	repoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	if err := h.formulaService.SyncRepository(uint(repoID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repository sync completed successfully"})
}

// UpdateRepository 更新仓库配置
func (h *FormulaHandler) UpdateRepository(c *gin.Context) {
	repoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	var req struct {
		Name      string `json:"name" binding:"required"`
		URL       string `json:"url" binding:"required"`
		Branch    string `json:"branch"`
		LocalPath string `json:"local_path"`
		IsActive  bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取现有仓库
	existingRepo, err := h.formulaService.GetRepository(uint(repoID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	// 更新仓库信息
	existingRepo.Name = req.Name
	existingRepo.URL = req.URL
	existingRepo.Branch = req.Branch
	existingRepo.LocalPath = req.LocalPath
	existingRepo.IsActive = req.IsActive

	if err := h.formulaService.UpdateRepository(existingRepo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"repository": existingRepo})
}

// GetRepository 获取单个仓库信息
func (h *FormulaHandler) GetRepository(c *gin.Context) {
	repoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	repo, err := h.formulaService.GetRepository(uint(repoID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"repository": repo})
}
