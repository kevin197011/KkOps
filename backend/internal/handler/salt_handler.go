package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kronos/backend/internal/salt"
)

type SaltHandler struct {
	saltClient *salt.Client
}

func NewSaltHandler(saltClient *salt.Client) *SaltHandler {
	return &SaltHandler{saltClient: saltClient}
}

// ExecuteCommand 执行 Salt 命令
func (h *SaltHandler) ExecuteCommand(c *gin.Context) {
	if h.saltClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Salt API not configured"})
		return
	}

	var req struct {
		Target   string        `json:"target" binding:"required"`   // 目标Minion，支持通配符
		Function string        `json:"function" binding:"required"` // Salt函数
		Args     []interface{} `json:"args"`                        // 函数参数
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.saltClient.ExecuteCommand(req.Target, req.Function, req.Args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

// TestPing 测试 Minion 连通性
func (h *SaltHandler) TestPing(c *gin.Context) {
	if h.saltClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Salt API not configured"})
		return
	}

	var req struct {
		Target string `json:"target" binding:"required"` // 目标Minion
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.saltClient.TestPing(req.Target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

// GetGrains 获取 Minion 的 Grains 信息
func (h *SaltHandler) GetGrains(c *gin.Context) {
	if h.saltClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Salt API not configured"})
		return
	}

	var req struct {
		Target string `json:"target" binding:"required"` // 目标Minion
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.saltClient.GetGrains(req.Target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

