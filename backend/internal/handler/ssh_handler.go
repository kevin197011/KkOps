package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kronos/backend/internal/models"
	"github.com/kronos/backend/internal/service"
)

type SSHHandler struct {
	connService   service.SSHConnectionService
	keyService    service.SSHKeyService
	sessionService service.SSHSessionService
}

func NewSSHHandler(
	connService service.SSHConnectionService,
	keyService service.SSHKeyService,
	sessionService service.SSHSessionService,
) *SSHHandler {
	return &SSHHandler{
		connService:    connService,
		keyService:     keyService,
		sessionService: sessionService,
	}
}

// SSH Connection Handlers
func (h *SSHHandler) CreateConnection(c *gin.Context) {
	var conn models.SSHConnection
	if err := c.ShouldBindJSON(&conn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.connService.CreateConnection(&conn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"connection": conn})
}

func (h *SSHHandler) GetConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid connection id"})
		return
	}

	conn, err := h.connService.GetConnection(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "connection not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"connection": conn})
}

func (h *SSHHandler) ListConnections(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if projectID := c.Query("project_id"); projectID != "" {
		if id, err := strconv.ParseUint(projectID, 10, 64); err == nil {
			filters["project_id"] = id
		}
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if hostID := c.Query("host_id"); hostID != "" {
		if id, err := strconv.ParseUint(hostID, 10, 64); err == nil {
			filters["host_id"] = id
		}
	}

	conns, total, err := h.connService.ListConnections(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"connections": conns,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
	})
}

func (h *SSHHandler) UpdateConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid connection id"})
		return
	}

	var conn models.SSHConnection
	if err := c.ShouldBindJSON(&conn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.connService.UpdateConnection(id, &conn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "connection updated successfully"})
}

func (h *SSHHandler) DeleteConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid connection id"})
		return
	}

	if err := h.connService.DeleteConnection(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "connection deleted successfully"})
}

// SSH Key Handlers
func (h *SSHHandler) CreateKey(c *gin.Context) {
	var key models.SSHKey
	if err := c.ShouldBindJSON(&key); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.keyService.CreateKey(&key); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"key": key})
}

func (h *SSHHandler) GetKey(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key id"})
		return
	}

	key, err := h.keyService.GetKey(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key})
}

func (h *SSHHandler) ListKeys(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	keys, total, err := h.keyService.ListKeys(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"keys":      keys,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *SSHHandler) UpdateKey(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key id"})
		return
	}

	var key models.SSHKey
	if err := c.ShouldBindJSON(&key); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.keyService.UpdateKey(id, &key); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "key updated successfully"})
}

func (h *SSHHandler) DeleteKey(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key id"})
		return
	}

	if err := h.keyService.DeleteKey(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "key deleted successfully"})
}

// SSH Session Handlers
func (h *SSHHandler) CreateSession(c *gin.Context) {
	var session models.SSHSession
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	session.UserID = userID.(uint64)

	if err := h.sessionService.CreateSession(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"session": session})
}

func (h *SSHHandler) GetSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	session, err := h.sessionService.GetSession(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"session": session})
}

func (h *SSHHandler) ListSessions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if userID := c.Query("user_id"); userID != "" {
		if id, err := strconv.ParseUint(userID, 10, 64); err == nil {
			filters["user_id"] = id
		}
	}
	if connectionID := c.Query("connection_id"); connectionID != "" {
		if id, err := strconv.ParseUint(connectionID, 10, 64); err == nil {
			filters["connection_id"] = id
		}
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	sessions, total, err := h.sessionService.ListSessions(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions":  sessions,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *SSHHandler) CloseSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	if err := h.sessionService.CloseSession(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session closed successfully"})
}

