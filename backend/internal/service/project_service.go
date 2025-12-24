package service

import (
	"errors"

	"github.com/kronos/backend/internal/models"
	"github.com/kronos/backend/internal/repository"
)

type ProjectService interface {
	CreateProject(name, description string, createdBy uint64) (*models.Project, error)
	GetProject(id uint64) (*models.Project, error)
	ListProjects(page, pageSize int, userID uint64, filters map[string]interface{}) ([]models.Project, int64, error)
	UpdateProject(id uint64, name, description, status string) (*models.Project, error)
	DeleteProject(id uint64) error
	AddMember(projectID, userID uint64, role string) error
	RemoveMember(projectID, userID uint64) error
	GetMembers(projectID uint64) ([]models.User, error)
	CheckUserAccess(projectID, userID uint64) (bool, error)
}

type projectService struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService(projectRepo repository.ProjectRepository) ProjectService {
	return &projectService{projectRepo: projectRepo}
}

func (s *projectService) CreateProject(name, description string, createdBy uint64) (*models.Project, error) {
	// 检查项目名是否已存在
	_, err := s.projectRepo.GetByName(name)
	if err == nil {
		return nil, errors.New("project name already exists")
	}

	project := &models.Project{
		Name:        name,
		Description: description,
		Status:      "active",
		CreatedBy:   createdBy,
	}

	err = s.projectRepo.Create(project)
	if err != nil {
		return nil, err
	}

	// 创建项目时，创建者自动成为owner
	err = s.projectRepo.AddMember(project.ID, createdBy, "owner")
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) GetProject(id uint64) (*models.Project, error) {
	return s.projectRepo.GetByID(id)
}

func (s *projectService) ListProjects(page, pageSize int, userID uint64, filters map[string]interface{}) ([]models.Project, int64, error) {
	if filters == nil {
		filters = make(map[string]interface{})
	}
	filters["user_id"] = userID // 只返回用户有权限的项目

	offset := (page - 1) * pageSize
	return s.projectRepo.List(offset, pageSize, filters)
}

func (s *projectService) UpdateProject(id uint64, name, description, status string) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		project.Name = name
	}
	if description != "" {
		project.Description = description
	}
	if status != "" {
		project.Status = status
	}

	err = s.projectRepo.Update(project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) DeleteProject(id uint64) error {
	return s.projectRepo.Delete(id)
}

func (s *projectService) AddMember(projectID, userID uint64, role string) error {
	return s.projectRepo.AddMember(projectID, userID, role)
}

func (s *projectService) RemoveMember(projectID, userID uint64) error {
	return s.projectRepo.RemoveMember(projectID, userID)
}

func (s *projectService) GetMembers(projectID uint64) ([]models.User, error) {
	return s.projectRepo.GetMembers(projectID)
}

func (s *projectService) CheckUserAccess(projectID, userID uint64) (bool, error) {
	return s.projectRepo.CheckUserAccess(projectID, userID)
}

