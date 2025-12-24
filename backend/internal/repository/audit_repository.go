package repository

import (
	"time"

	"github.com/kronos/backend/internal/models"
	"gorm.io/gorm"
)

type AuditRepository interface {
	Create(log *models.AuditLog) error
	GetByID(id uint64) (*models.AuditLog, error)
	List(offset, limit int, filters map[string]interface{}) ([]models.AuditLog, int64, error)
	GetByResource(resourceType string, resourceID uint64) ([]models.AuditLog, error)
	GetByUser(userID uint64, limit int) ([]models.AuditLog, error)
	CountByAction(action string, startTime, endTime time.Time) (int64, error)
	DeleteOldLogs(beforeTime time.Time) error
}

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *auditRepository) GetByID(id uint64) (*models.AuditLog, error) {
	var log models.AuditLog
	err := r.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *auditRepository) List(offset, limit int, filters map[string]interface{}) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := r.db.Model(&models.AuditLog{})

	// 应用过滤条件
	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if username, ok := filters["username"]; ok {
		query = query.Where("username LIKE ?", "%"+username.(string)+"%")
	}
	if action, ok := filters["action"]; ok {
		query = query.Where("action = ?", action)
	}
	if resourceType, ok := filters["resource_type"]; ok {
		query = query.Where("resource_type = ?", resourceType)
	}
	if resourceID, ok := filters["resource_id"]; ok {
		query = query.Where("resource_id = ?", resourceID)
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if startTime, ok := filters["start_time"]; ok {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok {
		query = query.Where("created_at <= ?", endTime)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *auditRepository) GetByResource(resourceType string, resourceID uint64) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

func (r *auditRepository) GetByUser(userID uint64, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

func (r *auditRepository) CountByAction(action string, startTime, endTime time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.AuditLog{}).
		Where("action = ? AND created_at >= ? AND created_at <= ?", action, startTime, endTime).
		Count(&count).Error
	return count, err
}

func (r *auditRepository) DeleteOldLogs(beforeTime time.Time) error {
	return r.db.Where("created_at < ?", beforeTime).Delete(&models.AuditLog{}).Error
}

