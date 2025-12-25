package repository

import (
	"github.com/kkops/backend/internal/models"
	"gorm.io/gorm"
)

type HostRepository interface {
	Create(host *models.Host) error
	GetByID(id uint64) (*models.Host, error)
	GetByHostname(hostname string) (*models.Host, error)
	GetBySaltMinionID(minionID string) (*models.Host, error)
	List(offset, limit int, filters map[string]interface{}) ([]models.Host, int64, error)
	Update(host *models.Host) error
	Delete(id uint64) error
	AddToGroup(hostID, groupID uint64) error
	RemoveFromGroup(hostID, groupID uint64) error
	AddTag(hostID, tagID uint64) error
	RemoveTag(hostID, tagID uint64) error
}

type hostRepository struct {
	db *gorm.DB
}

func NewHostRepository(db *gorm.DB) HostRepository {
	return &hostRepository{db: db}
}

func (r *hostRepository) Create(host *models.Host) error {
	// 清空 Tags，避免 GORM 自动处理 many2many 关联（标签通过单独的 API 端点管理）
	host.Tags = nil
	return r.db.Create(host).Error
}

func (r *hostRepository) GetByID(id uint64) (*models.Host, error) {
	var host models.Host
	err := r.db.Preload("Project").Preload("Groups").Preload("Tags").First(&host, id).Error
	return &host, err
}

func (r *hostRepository) GetByHostname(hostname string) (*models.Host, error) {
	var host models.Host
	err := r.db.Where("hostname = ?", hostname).First(&host).Error
	return &host, err
}

func (r *hostRepository) GetBySaltMinionID(minionID string) (*models.Host, error) {
	var host models.Host
	err := r.db.Where("salt_minion_id = ?", minionID).First(&host).Error
	return &host, err
}

func (r *hostRepository) List(offset, limit int, filters map[string]interface{}) ([]models.Host, int64, error) {
	var hosts []models.Host
	var total int64

	query := r.db.Model(&models.Host{})
	
	// 应用过滤器
	if projectID, ok := filters["project_id"].(uint64); ok && projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if hostname, ok := filters["hostname"].(string); ok && hostname != "" {
		query = query.Where("hostname LIKE ?", "%"+hostname+"%")
	}
	if ip, ok := filters["ip_address"].(string); ok && ip != "" {
		query = query.Where("ip_address LIKE ?", "%"+ip+"%")
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if environment, ok := filters["environment"].(string); ok && environment != "" {
		query = query.Where("environment = ?", environment)
	}
	if groupID, ok := filters["group_id"].(uint64); ok && groupID > 0 {
		query = query.Joins("INNER JOIN host_group_members ON hosts.id = host_group_members.host_id").
			Where("host_group_members.group_id = ?", groupID)
	}
	if tagID, ok := filters["tag_id"].(uint64); ok && tagID > 0 {
		query = query.Joins("INNER JOIN host_tag_assignments ON hosts.id = host_tag_assignments.host_id").
			Where("host_tag_assignments.tag_id = ?", tagID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Project").Preload("Groups").Preload("Tags").Offset(offset).Limit(limit).Find(&hosts).Error
	return hosts, total, err
}

func (r *hostRepository) Update(host *models.Host) error {
	// 清空 Tags，避免 GORM 自动处理 many2many 关联（标签通过单独的 API 端点管理）
	host.Tags = nil
	return r.db.Save(host).Error
}

func (r *hostRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Host{}, id).Error
}

func (r *hostRepository) AddToGroup(hostID, groupID uint64) error {
	return r.db.Create(&models.HostGroupMember{
		HostID:  hostID,
		GroupID: groupID,
	}).Error
}

func (r *hostRepository) RemoveFromGroup(hostID, groupID uint64) error {
	return r.db.Where("host_id = ? AND group_id = ?", hostID, groupID).
		Delete(&models.HostGroupMember{}).Error
}

func (r *hostRepository) AddTag(hostID, tagID uint64) error {
	// 先检查记录是否已存在
	var count int64
	err := r.db.Raw(
		"SELECT COUNT(*) FROM host_tag_assignments WHERE host_id = $1 AND tag_id = $2",
		hostID, tagID,
	).Scan(&count).Error
	if err != nil {
		return err
	}
	
	// 如果记录已存在，直接返回（幂等操作）
	if count > 0 {
		return nil
	}
	
	// 使用原生 SQL 插入，避免 GORM 字段映射问题
	// 使用 $1, $2 占位符（PostgreSQL 格式）
	result := r.db.Exec(
		"INSERT INTO host_tag_assignments (host_id, tag_id) VALUES ($1, $2)",
		hostID, tagID,
	)
	return result.Error
}

func (r *hostRepository) RemoveTag(hostID, tagID uint64) error {
	return r.db.Where("host_id = ? AND tag_id = ?", hostID, tagID).
		Delete(&models.HostTagAssignment{}).Error
}

