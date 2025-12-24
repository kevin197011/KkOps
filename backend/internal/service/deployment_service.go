package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/kronos/backend/internal/models"
	"github.com/kronos/backend/internal/repository"
	"github.com/kronos/backend/internal/salt"
)

type DeploymentConfigService interface {
	CreateConfig(config *models.DeploymentConfig) error
	GetConfig(id uint64) (*models.DeploymentConfig, error)
	ListConfigs(page, pageSize int, filters map[string]interface{}) ([]models.DeploymentConfig, int64, error)
	UpdateConfig(id uint64, config *models.DeploymentConfig) error
	DeleteConfig(id uint64) error
}

type deploymentConfigService struct {
	configRepo repository.DeploymentConfigRepository
}

func NewDeploymentConfigService(configRepo repository.DeploymentConfigRepository) DeploymentConfigService {
	return &deploymentConfigService{configRepo: configRepo}
}

func (s *deploymentConfigService) CreateConfig(config *models.DeploymentConfig) error {
	return s.configRepo.Create(config)
}

func (s *deploymentConfigService) GetConfig(id uint64) (*models.DeploymentConfig, error) {
	return s.configRepo.GetByID(id)
}

func (s *deploymentConfigService) ListConfigs(page, pageSize int, filters map[string]interface{}) ([]models.DeploymentConfig, int64, error) {
	offset := (page - 1) * pageSize
	return s.configRepo.List(offset, pageSize, filters)
}

func (s *deploymentConfigService) UpdateConfig(id uint64, config *models.DeploymentConfig) error {
	existing, err := s.configRepo.GetByID(id)
	if err != nil {
		return err
	}

	if config.Name != "" {
		existing.Name = config.Name
	}
	if config.ApplicationName != "" {
		existing.ApplicationName = config.ApplicationName
	}
	if config.Description != "" {
		existing.Description = config.Description
	}
	if config.SaltStateFiles != "" {
		existing.SaltStateFiles = config.SaltStateFiles
	}
	if config.TargetGroups != "" {
		existing.TargetGroups = config.TargetGroups
	}
	if config.TargetHosts != "" {
		existing.TargetHosts = config.TargetHosts
	}
	if config.Environment != "" {
		existing.Environment = config.Environment
	}
	if config.ConfigData != "" {
		existing.ConfigData = config.ConfigData
	}

	return s.configRepo.Update(existing)
}

func (s *deploymentConfigService) DeleteConfig(id uint64) error {
	return s.configRepo.Delete(id)
}

type DeploymentService interface {
	CreateDeployment(deployment *models.Deployment) error
	GetDeployment(id uint64) (*models.Deployment, error)
	ListDeployments(page, pageSize int, filters map[string]interface{}) ([]models.Deployment, int64, error)
	StartDeployment(configID uint64, version string, targetHosts []uint64, startedBy uint64) (*models.Deployment, error)
	UpdateDeploymentStatus(id uint64, status string, results map[string]interface{}, errorMsg string) error
	GetDeploymentStatus(id uint64) (*models.Deployment, error)
	RollbackDeployment(deploymentID uint64, startedBy uint64) (*models.Deployment, error)
}

type deploymentService struct {
	deploymentRepo repository.DeploymentRepository
	configRepo     repository.DeploymentConfigRepository
	hostRepo       repository.HostRepository
	saltClient     *salt.Client
}

func NewDeploymentService(
	deploymentRepo repository.DeploymentRepository,
	configRepo repository.DeploymentConfigRepository,
	hostRepo repository.HostRepository,
	saltClient *salt.Client,
) DeploymentService {
	return &deploymentService{
		deploymentRepo: deploymentRepo,
		configRepo:     configRepo,
		hostRepo:       hostRepo,
		saltClient:     saltClient,
	}
}

func (s *deploymentService) CreateDeployment(deployment *models.Deployment) error {
	return s.deploymentRepo.Create(deployment)
}

func (s *deploymentService) GetDeployment(id uint64) (*models.Deployment, error) {
	return s.deploymentRepo.GetByID(id)
}

