package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
	"github.com/kkops/backend/internal/salt"
)

type FormulaService interface {
	// Formula管理
	CreateFormula(formula *models.Formula) error
	GetFormula(id uint) (*models.Formula, error)
	ListFormulas(page, pageSize int, filters map[string]interface{}) ([]models.Formula, int64, error)
	UpdateFormula(formula *models.Formula) error
	DeleteFormula(id uint) error

	// Formula参数管理
	GetFormulaParameters(formulaID uint) ([]models.FormulaParameter, error)
	UpdateFormulaParameters(formulaID uint, params []models.FormulaParameter) error

	// Formula模板管理
	CreateFormulaTemplate(template *models.FormulaTemplate) error
	GetFormulaTemplates(formulaID uint) ([]models.FormulaTemplate, error)
	UpdateFormulaTemplate(template *models.FormulaTemplate) error
	DeleteFormulaTemplate(id uint) error

	// Formula部署
	CreateDeployment(deployment *models.FormulaDeployment) error
	ExecuteDeployment(deploymentID uint) error
	GetDeployment(id uint) (*models.FormulaDeployment, error)
	ListDeployments(page, pageSize int, filters map[string]interface{}) ([]models.FormulaDeployment, int64, error)
	CancelDeployment(id uint) error
	CleanupOldDeployments(beforeTime time.Time) (int64, error)

	// Formula仓库管理
	CreateRepository(repo *models.FormulaRepository) error
	GetRepository(id uint) (*models.FormulaRepository, error)
	UpdateRepository(repo *models.FormulaRepository) error
	ListRepositories(page, pageSize int, filters map[string]interface{}) ([]models.FormulaRepository, int64, error)
	SyncRepository(repoID uint) error

	// Formula发现
	DiscoverFormulas(repoPath string) ([]models.Formula, error)
}

type formulaService struct {
	formulaRepo repository.FormulaRepository
	hostRepo    repository.HostRepository
	sshKeyRepo  repository.SSHKeyRepository
	saltClient  *salt.Client
}

func NewFormulaService(formulaRepo repository.FormulaRepository, hostRepo repository.HostRepository, sshKeyRepo repository.SSHKeyRepository, saltClient *salt.Client) FormulaService {
	return &formulaService{
		formulaRepo: formulaRepo,
		hostRepo:    hostRepo,
		sshKeyRepo:  sshKeyRepo,
		saltClient:  saltClient,
	}
}

// Formula管理
func (s *formulaService) CreateFormula(formula *models.Formula) error {
	return s.formulaRepo.Create(formula)
}

func (s *formulaService) GetFormula(id uint) (*models.Formula, error) {
	return s.formulaRepo.GetByID(id)
}

func (s *formulaService) ListFormulas(page, pageSize int, filters map[string]interface{}) ([]models.Formula, int64, error) {
	return s.formulaRepo.List(page, pageSize, filters)
}

func (s *formulaService) UpdateFormula(formula *models.Formula) error {
	return s.formulaRepo.Update(formula)
}

func (s *formulaService) DeleteFormula(id uint) error {
	return s.formulaRepo.Delete(id)
}

// Formula参数管理
func (s *formulaService) GetFormulaParameters(formulaID uint) ([]models.FormulaParameter, error) {
	return s.formulaRepo.GetParameters(formulaID)
}

func (s *formulaService) UpdateFormulaParameters(formulaID uint, params []models.FormulaParameter) error {
	// 先删除现有的参数
	existingParams, err := s.formulaRepo.GetParameters(formulaID)
	if err != nil {
		return err
	}

	for _, param := range existingParams {
		if err := s.formulaRepo.DeleteParameter(param.ID); err != nil {
			return err
		}
	}

	// 创建新的参数
	for _, param := range params {
		param.FormulaID = formulaID
		if err := s.formulaRepo.CreateParameter(&param); err != nil {
			return err
		}
	}

	return nil
}

// Formula模板管理
func (s *formulaService) CreateFormulaTemplate(template *models.FormulaTemplate) error {
	return s.formulaRepo.CreateTemplate(template)
}

