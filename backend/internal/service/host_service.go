package service

import (
	"log"
	"time"

	"github.com/kronos/backend/internal/models"
	"github.com/kronos/backend/internal/repository"
	"github.com/kronos/backend/internal/salt"
)

type HostService interface {
	CreateHost(host *models.Host) error
	GetHost(id uint64) (*models.Host, error)
	ListHosts(page, pageSize int, filters map[string]interface{}) ([]models.Host, int64, error)
	UpdateHost(id uint64, host *models.Host) error
	DeleteHost(id uint64) error
	AddToGroup(hostID, groupID uint64) error
	RemoveFromGroup(hostID, groupID uint64) error
	AddTag(hostID, tagID uint64) error
	RemoveTag(hostID, tagID uint64) error
	SyncHostStatus(hostID uint64) error
	SyncAllHostsStatus() error
}

type hostService struct {
	hostRepo   repository.HostRepository
	saltClient *salt.Client
}

func NewHostService(hostRepo repository.HostRepository, saltClient *salt.Client) HostService {
	return &hostService{
		hostRepo:   hostRepo,
		saltClient: saltClient,
	}
}

func (s *hostService) CreateHost(host *models.Host) error {
	return s.hostRepo.Create(host)
}

func (s *hostService) GetHost(id uint64) (*models.Host, error) {
	return s.hostRepo.GetByID(id)
}

func (s *hostService) ListHosts(page, pageSize int, filters map[string]interface{}) ([]models.Host, int64, error) {
	offset := (page - 1) * pageSize
	return s.hostRepo.List(offset, pageSize, filters)
}

func (s *hostService) UpdateHost(id uint64, host *models.Host) error {
	existingHost, err := s.hostRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 更新字段
	if host.Hostname != "" {
		existingHost.Hostname = host.Hostname
	}
	if host.IPAddress != "" {
		existingHost.IPAddress = host.IPAddress
	}
	if host.SaltMinionID != "" {
		existingHost.SaltMinionID = host.SaltMinionID
	}
	if host.OSType != "" {
		existingHost.OSType = host.OSType
	}
	if host.OSVersion != "" {
		existingHost.OSVersion = host.OSVersion
	}
	if host.CPUCores != nil {
		existingHost.CPUCores = host.CPUCores
	}
	if host.MemoryGB != nil {
		existingHost.MemoryGB = host.MemoryGB
	}
	if host.DiskGB != nil {
		existingHost.DiskGB = host.DiskGB
	}
	if host.Status != "" {
		existingHost.Status = host.Status
	}
	if host.Description != "" {
		existingHost.Description = host.Description
	}

	return s.hostRepo.Update(existingHost)
}

func (s *hostService) DeleteHost(id uint64) error {
	return s.hostRepo.Delete(id)
}

func (s *hostService) AddToGroup(hostID, groupID uint64) error {
	return s.hostRepo.AddToGroup(hostID, groupID)
}

func (s *hostService) RemoveFromGroup(hostID, groupID uint64) error {
	return s.hostRepo.RemoveFromGroup(hostID, groupID)
}

func (s *hostService) AddTag(hostID, tagID uint64) error {
	return s.hostRepo.AddTag(hostID, tagID)
}

func (s *hostService) RemoveTag(hostID, tagID uint64) error {
	return s.hostRepo.RemoveTag(hostID, tagID)
}

// SyncHostStatus 同步单个主机的状态（通过 Salt）
func (s *hostService) SyncHostStatus(hostID uint64) error {
	if s.saltClient == nil {
		log.Println("Salt client not configured, skipping status sync")
		return nil
	}

	host, err := s.hostRepo.GetByID(hostID)
	if err != nil {
		return err
	}

	// 如果没有 Salt Minion ID，无法查询状态
	if host.SaltMinionID == "" {
		return nil
	}

	// 查询 Minion 状态
	isOnline, err := s.saltClient.GetMinionStatus(host.SaltMinionID)
	if err != nil {
		log.Printf("Failed to get minion status for %s: %v", host.SaltMinionID, err)
		// 查询失败时设置为 unknown，不返回错误
		host.Status = "unknown"
		host.LastSeenAt = nil
		return s.hostRepo.Update(host)
	}

	// 更新主机状态
	if isOnline {
		host.Status = "online"
		now := time.Now()
		host.LastSeenAt = &now
	} else {
		host.Status = "offline"
		host.LastSeenAt = nil
	}

	return s.hostRepo.Update(host)
}

// SyncAllHostsStatus 同步所有主机的状态（通过 Salt）
func (s *hostService) SyncAllHostsStatus() error {
	if s.saltClient == nil {
		log.Println("Salt client not configured, skipping status sync")
		return nil
	}

	// 获取所有有 Salt Minion ID 的主机
	hosts, _, err := s.hostRepo.List(0, 10000, map[string]interface{}{})
	if err != nil {
		return err
	}

	for _, host := range hosts {
		if host.SaltMinionID != "" {
			if err := s.SyncHostStatus(host.ID); err != nil {
				log.Printf("Failed to sync status for host %d: %v", host.ID, err)
				// 继续处理其他主机，不中断
			}
		}
	}

	return nil
}
