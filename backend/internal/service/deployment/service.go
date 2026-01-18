// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/utils"
)

// Service handles deployment management business logic
type Service struct {
	db     *gorm.DB
	config *config.Config
}

// NewService creates a new deployment service
func NewService(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{db: db, config: cfg}
}

// VersionSourceResponse represents the response from version source URL
type VersionSourceResponse struct {
	Versions []string `json:"versions"`
	Latest   string   `json:"latest"`
}

// CreateModuleRequest represents a request to create a deployment module
type CreateModuleRequest struct {
	ProjectID        uint   `json:"project_id" binding:"required"`
	EnvironmentID    *uint  `json:"environment_id"`
	TemplateID       *uint  `json:"template_id"` // 可选：关联执行模板
	Name             string `json:"name" binding:"required"`
	Description      string `json:"description"`
	VersionSourceURL string `json:"version_source_url"`
	DeployScript     string `json:"deploy_script"`
	ScriptType       string `json:"script_type"`
	Timeout          int    `json:"timeout"`
	AssetIDs         []uint `json:"asset_ids"`
}

// UpdateModuleRequest represents a request to update a deployment module
type UpdateModuleRequest struct {
	ProjectID        *uint  `json:"project_id"`
	EnvironmentID    *uint  `json:"environment_id"`
	TemplateID       *uint  `json:"template_id"` // 可选：关联执行模板
	Name             string `json:"name"`
	Description      string `json:"description"`
	VersionSourceURL string `json:"version_source_url"`
	DeployScript     string `json:"deploy_script"`
	ScriptType       string `json:"script_type"`
	Timeout          int    `json:"timeout"`
	AssetIDs         []uint `json:"asset_ids"`
}

// DeployRequest represents a request to execute deployment
type DeployRequest struct {
	Version  string `json:"version" binding:"required"`
	AssetIDs []uint `json:"asset_ids" binding:"required"`
}

// TemplateInfo represents basic template information
type TemplateInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ModuleResponse represents a deployment module response
type ModuleResponse struct {
	ID               uint          `json:"id"`
	ProjectID        uint          `json:"project_id"`
	ProjectName      string        `json:"project_name"`
	EnvironmentID    *uint         `json:"environment_id"`
	EnvironmentName  string        `json:"environment_name"`
	TemplateID       *uint         `json:"template_id"`
	TemplateName     string        `json:"template_name,omitempty"` // 模板名称（用于导出）
	Template         *TemplateInfo `json:"template,omitempty"`      // 关联的执行模板信息
	Name             string        `json:"name"`
	Description      string        `json:"description"`
	VersionSourceURL string        `json:"version_source_url"`
	DeployScript     string        `json:"deploy_script"`
	ScriptType       string        `json:"script_type"`
	Timeout          int           `json:"timeout"`
	AssetIDs         []uint        `json:"asset_ids"`
	CreatedBy        uint          `json:"created_by"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

// DeploymentResponse represents a deployment record response
type DeploymentResponse struct {
	ID           uint       `json:"id"`
	ModuleID     uint       `json:"module_id"`
	ModuleName   string     `json:"module_name"`
	ProjectName  string     `json:"project_name"`
	Version      string     `json:"version"`
	Status       string     `json:"status"`
	AssetIDs     []uint     `json:"asset_ids"`
	Output       string     `json:"output"`
	Error        string     `json:"error"`
	CreatedBy    uint       `json:"created_by"`
	CreatorName  string     `json:"creator_name"`
	StartedAt    *time.Time `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

// CreateModule creates a new deployment module
func (s *Service) CreateModule(req *CreateModuleRequest, userID uint) (*ModuleResponse, error) {
	// Convert asset IDs to comma-separated string
	assetIDsStr := ""
	if len(req.AssetIDs) > 0 {
		ids := make([]string, len(req.AssetIDs))
		for i, id := range req.AssetIDs {
			ids[i] = strconv.FormatUint(uint64(id), 10)
		}
		assetIDsStr = strings.Join(ids, ",")
	}

	scriptType := req.ScriptType
	deployScript := req.DeployScript

	// 如果关联了模板，从模板继承脚本内容和类型（如果未自定义）
	if req.TemplateID != nil && *req.TemplateID > 0 {
		var template model.TaskTemplate
		if err := s.db.First(&template, *req.TemplateID).Error; err == nil {
			if deployScript == "" {
				deployScript = template.Content
			}
			if scriptType == "" {
				scriptType = template.Type
			}
		}
	}

	if scriptType == "" {
		scriptType = "shell"
	}

	timeout := req.Timeout
	if timeout == 0 {
		timeout = 600
	}

	module := model.DeploymentModule{
		ProjectID:        req.ProjectID,
		EnvironmentID:    req.EnvironmentID,
		TemplateID:       req.TemplateID,
		Name:             req.Name,
		Description:      req.Description,
		VersionSourceURL: req.VersionSourceURL,
		DeployScript:     deployScript,
		ScriptType:       scriptType,
		Timeout:          timeout,
		AssetIDs:         assetIDsStr,
		CreatedBy:        userID,
	}

	if err := s.db.Create(&module).Error; err != nil {
		return nil, err
	}

	return s.GetModule(module.ID)
}

// GetModule retrieves a deployment module by ID
func (s *Service) GetModule(id uint) (*ModuleResponse, error) {
	var module model.DeploymentModule
	if err := s.db.Preload("Project").Preload("Environment").Preload("Template").First(&module, id).Error; err != nil {
		return nil, err
	}

	return s.moduleToResponse(&module), nil
}

// ListModules retrieves deployment modules with optional project filter
func (s *Service) ListModules(projectID *uint) ([]ModuleResponse, error) {
	var modules []model.DeploymentModule
	query := s.db.Preload("Project").Preload("Environment").Preload("Template")

	if projectID != nil && *projectID > 0 {
		query = query.Where("project_id = ?", *projectID)
	}

	if err := query.Order("created_at DESC").Find(&modules).Error; err != nil {
		// 如果表不存在，返回空数组而不是错误
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "no such table") {
			return []ModuleResponse{}, nil
		}
		return nil, err
	}

	result := make([]ModuleResponse, len(modules))
	for i, module := range modules {
		result[i] = *s.moduleToResponse(&module)
	}

	return result, nil
}

