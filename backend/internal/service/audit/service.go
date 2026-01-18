// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package audit

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/kkops/backend/internal/model"
	"gorm.io/gorm"
)

// Service 审计日志服务
type Service struct {
	db *gorm.DB
}

// NewService 创建审计服务实例
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateLogRequest 创建审计日志请求
type CreateLogRequest struct {
	UserID       uint
	Username     string
	Action       string
	Module       string
	ResourceID   *uint
	ResourceName string
	Detail       interface{} // 会被序列化为 JSON
	IPAddress    string
	UserAgent    string
	Status       string
	ErrorMsg     string
}

// ListRequest 查询审计日志请求
type ListRequest struct {
	Page      int
	PageSize  int
	UserID    *uint
	Username  string
	Module    string
	Action    string
	Status    string
	StartTime *time.Time
	EndTime   *time.Time
	Keyword   string
}

// ListResponse 查询审计日志响应
type ListResponse struct {
	Total int64             `json:"total"`
	Items []model.AuditLog `json:"items"`
}

// SensitiveFields 敏感字段列表，这些字段不会被记录到审计日志
var SensitiveFields = []string{
	"password",
	"password_hash",
	"old_password",
	"new_password",
	"confirm_password",
	"private_key",
	"passphrase",
	"token",
	"secret",
	"api_key",
}

// CreateLog 创建审计日志
func (s *Service) CreateLog(req *CreateLogRequest) error {
	// 序列化详情为 JSON
	var detailStr string
	if req.Detail != nil {
		// 过滤敏感字段
		filtered := filterSensitiveData(req.Detail)
		detailBytes, err := json.Marshal(filtered)
		if err != nil {
			detailStr = fmt.Sprintf("序列化失败: %v", err)
		} else {
			detailStr = string(detailBytes)
		}
	}

	log := &model.AuditLog{
		UserID:       req.UserID,
		Username:     req.Username,
		Action:       req.Action,
		Module:       req.Module,
		ResourceID:   req.ResourceID,
		ResourceName: req.ResourceName,
		Detail:       detailStr,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
		Status:       req.Status,
		ErrorMsg:     req.ErrorMsg,
		CreatedAt:    time.Now(),
	}

	return s.db.Create(log).Error
}

// ListLogs 查询审计日志
func (s *Service) ListLogs(req *ListRequest) (*ListResponse, error) {
	query := s.db.Model(&model.AuditLog{})

	// 应用筛选条件
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.Username != "" {
		query = query.Where("username ILIKE ?", "%"+req.Username+"%")
	}
	if req.Module != "" {
		query = query.Where("module = ?", req.Module)
	}
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.StartTime != nil {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != nil {
		query = query.Where("created_at <= ?", req.EndTime)
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("resource_name ILIKE ? OR detail ILIKE ?", keyword, keyword)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	var logs []model.AuditLog
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&logs).Error; err != nil {
		return nil, err
	}

	return &ListResponse{
		Total: total,
		Items: logs,
	}, nil
}

// GetLog 获取单条审计日志
func (s *Service) GetLog(id uint) (*model.AuditLog, error) {
	var log model.AuditLog
	if err := s.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// ExportLogs 导出审计日志
func (s *Service) ExportLogs(req *ListRequest, format string, writer io.Writer) error {
	query := s.db.Model(&model.AuditLog{})

	// 应用筛选条件（与 ListLogs 相同）
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.Username != "" {
		query = query.Where("username ILIKE ?", "%"+req.Username+"%")
	}
	if req.Module != "" {
		query = query.Where("module = ?", req.Module)
	}
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.StartTime != nil {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != nil {
		query = query.Where("created_at <= ?", req.EndTime)
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("resource_name ILIKE ? OR detail ILIKE ?", keyword, keyword)
	}

	var logs []model.AuditLog
	if err := query.Order("created_at DESC").Find(&logs).Error; err != nil {
		return err
	}

	switch format {
	case "csv":
		return s.exportCSV(logs, writer)
	case "json":
		return s.exportJSON(logs, writer)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportCSV 导出为 CSV 格式
func (s *Service) exportCSV(logs []model.AuditLog, writer io.Writer) error {
	w := csv.NewWriter(writer)
	defer w.Flush()

	// 写入表头
	headers := []string{"ID", "时间", "用户", "模块", "操作", "资源", "状态", "IP地址", "详情"}
	if err := w.Write(headers); err != nil {
		return err
	}

	// 写入数据
	for _, log := range logs {
		row := []string{
			fmt.Sprintf("%d", log.ID),
			log.CreatedAt.Format("2006-01-02 15:04:05"),
			log.Username,
			log.Module,
			log.Action,
			log.ResourceName,
			log.Status,
			log.IPAddress,
			log.Detail,
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// exportJSON 导出为 JSON 格式
func (s *Service) exportJSON(logs []model.AuditLog, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(logs)
}

// CleanupOldLogs 清理过期日志
func (s *Service) CleanupOldLogs(retentionDays int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	result := s.db.Where("created_at < ?", cutoff).Delete(&model.AuditLog{})
	return result.RowsAffected, result.Error
}

// GetModules 获取所有模块列表
func (s *Service) GetModules() []string {
	return []string{
		string(model.AuditModuleAuth),
		string(model.AuditModuleUser),
		string(model.AuditModuleRole),
		string(model.AuditModuleAsset),
		string(model.AuditModuleTask),
		string(model.AuditModuleTemplate),
		string(model.AuditModuleScheduled),
		string(model.AuditModuleDeployment),
		string(model.AuditModuleSSH),
		string(model.AuditModuleSSHKey),
		string(model.AuditModuleProject),
		string(model.AuditModuleEnv),
		string(model.AuditModuleCloud),
		string(model.AuditModuleTag),
	}
}

// GetActions 获取所有操作类型列表
func (s *Service) GetActions() []string {
	return []string{
		string(model.AuditActionCreate),
		string(model.AuditActionUpdate),
		string(model.AuditActionDelete),
		string(model.AuditActionExecute),
		string(model.AuditActionLogin),
		string(model.AuditActionLogout),
		string(model.AuditActionEnable),
		string(model.AuditActionDisable),
		string(model.AuditActionExport),
		string(model.AuditActionConnect),
	}
}

// filterSensitiveData 过滤敏感数据
func filterSensitiveData(data interface{}) interface{} {
	// 尝试转换为 map
	switch v := data.(type) {
	case map[string]interface{}:
		return filterMapSensitive(v)
	case map[string]string:
		result := make(map[string]interface{})
		for k, val := range v {
			result[k] = val
		}
		return filterMapSensitive(result)
	default:
		// 尝试 JSON 序列化后再解析
		bytes, err := json.Marshal(data)
		if err != nil {
			return data
		}
		var m map[string]interface{}
		if err := json.Unmarshal(bytes, &m); err != nil {
			return data
		}
		return filterMapSensitive(m)
	}
}

// filterMapSensitive 过滤 map 中的敏感字段
func filterMapSensitive(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		// 检查是否为敏感字段
		isSensitive := false
		lowerKey := strings.ToLower(k)
		for _, sf := range SensitiveFields {
			if strings.Contains(lowerKey, sf) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			result[k] = "[FILTERED]"
		} else if nested, ok := v.(map[string]interface{}); ok {
			result[k] = filterMapSensitive(nested)
		} else {
			result[k] = v
		}
	}
	return result
}
