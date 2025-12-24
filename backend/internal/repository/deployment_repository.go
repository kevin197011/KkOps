package repository

import (
	"github.com/kkops/backend/internal/models"
	"gorm.io/gorm"
)

type DeploymentConfigRepository interface {
	Create(config *models.DeploymentConfig) error
	GetByID(id uint64) (*models.DeploymentConfig, error)
	List(offset, limit int, filters map[string]interface{}) ([]models.DeploymentConfig, int64, error)
	Update(config *models.DeploymentConfig) error
	Delete(id uint64) error
}

type deploymentConfigRepository struct {
	db *gorm.DB
}

func NewDeploymentConfigRepository(db *gorm.DB) DeploymentConfigRepository {
	return &deploymentConfigRepository{db: db}
}

func (r *deploymentConfigRepository) Create(config *models.DeploymentConfig) error {
	return r.db.Create(config).Error
}

func (r *deploymentConfigRepository) GetByID(id uint64) (*models.DeploymentConfig, error) {
	var config models.DeploymentConfig
	err := r.db.Preload("Creator").Preload("Deployments").First(&config, id).Error
	return &config, err
}

func (r *deploymentConfigRepository) List(offset, limit int, filters map[string]interface{}) ([]models.DeploymentConfig, int64, error) {
	var configs []models.DeploymentConfig
	var total int64

	query := r.db.Model(&models.DeploymentConfig{})

	// 应用过滤器
	if applicationName, ok := filters["application_name"].(string); ok && applicationName != "" {
		query = query.Where("application_name LIKE ?", "%"+applicationName+"%")
	}
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if environment, ok := filters["environment"].(string); ok && environment != "" {
		query = query.Where("environment = ?", environment)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Creator").Offset(offset).Limit(limit).Find(&configs).Error
	return configs, total, err
}

func (r *deploymentConfigRepository) Update(config *models.DeploymentConfig) error {
	return r.db.Save(config).Error
}

func (r *deploymentConfigRepository) Delete(id uint64) error {
	return r.db.Delete(&models.DeploymentConfig{}, id).Error
}

type DeploymentRepository interface {
	Create(deployment *models.Deployment) error
	GetByID(id uint64) (*models.Deployment, error)
	List(offset, limit int, filters map[string]interface{}) ([]models.Deployment, int64, error)
	Update(deployment *models.Deployment) error
	GetByConfigID(configID uint64, limit int) ([]models.Deployment, error)
}

type deploymentRepository struct {
	db *gorm.DB
}

func NewDeploymentRepository(db *gorm.DB) DeploymentRepository {
	return &deploymentRepository{db: db}
}

func (r *deploymentRepository) Create(deployment *models.Deployment) error {
	return r.db.Create(deployment).Error
}

func (r *deploymentRepository) GetByID(id uint64) (*models.Deployment, error) {
	var deployment models.Deployment
	err := r.db.Preload("Config").Preload("Starter").Preload("RollbackFromDeployment").
		First(&deployment, id).Error
	return &deployment, err
}

func (r *deploymentRepository) List(offset, limit int, filters map[string]interface{}) ([]models.Deployment, int64, error) {
	var deployments []models.Deployment
	var total int64

	query := r.db.Model(&models.Deployment{})

	// 应用过滤器
	if configID, ok := filters["config_id"].(uint64); ok && configID > 0 {
		query = query.Where("config_id = ?", configID)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if version, ok := filters["version"].(string); ok && version != "" {
		query = query.Where("version = ?", version)
	}
	if startedBy, ok := filters["started_by"].(uint64); ok && startedBy > 0 {
		query = query.Where("started_by = ?", startedBy)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Config").Preload("Starter").Order("started_at DESC").
		Offset(offset).Limit(limit).Find(&deployments).Error
	return deployments, total, err
}

func (r *deploymentRepository) Update(deployment *models.Deployment) error {
	return r.db.Save(deployment).Error
}

func (r *deploymentRepository) GetByConfigID(configID uint64, limit int) ([]models.Deployment, error) {
	var deployments []models.Deployment
	err := r.db.Where("config_id = ?", configID).
		Order("started_at DESC").
		Limit(limit).
		Find(&deployments).Error
	return deployments, err
}

type DeploymentVersionRepository interface {
	Create(version *models.DeploymentVersion) error
	GetByID(id uint64) (*models.DeploymentVersion, error)
	List(offset, limit int, filters map[string]interface{}) ([]models.DeploymentVersion, int64, error)
	GetByApplication(applicationName string) ([]models.DeploymentVersion, error)
	GetByApplicationAndVersion(applicationName, version string) (*models.DeploymentVersion, error)
}

type deploymentVersionRepository struct {
	db *gorm.DB
}

func NewDeploymentVersionRepository(db *gorm.DB) DeploymentVersionRepository {
	return &deploymentVersionRepository{db: db}
}

func (r *deploymentVersionRepository) Create(version *models.DeploymentVersion) error {
	return r.db.Create(version).Error
}

func (r *deploymentVersionRepository) GetByID(id uint64) (*models.DeploymentVersion, error) {
	var version models.DeploymentVersion
	err := r.db.Preload("Creator").First(&version, id).Error
	return &version, err
}

func (r *deploymentVersionRepository) List(offset, limit int, filters map[string]interface{}) ([]models.DeploymentVersion, int64, error) {
	var versions []models.DeploymentVersion
	var total int64

	query := r.db.Model(&models.DeploymentVersion{})

	// 应用过滤器
	if applicationName, ok := filters["application_name"].(string); ok && applicationName != "" {
		query = query.Where("application_name = ?", applicationName)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Creator").Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&versions).Error
	return versions, total, err
}

func (r *deploymentVersionRepository) GetByApplication(applicationName string) ([]models.DeploymentVersion, error) {
	var versions []models.DeploymentVersion
	err := r.db.Where("application_name = ?", applicationName).
		Order("created_at DESC").
		Find(&versions).Error
	return versions, err
}

func (r *deploymentVersionRepository) GetByApplicationAndVersion(applicationName, version string) (*models.DeploymentVersion, error) {
	var v models.DeploymentVersion
	err := r.db.Where("application_name = ? AND version = ?", applicationName, version).
		First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