// UpdateModule updates a deployment module
func (s *Service) UpdateModule(id uint, req *UpdateModuleRequest) (*ModuleResponse, error) {
	var module model.DeploymentModule
	if err := s.db.First(&module, id).Error; err != nil {
		return nil, err
	}

	if req.ProjectID != nil {
		module.ProjectID = *req.ProjectID
	}
	if req.EnvironmentID != nil {
		module.EnvironmentID = req.EnvironmentID
	}
	// 处理模板关联
	if req.TemplateID != nil {
		if *req.TemplateID == 0 {
			// 清除模板关联
			module.TemplateID = nil
		} else {
			module.TemplateID = req.TemplateID
			// 如果没有自定义脚本，从模板继承
			if req.DeployScript == "" {
				var template model.TaskTemplate
				if err := s.db.First(&template, *req.TemplateID).Error; err == nil {
					module.DeployScript = template.Content
					if req.ScriptType == "" {
						module.ScriptType = template.Type
					}
				}
			}
		}
	}
	if req.Name != "" {
		module.Name = req.Name
	}
	if req.Description != "" {
		module.Description = req.Description
	}
	if req.VersionSourceURL != "" {
		module.VersionSourceURL = req.VersionSourceURL
	}
	if req.DeployScript != "" {
		module.DeployScript = req.DeployScript
	}
	if req.ScriptType != "" {
		module.ScriptType = req.ScriptType
	}
	if req.Timeout > 0 {
		module.Timeout = req.Timeout
	}
	if req.AssetIDs != nil {
		ids := make([]string, len(req.AssetIDs))
		for i, id := range req.AssetIDs {
			ids[i] = strconv.FormatUint(uint64(id), 10)
		}
		module.AssetIDs = strings.Join(ids, ",")
	}

	if err := s.db.Save(&module).Error; err != nil {
		return nil, err
	}

	return s.GetModule(id)
}

// DeleteModule deletes a deployment module
func (s *Service) DeleteModule(id uint) error {
	return s.db.Delete(&model.DeploymentModule{}, id).Error
}

