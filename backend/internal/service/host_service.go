package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
	"github.com/kkops/backend/internal/salt"
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
	DiscoverMinions() ([]salt.MinionInfo, error)
	AutoAssignMinion(hostID uint64) (*string, error)
}

type hostService struct {
	hostRepo   repository.HostRepository
	envRepo    repository.EnvironmentRepository
	saltClient *salt.Client
}

func NewHostService(hostRepo repository.HostRepository, envRepo repository.EnvironmentRepository, saltClient *salt.Client) HostService {
	return &hostService{
		hostRepo:   hostRepo,
		envRepo:    envRepo,
		saltClient: saltClient,
	}
}

// validateEnvironment 验证环境字段值（从数据库动态获取有效值）
func (s *hostService) validateEnvironment(env string) bool {
	if env == "" {
		return true // 环境字段是可选的
	}
	// 从数据库获取所有有效的环境名称
	environments, _, err := s.envRepo.List(0, 1000)
	if err != nil {
		// 如果查询失败，回退到默认值
		validEnvs := []string{"dev", "uat", "staging", "prod", "test"}
		for _, validEnv := range validEnvs {
			if env == validEnv {
				return true
			}
		}
		return false
	}
	for _, e := range environments {
		if e.Name == env {
			return true
		}
	}
	return false
}

