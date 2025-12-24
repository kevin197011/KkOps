package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/service"
)

type ProjectHandler struct {
	projectService service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID, _ := c.Get("user_id")
	createdBy := userID.(uint64)

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.projectService.CreateProject(req.Name, req.Description, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"project": project})
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	userID, _ := c.Get("user_id")
	hasAccess, err := h.projectService.CheckUserAccess(id, userID.(uint64))
	if err != nil || !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	project, err := h.projectService.GetProject(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	userID, _ := c.Get("user_id")

	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	projects, total, err := h.projectService.ListProjects(page, pageSize, userID.(uint64), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects":  projects,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	userID, _ := c.Get("user_id")
	hasAccess, err := h.projectService.CheckUserAccess(id, userID.(uint64))
	if err != nil || !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.projectService.UpdateProject(id, req.Name, req.Description, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	userID, _ := c.Get("user_id")
	hasAccess, err := h.projectService.CheckUserAccess(id, userID.(uint64))
	if err != nil || !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	err = h.projectService.DeleteProject(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project deleted successfully"})
}

func (h *ProjectHandler) AddMember(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.projectService.AddMember(projectID, req.UserID, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member added successfully"})
}

func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.projectService.RemoveMember(projectID, req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed successfully"})
}

func (h *ProjectHandler) GetMembers(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	members, err := h.projectService.GetMembers(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