func (s *formulaService) GetFormulaTemplates(formulaID uint) ([]models.FormulaTemplate, error) {
	return s.formulaRepo.GetTemplates(formulaID)
}

func (s *formulaService) UpdateFormulaTemplate(template *models.FormulaTemplate) error {
	return s.formulaRepo.UpdateTemplate(template)
}

func (s *formulaService) DeleteFormulaTemplate(id uint) error {
	return s.formulaRepo.DeleteTemplate(id)
}

// Formula部署
func (s *formulaService) CreateDeployment(deployment *models.FormulaDeployment) error {
	deployment.Status = "pending"
	return s.formulaRepo.CreateDeployment(deployment)
}

func (s *formulaService) ExecuteDeployment(deploymentID uint) error {
	deployment, err := s.formulaRepo.GetDeploymentByID(deploymentID)
	if err != nil {
		return err
	}

	if deployment.Status != "pending" {
		return fmt.Errorf("deployment is not in pending status")
	}

	// 更新部署状态
	now := time.Now()
	deployment.Status = "running"
	deployment.StartedAt = &now
	if err := s.formulaRepo.UpdateDeployment(deployment); err != nil {
		return err
	}

	// 异步执行部署
	go s.executeDeploymentAsync(deployment)

	return nil
}

func (s *formulaService) executeDeploymentAsync(deployment *models.FormulaDeployment) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Formula deployment panic: %v", r)
			deployment.Status = "failed"
			deployment.ErrorMessage = fmt.Sprintf("Panic: %v", r)
			s.formulaRepo.UpdateDeployment(deployment)
		}
	}()

	// 获取Formula信息
	formula, err := s.formulaRepo.GetByID(deployment.FormulaID)
	if err != nil {
		deployment.Status = "failed"
		deployment.ErrorMessage = fmt.Sprintf("Failed to get formula: %v", err)
		s.formulaRepo.UpdateDeployment(deployment)
		return
	}

	// 解析目标主机
	var targetHosts []string
	if err := json.Unmarshal(deployment.TargetHosts, &targetHosts); err != nil {
		deployment.Status = "failed"
		deployment.ErrorMessage = fmt.Sprintf("Failed to parse target hosts: %v", err)
		s.formulaRepo.UpdateDeployment(deployment)
		return
	}

	// 构建Salt target
	var target string
	if len(targetHosts) == 1 {
		target = targetHosts[0]
	} else {
		target = strings.Join(targetHosts, ",")
	}

	// 解析Pillar数据
	var pillarData map[string]interface{}
	if deployment.PillarData != nil {
		if err := json.Unmarshal(deployment.PillarData, &pillarData); err != nil {
			deployment.Status = "failed"
			deployment.ErrorMessage = fmt.Sprintf("Failed to parse pillar data: %v", err)
			s.formulaRepo.UpdateDeployment(deployment)
			return
		}
	}

	// 使用 formula 的 path 作为 state 名称（如 runtime/java），如果 path 为空则使用 name
	stateName := formula.Path
	if stateName == "" {
		stateName = formula.Name
	}

	log.Printf("Executing state.apply: state=%s, target=%s, pillar=%v", stateName, target, pillarData)

	// 执行state.apply
	result, err := s.saltClient.ExecuteState(stateName, target, pillarData)
	if err != nil {
		deployment.Status = "failed"
		deployment.ErrorMessage = fmt.Sprintf("Failed to execute state: %v", err)
		s.formulaRepo.UpdateDeployment(deployment)
		return
	}

	// 处理结果
	results := make(map[string]interface{})
	successCount := 0
	failedCount := 0

	log.Printf("Salt API result: %+v", result)

	if resultData, ok := result["return"].([]interface{}); ok && len(resultData) > 0 {
		if minions, ok := resultData[0].(map[string]interface{}); ok {
			for minionID, minionResult := range minions {
				results[minionID] = minionResult

				// Salt state.apply 返回格式: {minion_id: {state_id: {result: bool, ...}, ...}}
				if minionResultMap, ok := minionResult.(map[string]interface{}); ok {
					allSuccess := true
					hasStates := false

					for stateID, stateResult := range minionResultMap {
						// 跳过非 state 结果的字段
						if stateID == "retcode" {
							continue
						}

						hasStates = true
						if stateResultMap, ok := stateResult.(map[string]interface{}); ok {
							if success, ok := stateResultMap["result"].(bool); ok && !success {
								allSuccess = false
								log.Printf("State %s failed on %s", stateID, minionID)
							}
						}
					}

					if hasStates {
						if allSuccess {
							successCount++
							log.Printf("Minion %s: all states succeeded", minionID)
						} else {
							failedCount++
							log.Printf("Minion %s: some states failed", minionID)
						}
					} else {
						// 没有 state 结果，可能是错误
						failedCount++
						log.Printf("Minion %s: no state results found", minionID)
					}
				} else if errStr, ok := minionResult.(string); ok {
					// 错误信息是字符串
					failedCount++
					log.Printf("Minion %s error: %s", minionID, errStr)
				} else {
					failedCount++
					log.Printf("Minion %s: unexpected result format", minionID)
				}
			}
		}
	} else {
		log.Printf("No return data in Salt API result")
	}

	// 更新部署结果
	resultsJSON, _ := json.Marshal(results)
	now := time.Now()
	deployment.Status = "completed"
	deployment.CompletedAt = &now
	deployment.SuccessCount = successCount
	deployment.FailedCount = failedCount
	deployment.Results = resultsJSON

	if deployment.StartedAt != nil {
		duration := time.Since(*deployment.StartedAt)
		deployment.Duration = int(duration.Seconds())
	}

	s.formulaRepo.UpdateDeployment(deployment)
}

