package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kronos/backend/internal/models"
	"github.com/kronos/backend/internal/service"
)

type HostTagHandler struct {
	tagService service.HostTagService
}

func NewHostTagHandler(tagService service.HostTagService) *HostTagHandler {
	return &HostTagHandler{tagService: tagService}
}

func (h *HostTagHandler) CreateTag(c *gin.Context) {
	var tag models.HostTag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.tagService.CreateTag(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"tag": tag})
}

func (h *HostTagHandler) GetTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
		return
	}

	tag, err := h.tagService.GetTag(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tag": tag})
}

func (h *HostTagHandler) ListTags(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}

	tags, total, err := h.tagService.ListTags(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tags":      tags,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *HostTagHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
		return
	}

	var tag models.HostTag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.tagService.UpdateTag(id, &tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tag updated successfully"})
}

func (h *HostTagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
		return
	}

	if err := h.tagService.DeleteTag(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tag deleted successfully"})
}

