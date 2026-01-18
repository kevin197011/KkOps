// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package authorization

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
)

// Service 授权服务，处理用户资产访问权限（纯角色授权模型）
type Service struct {
	db *gorm.DB
}

// NewService 创建授权服务实例
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// IsAdmin 检查用户是否为管理员
// 管理员角色的 is_admin 字段为 true
func (s *Service) IsAdmin(userID uint) (bool, error) {
	var user model.User
	if err := s.db.Preload("Roles").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	for _, role := range user.Roles {
		if role.IsAdmin {
			return true, nil
		}
	}

	return false, nil
}

// HasAssetAccess 检查用户是否有资产访问权限
// 检查顺序：1.管理员 2.角色授权(role_assets)
func (s *Service) HasAssetAccess(userID, assetID uint) (bool, error) {
	// 1. 先检查是否为管理员
	isAdmin, err := s.IsAdmin(userID)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}

	// 2. 检查 role_assets 表（角色授权，通过 user_roles 关联）
	var roleCount int64
	err = s.db.Model(&model.RoleAsset{}).
		Joins("JOIN user_roles ON user_roles.role_id = role_assets.role_id").
		Where("user_roles.user_id = ? AND role_assets.asset_id = ?", userID, assetID).
		Count(&roleCount).Error
	if err != nil {
		return false, err
	}

	return roleCount > 0, nil
}

// HasMultipleAssetAccess 批量检查用户对多个资产的访问权限
// 返回用户有权限访问的资产ID列表（通过角色授权）
func (s *Service) HasMultipleAssetAccess(userID uint, assetIDs []uint) ([]uint, error) {
	if len(assetIDs) == 0 {
		return []uint{}, nil
	}

	// 先检查是否为管理员
	isAdmin, err := s.IsAdmin(userID)
	if err != nil {
		return nil, err
	}
	if isAdmin {
		return assetIDs, nil // 管理员有权限访问所有
	}

	// 检查 role_assets 表（角色授权）
	var roleAssetIDs []uint
	err = s.db.Model(&model.RoleAsset{}).
		Select("DISTINCT role_assets.asset_id").
		Joins("JOIN user_roles ON user_roles.role_id = role_assets.role_id").
		Where("user_roles.user_id = ? AND role_assets.asset_id IN ?", userID, assetIDs).
		Pluck("role_assets.asset_id", &roleAssetIDs).Error
	if err != nil {
		return nil, err
	}

	return roleAssetIDs, nil
}

// GetUserAssetIDs 获取用户已授权的资产ID列表
// 管理员返回 nil（表示可以访问所有），普通用户返回角色授权的资产ID列表
func (s *Service) GetUserAssetIDs(userID uint) ([]uint, bool, error) {
	// 先检查是否为管理员
	isAdmin, err := s.IsAdmin(userID)
	if err != nil {
		return nil, false, err
	}
	if isAdmin {
		return nil, true, nil // 管理员可以访问所有资产
	}

	// 获取角色授权的资产ID (role_assets through user_roles)
	var roleAssetIDs []uint
	err = s.db.Model(&model.RoleAsset{}).
		Select("DISTINCT role_assets.asset_id").
		Joins("JOIN user_roles ON user_roles.role_id = role_assets.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("role_assets.asset_id", &roleAssetIDs).Error
	if err != nil {
		return nil, false, err
	}

	return roleAssetIDs, false, nil
}

// ==================== 角色资产授权管理 ====================

// RoleAssetInfo 角色已授权资产的简要信息
type RoleAssetInfo struct {
	ID       uint   `json:"id"`
	HostName string `json:"hostname"`
	IP       string `json:"ip"`
}

// GetRoleAssets 获取角色已授权的资产列表
func (s *Service) GetRoleAssets(roleID uint) ([]RoleAssetInfo, error) {
	var roleAssets []model.RoleAsset
	err := s.db.Where("role_id = ?", roleID).Find(&roleAssets).Error
	if err != nil {
		return nil, err
	}

	if len(roleAssets) == 0 {
		return []RoleAssetInfo{}, nil
	}

	// 获取资产ID列表
	assetIDs := make([]uint, len(roleAssets))
	for i, ra := range roleAssets {
		assetIDs[i] = ra.AssetID
	}

	// 批量查询资产信息
	var assets []model.Asset
	err = s.db.Where("id IN ?", assetIDs).Find(&assets).Error
	if err != nil {
		return nil, err
	}

	result := make([]RoleAssetInfo, len(assets))
	for i, asset := range assets {
		result[i] = RoleAssetInfo{
			ID:       asset.ID,
			HostName: asset.HostName,
			IP:       asset.IP,
		}
	}

	return result, nil
}

// GetRoleAssetIDs 获取角色已授权的资产ID列表
func (s *Service) GetRoleAssetIDs(roleID uint) ([]uint, error) {
	var assetIDs []uint
	err := s.db.Model(&model.RoleAsset{}).
		Where("role_id = ?", roleID).
		Pluck("asset_id", &assetIDs).Error
	if err != nil {
		return nil, err
	}
	return assetIDs, nil
}

// GetRoleAssetCount 获取角色已授权的资产数量
func (s *Service) GetRoleAssetCount(roleID uint) (int64, error) {
	var count int64
	err := s.db.Model(&model.RoleAsset{}).
		Where("role_id = ?", roleID).
		Count(&count).Error
	return count, err
}

// GrantRoleAssets 为角色授权资产
func (s *Service) GrantRoleAssets(roleID uint, assetIDs []uint, createdBy uint) (int, error) {
	if len(assetIDs) == 0 {
		return 0, nil
	}

	var created int
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, assetID := range assetIDs {
			// 检查是否已存在
			var count int64
			tx.Model(&model.RoleAsset{}).
				Where("role_id = ? AND asset_id = ?", roleID, assetID).
				Count(&count)
			if count > 0 {
				continue // 已存在，跳过
			}

			// 创建授权记录
			ra := model.RoleAsset{
				RoleID:    roleID,
				AssetID:   assetID,
				CreatedBy: createdBy,
			}
			if err := tx.Create(&ra).Error; err != nil {
				return err
			}
			created++
		}
		return nil
	})

	return created, err
}