func (s *formulaService) GetDeployment(id uint) (*models.FormulaDeployment, error) {
	return s.formulaRepo.GetDeploymentByID(id)
}

func (s *formulaService) ListDeployments(page, pageSize int, filters map[string]interface{}) ([]models.FormulaDeployment, int64, error) {
	return s.formulaRepo.ListDeployments(page, pageSize, filters)
}

func (s *formulaService) CancelDeployment(id uint) error {
	deployment, err := s.formulaRepo.GetDeploymentByID(id)
	if err != nil {
		return err
	}

	if deployment.Status != "running" {
		return fmt.Errorf("deployment is not running")
	}

	deployment.Status = "cancelled"
	return s.formulaRepo.UpdateDeployment(deployment)
}

func (s *formulaService) CleanupOldDeployments(beforeTime time.Time) (int64, error) {
	return s.formulaRepo.CleanupOldDeployments(beforeTime)
}

// Formula仓库管理
func (s *formulaService) CreateRepository(repo *models.FormulaRepository) error {
	return s.formulaRepo.CreateRepository(repo)
}

func (s *formulaService) GetRepository(id uint) (*models.FormulaRepository, error) {
	return s.formulaRepo.GetRepositoryByID(id)
}

func (s *formulaService) UpdateRepository(repo *models.FormulaRepository) error {
	return s.formulaRepo.UpdateRepository(repo)
}

func (s *formulaService) ListRepositories(page, pageSize int, filters map[string]interface{}) ([]models.FormulaRepository, int64, error) {
	return s.formulaRepo.ListRepositories(page, pageSize, filters)
}

