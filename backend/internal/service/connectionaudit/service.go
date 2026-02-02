// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package connectionaudit

import (
	"time"

	"github.com/kkops/backend/internal/model"
	"gorm.io/gorm"
)

// Service 审计连线服务
type Service struct {
	db *gorm.DB
}

// NewService 创建审计连线服务实例
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateRequest 创建连线记录请求
type CreateRequest struct {
	UserID              uint
	Username            string
	AssetID             uint
	AssetHostname       string
	StartedAt           time.Time
	EndedAt             time.Time
	Transcript          string
	TranscriptTruncated bool
}

// ListRequest 查询连线记录请求
type ListRequest struct {
	Page      int
	PageSize  int
	UserID    *uint
	AssetID   *uint
	StartTime *time.Time
	EndTime   *time.Time
}

// ListResponse 查询连线记录响应（列表不含 Transcript）
type ListResponse struct {
	Total int64                       `json:"total"`
	Data  []model.SSHConnectionRecord `json:"data"`
}

// MaxTranscriptSize 单条录像 Transcript 最大字节数（1MB），超出则截断并标记
const MaxTranscriptSize = 1024 * 1024

// Create 创建连线记录
func (s *Service) Create(req *CreateRequest) error {
	rec := &model.SSHConnectionRecord{
		UserID:              req.UserID,
		Username:            req.Username,
		AssetID:             req.AssetID,
		AssetHostname:       req.AssetHostname,
		StartedAt:           req.StartedAt,
		EndedAt:             req.EndedAt,
		DurationSeconds:     int64(req.EndedAt.Sub(req.StartedAt).Seconds()),
		Transcript:          req.Transcript,
		TranscriptTruncated: req.TranscriptTruncated,
	}
	return s.db.Create(rec).Error
}

// List 分页列表（不含 Transcript 大字段）
func (s *Service) List(req *ListRequest) (*ListResponse, error) {
	query := s.db.Model(&model.SSHConnectionRecord{})

	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.AssetID != nil {
		query = query.Where("asset_id = ?", *req.AssetID)
	}
	if req.StartTime != nil {
		query = query.Where("started_at >= ?", req.StartTime)
	}
	if req.EndTime != nil {
		query = query.Where("ended_at <= ?", req.EndTime)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []model.SSHConnectionRecord
	offset := (req.Page - 1) * req.PageSize
	if err := query.Omit("Transcript").Order("started_at DESC").Offset(offset).Limit(req.PageSize).Find(&list).Error; err != nil {
		return nil, err
	}

	return &ListResponse{Total: total, Data: list}, nil
}

// GetByID 按 ID 获取单条（含 Transcript）
func (s *Service) GetByID(id uint) (*model.SSHConnectionRecord, error) {
	var rec model.SSHConnectionRecord
	if err := s.db.First(&rec, id).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}
