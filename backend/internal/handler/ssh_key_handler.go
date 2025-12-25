package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/service"
)

type SSHKeyHandler struct {
	sshKeyService service.SSHKeyService
}

func NewSSHKeyHandler(sshKeyService service.SSHKeyService) *SSHKeyHandler {
	return &SSHKeyHandler{sshKeyService: sshKeyService}
}

type CreateSSHKeyRequest struct {
	Name            string `json:"name" binding:"required"`
	Username        string `json:"username"` // 可选：SSH用户名
	PrivateKeyContent string `json:"private_key_content" binding:"required"`
}

type UpdateSSHKeyRequest struct {
	Name            string `json:"name" binding:"required"`
	Username        string `json:"username"`                    // 可选：SSH用户名
	PrivateKeyContent string `json:"private_key_content"`       // 可选：私钥内容（如果提供，将更新私钥、公钥、指纹和密钥类型）
}

func (h *SSHKeyHandler) CreateSSHKey(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req CreateSSHKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key, err := h.sshKeyService.CreateKey(userID.(uint64), req.Name, req.Username, req.PrivateKeyContent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ssh_key": key})
}

func (h *SSHKeyHandler) ListSSHKeys(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	keys, total, err := h.sshKeyService.ListKeys(userID.(uint64), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ssh_keys": keys,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func (h *SSHKeyHandler) GetSSHKey(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key id"})
		return
	}

	key, err := h.sshKeyService.GetKey(id, userID.(uint64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ssh_key": key})
}

func (h *SSHKeyHandler) UpdateSSHKey(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key id"})
		return
	}

	var req UpdateSSHKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key, err := h.sshKeyService.UpdateKey(id, userID.(uint64), req.Name, req.Username, req.PrivateKeyContent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ssh_key": key})
}

func (h *SSHKeyHandler) DeleteSSHKey(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key id"})
		return
	}

	if err := h.sshKeyService.DeleteKey(id, userID.(uint64)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SSH key deleted successfully"})
}

func (h *SSHKeyHandler) GetSSHKeyPrivateKey(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key id"})
		return
	}

	privateKey, err := h.sshKeyService.GetDecryptedPrivateKey(id, userID.(uint64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"private_key": privateKey})
}

