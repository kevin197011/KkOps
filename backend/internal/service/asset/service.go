// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package asset

import (
	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
)

// Service handles asset management business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new asset service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateAssetRequest represents a request to create an asset
type CreateAssetRequest struct {
	HostName        string `json:"hostName" binding:"required"`
	ProjectID       *uint  `json:"project_id"`
	CloudPlatformID *uint  `json:"cloud_platform_id"`
	EnvironmentID   *uint  `json:"environment_id"`
	IP            string `json:"ip"`
	SSHPort       int    `json:"ssh_port"`
	SSHKeyID      *uint  `json:"ssh_key_id"`
	SSHUser       string `json:"ssh_user"`
	CPU           string `json:"cpu"`
	Memory        string `json:"memory"`
	Disk          string `json:"disk"`
	Status        string `json:"status"`
	Description   string `json:"description"`
	TagIDs        []uint `json:"tag_ids"`
}

// UpdateAssetRequest represents a request to update an asset
type UpdateAssetRequest struct {
	HostName        string `json:"hostName"`
	ProjectID       *uint  `json:"project_id"`
	CloudPlatformID *uint  `json:"cloud_platform_id"`
	EnvironmentID   *uint  `json:"environment_id"`
	IP            string `json:"ip"`
	SSHPort       int    `json:"ssh_port"`
	SSHKeyID      *uint  `json:"ssh_key_id"`
	SSHUser       string `json:"ssh_user"`
	CPU           string `json:"cpu"`
	Memory        string `json:"memory"`
	Disk          string `json:"disk"`
	Status        string `json:"status"`
	Description   string `json:"description"`
	TagIDs        []uint `json:"tag_ids"`
}