func (s *formulaService) SyncRepository(repoID uint) error {
	repo, err := s.formulaRepo.GetRepositoryByID(repoID)
	if err != nil {
		return err
	}

	log.Printf("Syncing repository: ID=%d, Name=%s, CreatedBy=%d", repo.ID, repo.Name, repo.CreatedBy)

	// 1. 首先同步到 Salt Master 的 /srv/salt 目录
	if err := s.syncToSaltMaster(repo); err != nil {
		log.Printf("Warning: Failed to sync to Salt Master: %v", err)
		// 不返回错误，继续本地同步
	}

	// 2. 同步仓库到本地（用于扫描 Formula）
	repoPath := repo.LocalPath
	if repoPath == "" {
		repoPath = fmt.Sprintf("/tmp/salt-formulas-%d", repo.ID)
	}

	// 确保目录存在
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// 执行Git同步
	if err := s.syncGitRepository(repo, repoPath); err != nil {
		return fmt.Errorf("failed to sync git repository: %v", err)
	}

	// 扫描Formula
	formulas, err := s.DiscoverFormulas(repoPath)
	if err != nil {
		return fmt.Errorf("failed to discover formulas: %v", err)
	}

	// 保存发现的Formula到数据库
	for _, formula := range formulas {
		// 最后验证：确保Formula name不为空
		if formula.Name == "" {
			log.Printf("ERROR: Skipping formula with empty name from path: %s", formula.Path)
			continue
		}

		// 创建Formula的副本，避免并发修改
		formulaCopy := formula
		formulaCopy.Repository = repo.URL

		// 确保created_by不为0，使用仓库创建者ID，如果为0或无效则使用系统默认ID
		if repo.CreatedBy == 0 || repo.CreatedBy > 1000000 { // 防止无效的ID
			formulaCopy.CreatedBy = 1 // 系统默认用户ID（admin）
		} else {
			formulaCopy.CreatedBy = repo.CreatedBy
		}

		existing, err := s.formulaRepo.GetByName(formulaCopy.Name)
		if err != nil {
			if err.Error() != "record not found" {
				log.Printf("Error checking existing formula %s: %v", formulaCopy.Name, err)
				continue
			}
			// record not found - 创建新Formula
			log.Printf("Creating new formula: %s (created_by=%d)", formulaCopy.Name, formulaCopy.CreatedBy)
			if err := s.formulaRepo.Create(&formulaCopy); err != nil {
				log.Printf("Failed to create formula %s: %v", formulaCopy.Name, err)
			}
		} else if existing != nil && existing.ID > 0 {
			// 更新现有Formula
			existing.Description = formulaCopy.Description
			existing.Category = formulaCopy.Category
			existing.Version = formulaCopy.Version
			existing.Path = formulaCopy.Path
			existing.Metadata = formulaCopy.Metadata
			existing.UpdatedAt = time.Now()
			log.Printf("Updating existing formula: %s (id=%d)", formulaCopy.Name, existing.ID)
			if err := s.formulaRepo.Update(existing); err != nil {
				log.Printf("Failed to update formula %s: %v", formulaCopy.Name, err)
			}
		} else {
			// existing 为 nil 或 ID 为 0 - 创建新Formula
			log.Printf("Creating new formula: %s (created_by=%d)", formulaCopy.Name, formulaCopy.CreatedBy)
			if err := s.formulaRepo.Create(&formulaCopy); err != nil {
				log.Printf("Failed to create formula %s: %v", formulaCopy.Name, err)
			}
		}
	}

	// 更新最后同步时间
	now := time.Now()
	repo.LastSyncAt = &now
	return s.formulaRepo.UpdateRepository(repo)
}

// Formula发现 - 扫描根目录下的一级目录（扁平结构）
// 目录命名格式: {category}_{name}，如 app_webapp, base_users, runtime_java
func (s *formulaService) DiscoverFormulas(repoPath string) ([]models.Formula, error) {
	var formulas []models.Formula

	log.Printf("Starting Formula discovery in path: %s", repoPath)

	// 检查目录是否存在
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("repository path does not exist: %s", repoPath)
	}

	// 只扫描根目录下的一级目录
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, entry := range entries {
		// 跳过非目录和隐藏目录
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		dirName := entry.Name()
		dirPath := filepath.Join(repoPath, dirName)

		// 检查是否是Formula目录（包含init.sls）
		initPath := filepath.Join(dirPath, "init.sls")
		if _, err := os.Stat(initPath); err != nil {
			log.Printf("Skipping directory %s: no init.sls found", dirName)
			continue
		}

		// 确保Formula名称有效
		if dirName == "" || dirName == "." || dirName == ".." {
			log.Printf("Skipping invalid formula directory: %s", dirName)
			continue
		}

		formula := models.Formula{
			Name:        dirName,
			Path:        dirName, // 扁平结构，path 就是目录名
			Description: fmt.Sprintf("Auto-discovered formula: %s", dirName),
			Category:    s.inferCategory(dirName),
			IsActive:    true,
		}

		log.Printf("Discovered Formula: %s (category: %s, path: %s)", formula.Name, formula.Category, formula.Path)

		// 解析元数据
		if metadata, err := s.parseFormulaMetadata(dirPath); err == nil {
			metadataJSON, _ := json.Marshal(metadata)
			formula.Metadata = metadataJSON
		}

		formulas = append(formulas, formula)
	}

	log.Printf("Formula discovery completed, found %d formulas", len(formulas))
	return formulas, nil
}

