// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package operationtool

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"github.com/kkops/backend/internal/model"
)

// Service 运维工具服务
type Service struct {
	db *gorm.DB
}

// NewService 创建运维工具服务实例
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// ListOptions 列表查询选项
type ListOptions struct {
	Category *string
	Enabled  *bool
}

// CreateToolRequest 创建工具请求
type CreateToolRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Icon        string `json:"icon"`
	URL         string `json:"url" binding:"required"`
	OrderIndex  int    `json:"order"`
	Enabled     bool   `json:"enabled"`
}

// UpdateToolRequest 更新工具请求
type UpdateToolRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Icon        *string `json:"icon"`
	URL         *string `json:"url"`
	OrderIndex  *int    `json:"order"`
	Enabled     *bool   `json:"enabled"`
}

// ValidateURL 验证 URL 格式
func validateURL(url string) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("URL must start with http:// or https://")
	}
	return nil
}

// List 获取工具列表
func (s *Service) List(opts *ListOptions) ([]model.OperationTool, error) {
	var tools []model.OperationTool
	query := s.db.Model(&model.OperationTool{})

	// 按分类过滤
	if opts != nil && opts.Category != nil && *opts.Category != "" {
		query = query.Where("category = ?", *opts.Category)
	}

	// 按启用状态过滤（默认只返回启用的）
	if opts != nil && opts.Enabled != nil {
		query = query.Where("enabled = ?", *opts.Enabled)
	} else if opts == nil || opts.Enabled == nil {
		// 默认只返回启用的工具
		query = query.Where("enabled = ?", true)
	}

	// 排序：按 order_index 升序，然后按 id 升序
	err := query.Order("order_index ASC, id ASC").Find(&tools).Error
	if err != nil {
		return nil, fmt.Errorf("获取工具列表失败: %w", err)
	}

	return tools, nil
}

// Get 获取单个工具
func (s *Service) Get(id uint) (*model.OperationTool, error) {
	var tool model.OperationTool
	err := s.db.First(&tool, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("工具不存在")
		}
		return nil, fmt.Errorf("获取工具失败: %w", err)
	}
	return &tool, nil
}

// Create 创建工具
func (s *Service) Create(req *CreateToolRequest) (*model.OperationTool, error) {
	// 验证必填字段
	if req.Name == "" {
		return nil, errors.New("工具名称不能为空")
	}
	if req.URL == "" {
		return nil, errors.New("工具 URL 不能为空")
	}

	// 验证 URL 格式
	if err := validateURL(req.URL); err != nil {
		return nil, err
	}

	tool := &model.OperationTool{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Icon:        req.Icon,
		URL:         req.URL,
		OrderIndex:  req.OrderIndex,
		Enabled:     req.Enabled,
	}

	if err := s.db.Create(tool).Error; err != nil {
		return nil, fmt.Errorf("创建工具失败: %w", err)
	}

	return tool, nil
}

// Update 更新工具
func (s *Service) Update(id uint, req *UpdateToolRequest) (*model.OperationTool, error) {
	tool, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != nil {
		tool.Name = *req.Name
	}
	if req.Description != nil {
		tool.Description = *req.Description
	}
	if req.Category != nil {
		tool.Category = *req.Category
	}
	if req.Icon != nil {
		tool.Icon = *req.Icon
	}
	if req.URL != nil {
		// 验证 URL 格式
		if err := validateURL(*req.URL); err != nil {
			return nil, err
		}
		tool.URL = *req.URL
	}
	if req.OrderIndex != nil {
		tool.OrderIndex = *req.OrderIndex
	}
	if req.Enabled != nil {
		tool.Enabled = *req.Enabled
	}

	if err := s.db.Save(tool).Error; err != nil {
		return nil, fmt.Errorf("更新工具失败: %w", err)
	}

	return tool, nil
}

// Delete 删除工具
func (s *Service) Delete(id uint) error {
	result := s.db.Delete(&model.OperationTool{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除工具失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("工具不存在")
	}
	return nil
}
