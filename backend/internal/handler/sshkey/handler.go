// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package sshkey

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/service/sshkey"
)

// Handler handles SSH key management HTTP requests
type Handler struct {
	service *sshkey.Service
}

// NewHandler creates a new SSH key handler
func NewHandler(service *sshkey.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ListSSHKeys lists all SSH keys for the current user
// @Summary List SSH keys
// @Description Get list of all SSH keys for the current user
// @Tags ssh-keys
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Response with data array"
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/ssh/keys [get]
func (h *Handler) ListSSHKeys(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	keys, err := h.service.ListSSHKeys(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": keys})
}

// GetSSHKey gets an SSH key by ID
// @Summary Get SSH key
// @Description Get SSH key by ID for the current user
// @Tags ssh-keys
// @Produce json
// @Security BearerAuth
// @Param id path int true "SSH Key ID"
// @Success 200 {object} map[string]interface{} "Response with data object"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/ssh/keys/{id} [get]
func (h *Handler) GetSSHKey(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key ID"})
		return
	}

	key, err := h.service.GetSSHKey(userID.(uint), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SSH key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": key})
}

// CreateSSHKey creates a new SSH key
// @Summary Create SSH key
// @Description Create a new SSH key for the current user
// @Tags ssh-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sshkey.CreateSSHKeyRequest true "Create SSH key request"
// @Success 201 {object} map[string]interface{} "Response with data object"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/ssh/keys [post]
func (h *Handler) CreateSSHKey(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req sshkey.CreateSSHKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key, err := h.service.CreateSSHKey(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": key})
}

// UpdateSSHKey updates an SSH key
// @Summary Update SSH key
// @Description Update an SSH key for the current user
// @Tags ssh-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "SSH Key ID"
// @Param request body sshkey.UpdateSSHKeyRequest true "Update SSH key request"
// @Success 200 {object} map[string]interface{} "Response with data object"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/ssh/keys/{id} [put]
func (h *Handler) UpdateSSHKey(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key ID"})
		return
	}

	var req sshkey.UpdateSSHKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key, err := h.service.UpdateSSHKey(userID.(uint), uint(id), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SSH key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": key})
}

// DeleteSSHKey deletes an SSH key
// @Summary Delete SSH key
// @Description Delete an SSH key for the current user
// @Tags ssh-keys
// @Security BearerAuth
// @Param id path int true "SSH Key ID"
// @Success 200 {object} map[string]string "Response with message"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/ssh/keys/{id} [delete]
func (h *Handler) DeleteSSHKey(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key ID"})
		return
	}

	if err := h.service.DeleteSSHKey(userID.(uint), uint(id)); err != nil {
		if err.Error() == "SSH key not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SSH key deleted successfully"})
}

// TestSSHKey tests an SSH key connection
// @Summary Test SSH key connection
// @Description Test SSH key connection to a host
// @Tags ssh-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "SSH Key ID"
// @Param request body sshkey.TestSSHKeyRequest true "Test SSH key request"
// @Success 200 {object} map[string]string "Response with message"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/ssh/keys/{id}/test [post]
func (h *Handler) TestSSHKey(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key ID"})
		return
	}

	var req sshkey.TestSSHKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.TestSSHKey(userID.(uint), uint(id), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SSH connection test successful"})
}
