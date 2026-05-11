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
	LoginUsername       string // 操作用户：KkOps 登录用户
	Username            string // 连线用户：SSH 登录名（如 root）
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

// Create 创建连线记录，返回记录 ID
func (s *Service) Create(req *CreateRequest) (uint, error) {
	rec := &model.SSHConnectionRecord{
		UserID:              req.UserID,
		LoginUsername:       req.LoginUsername,
		Username:            req.Username,
		AssetID:             req.AssetID,
		AssetHostname:       req.AssetHostname,
		StartedAt:           req.StartedAt,
		EndedAt:             req.EndedAt,
		DurationSeconds:     int64(req.EndedAt.Sub(req.StartedAt).Seconds()),
		Transcript:          req.Transcript,
		TranscriptTruncated: req.TranscriptTruncated,
	}
	if err := s.db.Create(rec).Error; err != nil {
		return 0, err
	}
	return rec.ID, nil
}

// UpdateTranscript 仅更新录像内容（用于进行中连线的实时写入，便于「查看」时看到当前已录内容）
func (s *Service) UpdateTranscript(id uint, transcript string, truncated bool) error {
	return s.db.Model(&model.SSHConnectionRecord{}).Where("id = ?", id).Updates(map[string]interface{}{
		"transcript":           transcript,
		"transcript_truncated": truncated,
	}).Error
}

// UpdateOnDisconnect 在断开连接时更新记录（结束时间、时长、录像）
func (s *Service) UpdateOnDisconnect(id uint, startedAt, endedAt time.Time, transcript string, truncated bool) error {
	duration := int64(endedAt.Sub(startedAt).Seconds())
	return s.db.Model(&model.SSHConnectionRecord{}).Where("id = ?", id).Updates(map[string]interface{}{
		"ended_at":             endedAt,
		"duration_seconds":     duration,
		"transcript":           transcript,
		"transcript_truncated": truncated,
	}).Error
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
		// Include active connections (ended_at = started_at) so they always show
		query = query.Where("ended_at <= ? OR ended_at = started_at", req.EndTime)
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

	// Backfill login_username from users for records that have it empty (e.g. old data or migration)
	var needBackfill []uint
	for i := range list {
		if list[i].LoginUsername == "" && list[i].UserID != 0 {
			needBackfill = append(needBackfill, list[i].UserID)
		}
	}
	if len(needBackfill) > 0 {
		seen := make(map[uint]struct{})
		unique := make([]uint, 0, len(needBackfill))
		for _, id := range needBackfill {
			if _, ok := seen[id]; !ok {
				seen[id] = struct{}{}
				unique = append(unique, id)
			}
		}
		var users []struct {
			ID       uint   `gorm:"column:id"`
			Username string `gorm:"column:username"`
		}
		if err := s.db.Table("users").Select("id, username").Where("id IN ?", unique).Find(&users).Error; err == nil {
			userMap := make(map[uint]string)
			for _, u := range users {
				userMap[u.ID] = u.Username
			}
			for i := range list {
				if list[i].LoginUsername == "" && list[i].UserID != 0 {
					if name, ok := userMap[list[i].UserID]; ok {
						list[i].LoginUsername = name
					}
				}
			}
		}
	}

	return &ListResponse{Total: total, Data: list}, nil
}

// GetByID 按 ID 获取单条（含 Transcript）
func (s *Service) GetByID(id uint) (*model.SSHConnectionRecord, error) {
	var rec model.SSHConnectionRecord
	if err := s.db.First(&rec, id).Error; err != nil {
		return nil, err
	}
	if rec.LoginUsername == "" && rec.UserID != 0 {
		var u model.User
		if err := s.db.Select("username").First(&u, rec.UserID).Error; err == nil {
			rec.LoginUsername = u.Username
		}
	}
	return &rec, nil
}