func (s *deploymentService) ListDeployments(page, pageSize int, filters map[string]interface{}) ([]models.Deployment, int64, error) {
	offset := (page - 1) * pageSize
	return s.deploymentRepo.List(offset, pageSize, filters)
}

func (s *deploymentService) StartDeployment(configID uint64, version string, targetHosts []uint64, startedBy uint64) (*models.Deployment, error) {
	// 获取配置
	config, err := s.configRepo.GetByID(configID)
	if err != nil {
		return nil, err
	}

	// 验证目标主机
	if len(targetHosts) == 0 {
		return nil, errors.New("target hosts cannot be empty")
	}

	// 构建目标主机列表（Salt Minion ID）
	var minionIDs []string
	for _, hostID := range targetHosts {
		host, err := s.hostRepo.GetByID(hostID)
		if err != nil {
			return nil, err
		}
		if host.SaltMinionID == "" {
			return nil, errors.New("host does not have Salt Minion ID")
		}
		minionIDs = append(minionIDs, host.SaltMinionID)
	}

	// 创建部署记录
	deployment := &models.Deployment{
		ConfigID:      configID,
		Version:       version,
		Status:        "pending",
		StartedBy:     startedBy,
		StartedAt:     time.Now(),
		TargetHosts:   convertUint64ArrayToString(targetHosts),
		IsRollback:    false,
	}

	if err := s.deploymentRepo.Create(deployment); err != nil {
		return nil, err
	}

	// 如果 Salt 客户端可用，执行部署
	if s.saltClient != nil {
		// 构建 Salt target（使用逗号分隔的 Minion ID）
		target := ""
		for i, id := range minionIDs {
			if i > 0 {
				target += ","
			}
			target += id
		}

		// 执行 Salt state.apply
		// 注意：这里简化处理，实际应该解析 SaltStateFiles 数组
		state := "apache" // 默认 state，实际应该从配置中获取
		if config.SaltStateFiles != "" {
			// 解析 SaltStateFiles（假设是 JSON 数组格式）
			var states []string
			if err := json.Unmarshal([]byte(config.SaltStateFiles), &states); err == nil && len(states) > 0 {
				state = states[0]
			}
		}

		// 异步执行部署（在实际实现中应该使用 goroutine）
		go func() {
			s.executeDeployment(deployment.ID, target, state)
		}()

		// 更新状态为 running
		deployment.Status = "running"
		s.deploymentRepo.Update(deployment)
	}

	return deployment, nil
}

func (s *deploymentService) executeDeployment(deploymentID uint64, target string, state string) {
	deployment, err := s.deploymentRepo.GetByID(deploymentID)
	if err != nil {
		return
	}

	// 执行 Salt state.apply
	result, err := s.saltClient.StateApply(target, state, nil)
	if err != nil {
		deployment.Status = "failed"
		deployment.ErrorMessage = err.Error()
		deployment.CompletedAt = &[]time.Time{time.Now()}[0]
		s.deploymentRepo.Update(deployment)
		return
	}

	// 解析结果
	resultsJSON, _ := json.Marshal(result)
	deployment.Results = string(resultsJSON)

	// 检查结果是否成功
	allSuccess := true
	for _, v := range result {
		if resultMap, ok := v.(map[string]interface{}); ok {
			if resultValue, exists := resultMap["result"]; exists {
				if success, ok := resultValue.(bool); ok && !success {
					allSuccess = false
					break
				}
			}
		}
	}

	if allSuccess {
		deployment.Status = "completed"
	} else {
		deployment.Status = "failed"
		deployment.ErrorMessage = "some hosts failed"
	}

	now := time.Now()
	deployment.CompletedAt = &now
	duration := int(now.Sub(deployment.StartedAt).Seconds())
	deployment.DurationSeconds = &duration

	s.deploymentRepo.Update(deployment)
}

