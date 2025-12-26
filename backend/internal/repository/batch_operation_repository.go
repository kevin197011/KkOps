package repository

import (
	"encoding/json"
	"time"

	"github.com/kkops/backend/internal/models"
	"gorm.io/gorm"
)

type BatchOperationRepository struct {
	db *gorm.DB
}

func NewBatchOperationRepository(db *gorm.DB) *BatchOperationRepository {
	return &BatchOperationRepository{db: db}
}

// GetDB 获取数据库连接（用于复杂查询）
func (r *BatchOperationRepository) GetDB() *gorm.DB {
	return r.db
}

// Create 创建批量操作
func (r *BatchOperationRepository) Create(operation *models.BatchOperation) error {
	return r.db.Create(operation).Error
}

// Get 获取批量操作
func (r *BatchOperationRepository) Get(id uint) (*models.BatchOperation, error) {
	var operation models.BatchOperation
	err := r.db.Preload("StartedByUser").First(&operation, id).Error
	if err != nil {
		return nil, err
	}
	return &operation, nil
}

// List 列出批量操作
func (r *BatchOperationRepository) List(page, pageSize int, filters map[string]interface{}) ([]models.BatchOperation, int64, error) {
	var operations []models.BatchOperation
	var total int64

	query := r.db.Model(&models.BatchOperation{})

	// 默认只查询最近1个月的数据，除非明确指定了时间范围
	hasTimeFilter := false
	if _, ok := filters["start_time"]; ok {
		hasTimeFilter = true
	}
	if _, ok := filters["end_time"]; ok {
		hasTimeFilter = true
	}

	// 如果没有指定时间过滤器，默认查询最近1个月
	if !hasTimeFilter {
		oneMonthAgo := time.Now().AddDate(0, -1, 0)
		query = query.Where("started_at >= ?", oneMonthAgo)
	}

	// 应用过滤器
	if startedBy, ok := filters["started_by"].(uint); ok {
		query = query.Where("started_by = ?", startedBy)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if startTime, ok := filters["start_time"].(time.Time); ok {
		query = query.Where("started_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"].(time.Time); ok {
		query = query.Where("started_at <= ?", endTime)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Preload("StartedByUser").
		Order("started_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&operations).Error

	return operations, total, err
}

// Update 更新批量操作
func (r *BatchOperationRepository) Update(operation *models.BatchOperation) error {
	return r.db.Save(operation).Error
}

// UpdateStatus 更新操作状态
func (r *BatchOperationRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.BatchOperation{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdateResults 更新操作结果
func (r *BatchOperationRepository) UpdateResults(id uint, results map[string]interface{}) error {
	resultsJSON, _ := json.Marshal(results)
	updates := map[string]interface{}{
		"results": string(resultsJSON),
	}
	if results["status"] != nil {
		updates["status"] = results["status"]
	}
	if results["completed_at"] != nil {
		updates["completed_at"] = results["completed_at"]
	}
	if results["duration_seconds"] != nil {
		updates["duration_seconds"] = results["duration_seconds"]
	}
	if results["success_count"] != nil {
		updates["success_count"] = results["success_count"]
	}
	if results["failed_count"] != nil {
		updates["failed_count"] = results["failed_count"]
	}
	if results["error_message"] != nil {
		updates["error_message"] = results["error_message"]
	}
	if results["salt_job_id"] != nil {
		updates["salt_job_id"] = results["salt_job_id"]
	}

	return r.db.Model(&models.BatchOperation{}).
		Where("id = ?", id).
		Updates(updates).Error
}

