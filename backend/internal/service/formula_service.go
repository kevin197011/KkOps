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
	saltClient  *salt.Client
}

func NewFormulaService(formulaRepo repository.FormulaRepository, saltClient *salt.Client) FormulaService {
	return &formulaService{
		formulaRepo: formulaRepo,
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

	// 执行state.apply
	result, err := s.saltClient.ExecuteState(formula.Name, target, pillarData)
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

	if resultData, ok := result["return"].([]interface{}); ok && len(resultData) > 0 {
		if minions, ok := resultData[0].(map[string]interface{}); ok {
			for minionID, minionResult := range minions {
				results[minionID] = minionResult

				// 检查结果是否成功
				if minionResultMap, ok := minionResult.(map[string]interface{}); ok {
					if resultList, ok := minionResultMap["result"].([]interface{}); ok {
						allSuccess := true
						for _, r := range resultList {
							if resultMap, ok := r.(map[string]interface{}); ok {
								if success, ok := resultMap["result"].(bool); ok && !success {
									allSuccess = false
									break
								}
							}
						}
						if allSuccess {
							successCount++
						} else {
							failedCount++
						}
					} else {
						failedCount++
					}
				} else {
					failedCount++
				}
			}
		}
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

	// 同步仓库（克隆或更新Git仓库）
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
		if err != nil && err.Error() != "record not found" {
			log.Printf("Error checking existing formula %s: %v", formulaCopy.Name, err)
			continue
		}

		if existing != nil {
			// 更新现有Formula
			existing.Description = formulaCopy.Description
			existing.Category = formulaCopy.Category
			existing.Version = formulaCopy.Version
			existing.Path = formulaCopy.Path
			existing.Metadata = formulaCopy.Metadata
			existing.UpdatedAt = time.Now()
			log.Printf("Updating existing formula: %s", formulaCopy.Name)
			if err := s.formulaRepo.Update(existing); err != nil {
				log.Printf("Failed to update formula %s: %v", formulaCopy.Name, err)
				continue
			}
		} else {
			// 创建新Formula
			log.Printf("Creating new formula: %s (created_by=%d)", formulaCopy.Name, formulaCopy.CreatedBy)
			if err := s.formulaRepo.Create(&formulaCopy); err != nil {
				log.Printf("Failed to create formula %s: %v", formulaCopy.Name, err)
				continue // 继续处理其他Formula，不要因为一个失败而停止整个同步
			}
		}
	}

	// 更新最后同步时间
	now := time.Now()
	repo.LastSyncAt = &now
	return s.formulaRepo.UpdateRepository(repo)
}

// Formula发现
func (s *formulaService) DiscoverFormulas(repoPath string) ([]models.Formula, error) {
	var formulas []models.Formula

	log.Printf("Starting Formula discovery in path: %s", repoPath)

	// 检查目录是否存在
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("repository path does not exist: %s", repoPath)
	}

	// 扫描目录结构
	log.Printf("Starting directory walk from: %s", repoPath)
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error walking path %s: %v", path, err)
			return err
		}

		// 跳过根目录本身
		if path == repoPath {
			return nil
		}

		// 检查是否是Formula目录（包含init.sls）
		if info.IsDir() {
			initPath := filepath.Join(path, "init.sls")
			if _, err := os.Stat(initPath); err == nil {
				// 发现Formula
				relPath, _ := filepath.Rel(repoPath, path)
				formulaName := filepath.Base(path)

				log.Printf("Processing potential formula - path: %s, base: %s, rel: %s", path, formulaName, relPath)

				// 确保Formula名称不为空，且不包含特殊字符
				if formulaName == "" || formulaName == "." || formulaName == ".." ||
				   strings.Contains(formulaName, "/") || strings.Contains(formulaName, "\\") ||
				   strings.HasPrefix(formulaName, ".") {
					log.Printf("Skipping invalid formula path: %s (basename: '%s')", path, formulaName)
					return nil
				}

				// 确保相对路径不为空
				if relPath == "" || relPath == "." {
					log.Printf("Skipping root formula path: %s", path)
					return nil
				}

				formula := models.Formula{
					Name:        formulaName,
					Path:        relPath,
					Description: fmt.Sprintf("Auto-discovered formula: %s", formulaName),
					Category:    s.inferCategory(relPath),
					IsActive:    true,
				}

				log.Printf("Discovered Formula: %s (category: %s, path: %s, full_path: %s)", formula.Name, formula.Category, formula.Path, path)

				// 解析元数据
				if metadata, err := s.parseFormulaMetadata(path); err == nil {
					metadataJSON, _ := json.Marshal(metadata)
					formula.Metadata = metadataJSON
				}

				// 最后验证Formula数据
				if formula.Name == "" {
					log.Printf("Warning: Formula name is empty for path: %s", path)
					return nil
				}

				formulas = append(formulas, formula)
			}
		}

		return nil
	})

	log.Printf("Formula discovery completed, found %d formulas", len(formulas))
	return formulas, err
}

// 推断Formula分类
func (s *formulaService) inferCategory(relPath string) string {
	parts := strings.Split(relPath, string(filepath.Separator))
	if len(parts) >= 2 {
		parentDir := parts[0]
		switch parentDir {
		case "base":
			return "base"
		case "middleware":
			return "middleware"
		case "runtime":
			return "runtime"
		case "app":
			return "app"
		default:
			return "custom"
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
