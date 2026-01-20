// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package role

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/service/role"
)

// Handler handles role and permission management HTTP requests
type Handler struct {
	service *role.Service
}

// NewHandler creates a new role handler
func NewHandler(service *role.Service) *Handler {
	return &Handler{service: service}
}

// CreateRole handles role creation
// @Summary Create role
// @Description Create a new role
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body role.CreateRoleRequest true "Create role request"
// @Success 201 {object} role.RoleResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/roles [post]
func (h *Handler) CreateRole(c *gin.Context) {
	var req role.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateRole(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetRole handles role retrieval
// @Summary Get role
// @Description Get role by ID
// @Tags roles
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} role.RoleResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/roles/{id} [get]
func (h *Handler) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	resp, err := h.service.GetRole(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListRoles handles role list retrieval
// @Summary List roles
// @Description Get list of all roles
// @Tags roles
// @Produce json
// @Security BearerAuth
// @Success 200 {array} role.RoleResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/roles [get]
func (h *Handler) ListRoles(c *gin.Context) {
	roles, err := h.service.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// UpdateRole handles role update
// @Summary Update role
// @Description Update role information
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param request body role.UpdateRoleRequest true "Update role request"
// @Success 200 {object} role.RoleResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/roles/{id} [put]
func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	var req role.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.UpdateRole(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteRole handles role deletion
// @Summary Delete role
// @Description Delete a role by ID
// @Tags roles
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/roles/{id} [delete]
func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	if err := h.service.DeleteRole(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignRoleToUser handles role assignment to user
// @Summary Assign role to user
// @Description Assign a role to a user (deprecated, use /users/{id}/roles instead)
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body object true "Role assignment request" SchemaExample({"role_id": 1})
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /api/v1/users/{id}/roles [post]
func (h *Handler) AssignRoleToUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignRoleToUser(uint(userID), req.RoleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignPermissionToRole handles permission assignment to role
// @Summary Assign permission to role
// @Description Assign a permission to a role
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param request body object true "Permission assignment request" SchemaExample({"permission_id": 1})
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /api/v1/roles/{id}/permissions [post]
func (h *Handler) AssignPermissionToRole(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	var req struct {
		PermissionID uint `json:"permission_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignPermissionToRole(uint(roleID), req.PermissionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListPermissions handles permission list retrieval
// @Summary List permissions
// @Description Get list of all permissions
// @Tags roles
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Permission
// @Failure 500 {object} map[string]string
// @Router /api/v1/permissions [get]
func (h *Handler) ListPermissions(c *gin.Context) {
	permissions, err := h.service.ListPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// GetRolePermissions handles role permissions retrieval
// @Summary Get role permissions
// @Description Get list of permissions for a role
// @Tags roles
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {array} model.Permission
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/roles/{id}/permissions [get]
func (h *Handler) GetRolePermissions(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	permissions, err := h.service.GetRolePermissions(uint(roleID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// RemovePermissionFromRole handles permission removal from role
// @Summary Remove permission from role
// @Description Remove a permission from a role
// @Tags roles
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param permission_id path int true "Permission ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /api/v1/roles/{id}/permissions/{permission_id} [delete]
func (h *Handler) RemovePermissionFromRole(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	permissionID, err := strconv.ParseUint(c.Param("permission_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission ID"})
		return
	}

	if err := h.service.RemovePermissionFromRole(uint(roleID), uint(permissionID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