// GetVersions fetches versions from the module's version source URL
func (s *Service) GetVersions(moduleID uint) (*VersionSourceResponse, error) {
	var module model.DeploymentModule
	if err := s.db.First(&module, moduleID).Error; err != nil {
		return nil, err
	}

	if module.VersionSourceURL == "" {
		return nil, fmt.Errorf("version source URL not configured")
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(module.VersionSourceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("version source returned status %d", resp.StatusCode)
	}

	var versionResp VersionSourceResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
		return nil, fmt.Errorf("failed to parse version response: %w", err)
	}

	return &versionResp, nil
}

// Deploy executes deployment for a module
func (s *Service) Deploy(moduleID uint, req *DeployRequest, userID uint) (*DeploymentResponse, error) {
	var module model.DeploymentModule
	if err := s.db.Preload("Project").Preload("Environment").First(&module, moduleID).Error; err != nil {
		return nil, err
	}

	// Convert asset IDs to comma-separated string
	assetIDsStr := ""
	if len(req.AssetIDs) > 0 {
		ids := make([]string, len(req.AssetIDs))
		for i, id := range req.AssetIDs {
			ids[i] = strconv.FormatUint(uint64(id), 10)
		}
		assetIDsStr = strings.Join(ids, ",")
	}

	// Create deployment record
	deployment := model.Deployment{
		ModuleID:  moduleID,
		Version:   req.Version,
		Status:    "pending",
		AssetIDs:  assetIDsStr,
		CreatedBy: userID,
	}

	if err := s.db.Create(&deployment).Error; err != nil {
		return nil, err
	}

	// Execute deployment asynchronously
	go s.executeDeployment(&deployment, &module, req.AssetIDs)

	return s.GetDeployment(deployment.ID)
}

// executeDeployment performs the actual deployment execution
func (s *Service) executeDeployment(deployment *model.Deployment, module *model.DeploymentModule, assetIDs []uint) {
	// Update status to running
	now := time.Now()
	deployment.Status = "running"
	deployment.StartedAt = &now
	s.db.Save(deployment)

	// Get project and environment names for variable replacement
	projectName := ""
	if module.Project != nil {
		projectName = module.Project.Name
	}
	environmentName := ""
	if module.Environment != nil {
		environmentName = module.Environment.Name
	}

	// Prepare script with variable replacement
	script := module.DeployScript
	script = strings.ReplaceAll(script, "${VERSION}", deployment.Version)
	script = strings.ReplaceAll(script, "${MODULE_NAME}", module.Name)
	script = strings.ReplaceAll(script, "${PROJECT_NAME}", projectName)
	script = strings.ReplaceAll(script, "${ENVIRONMENT_NAME}", environmentName)

	// Collect outputs
	var outputs []string
	var errors []string
	allSuccess := true

	// Execute on each asset
	for _, assetID := range assetIDs {
		var asset model.Asset
		if err := s.db.Preload("SSHKey").First(&asset, assetID).Error; err != nil {
			errors = append(errors, fmt.Sprintf("[%s] Failed to get asset: %v", asset.HostName, err))
			allSuccess = false
			continue
		}

		// Execute script via SSH
		output, err := s.executeScriptOnAsset(&asset, script, module.Timeout)
		if err != nil {
			errors = append(errors, fmt.Sprintf("[%s] %v", asset.HostName, err))
			allSuccess = false
		}
		if output != "" {
			outputs = append(outputs, fmt.Sprintf("=== %s (%s) ===\n%s", asset.HostName, asset.IP, output))
		}
	}

	// Update deployment record
	finishedAt := time.Now()
	deployment.FinishedAt = &finishedAt
	deployment.Output = strings.Join(outputs, "\n\n")
	deployment.Error = strings.Join(errors, "\n")

	if allSuccess {
		deployment.Status = "success"
	} else {
		deployment.Status = "failed"
	}

	s.db.Save(deployment)
}

