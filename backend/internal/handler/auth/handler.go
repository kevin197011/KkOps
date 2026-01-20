// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package auth

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/service/auth"
)

// Handler handles authentication HTTP requests
type Handler struct {
	service *auth.Service
}

// NewHandler creates a new authentication handler
func NewHandler(service *auth.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login request"
// @Success 200 {object} auth.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMe returns current user information
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} auth.UserInfo
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/me [get]
func (h *Handler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userInfo, err := h.service.GetCurrentUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

// Logout handles user logout
// @Summary User logout
// @Description Logout the current authenticated user (JWT is stateless, so logout is primarily client-side)
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	// For JWT-based authentication, logout is primarily client-side (token deletion)
	// This endpoint provides API consistency and can be used for logging/logging out events
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// ChangePassword handles password change for current user
// @Summary Change password
// @Description Change password for the current authenticated user
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body auth.ChangePasswordRequest true "Change password request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/change-password [post]
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req auth.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ChangePassword(userID.(uint), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// CreateAPIToken handles API token creation
// @Summary Create API token
// @Description Create a new API token for a user (token is only shown once)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body auth.CreateAPITokenRequest true "Create API token request"
// @Success 201 {object} auth.APITokenResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/users/{id}/tokens [post]
func (h *Handler) CreateAPIToken(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req auth.CreateAPITokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateAPIToken(uint(userID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// ListAPITokens handles API token list retrieval
// @Summary List API tokens
// @Description Get list of API tokens for a user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {array} auth.APITokenResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/{id}/tokens [get]
func (h *Handler) ListAPITokens(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	tokens, err := h.service.ListAPITokens(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// GetAPIToken handles API token retrieval
// @Summary Get API token
// @Description Get full API token by ID (for the token owner)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param token_id path int true "Token ID"
// @Success 200 {object} auth.APITokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/{id}/tokens/{token_id} [get]
func (h *Handler) GetAPIToken(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	tokenID, err := strconv.ParseUint(c.Param("token_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token ID"})
		return
	}

	// Verify the requesting user matches the token owner
	requestingUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if requestingUserID.(uint) != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	resp, err := h.service.GetAPIToken(uint(userID), uint(tokenID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteAPIToken handles API token deletion
// @Summary Delete API token
// @Description Delete an API token by ID
// @Tags users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param token_id path int true "Token ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/{id}/tokens/{token_id} [delete]
func (h *Handler) DeleteAPIToken(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	tokenID, err := strconv.ParseUint(c.Param("token_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token ID"})
		return
	}

	// Verify the requesting user matches the token owner
	requestingUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if requestingUserID.(uint) != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.service.DeleteAPIToken(uint(userID), uint(tokenID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
