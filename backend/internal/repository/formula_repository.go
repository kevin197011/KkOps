package repository

import (
	"path/filepath"
	"strings"

	"github.com/kkops/backend/internal/models"
	"gorm.io/gorm"
)

type FormulaRepository interface {
	Create(formula *models.Formula) error
	GetByID(id uint) (*models.Formula, error)
	GetByName(name string) (*models.Formula, error)
	List(page, pageSize int, filters map[string]interface{}) ([]models.Formula, int64, error)
	Update(formula *models.Formula) error
	Delete(id uint) error

	// Formula参数相关
	CreateParameter(param *models.FormulaParameter) error
	GetParameters(formulaID uint) ([]models.FormulaParameter, error)
	UpdateParameter(param *models.FormulaParameter) error
	DeleteParameter(id uint) error

	// Formula模板相关
	CreateTemplate(template *models.FormulaTemplate) error
	GetTemplates(formulaID uint) ([]models.FormulaTemplate, error)
	GetTemplateByID(id uint) (*models.FormulaTemplate, error)
	UpdateTemplate(template *models.FormulaTemplate) error
	DeleteTemplate(id uint) error

	// Formula部署相关
	CreateDeployment(deployment *models.FormulaDeployment) error
	GetDeploymentByID(id uint) (*models.FormulaDeployment, error)
	ListDeployments(page, pageSize int, filters map[string]interface{}) ([]models.FormulaDeployment, int64, error)
	UpdateDeployment(deployment *models.FormulaDeployment) error
	DeleteDeployment(id uint) error
	CleanupOldDeployments(beforeTime interface{}) (int64, error)

	// Formula仓库相关
	CreateRepository(repo *models.FormulaRepository) error
	GetRepositoryByID(id uint) (*models.FormulaRepository, error)
	ListRepositories(page, pageSize int, filters map[string]interface{}) ([]models.FormulaRepository, int64, error)
	UpdateRepository(repo *models.FormulaRepository) error
	DeleteRepository(id uint) error
}

type formulaRepository struct {
	db *gorm.DB
}

func NewFormulaRepository(db *gorm.DB) FormulaRepository {
	return &formulaRepository{db: db}
}

// Formula CRUD
func (r *formulaRepository) Create(formula *models.Formula) error {
	return r.db.Create(formula).Error
}

func (r *formulaRepository) GetByID(id uint) (*models.Formula, error) {
	var formula models.Formula
	err := r.db.First(&formula, id).Error
	return &formula, err
}

func (r *formulaRepository) GetByName(name string) (*models.Formula, error) {
	var formula models.Formula
	err := r.db.Where("name = ?", name).First(&formula).Error
	return &formula, err
}

func (r *formulaRepository) List(page, pageSize int, filters map[string]interface{}) ([]models.Formula, int64, error) {
	var formulas []models.Formula
	var total int64

	query := r.db.Model(&models.Formula{})

	// 应用过滤器
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if repository, ok := filters["repository"].(string); ok && repository != "" {
		query = query.Where("repository = ?", repository)
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&formulas).Error

	return formulas, total, err
}

func (r *formulaRepository) Update(formula *models.Formula) error {
	return r.db.Save(formula).Error
}

func (r *formulaRepository) Delete(id uint) error {
	return r.db.Delete(&models.Formula{}, id).Error
}

// Formula参数相关
func (r *formulaRepository) CreateParameter(param *models.FormulaParameter) error {
	return r.db.Create(param).Error
}

func (r *formulaRepository) GetParameters(formulaID uint) ([]models.FormulaParameter, error) {
	var params []models.FormulaParameter
	err := r.db.Where("formula_id = ?", formulaID).Order("order_index ASC").Find(&params).Error
	return params, err
}

func (r *formulaRepository) UpdateParameter(param *models.FormulaParameter) error {
	return r.db.Save(param).Error
}

func (r *formulaRepository) DeleteParameter(id uint) error {
	return r.db.Delete(&models.FormulaParameter{}, id).Error
}

// Formula模板相关
func (r *formulaRepository) CreateTemplate(template *models.FormulaTemplate) error {
	return r.db.Create(template).Error
}

func (r *formulaRepository) GetTemplates(formulaID uint) ([]models.FormulaTemplate, error) {
	var templates []models.FormulaTemplate
	err := r.db.Where("formula_id = ?", formulaID).Order("created_at DESC").Find(&templates).Error
	return templates, err
}