// executeScriptOnAsset executes a script on a single asset via SSH
func (s *Service) executeScriptOnAsset(asset *model.Asset, script string, timeout int) (string, error) {
	if asset.SSHKey == nil {
		return "", fmt.Errorf("no SSH key configured for asset")
	}

	sshUser := asset.SSHUser
	if sshUser == "" {
		sshUser = asset.SSHKey.SSHUser
	}
	if sshUser == "" {
		sshUser = "root"
	}

	sshPort := asset.SSHPort
	if sshPort == 0 {
		sshPort = 22
	}

	// Decrypt private key
	privateKeyBytes, err := utils.Decrypt(asset.SSHKey.PrivateKey, s.config.Encryption.Key)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt private key: %w", err)
	}

	// Decrypt passphrase if exists
	var passphraseBytes []byte
	if asset.SSHKey.Passphrase != "" {
		passphraseBytes, err = utils.Decrypt(asset.SSHKey.Passphrase, s.config.Encryption.Key)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt passphrase: %w", err)
		}
	}

	// Create SSH client
	var client *utils.SSHClient

	if len(passphraseBytes) > 0 {
		client, err = utils.NewSSHClientWithPassphrase(
			asset.IP,
			sshPort,
			sshUser,
			privateKeyBytes,
			passphraseBytes,
			30*time.Second,
		)
	} else {
		client, err = utils.NewSSHClient(
			asset.IP,
			sshPort,
			sshUser,
			privateKeyBytes,
			30*time.Second,
		)
	}
	if err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	// Execute command with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	output, _, err := client.ExecuteCommandWithTimeout(ctx, script)
	return output, err
}

// GetDeployment retrieves a deployment record by ID
func (s *Service) GetDeployment(id uint) (*DeploymentResponse, error) {
	var deployment model.Deployment
	if err := s.db.Preload("Module.Project").Preload("Creator").First(&deployment, id).Error; err != nil {
		return nil, err
	}

	return s.deploymentToResponse(&deployment), nil
}

