package repository

import (
	"github.com/kkops/backend/internal/models"
	"gorm.io/gorm"
)

type HostGroupRepository interface {
	Create(group *models.HostGroup) error
	GetByID(id uint64) (*models.HostGroup, error)
	List(offset, limit int, filters map[string]interface{}) ([]models.HostGroup, int64, error)
	Update(group *models.HostGroup) error
	Delete(id uint64) error
}

type hostGroupRepository struct {
	db *gorm.DB
}

func NewHostGroupRepository(db *gorm.DB) HostGroupRepository {
	return &hostGroupRepository{db: db}
}

func (r *hostGroupRepository) Create(group *models.HostGroup) error {
	return r.db.Create(group).Error
}

func (r *hostGroupRepository) GetByID(id uint64) (*models.HostGroup, error) {
	var group models.HostGroup
	err := r.db.Preload("Project").Preload("Hosts").First(&group, id).Error
	return &group, err
}

func (r *hostGroupRepository) List(offset, limit int, filters map[string]interface{}) ([]models.HostGroup, int64, error) {
	var groups []models.HostGroup
	var total int64

	query := r.db.Model(&models.HostGroup{})

	// 应用过滤器
	if projectID, ok := filters["project_id"].(uint64); ok && projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Project").Preload("Hosts").Offset(offset).Limit(limit).Find(&groups).Error
	return groups, total, err
}

func (r *hostGroupRepository) Update(group *models.HostGroup) error {
	return r.db.Save(group).Error
}

func (r *hostGroupRepository) Delete(id uint64) error {
	return r.db.Delete(&models.HostGroup{}, id).Error
}

