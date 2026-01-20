// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package role

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kkops/backend/internal/service/authorization"
)

// AssetHandler 角色资产授权处理器
type AssetHandler struct {
	authzService *authorization.Service
}

// NewAssetHandler 创建角色资产授权处理器
func NewAssetHandler(authzService *authorization.Service) *AssetHandler {
	return &AssetHandler{authzService: authzService}
}

// GetRoleAssets 获取角色授权的资产列表
// @Summary 获取角色授权的资产列表
// @Tags Roles
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/roles/{id}/assets [get]
func (h *AssetHandler) GetRoleAssets(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	assets, err := h.authzService.GetRoleAssets(uint(roleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取角色授权资产失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assets})
}

// GrantRoleAssetsRequest 角色资产授权请求
type GrantRoleAssetsRequest struct {
	AssetIDs []uint `json:"asset_ids" binding:"required"`
}

// GrantRoleAssets 为角色授权资产
// @Summary 为角色授权资产
// @Tags Roles
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param body body GrantRoleAssetsRequest true "资产ID列表"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/roles/{id}/assets [post]
func (h *AssetHandler) GrantRoleAssets(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var req GrantRoleAssetsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy := c.MustGet("user_id").(uint)
	created, err := h.authzService.GrantRoleAssets(uint(roleID), req.AssetIDs, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "授权失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "授权成功",
		"created": created,
	})
}

// RevokeRoleAssetsRequest 撤销角色资产授权请求
type RevokeRoleAssetsRequest struct {
	AssetIDs []uint `json:"asset_ids" binding:"required"`
}

// RevokeRoleAssets 批量撤销角色的资产授权
// @Summary 批量撤销角色的资产授权
// @Tags Roles
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param body body RevokeRoleAssetsRequest true "资产ID列表"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/roles/{id}/assets [delete]
func (h *AssetHandler) RevokeRoleAssets(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var req RevokeRoleAssetsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	revoked, err := h.authzService.RevokeRoleAssets(uint(roleID), req.AssetIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "撤销授权失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "撤销成功",
		"revoked": revoked,
	})
}

// RevokeSingleRoleAsset 撤销角色对单个资产的授权
// @Summary 撤销角色对单个资产的授权
// @Tags Roles
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param asset_id path int true "资产ID"
// @Success 204 "No Content"
// @Router /api/v1/roles/{id}/assets/{asset_id} [delete]
func (h *AssetHandler) RevokeSingleRoleAsset(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	assetID, err := strconv.ParseUint(c.Param("asset_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的资产ID"})
		return
	}

	if err := h.authzService.RevokeSingleRoleAsset(uint(roleID), uint(assetID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "撤销授权失败"})
		return
	}

	c.Status(http.StatusNoContent)
}
