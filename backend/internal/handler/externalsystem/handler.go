// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package externalsystem

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	authService "github.com/kkops/backend/internal/service/auth"
	"github.com/kkops/backend/internal/service/externalsystem"
	"github.com/kkops/backend/internal/service/rbac"
)

// Handler handles external system CRUD and launch
type Handler struct {
	svc     *externalsystem.Service
	authSvc *authService.Service
	rbacSvc *rbac.Service
}

// NewHandler creates a new handler
func NewHandler(svc *externalsystem.Service, authSvc *authService.Service, rbacSvc *rbac.Service) *Handler {
	return &Handler{svc: svc, authSvc: authSvc, rbacSvc: rbacSvc}
}

// List returns list of external systems (enabled-only for non-admin listing; secret never returned)
func (h *Handler) List(c *gin.Context) {
	enabledOnly := c.Query("enabled") != "false"
	list, err := h.svc.List(enabledOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// Get returns one external system by id (secret never returned)
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	sys, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	// Do not expose Secret
	out := *sys
	out.Secret = ""
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// Create creates an external system (admin only)
func (h *Handler) Create(c *gin.Context) {
	var req externalsystem.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sys, err := h.svc.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out := *sys
	out.Secret = ""
	c.JSON(http.StatusCreated, gin.H{"data": out})
}

// Update updates an external system (admin only)
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req externalsystem.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sys, err := h.svc.Update(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out := *sys
	out.Secret = ""
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// Delete deletes an external system (admin only)
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// LaunchResponse is returned for launch (redirect URL so frontend can open in new tab or redirect)
type LaunchResponse struct {
	RedirectURL string `json:"redirect_url"`
}

// Launch builds SSO redirect URL with JWT and returns it (frontend redirects or opens in new tab)
func (h *Handler) Launch(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userInfo, err := h.authSvc.GetCurrentUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	permissions, err := h.rbacSvc.GetUserPermissionList(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get permissions"})
		return
	}

	redirectURL, err := h.svc.Launch(uint(id), userID, userInfo.Username, userInfo.Email, userInfo.RealName, userInfo.Roles, permissions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LaunchResponse{RedirectURL: redirectURL})
}