// ListDeployments retrieves deployment records with optional module filter
func (s *Service) ListDeployments(moduleID *uint, page, pageSize int) ([]DeploymentResponse, int64, error) {
	var deployments []model.Deployment
	var total int64

	query := s.db.Model(&model.Deployment{})

	if moduleID != nil && *moduleID > 0 {
		query = query.Where("module_id = ?", *moduleID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Module.Project").Preload("Creator").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&deployments).Error; err != nil {
		return nil, 0, err
	}

	result := make([]DeploymentResponse, len(deployments))
	for i, d := range deployments {
		result[i] = *s.deploymentToResponse(&d)
	}

	return result, total, nil
}

// CancelDeployment cancels a running deployment
func (s *Service) CancelDeployment(id uint) error {
	var deployment model.Deployment
	if err := s.db.First(&deployment, id).Error; err != nil {
		return err
	}

	if deployment.Status != "pending" && deployment.Status != "running" {
		return fmt.Errorf("deployment is not running")
	}

	now := time.Now()
	deployment.Status = "cancelled"
	deployment.FinishedAt = &now

	return s.db.Save(&deployment).Error
}

// Helper functions

func (s *Service) moduleToResponse(m *model.DeploymentModule) *ModuleResponse {
	assetIDs := parseAssetIDs(m.AssetIDs)
	projectName := ""
	if m.Project != nil {
		projectName = m.Project.Name
	}
	environmentName := ""
	if m.Environment != nil {
		environmentName = m.Environment.Name
	}

	// 构建模板信息
	var templateInfo *TemplateInfo
	if m.Template != nil {
		templateInfo = &TemplateInfo{
			ID:   m.Template.ID,
			Name: m.Template.Name,
			Type: m.Template.Type,
		}
	}

	// 获取模板名称
	templateName := ""
	if m.Template != nil {
		templateName = m.Template.Name
	}

	return &ModuleResponse{
		ID:               m.ID,
		ProjectID:        m.ProjectID,
		ProjectName:      projectName,
		EnvironmentID:    m.EnvironmentID,
		EnvironmentName:  environmentName,
		TemplateID:       m.TemplateID,
		TemplateName:     templateName,
		Template:         templateInfo,
		Name:             m.Name,
		Description:      m.Description,
		VersionSourceURL: m.VersionSourceURL,
		DeployScript:     m.DeployScript,
		ScriptType:       m.ScriptType,
		Timeout:          m.Timeout,
		AssetIDs:         assetIDs,
		CreatedBy:        m.CreatedBy,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func (s *Service) deploymentToResponse(d *model.Deployment) *DeploymentResponse {
	assetIDs := parseAssetIDs(d.AssetIDs)
	moduleName := ""
	projectName := ""
	if d.Module != nil {
		moduleName = d.Module.Name
		if d.Module.Project != nil {
			projectName = d.Module.Project.Name
		}
	}
	creatorName := ""
	if d.Creator.ID > 0 {
		creatorName = d.Creator.Username
	}

	return &DeploymentResponse{
		ID:          d.ID,
		ModuleID:    d.ModuleID,
		ModuleName:  moduleName,
		ProjectName: projectName,
		Version:     d.Version,
		Status:      d.Status,
		AssetIDs:    assetIDs,
		Output:      d.Output,
		Error:       d.Error,
		CreatedBy:   d.CreatedBy,
		CreatorName: creatorName,
		StartedAt:   d.StartedAt,
		FinishedAt:  d.FinishedAt,
		CreatedAt:   d.CreatedAt,
	}
}

func parseAssetIDs(s string) []uint {
	if s == "" {
		return []uint{}
	}

	parts := strings.Split(s, ",")
	result := make([]uint, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.ParseUint(p, 10, 32)
		if err == nil {
			result = append(result, uint(id))
		}
	}
	return result
}

// FindProjectByName 根据项目名称查找项目 ID
func (s *Service) FindProjectByName(name string) (uint, error) {
	var project model.Project
	if err := s.db.Where("name = ?", name).First(&project).Error; err != nil {
		return 0, fmt.Errorf("项目不存在: %s", name)
	}
	return project.ID, nil
}

// FindEnvironmentByName 根据环境名称查找环境 ID
func (s *Service) FindEnvironmentByName(name string) (uint, error) {
	var env model.Environment
	if err := s.db.Where("name = ?", name).First(&env).Error; err != nil {
		return 0, fmt.Errorf("环境不存在: %s", name)
	}
	return env.ID, nil
}

// FindTemplateByName 根据模板名称查找模板 ID
func (s *Service) FindTemplateByName(name string) (uint, error) {
	var template model.TaskTemplate
	if err := s.db.Where("name = ?", name).First(&template).Error; err != nil {
		return 0, fmt.Errorf("模板不存在: %s", name)
	}
	return template.ID, nil
}

// ModuleExistsByName 检查同名模块是否存在
func (s *Service) ModuleExistsByName(name string, projectID uint) (bool, uint) {
	var module model.DeploymentModule
	query := s.db.Where("name = ?", name)
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if err := query.First(&module).Error; err != nil {
		return false, 0
	}
	return true, module.ID
}

// GetAssetHostnames 根据资产ID列表获取主机名列表
func (s *Service) GetAssetHostnames(assetIDs []uint) []string {
	if len(assetIDs) == 0 {
		return []string{}
	}

	var assets []model.Asset
	s.db.Where("id IN ?", assetIDs).Find(&assets)

	hostnames := make([]string, 0, len(assets))
	for _, asset := range assets {
		// 优先使用主机名，没有则使用 IP
		if asset.HostName != "" {
			hostnames = append(hostnames, asset.HostName)
		} else if asset.IP != "" {
			hostnames = append(hostnames, asset.IP)
		}
	}
	return hostnames
}

// FindAssetsByHostnames 根据主机名或IP列表查找资产ID
func (s *Service) FindAssetsByHostnames(hostnames []string) ([]uint, []string) {
	if len(hostnames) == 0 {
		return []uint{}, []string{}
	}

	var assets []model.Asset
	// 同时匹配主机名和 IP
	s.db.Where("host_name IN ? OR ip IN ?", hostnames, hostnames).Find(&assets)

	// 创建已找到的主机名/IP 映射
	foundSet := make(map[string]bool)
	assetIDs := make([]uint, 0, len(assets))
	for _, asset := range assets {
		assetIDs = append(assetIDs, asset.ID)
		if asset.HostName != "" {
			foundSet[asset.HostName] = true
		}
		if asset.IP != "" {
			foundSet[asset.IP] = true
		}
	}

	// 找出未找到的主机名
	missingHosts := make([]string, 0)
	for _, h := range hostnames {
		if !foundSet[h] {
			missingHosts = append(missingHosts, h)
		}
	}

	return assetIDs, missingHosts
}
