package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/service"
)

type HostHandler struct {
	hostService service.HostService
}

func NewHostHandler(hostService service.HostService) *HostHandler {
	return &HostHandler{hostService: hostService}
}

func (h *HostHandler) CreateHost(c *gin.Context) {
	var host models.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.hostService.CreateHost(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"host": host})
}

func (h *HostHandler) GetHost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid host id"})
		return
	}

	host, err := h.hostService.GetHost(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "host not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"host": host})
}

func (h *HostHandler) ListHosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if projectID := c.Query("project_id"); projectID != "" {
		if id, err := strconv.ParseUint(projectID, 10, 64); err == nil {
			filters["project_id"] = id
		}
	}
	if hostname := c.Query("hostname"); hostname != "" {
		filters["hostname"] = hostname
	}
	if ip := c.Query("ip_address"); ip != "" {
		filters["ip_address"] = ip
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if environment := c.Query("environment"); environment != "" {
		filters["environment"] = environment
	}
	if groupID := c.Query("group_id"); groupID != "" {
		if id, err := strconv.ParseUint(groupID, 10, 64); err == nil {
			filters["group_id"] = id
		}
	}
	if tagID := c.Query("tag_id"); tagID != "" {
		if id, err := strconv.ParseUint(tagID, 10, 64); err == nil {
			filters["tag_id"] = id
		}
	}

	hosts, total, err := h.hostService.ListHosts(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hosts":     hosts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *HostHandler) UpdateHost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid host id"})
		return
	}

	var host models.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.hostService.UpdateHost(id, &host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "host updated successfully"})
}

func (h *HostHandler) DeleteHost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid host id"})
		return
	}

	if err := h.hostService.DeleteHost(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "host deleted successfully"})
}

func (h *HostHandler) AddToGroup(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid host id"})
		return
	}

	var req struct {
		GroupID uint64 `json:"group_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.hostService.AddToGroup(hostID, req.GroupID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "host added to group successfully"})
}

func (h *HostHandler) RemoveFromGroup(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid host id"})
		return
	}

	var req struct {
		GroupID uint64 `json:"group_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.hostService.RemoveFromGroup(hostID, req.GroupID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "host removed from group successfully"})
}

func (h *HostHandler) AddTag(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid host id"})
		return
	}

	var req struct {
		TagID uint64 `json:"tag_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.hostService.AddTag(hostID, req.TagID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tag added successfully"})
}

func (h *HostHandler) RemoveTag(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid host id"})
		return
	}

	var req struct {
		TagID uint64 `json:"tag_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.hostService.RemoveTag(hostID, req.TagID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tag removed successfully"})
}

// SyncHostStatus 同步主机状态（通过Salt）
func (h *HostHandler) SyncHostStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid host id"})
		return
	}

	if err := h.hostService.SyncHostStatus(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回更新后的主机信息
	host, err := h.hostService.GetHost(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get updated host"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "host status synced successfully",
		"host":    host,
	})
}

// SyncAllHostsStatus 同步所有主机状态（通过Salt）
func (h *HostHandler) SyncAllHostsStatus(c *gin.Context) {
	if err := h.hostService.SyncAllHostsStatus(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all hosts status synced successfully"})
}