// 推断Formula分类 - 从目录名解析分类
// 目录命名格式: {category}_{name}，如 app_webapp, base_users, runtime_java
// 支持的分类前缀: app_, base_, middleware_, runtime_
func (s *formulaService) inferCategory(dirName string) string {
	// 检查是否以已知分类前缀开头
	prefixes := map[string]string{
		"app_":        "app",
		"base_":       "base",
		"middleware_": "middleware",
		"runtime_":    "runtime",
	}

	for prefix, category := range prefixes {
		if strings.HasPrefix(dirName, prefix) {
			return category
		}
	}

	return "custom"
}

// 解析Formula元数据
func (s *formulaService) parseFormulaMetadata(formulaPath string) (map[string]interface{}, error) {
	metadata := make(map[string]interface{})

	// 检查是否存在map.jinja
	mapPath := filepath.Join(formulaPath, "map.jinja")
	if _, err := os.Stat(mapPath); err == nil {
		metadata["has_map"] = true
	}

	// 检查子目录
	entries, err := os.ReadDir(formulaPath)
	if err != nil {
		return metadata, err
	}

	subdirs := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			subdirs = append(subdirs, entry.Name())
		}
	}
	metadata["subdirs"] = subdirs

	return metadata, nil
}

// syncGitRepository 同步Git仓库
func (s *formulaService) syncGitRepository(repo *models.FormulaRepository, repoPath string) error {
	// 检查是否已经存在Git仓库
	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		// Git仓库已存在，执行pull
		log.Printf("Git repository exists at %s, pulling latest changes", repoPath)
		return s.gitPull(repoPath, repo.Branch)
	} else {
		// Git仓库不存在，执行clone
		log.Printf("Cloning git repository %s to %s", repo.URL, repoPath)
		return s.gitClone(repo.URL, repoPath, repo.Branch)
	}
}

// gitClone 执行Git clone操作
func (s *formulaService) gitClone(url, path, branch string) error {
	args := []string{"clone", "--depth", "1", "--branch", branch, url, path}
	cmd := exec.Command("git", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %v, output: %s", err, string(output))
	}

	// 设置默认的pull策略为merge，避免将来出现分歧问题
	configCmd := exec.Command("git", "config", "pull.rebase", "false")
	configCmd.Dir = path
	if configOutput, configErr := configCmd.CombinedOutput(); configErr != nil {
		log.Printf("Warning: failed to set git pull config: %v, output: %s", configErr, string(configOutput))
		// 不返回错误，因为这不是致命问题
	}

	log.Printf("Successfully cloned repository to %s", path)
	return nil
}

// gitPull 执行Git pull操作
func (s *formulaService) gitPull(path, branch string) error {
	// 确保pull策略设置为merge
	configCmd := exec.Command("git", "config", "pull.rebase", "false")
	configCmd.Dir = path
	if configOutput, configErr := configCmd.CombinedOutput(); configErr != nil {
		log.Printf("Warning: failed to set git pull config: %v, output: %s", configErr, string(configOutput))
	}

	// 切换到指定分支
	checkoutCmd := exec.Command("git", "checkout", branch)
	checkoutCmd.Dir = path
	if output, err := checkoutCmd.CombinedOutput(); err != nil {
		log.Printf("Warning: git checkout failed: %v, output: %s", err, string(output))
	}

	// 执行pull，使用merge策略避免分歧问题
	pullCmd := exec.Command("git", "pull", "--no-rebase", "origin", branch)
	pullCmd.Dir = path

	output, err := pullCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %v, output: %s", err, string(output))
	}

	log.Printf("Successfully pulled latest changes to %s", path)
	return nil
}