// AssetResponse represents an asset response
type AssetResponse struct {
	ID              uint               `json:"id"`
	HostName        string             `json:"hostName"`
	ProjectID       *uint              `json:"project_id"`
	CloudPlatformID *uint              `json:"cloud_platform_id"`
	CloudPlatform   *CloudPlatformInfo `json:"cloud_platform,omitempty"`
	EnvironmentID   *uint              `json:"environment_id"`
	IP            string    `json:"ip"`
	SSHPort       int       `json:"ssh_port"`
	SSHKeyID      *uint     `json:"ssh_key_id"`
	SSHUser       string    `json:"ssh_user"`
	CPU           string    `json:"cpu"`
	Memory        string    `json:"memory"`
	Disk          string    `json:"disk"`
	Status        string    `json:"status"`
	Description   string    `json:"description"`
	Tags          []TagInfo `json:"tags"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

// TagInfo represents tag information in asset response
type TagInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// CloudPlatformInfo represents cloud platform information in asset response
type CloudPlatformInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateAsset creates a new asset
func (s *Service) CreateAsset(req *CreateAssetRequest) (*AssetResponse, error) {
	asset := model.Asset{
		HostName:        req.HostName,
		ProjectID:       req.ProjectID,
		CloudPlatformID: req.CloudPlatformID,
		EnvironmentID:   req.EnvironmentID,
		IP:            req.IP,
		SSHPort:       req.SSHPort,
		SSHKeyID:      req.SSHKeyID,
		SSHUser:       req.SSHUser,
		CPU:           req.CPU,
		Memory:        req.Memory,
		Disk:          req.Disk,
		Status:        req.Status,
		Description:   req.Description,
	}

	if asset.Status == "" {
		asset.Status = "active"
	}
	if asset.SSHPort == 0 {
		asset.SSHPort = 22
	}

	// Create asset
	if err := s.db.Create(&asset).Error; err != nil {
		return nil, err
	}

	// Assign tags if provided
	if len(req.TagIDs) > 0 {
		var tags []model.Tag
		if err := s.db.Where("id IN ?", req.TagIDs).Find(&tags).Error; err == nil {
			s.db.Model(&asset).Association("Tags").Replace(tags)
		}
	}

	return s.getAssetResponse(asset.ID)
}

// GetAsset retrieves an asset by ID
func (s *Service) GetAsset(id uint) (*AssetResponse, error) {
	return s.getAssetResponse(id)
}

// ListAssetsFilter 资产列表查询过滤条件
type ListAssetsFilter struct {
	Page            int
	PageSize        int
	ProjectID       *uint
	CloudPlatformID *uint
	EnvironmentID   *uint
	Status          string
	IP              string
	SSHKeyID        *uint
	TagIDs          []uint
	Search          string // Search in hostName, IP, description
	// 权限过滤
	AllowedAssetIDs []uint // 允许访问的资产ID列表，nil 表示不限制（管理员）
	IsAdmin         bool   // 是否为管理员，管理员可以访问所有资产
}

// ListAssets retrieves assets with filtering and pagination
// 支持权限过滤：非管理员用户只能看到已授权的资产
func (s *Service) ListAssets(filter ListAssetsFilter) ([]AssetResponse, int64, error) {
	query := s.db.Model(&model.Asset{})

	// 权限过滤：非管理员需要限制在允许的资产范围内
	if !filter.IsAdmin {
		if len(filter.AllowedAssetIDs) == 0 {
			// 非管理员且没有任何授权，返回空结果
			return []AssetResponse{}, 0, nil
		}
		query = query.Where("id IN ?", filter.AllowedAssetIDs)
	}

	// Apply filters
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.CloudPlatformID != nil {
		query = query.Where("cloud_platform_id = ?", *filter.CloudPlatformID)
	}
	if filter.EnvironmentID != nil {
		query = query.Where("environment_id = ?", *filter.EnvironmentID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.IP != "" {
		query = query.Where("ip LIKE ?", "%"+filter.IP+"%")
	}
	if filter.SSHKeyID != nil {
		query = query.Where("ssh_key_id = ?", *filter.SSHKeyID)
	}

	// Search
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("host_name LIKE ? OR ip LIKE ? OR description LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Tag filter (requires join)
	if len(filter.TagIDs) > 0 {
		query = query.Joins("JOIN asset_tags ON asset_tags.asset_id = assets.id").
			Where("asset_tags.tag_id IN ?", filter.TagIDs).
			Group("assets.id")
	}

	// Count total
	var total int64
	query.Count(&total)

	// Pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}
	offset := (filter.Page - 1) * filter.PageSize

	// Fetch assets
	var assets []model.Asset
	if err := query.Preload("Tags").Preload("Environment").Offset(offset).Limit(filter.PageSize).Find(&assets).Error; err != nil {
		return nil, 0, err
	}

	// Convert to response
	result := make([]AssetResponse, len(assets))
	for i, asset := range assets {
		result[i] = s.assetToResponse(asset)
	}

	return result, total, nil
}

// UpdateAsset updates an asset
func (s *Service) UpdateAsset(id uint, req *UpdateAssetRequest) (*AssetResponse, error) {
	var asset model.Asset
	if err := s.db.Preload("Tags").First(&asset, id).Error; err != nil {
		return nil, err
	}

	// Update fields
	if req.HostName != "" {
		asset.HostName = req.HostName
	}
	if req.ProjectID != nil {
		asset.ProjectID = req.ProjectID
	}
	if req.CloudPlatformID != nil {
		asset.CloudPlatformID = req.CloudPlatformID
	}
	if req.EnvironmentID != nil {
		asset.EnvironmentID = req.EnvironmentID
	}
	if req.IP != "" {
		asset.IP = req.IP
	}
	if req.SSHPort != 0 {
		asset.SSHPort = req.SSHPort
	}
	if req.SSHKeyID != nil {
		asset.SSHKeyID = req.SSHKeyID
	}
	if req.SSHUser != "" {
		asset.SSHUser = req.SSHUser
	}
	if req.CPU != "" {
		asset.CPU = req.CPU
	}
	if req.Memory != "" {
		asset.Memory = req.Memory
	}
	if req.Disk != "" {
		asset.Disk = req.Disk
	}
	if req.Status != "" {
		asset.Status = req.Status
	}
	if req.Description != "" {
		asset.Description = req.Description
	}

	if err := s.db.Save(&asset).Error; err != nil {
		return nil, err
	}

	// Update tags if provided
	if req.TagIDs != nil {
		var tags []model.Tag
		if len(req.TagIDs) > 0 {
			s.db.Where("id IN ?", req.TagIDs).Find(&tags)
		}
		s.db.Model(&asset).Association("Tags").Replace(tags)
	}

	return s.getAssetResponse(id)
}

// DeleteAsset deletes an asset
func (s *Service) DeleteAsset(id uint) error {
	return s.db.Delete(&model.Asset{}, id).Error
}

// getAssetResponse retrieves an asset and converts it to response
func (s *Service) getAssetResponse(id uint) (*AssetResponse, error) {
	var asset model.Asset
	if err := s.db.Preload("Tags").Preload("Environment").First(&asset, id).Error; err != nil {
		return nil, err
	}
	response := s.assetToResponse(asset)
	return &response, nil
}

// assetToResponse converts an asset model to response
func (s *Service) assetToResponse(asset model.Asset) AssetResponse {
	tags := make([]TagInfo, len(asset.Tags))
	for i, tag := range asset.Tags {
		tags[i] = TagInfo{
			ID:    tag.ID,
			Name:  tag.Name,
			Color: tag.Color,
		}
	}

	var cloudPlatformInfo *CloudPlatformInfo
	if asset.CloudPlatform != nil {
		cloudPlatformInfo = &CloudPlatformInfo{
			ID:          asset.CloudPlatform.ID,
			Name:        asset.CloudPlatform.Name,
			Description: asset.CloudPlatform.Description,
		}
	}

	return AssetResponse{
		ID:              asset.ID,
		HostName:        asset.HostName,
		ProjectID:       asset.ProjectID,
		CloudPlatformID: asset.CloudPlatformID,
		CloudPlatform:   cloudPlatformInfo,
		EnvironmentID:   asset.EnvironmentID,
		IP:            asset.IP,
		SSHPort:       asset.SSHPort,
		SSHKeyID:      asset.SSHKeyID,
		SSHUser:       asset.SSHUser,
		CPU:           asset.CPU,
		Memory:        asset.Memory,
		Disk:          asset.Disk,
		Status:        asset.Status,
		Description:   asset.Description,
		Tags:          tags,
		CreatedAt:     asset.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     asset.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