func (s *hostService) CreateHost(host *models.Host) error {
	// 检查IP地址是否已被使用
	if existingHost, err := s.hostRepo.GetByIPAddress(host.IPAddress); err == nil && existingHost != nil {
		return fmt.Errorf("IP address %s is already in use by host: %s (ID: %d)", host.IPAddress, existingHost.Hostname, existingHost.ID)
	}

	// 验证环境字段
	if host.Environment != "" && !s.validateEnvironment(host.Environment) {
		return fmt.Errorf("invalid environment value: %s", host.Environment)
	}
	// 确保 Metadata 字段是有效的 JSON（jsonb 类型不能是空字符串）
	if host.Metadata == "" {
		host.Metadata = "{}"
	}
	// 如果 SaltMinionID 是空字符串，设置为 nil（NULL），避免唯一约束冲突
	if host.SaltMinionID != nil && *host.SaltMinionID == "" {
		host.SaltMinionID = nil
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
	if host.Environment != "" && !s.validateEnvironment(host.Environment) {
		return fmt.Errorf("invalid environment value: %s", host.Environment)
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
		// 检查IP地址是否已被其他主机使用
		if existingIPHost, err := s.hostRepo.GetByIPAddress(host.IPAddress); err == nil && existingIPHost != nil && existingIPHost.ID != id {
			return fmt.Errorf("IP address %s is already in use by host: %s (ID: %d)", host.IPAddress, existingIPHost.Hostname, existingIPHost.ID)
		}
		existingHost.IPAddress = host.IPAddress
	}
	if host.SaltMinionID != nil {
		// 如果 SaltMinionID 是空字符串，设置为 nil（NULL），避免唯一约束冲突
		if *host.SaltMinionID == "" {
			existingHost.SaltMinionID = nil
		} else {
		existingHost.SaltMinionID = host.SaltMinionID
		}
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
	// 更新云平台字段（支持设置和清空）
	existingHost.CloudPlatformID = host.CloudPlatformID

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

	// 如果没有 Salt Minion ID，尝试自动匹配
	if host.SaltMinionID == nil || *host.SaltMinionID == "" {
		if minionID, err := s.AutoAssignMinion(hostID); err != nil {
			log.Printf("Failed to auto-assign minion for host %d: %v", hostID, err)
			// 自动匹配失败，设置状态为unknown
			host.Status = "unknown"
			host.LastSeenAt = nil
			return s.hostRepo.Update(host)
		} else if minionID != nil {
			log.Printf("Successfully auto-assigned minion ID %s for host %d", *minionID, hostID)
			// 重新获取更新后的主机信息
			if updatedHost, err := s.hostRepo.GetByID(hostID); err == nil {
				host = updatedHost
			} else {
				return err
			}
		} else {
			// 没有找到匹配的minion，设置状态为unknown
			host.Status = "unknown"
			host.LastSeenAt = nil
			return s.hostRepo.Update(host)
		}
	}

	// 查询 Minion 状态
	isOnline, err := s.saltClient.GetMinionStatus(*host.SaltMinionID)
	if err != nil {
		log.Printf("Failed to get minion status for %s: %v", *host.SaltMinionID, err)
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
		grains, err := s.saltClient.GetGrains(*host.SaltMinionID)
		if err == nil {
			s.fillHostFromGrains(host, grains)
		} else {
			log.Printf("Failed to get grains for %s: %v", *host.SaltMinionID, err)
		}

		// 获取磁盘使用情况
		diskUsage, err := s.saltClient.GetDiskUsage(*host.SaltMinionID)
		if err == nil {
			s.fillHostDiskInfo(host, diskUsage)
		} else {
			log.Printf("Failed to get disk usage for %s: %v", *host.SaltMinionID, err)
		}

		// 获取 SSH 端口配置
		sshPort, err := s.getSSHPortFromSalt(*host.SaltMinionID)
		if err == nil && sshPort > 0 {
			host.SSHPort = sshPort
		} else {
			log.Printf("Failed to get SSH port for %s: %v, using default 22", *host.SaltMinionID, err)
			// 使用默认端口
			host.SSHPort = 22
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
	if mg, ok := grains[*host.SaltMinionID].(map[string]interface{}); ok {
		minionGrains = mg
	} else {
		// 如果没有找到，尝试遍历所有 key
		for key, value := range grains {
			if key == *host.SaltMinionID {
				if mg, ok := value.(map[string]interface{}); ok {
					minionGrains = mg
					break
				}
			}
		}
	}
	
	if minionGrains == nil {
		log.Printf("Failed to extract grains data for minion %s", *host.SaltMinionID)
		return
	}

	// 填充主机名（始终更新）
	if hostname, ok := minionGrains["host"].(string); ok && hostname != "" {
		host.Hostname = hostname
	} else if fqdn, ok := minionGrains["fqdn"].(string); ok && fqdn != "" {
		host.Hostname = fqdn
	}

	// 填充操作系统信息（始终更新）
	if osType, ok := minionGrains["os"].(string); ok && osType != "" {
		host.OSType = osType
	}
	if osVersion, ok := minionGrains["osrelease"].(string); ok && osVersion != "" {
		host.OSVersion = osVersion
	}

	// 填充 CPU 核心数（始终更新）
	if numCpus, ok := minionGrains["num_cpus"].(float64); ok && numCpus > 0 {
		cpuCores := int(numCpus)
		host.CPUCores = &cpuCores
	}

	// 填充内存信息（MB -> GB，始终更新）
	if memTotal, ok := minionGrains["mem_total"].(float64); ok && memTotal > 0 {
		// Salt 返回的内存是 MB，转换为 GB
		memoryGB := memTotal / 1024.0
		host.MemoryGB = &memoryGB
	}

	// 填充磁盘信息 - 尝试多种方式获取
	// 方式1: 从 disks grain 获取磁盘列表，然后计算总大小
	// 方式2: 如果有 disk_usage grain
	// 注意：Salt 默认 grains 不包含磁盘总大小，需要通过 disk.usage 模块获取
	// 这里我们暂时跳过磁盘信息，因为需要额外的 Salt 命令调用

	// 填充 Salt 版本（始终更新）
	if saltVersion, ok := minionGrains["saltversion"].(string); ok && saltVersion != "" {
		host.SaltVersion = saltVersion
	}

		// 填充 IP 地址（始终更新，优先选择私有网络 IP）
		if ipv4, ok := minionGrains["ipv4"].([]interface{}); ok && len(ipv4) > 0 {
			// 优先选择私有网络 IP（192.168.x.x, 10.x.x.x, 172.16-31.x.x）
			var privateIP, publicIP string
			for _, ip := range ipv4 {
				if ipStr, ok := ip.(string); ok && ipStr != "127.0.0.1" {
					if isPrivateIP(ipStr) {
						privateIP = ipStr
					} else {
						publicIP = ipStr
					}
				}
			}
			// 优先使用私有 IP
			if privateIP != "" {
				host.IPAddress = privateIP
			} else if publicIP != "" {
				host.IPAddress = publicIP
			}
		}
}

// isPrivateIP 判断是否为私有 IP 地址
func isPrivateIP(ip string) bool {
	// 简单判断私有 IP 范围
	return len(ip) > 0 && (
		ip[:3] == "10." ||
		ip[:4] == "172." ||
		ip[:8] == "192.168.")
}

// getSSHPortFromSalt 通过 Salt 命令获取 SSH 端口
func (s *hostService) getSSHPortFromSalt(minionID string) (int, error) {
	// 使用 Salt cmd.run 执行命令获取 SSH 端口
	// 命令: grep -oP '^Port \K\d+' /etc/ssh/sshd_config || echo 22
	cmd := `grep -oP '^Port \K\d+' /etc/ssh/sshd_config 2>/dev/null || echo 22`

	result, err := s.saltClient.ExecuteCommand(minionID, "cmd.run", []interface{}{cmd})
	if err != nil {
		return 0, fmt.Errorf("failed to execute SSH port check command: %w", err)
	}

	// 解析结果
	if minionResult, ok := result[minionID]; ok {
		if output, ok := minionResult.(string); ok {
			// 清理输出（去除换行符等）
			output = strings.TrimSpace(output)

			// 尝试转换为整数
			if port, err := strconv.Atoi(output); err == nil && port > 0 && port < 65536 {
				return port, nil
			}
		}
	}

	// 如果解析失败，返回默认端口
	return 22, nil
}

// getSSHPortFromGrains 从 Salt Grains 中获取 SSH 端口（备用方法）
func (s *hostService) getSSHPortFromGrains(grains map[string]interface{}) int {
	// 方法1: 检查 sshd grains（如果安装了sshd模块）
	if sshd, ok := grains["sshd"].(map[string]interface{}); ok {
		if port, ok := sshd["Port"].(float64); ok {
			return int(port)
		}
	}

	// 返回默认端口
	return 22
}

// fillHostDiskInfo 从 Salt disk.usage 结果填充磁盘信息
func (s *hostService) fillHostDiskInfo(host *models.Host, diskUsage map[string]interface{}) {
	// disk.usage 返回格式: { "minion-id": { "/": { "1K-blocks": xxx, "available": xxx, ... }, "/boot": {...} } }
	var minionDiskData map[string]interface{}

	// 尝试获取 minion ID 对应的数据
	if data, ok := diskUsage[*host.SaltMinionID].(map[string]interface{}); ok {
		minionDiskData = data
	} else {
		// 如果没有找到，尝试遍历
		for _, value := range diskUsage {
			if data, ok := value.(map[string]interface{}); ok {
				minionDiskData = data
				break
			}
		}
	}

	if minionDiskData == nil {
		return
	}

	// 计算根分区或最大分区的磁盘大小
	var totalDiskKB float64

	// 优先获取根分区 "/"
	if rootData, ok := minionDiskData["/"].(map[string]interface{}); ok {
		if blocks, ok := rootData["1K-blocks"].(string); ok {
			if kb, err := parseFloat(blocks); err == nil {
				totalDiskKB = kb
			}
		}
	}

	// 如果没有根分区数据，累加所有分区
	if totalDiskKB == 0 {
		for mountPoint, data := range minionDiskData {
			// 跳过一些特殊挂载点
			if mountPoint == "/dev" || mountPoint == "/run" || mountPoint == "/sys" ||
				mountPoint == "/proc" || mountPoint == "/dev/shm" {
				continue
			}
			if diskData, ok := data.(map[string]interface{}); ok {
				if blocks, ok := diskData["1K-blocks"].(string); ok {
					if kb, err := parseFloat(blocks); err == nil && kb > totalDiskKB {
						totalDiskKB = kb
					}
				}
			}
		}
	}

	if totalDiskKB > 0 {
		// 转换为 GB (1K-blocks -> GB)
		diskGB := totalDiskKB / 1024.0 / 1024.0
		host.DiskGB = &diskGB
	}
}

// parseFloat 解析字符串为浮点数
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// SyncAllHostsStatus 并发同步所有主机的状态（通过 Salt），最多100个并发
func (s *hostService) SyncAllHostsStatus() error {
	if s.saltClient == nil {
		log.Println("Salt client not configured, skipping status sync")
		return nil
	}

	// 获取所有主机
	hosts, _, err := s.hostRepo.List(0, 10000, map[string]interface{}{})
	if err != nil {
		return err
	}

	const maxConcurrency = 100
	semaphore := make(chan struct{}, maxConcurrency)
	errorChan := make(chan error, len(hosts))

	// 并发处理每个主机
	for i := range hosts {
		host := &hosts[i] // 使用指针避免闭包问题

		go func(h *models.Host) {
			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 如果主机没有 Salt Minion ID，尝试自动匹配
			if h.SaltMinionID == nil || *h.SaltMinionID == "" {
				if minionID, err := s.AutoAssignMinion(h.ID); err != nil {
					log.Printf("Failed to auto-assign minion for host %d: %v", h.ID, err)
					// 自动匹配失败，设置状态为unknown
					h.Status = "unknown"
					h.LastSeenAt = nil
					if updateErr := s.hostRepo.Update(h); updateErr != nil {
						log.Printf("Failed to update host status for %d: %v", h.ID, updateErr)
					}
					errorChan <- fmt.Errorf("auto-assign failed for host %d: %v", h.ID, err)
					return
				} else if minionID != nil {
					log.Printf("Successfully auto-assigned minion ID %s for host %d", *minionID, h.ID)
					// 重新获取更新后的主机信息
					if updatedHost, err := s.hostRepo.GetByID(h.ID); err == nil {
						*h = *updatedHost
					}
				}
			}

			// 同步主机状态（现在应该有Salt Minion ID了）
			if h.SaltMinionID != nil && *h.SaltMinionID != "" {
				if err := s.SyncHostStatus(h.ID); err != nil {
					log.Printf("Failed to sync status for host %d: %v", h.ID, err)
					errorChan <- fmt.Errorf("sync failed for host %d: %v", h.ID, err)
					return
				}
			}

			errorChan <- nil // 成功
		}(host)
	}

	// 等待所有goroutines完成
	var errors []error
	for i := 0; i < len(hosts); i++ {
		if err := <-errorChan; err != nil {
			errors = append(errors, err)
		}
	}

	// 如果有错误，记录但不中断整个同步过程
	if len(errors) > 0 {
		log.Printf("SyncAllHostsStatus completed with %d errors out of %d hosts", len(errors), len(hosts))
		for _, err := range errors {
			log.Printf("Sync error: %v", err)
		}
	} else {
		log.Printf("SyncAllHostsStatus completed successfully for %d hosts", len(hosts))
	}

	return nil
}


// DiscoverMinions 发现所有 Salt Minion
func (s *hostService) DiscoverMinions() ([]salt.MinionInfo, error) {
	if s.saltClient == nil {
		return nil, fmt.Errorf("Salt client not configured")
	}

	return s.saltClient.ListMinions()
}

// AutoAssignMinion 自动为主机分配 Salt Minion ID（通过 IP 地址匹配）
// 匹配成功后会直接使用 /minions 返回的 grains 数据更新主机信息，无需额外 API 调用
func (s *hostService) AutoAssignMinion(hostID uint64) (*string, error) {
	if s.saltClient == nil {
		return nil, fmt.Errorf("Salt client not configured")
	}

	host, err := s.hostRepo.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("failed to get host: %w", err)
	}

	// 如果已经有 Minion ID，直接返回
	if host.SaltMinionID != nil && *host.SaltMinionID != "" {
		return host.SaltMinionID, nil
	}

	// 获取所有 Minion（/minions 端点已返回完整的 grains 数据）
	minions, err := s.saltClient.ListMinions()
	if err != nil {
		return nil, fmt.Errorf("failed to list minions: %w", err)
	}

	var matchedMinion *salt.MinionInfo

	// 通过 IP 地址匹配
	for i := range minions {
		if minions[i].IPAddress == host.IPAddress {
			matchedMinion = &minions[i]
			log.Printf("Auto-assigned minion %s to host %d (IP: %s)", matchedMinion.ID, hostID, host.IPAddress)
			break
		}
	}

	// 如果 IP 没匹配到，尝试通过主机名匹配
	if matchedMinion == nil {
		for i := range minions {
			if minions[i].Hostname == host.Hostname || minions[i].ID == host.Hostname {
				matchedMinion = &minions[i]
				log.Printf("Auto-assigned minion %s to host %d (hostname: %s)", matchedMinion.ID, hostID, host.Hostname)
				break
			}
		}
	}

	// 没有找到匹配的 Minion
	if matchedMinion == nil {
		return nil, nil
	}

	// 使用 grains 数据直接更新主机信息
	host.SaltMinionID = &matchedMinion.ID
	host.Status = "online"
	now := time.Now()
	host.LastSeenAt = &now

	// 从 MinionInfo 填充主机信息
	if matchedMinion.Hostname != "" {
		host.Hostname = matchedMinion.Hostname
	}
	if matchedMinion.OS != "" {
		host.OSType = matchedMinion.OS
	}
	if matchedMinion.OSRelease != "" {
		host.OSVersion = matchedMinion.OSRelease
	}
	if matchedMinion.NumCPUs > 0 {
		host.CPUCores = &matchedMinion.NumCPUs
	}
	if matchedMinion.MemTotalMB > 0 {
		// 转换 MB 为 GB
		memoryGB := matchedMinion.MemTotalMB / 1024.0
		host.MemoryGB = &memoryGB
	}
	if matchedMinion.SaltVersion != "" {
		host.SaltVersion = matchedMinion.SaltVersion
	}
	if matchedMinion.IPAddress != "" {
		host.IPAddress = matchedMinion.IPAddress
	}

	// 磁盘信息需要额外调用 disk.usage，这里先跳过
	// 用户可以点击"同步"按钮获取完整信息（包括磁盘）

	if err := s.hostRepo.Update(host); err != nil {
		return nil, fmt.Errorf("failed to update host: %w", err)
	}

	return &matchedMinion.ID, nil
}
