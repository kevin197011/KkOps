package repository

import (
	"github.com/kkops/backend/internal/models"
	"gorm.io/gorm"
)

type CommandTemplateRepository struct {
	db *gorm.DB
}

func NewCommandTemplateRepository(db *gorm.DB) *CommandTemplateRepository {
	return &CommandTemplateRepository{db: db}
}

// Create 创建命令模板
func (r *CommandTemplateRepository) Create(template *models.CommandTemplate) error {
	return r.db.Create(template).Error
}

// Get 获取命令模板
func (r *CommandTemplateRepository) Get(id uint) (*models.CommandTemplate, error) {
	var template models.CommandTemplate
	err := r.db.Preload("CreatedByUser").First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// List 列出命令模板
func (r *CommandTemplateRepository) List(filters map[string]interface{}) ([]models.CommandTemplate, error) {
	var templates []models.CommandTemplate

	query := r.db.Model(&models.CommandTemplate{})

	// 应用过滤器
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if createdBy, ok := filters["created_by"].(uint); ok {
		query = query.Where("created_by = ?", createdBy)
	}
	if isPublic, ok := filters["is_public"].(bool); ok {
		if isPublic {
			query = query.Where("is_public = ?", true)
		} else {
			// 如果 is_public 为 false，查询用户自己的模板或公开模板
			if createdBy, ok := filters["created_by"].(uint); ok {
				query = query.Where("is_public = ? OR created_by = ?", true, createdBy)
			} else {
				query = query.Where("is_public = ?", true)
			}
		}
	}

	err := query.Preload("CreatedByUser").
		Order("category ASC, usage_count DESC, created_at DESC").
		Find(&templates).Error

	return templates, err
}

// Update 更新命令模板
func (r *CommandTemplateRepository) Update(template *models.CommandTemplate) error {
	return r.db.Save(template).Error
}

// Delete 删除命令模板
func (r *CommandTemplateRepository) Delete(id uint) error {
	return r.db.Delete(&models.CommandTemplate{}, id).Error
}

// IncrementUsageCount 增加使用次数
func (r *CommandTemplateRepository) IncrementUsageCount(id uint) error {
	return r.db.Model(&models.CommandTemplate{}).
		Where("id = ?", id).
		Update("usage_count", gorm.Expr("usage_count + 1")).Error
}