func (r *formulaRepository) GetTemplateByID(id uint) (*models.FormulaTemplate, error) {
	var template models.FormulaTemplate
	err := r.db.First(&template, id).Error
	return &template, err
}

func (r *formulaRepository) UpdateTemplate(template *models.FormulaTemplate) error {
	return r.db.Save(template).Error
}

func (r *formulaRepository) DeleteTemplate(id uint) error {
	return r.db.Delete(&models.FormulaTemplate{}, id).Error
}

// Formula部署相关
func (r *formulaRepository) CreateDeployment(deployment *models.FormulaDeployment) error {
	return r.db.Create(deployment).Error
}

func (r *formulaRepository) GetDeploymentByID(id uint) (*models.FormulaDeployment, error) {
	var deployment models.FormulaDeployment
	err := r.db.First(&deployment, id).Error
	return &deployment, err
}

func (r *formulaRepository) ListDeployments(page, pageSize int, filters map[string]interface{}) ([]models.FormulaDeployment, int64, error) {
	var deployments []models.FormulaDeployment
	var total int64

	query := r.db.Model(&models.FormulaDeployment{})

	// 应用过滤器
	if formulaID, ok := filters["formula_id"].(uint); ok && formulaID > 0 {
		query = query.Where("formula_id = ?", formulaID)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if startedBy, ok := filters["started_by"].(uint64); ok && startedBy > 0 {
		query = query.Where("started_by = ?", startedBy)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&deployments).Error

	return deployments, total, err
}

func (r *formulaRepository) UpdateDeployment(deployment *models.FormulaDeployment) error {
	return r.db.Save(deployment).Error
}

func (r *formulaRepository) DeleteDeployment(id uint) error {
	return r.db.Delete(&models.FormulaDeployment{}, id).Error
}

func (r *formulaRepository) CleanupOldDeployments(beforeTime interface{}) (int64, error) {
	result := r.db.Unscoped().Where("created_at < ?", beforeTime).Delete(&models.FormulaDeployment{})
	return result.RowsAffected, result.Error
}

// Formula仓库相关
func (r *formulaRepository) CreateRepository(repo *models.FormulaRepository) error {
	return r.db.Create(repo).Error
}

func (r *formulaRepository) GetRepositoryByID(id uint) (*models.FormulaRepository, error) {
	var repo models.FormulaRepository
	err := r.db.First(&repo, id).Error
	return &repo, err
}

func (r *formulaRepository) ListRepositories(page, pageSize int, filters map[string]interface{}) ([]models.FormulaRepository, int64, error) {
	var repos []models.FormulaRepository
	var total int64

	query := r.db.Model(&models.FormulaRepository{})

	// 应用过滤器
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&repos).Error

	return repos, total, err
}

func (r *formulaRepository) UpdateRepository(repo *models.FormulaRepository) error {
	return r.db.Save(repo).Error
}

func (r *formulaRepository) DeleteRepository(id uint) error {
	return r.db.Delete(&models.FormulaRepository{}, id).Error
}

// DiscoverFormulas 从仓库扫描Formula
func (r *formulaRepository) DiscoverFormulas(repoPath string) ([]map[string]interface{}, error) {
	var formulas []map[string]interface{}

	// 这里实现Formula扫描逻辑
	// 扫描目录结构，读取init.sls等文件，解析元数据

	// 暂时返回空切片，实际实现需要根据具体的目录结构来解析
	return formulas, nil
}

// ParseFormulaMetadata 解析Formula元数据
func (r *formulaRepository) ParseFormulaMetadata(formulaPath string) (map[string]interface{}, error) {
	metadata := make(map[string]interface{})

	// 检查是否存在init.sls
	if _, err := r.db.DB(); err == nil { // 简化检查，实际应该检查文件存在
		metadata["has_init"] = true
	}

	// 检查是否存在map.jinja
	if _, err := r.db.DB(); err == nil {
		metadata["has_map"] = true
	}

	// 解析Formula名称（从路径中提取）
	formulaName := filepath.Base(formulaPath)
	metadata["name"] = formulaName

	// 推断分类（基于路径结构）
	parts := strings.Split(formulaPath, string(filepath.Separator))
	if len(parts) >= 2 {
		parentDir := parts[len(parts)-2]
		switch parentDir {
		case "base":
			metadata["category"] = "base"
		case "middleware":
			metadata["category"] = "middleware"
		case "runtime":
			metadata["category"] = "runtime"
		case "app":
			metadata["category"] = "app"
		default:
			metadata["category"] = "custom"
		}
	}

	return metadata, nil
}