// syncToSaltMaster 通过 rsync 同步仓库到 Salt Master
func (s *formulaService) syncToSaltMaster(repo *models.FormulaRepository) error {
	log.Printf("Syncing repository to Salt Master via rsync: %s", repo.URL)

	// 查找 Salt Master 主机信息
	saltMasterHost, err := s.hostRepo.GetByHostname("salt-master")
	if err != nil {
		return fmt.Errorf("failed to find salt-master host: %v", err)
	}

	if saltMasterHost == nil {
		return fmt.Errorf("salt-master host not found in database")
	}

	// 获取 SSH 端口
	sshPort := saltMasterHost.SSHPort
	if sshPort == 0 {
		sshPort = 22
	}

	// 获取 SSH 用户名（从 SSH Key 或默认 root）
	sshUser := "root"
	var sshKeyPath string

	if saltMasterHost.SSHKeyID != nil && *saltMasterHost.SSHKeyID > 0 {
		sshKey, err := s.sshKeyRepo.GetByID(*saltMasterHost.SSHKeyID)
		if err != nil {
			log.Printf("Warning: failed to get SSH key: %v, will try without key", err)
		} else if sshKey != nil {
			// 如果 SSH Key 有指定用户名，使用它
			if sshKey.Username != "" {
				sshUser = sshKey.Username
			}

			// 将私钥写入临时文件
			tmpKeyFile, err := os.CreateTemp("", "ssh-key-*")
			if err != nil {
				return fmt.Errorf("failed to create temp key file: %v", err)
			}
			defer os.Remove(tmpKeyFile.Name())

			if _, err := tmpKeyFile.WriteString(sshKey.PrivateKey); err != nil {
				tmpKeyFile.Close()
				return fmt.Errorf("failed to write private key: %v", err)
			}
			tmpKeyFile.Close()

			// 设置正确的权限
			if err := os.Chmod(tmpKeyFile.Name(), 0600); err != nil {
				return fmt.Errorf("failed to set key file permissions: %v", err)
			}

			sshKeyPath = tmpKeyFile.Name()
		}
	}

	// 本地仓库路径
	localPath := repo.LocalPath
	if localPath == "" {
		localPath = fmt.Sprintf("/tmp/salt-formulas-%d", repo.ID)
	}

	// 确保本地仓库已同步
	if err := s.syncGitRepository(repo, localPath); err != nil {
		return fmt.Errorf("failed to sync local repository: %v", err)
	}

	// Salt Master 上的目标路径
	saltPath := "/srv/salt/"

	// 构建 rsync 命令
	// rsync -avz --delete -e "ssh -p PORT -i KEYFILE -o StrictHostKeyChecking=no" LOCAL_PATH/ user@IP:/srv/salt/
	sshCmd := fmt.Sprintf("ssh -p %d -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null", sshPort)
	if sshKeyPath != "" {
		sshCmd = fmt.Sprintf("%s -i %s", sshCmd, sshKeyPath)
	}

	rsyncArgs := []string{
		"-avz",
		"--delete",
		"--exclude", ".git",
		"-e", sshCmd,
		localPath + "/",
		fmt.Sprintf("%s@%s:%s", sshUser, saltMasterHost.IPAddress, saltPath),
	}

	log.Printf("Executing rsync command: rsync %v", rsyncArgs)

	cmd := exec.Command("rsync", rsyncArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("rsync failed: %v, output: %s", err, string(output))
	}

	log.Printf("Rsync output: %s", string(output))
	log.Printf("Successfully synced repository to Salt Master via rsync")
	return nil
}
