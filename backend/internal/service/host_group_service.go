package service

import (
	"github.com/kronos/backend/internal/models"
	"github.com/kronos/backend/internal/repository"
)

type HostGroupService interface {
	CreateGroup(group *models.HostGroup) error
	GetGroup(id uint64) (*models.HostGroup, error)
	ListGroups(page, pageSize int, filters map[string]interface{}) ([]models.HostGroup, int64, error)
	UpdateGroup(id uint64, group *models.HostGroup) error
	DeleteGroup(id uint64) error
}

type hostGroupService struct {
	groupRepo repository.HostGroupRepository
}

func NewHostGroupService(groupRepo repository.HostGroupRepository) HostGroupService {
	return &hostGroupService{groupRepo: groupRepo}
}

func (s *hostGroupService) CreateGroup(group *models.HostGroup) error {
	return s.groupRepo.Create(group)
}

func (s *hostGroupService) GetGroup(id uint64) (*models.HostGroup, error) {
	return s.groupRepo.GetByID(id)
}

func (s *hostGroupService) ListGroups(page, pageSize int, filters map[string]interface{}) ([]models.HostGroup, int64, error) {
	offset := (page - 1) * pageSize
	return s.groupRepo.List(offset, pageSize, filters)
}

func (s *hostGroupService) UpdateGroup(id uint64, group *models.HostGroup) error {
	existingGroup, err := s.groupRepo.GetByID(id)
	if err != nil {
		return err
	}

	if group.Name != "" {
		existingGroup.Name = group.Name
	}
	if group.Description != "" {
		existingGroup.Description = group.Description
	}

	return s.groupRepo.Update(existingGroup)
}

func (s *hostGroupService) DeleteGroup(id uint64) error {
	return s.groupRepo.Delete(id)
}

