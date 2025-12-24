package repository

import (
	"github.com/kronos/backend/internal/models"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *models.Project) error
	GetByID(id uint64) (*models.Project, error)
	GetByName(name string) (*models.Project, error)
	List(offset, limit int, filters map[string]interface{}) ([]models.Project, int64, error)
	Update(project *models.Project) error
	Delete(id uint64) error
	AddMember(projectID, userID uint64, role string) error
	RemoveMember(projectID, userID uint64) error
	GetMembers(projectID uint64) ([]models.User, error)
	GetUserProjects(userID uint64) ([]models.Project, error)
	CheckUserAccess(projectID, userID uint64) (bool, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) GetByID(id uint64) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Creator").First(&project, id).Error
	return &project, err
}

func (r *projectRepository) GetByName(name string) (*models.Project, error) {
	var project models.Project
	err := r.db.Where("name = ?", name).First(&project).Error
	return &project, err
}

func (r *projectRepository) List(offset, limit int, filters map[string]interface{}) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := r.db.Model(&models.Project{})

	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if userID, ok := filters["user_id"].(uint64); ok && userID > 0 {
		// 只返回用户有权限访问的项目
		query = query.Joins("INNER JOIN project_members ON projects.id = project_members.project_id").
			Where("project_members.user_id = ?", userID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Creator").Offset(offset).Limit(limit).Find(&projects).Error
	return projects, total, err
}

func (r *projectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Project{}, id).Error
}

func (r *projectRepository) AddMember(projectID, userID uint64, role string) error {
	return r.db.Create(&models.ProjectMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}).Error
}

func (r *projectRepository) RemoveMember(projectID, userID uint64) error {
	return r.db.Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&models.ProjectMember{}).Error
}

func (r *projectRepository) GetMembers(projectID uint64) ([]models.User, error) {
	var users []models.User
	err := r.db.Table("users").
		Joins("INNER JOIN project_members ON users.id = project_members.user_id").
		Where("project_members.project_id = ?", projectID).
		Find(&users).Error
	return users, err
}

func (r *projectRepository) GetUserProjects(userID uint64) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Table("projects").
		Joins("INNER JOIN project_members ON projects.id = project_members.project_id").
		Where("project_members.user_id = ?", userID).
		Find(&projects).Error
	return projects, err
}

func (r *projectRepository) CheckUserAccess(projectID, userID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}

