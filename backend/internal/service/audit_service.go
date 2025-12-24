package service

import (
	"encoding/json"

	"github.com/kronos/backend/internal/models"
	"github.com/kronos/backend/internal/repository"
)

type AuditService interface {
	CreateLog(log *models.AuditLog) error
	GetLog(id uint64) (*models.AuditLog, error)
	ListLogs(page, pageSize int, filters map[string]interface{}) ([]models.AuditLog, int64, error)
	GetLogsByResource(resourceType string, resourceID uint64) ([]models.AuditLog, error)
	GetLogsByUser(userID uint64, limit int) ([]models.AuditLog, error)
	LogOperation(userID *uint64, username, action, resourceType string, resourceID *uint64, resourceName string, ipAddress, userAgent string, requestData, responseData, beforeData, afterData interface{}, status string, errorMsg string, durationMs int) error
}

type auditService struct {
	auditRepo repository.AuditRepository
}

func NewAuditService(auditRepo repository.AuditRepository) AuditService {
	return &auditService{auditRepo: auditRepo}
}

func (s *auditService) CreateLog(log *models.AuditLog) error {
	return s.auditRepo.Create(log)
}

func (s *auditService) GetLog(id uint64) (*models.AuditLog, error) {
	return s.auditRepo.GetByID(id)
}

func (s *auditService) ListLogs(page, pageSize int, filters map[string]interface{}) ([]models.AuditLog, int64, error) {
	offset := (page - 1) * pageSize
	return s.auditRepo.List(offset, pageSize, filters)
}

func (s *auditService) GetLogsByResource(resourceType string, resourceID uint64) ([]models.AuditLog, error) {
	return s.auditRepo.GetByResource(resourceType, resourceID)
}

func (s *auditService) GetLogsByUser(userID uint64, limit int) ([]models.AuditLog, error) {
	return s.auditRepo.GetByUser(userID, limit)
}

// LogOperation 记录操作日志（便捷方法）
func (s *auditService) LogOperation(
	userID *uint64,
	username, action, resourceType string,
	resourceID *uint64,
	resourceName string,
	ipAddress, userAgent string,
	requestData, responseData, beforeData, afterData interface{},
	status string,
	errorMsg string,
	durationMs int,
) error {
	log := &models.AuditLog{
		UserID:       userID,
		Username:     username,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ResourceName: resourceName,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Status:       status,
		ErrorMessage: errorMsg,
		DurationMs:   durationMs,
	}

	// 序列化 JSON 数据
	// 对于 jsonb 类型，GORM 需要有效的 JSON 字符串
	// 空字符串会导致 PostgreSQL 解析错误，所以使用 "{}" 作为默认值
	log.RequestData = serializeJSONData(requestData)
	log.ResponseData = serializeJSONData(responseData)
	log.BeforeData = serializeJSONData(beforeData)
	log.AfterData = serializeJSONData(afterData)

	return s.auditRepo.Create(log)
}

// serializeJSONData 序列化数据为有效的 JSON 字符串
// 对于 jsonb 类型，PostgreSQL 需要有效的 JSON，空字符串会导致错误
func serializeJSONData(data interface{}) string {
	if data == nil {
		return "{}" // 使用空对象而不是空字符串
	}

	// 如果是字符串类型
	if str, ok := data.(string); ok {
		if str == "" {
			return "{}" // 空字符串转换为空对象
		}
		// 验证是否是有效的 JSON
		var test interface{}
		if err := json.Unmarshal([]byte(str), &test); err == nil {
			// 已经是有效的 JSON，直接返回
			return str
		}
		// 不是有效的 JSON，需要编码为 JSON 字符串
		if data, err := json.Marshal(str); err == nil {
			return string(data)
		}
		return "{}"
	}

	// 其他类型，直接 JSON 编码
	if data, err := json.Marshal(data); err == nil {
		return string(data)
	}
	return "{}"
}