func (s *deploymentService) UpdateDeploymentStatus(id uint64, status string, results map[string]interface{}, errorMsg string) error {
	deployment, err := s.deploymentRepo.GetByID(id)
	if err != nil {
		return err
	}

	deployment.Status = status
	if results != nil {
		resultsJSON, _ := json.Marshal(results)
		deployment.Results = string(resultsJSON)
	}
	if errorMsg != "" {
		deployment.ErrorMessage = errorMsg
	}
	if status == "completed" || status == "failed" || status == "cancelled" {
		now := time.Now()
		deployment.CompletedAt = &now
		if deployment.StartedAt.Unix() > 0 {
			duration := int(now.Sub(deployment.StartedAt).Seconds())
			deployment.DurationSeconds = &duration
		}
	}

	return s.deploymentRepo.Update(deployment)
}

func (s *deploymentService) GetDeploymentStatus(id uint64) (*models.Deployment, error) {
	return s.deploymentRepo.GetByID(id)
}

func (s *deploymentService) RollbackDeployment(deploymentID uint64, startedBy uint64) (*models.Deployment, error) {
	// 获取要回滚的部署
	originalDeployment, err := s.deploymentRepo.GetByID(deploymentID)
	if err != nil {
		return nil, err
	}

	if originalDeployment.Status != "completed" {
		return nil, errors.New("can only rollback completed deployments")
	}

	// 获取配置
	config, err := s.configRepo.GetByID(originalDeployment.ConfigID)
	if err != nil {
		return nil, err
	}

	// 获取上一个成功的部署版本
	previousDeployments, err := s.deploymentRepo.GetByConfigID(config.ID, 10)
	if err != nil {
		return nil, err
	}

	var previousVersion string
	for _, dep := range previousDeployments {
		if dep.ID != deploymentID && dep.Status == "completed" {
			previousVersion = dep.Version
			break
		}
	}

	if previousVersion == "" {
		return nil, errors.New("no previous version found for rollback")
	}

	// 解析目标主机
	targetHosts := parseStringToUint64Array(originalDeployment.TargetHosts)

	// 创建回滚部署
	rollbackDeployment, err := s.StartDeployment(config.ID, previousVersion, targetHosts, startedBy)
	if err != nil {
		return nil, err
	}

	rollbackDeployment.IsRollback = true
	rollbackDeployment.RollbackFromDeploymentID = &deploymentID
	s.deploymentRepo.Update(rollbackDeployment)

	return rollbackDeployment, nil
}

type DeploymentVersionService interface {
	CreateVersion(version *models.DeploymentVersion) error
	GetVersion(id uint64) (*models.DeploymentVersion, error)
	ListVersions(page, pageSize int, filters map[string]interface{}) ([]models.DeploymentVersion, int64, error)
	GetVersionsByApplication(applicationName string) ([]models.DeploymentVersion, error)
}

type deploymentVersionService struct {
	versionRepo repository.DeploymentVersionRepository
}

func NewDeploymentVersionService(versionRepo repository.DeploymentVersionRepository) DeploymentVersionService {
	return &deploymentVersionService{versionRepo: versionRepo}
}

func (s *deploymentVersionService) CreateVersion(version *models.DeploymentVersion) error {
	return s.versionRepo.Create(version)
}

func (s *deploymentVersionService) GetVersion(id uint64) (*models.DeploymentVersion, error) {
	return s.versionRepo.GetByID(id)
}

func (s *deploymentVersionService) ListVersions(page, pageSize int, filters map[string]interface{}) ([]models.DeploymentVersion, int64, error) {
	offset := (page - 1) * pageSize
	return s.versionRepo.List(offset, pageSize, filters)
}

func (s *deploymentVersionService) GetVersionsByApplication(applicationName string) ([]models.DeploymentVersion, error) {
	return s.versionRepo.GetByApplication(applicationName)
}

// 辅助函数
func convertUint64ArrayToString(arr []uint64) string {
	data, _ := json.Marshal(arr)
	return string(data)
}

func parseStringToUint64Array(s string) []uint64 {
	var arr []uint64
	json.Unmarshal([]byte(s), &arr)
	return arr
}

