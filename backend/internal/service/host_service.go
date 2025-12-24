package service

import (
	"fmt"
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

// validateEnvironment 验证环境字段值
func validateEnvironment(env string) bool {
	if env == "" {
		return true // 环境字段是可选的
	}
	validEnvs := []string{"dev", "uat", "staging", "prod", "test"}
	for _, validEnv := range validEnvs {
		if env == validEnv {
			return true
		}
	}
	return false
}

func (s *hostService) CreateHost(host *models.Host) error {
	// 验证环境字段
	if host.Environment != "" && !validateEnvironment(host.Environment) {
		return fmt.Errorf("invalid environment value: %s. Valid values are: dev, uat, staging, prod, test", host.Environment)
	}
	// 确保 Metadata 字段是有效的 JSON（jsonb 类型不能是空字符串）
	if host.Metadata == "" {
		host.Metadata = "{}"
	}
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
	// 验证环境字段
	if host.Environment != "" && !validateEnvironment(host.Environment) {
		return fmt.Errorf("invalid environment value: %s. Valid values are: dev, uat, staging, prod, test", host.Environment)
	}
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
	// 更新环境字段（支持设置为空字符串）
	if host.Environment != "" {
		existingHost.Environment = host.Environment
	} else {
		// 允许清空环境字段（当传入空字符串时）
		existingHost.Environment = ""
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

// SyncHostStatus 同步单个主机的状态和详细信息（通过 Salt）
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

		// 获取 Grains 信息来填充主机详细信息
		grains, err := s.saltClient.GetGrains(host.SaltMinionID)
		if err == nil {
			s.fillHostFromGrains(host, grains)
		} else {
			log.Printf("Failed to get grains for %s: %v", host.SaltMinionID, err)
		}
	} else {
		host.Status = "offline"
		host.LastSeenAt = nil
	}

	return s.hostRepo.Update(host)
}

// fillHostFromGrains 从 Salt Grains 填充主机信息
func (s *hostService) fillHostFromGrains(host *models.Host, grains map[string]interface{}) {
	// Salt API 返回的格式是: { "minion-id": { grains data } }
	// 获取 Minion 的 Grains 数据
	var minionGrains map[string]interface{}
	
	// 尝试直接获取 minion ID 对应的数据
	if mg, ok := grains[host.SaltMinionID].(map[string]interface{}); ok {
		minionGrains = mg
	} else {
		// 如果没有找到，尝试遍历所有 key
		for key, value := range grains {
			if key == host.SaltMinionID {
				if mg, ok := value.(map[string]interface{}); ok {
					minionGrains = mg
					break
				}
			}
		}
	}
	
	if minionGrains == nil {
		log.Printf("Failed to extract grains data for minion %s", host.SaltMinionID)
		return
	}

	// 填充操作系统信息
	if osType, ok := minionGrains["os"].(string); ok && osType != "" {
		host.OSType = osType
	}
	if osVersion, ok := minionGrains["osrelease"].(string); ok && osVersion != "" {
		host.OSVersion = osVersion
	}

	// 填充 CPU 核心数
	if numCpus, ok := minionGrains["num_cpus"].(float64); ok && numCpus > 0 {
		cpuCores := int(numCpus)
		host.CPUCores = &cpuCores
	}

	// 填充内存信息（GB）
	if memTotal, ok := minionGrains["mem_total"].(float64); ok && memTotal > 0 {
		// Salt 返回的内存通常是 MB，转换为 GB
		memoryGB := memTotal / 1024.0
		host.MemoryGB = &memoryGB
	}

	// 填充磁盘信息（GB）
	if diskTotal, ok := minionGrains["disk_total"].(float64); ok && diskTotal > 0 {
		// Salt 返回的磁盘通常是 MB，转换为 GB
		diskGB := diskTotal / 1024.0
		host.DiskGB = &diskGB
	}

	// 填充 Salt 版本
	if saltVersion, ok := minionGrains["saltversion"].(string); ok && saltVersion != "" {
		host.SaltVersion = saltVersion
	}

	// 填充主机名（如果为空）
	if host.Hostname == "" {
		if hostname, ok := minionGrains["host"].(string); ok && hostname != "" {
			host.Hostname = hostname
		} else if fqdn, ok := minionGrains["fqdn"].(string); ok && fqdn != "" {
			host.Hostname = fqdn
		}
	}

	// 填充 IP 地址（如果为空）
	if host.IPAddress == "" {
		if ipv4, ok := minionGrains["ipv4"].([]interface{}); ok && len(ipv4) > 0 {
			// 取第一个非回环地址
			for _, ip := range ipv4 {
				if ipStr, ok := ip.(string); ok && ipStr != "127.0.0.1" && ipStr != "::1" {
					host.IPAddress = ipStr
					break
				}
			}
		}
	}
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