// RevokeRoleAssets 批量撤销角色的资产授权
func (s *Service) RevokeRoleAssets(roleID uint, assetIDs []uint) (int, error) {
	if len(assetIDs) == 0 {
		return 0, nil
	}

	result := s.db.Where("role_id = ? AND asset_id IN ?", roleID, assetIDs).
		Delete(&model.RoleAsset{})
	return int(result.RowsAffected), result.Error
}

// RevokeSingleRoleAsset 撤销角色的单个资产授权
func (s *Service) RevokeSingleRoleAsset(roleID, assetID uint) error {
	return s.db.Where("role_id = ? AND asset_id = ?", roleID, assetID).
		Delete(&model.RoleAsset{}).Error
}

// ==================== 用户角色管理 ====================

// UserRoleInfo 用户角色信息
type UserRoleInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsAdmin     bool   `json:"is_admin"`
}

// GetUserRoles 获取用户的角色列表
func (s *Service) GetUserRoles(userID uint) ([]UserRoleInfo, error) {
	var user model.User
	err := s.db.Preload("Roles").First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []UserRoleInfo{}, nil
		}
		return nil, err
	}

	result := make([]UserRoleInfo, len(user.Roles))
	for i, role := range user.Roles {
		result[i] = UserRoleInfo{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			IsAdmin:     role.IsAdmin,
		}
	}

	return result, nil
}

// GetUserRoleIDs 获取用户的角色ID列表
func (s *Service) GetUserRoleIDs(userID uint) ([]uint, error) {
	var roleIDs []uint
	err := s.db.Model(&model.UserRole{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs).Error
	if err != nil {
		return nil, err
	}
	return roleIDs, nil
}

// SetUserRoles 设置用户的角色（全量替换）
func (s *Service) SetUserRoles(userID uint, roleIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除现有角色关联
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		// 添加新的角色关联
		for _, roleID := range roleIDs {
			userRole := model.UserRole{
				UserID: userID,
				RoleID: roleID,
			}
			if err := tx.Create(&userRole).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// AddUserRole 为用户添加单个角色
func (s *Service) AddUserRole(userID, roleID uint) error {
	// 检查是否已存在
	var count int64
	s.db.Model(&model.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count)
	if count > 0 {
		return nil // 已存在，无需添加
	}

	userRole := model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return s.db.Create(&userRole).Error
}

// RemoveUserRole 移除用户的单个角色
func (s *Service) RemoveUserRole(userID, roleID uint) error {
	return s.db.Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&model.UserRole{}).Error
}
