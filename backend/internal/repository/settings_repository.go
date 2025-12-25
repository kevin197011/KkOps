package repository

import (
	"github.com/kkops/backend/internal/models"
	"gorm.io/gorm"
)

type SettingsRepository interface {
	GetByKey(key string) (*models.SystemSettings, error)
	GetByCategory(category string) ([]models.SystemSettings, error)
	GetAll() ([]models.SystemSettings, error)
	CreateOrUpdate(setting *models.SystemSettings) error
	Delete(key string) error
}

type settingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) SettingsRepository {
	return &settingsRepository{db: db}
}

func (r *settingsRepository) GetByKey(key string) (*models.SystemSettings, error) {
	var setting models.SystemSettings
	err := r.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *settingsRepository) GetByCategory(category string) ([]models.SystemSettings, error) {
	var settings []models.SystemSettings
	err := r.db.Where("category = ?", category).Find(&settings).Error
	return settings, err
}

func (r *settingsRepository) GetAll() ([]models.SystemSettings, error) {
	var settings []models.SystemSettings
	err := r.db.Find(&settings).Error
	return settings, err
}

func (r *settingsRepository) CreateOrUpdate(setting *models.SystemSettings) error {
	// 使用 ON CONFLICT 或先查找再更新
	var existing models.SystemSettings
	err := r.db.Where("key = ?", setting.Key).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		return r.db.Create(setting).Error
	} else if err != nil {
		return err
	}
	
	// 更新现有记录
	existing.Value = setting.Value
	existing.Category = setting.Category
	existing.Description = setting.Description
	existing.UpdatedBy = setting.UpdatedBy
	existing.UpdatedAt = setting.UpdatedAt
	return r.db.Save(&existing).Error
}

func (r *settingsRepository) Delete(key string) error {
	return r.db.Where("key = ?", key).Delete(&models.SystemSettings{}).Error
}

