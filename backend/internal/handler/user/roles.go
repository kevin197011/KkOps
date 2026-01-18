// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/service/authorization"
)

// RoleHandler 用户角色管理处理器
type RoleHandler struct {
	authzService *authorization.Service
}

// NewRoleHandler 创建用户角色管理处理器
func NewRoleHandler(authzService *authorization.Service) *RoleHandler {
	return &RoleHandler{authzService: authzService}
}

// GetUserRoles 获取用户的角色列表
// @Summary 获取用户的角色列表
// @Tags Users
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id}/roles [get]
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	roles, err := h.authzService.GetUserRoles(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// SetUserRolesRequest 设置用户角色请求
type SetUserRolesRequest struct {
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

// SetUserRoles 设置用户的角色（全量替换）
// @Summary 设置用户的角色
// @Tags Users
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param body body SetUserRolesRequest true "角色ID列表"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id}/roles [post]
func (h *RoleHandler) SetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req SetUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authzService.SetUserRoles(uint(userID), req.RoleIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "设置用户角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设置成功"})
}

// RemoveUserRole 移除用户的某个角色
// @Summary 移除用户的某个角色
// @Tags Users
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param role_id path int true "角色ID"
// @Success 204 "No Content"
// @Router /users/{id}/roles/{role_id} [delete]
func (h *RoleHandler) RemoveUserRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	roleID, err := strconv.ParseUint(c.Param("role_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	if err := h.authzService.RemoveUserRole(uint(userID), uint(roleID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移除角色失败"})
		return
	}

	c.Status(http.StatusNoContent)
}
