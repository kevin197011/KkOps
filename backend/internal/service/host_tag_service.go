package service

import (
	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
)

type HostTagService interface {
	CreateTag(tag *models.HostTag) error
	GetTag(id uint64) (*models.HostTag, error)
	ListTags(page, pageSize int, filters map[string]interface{}) ([]models.HostTag, int64, error)
	UpdateTag(id uint64, tag *models.HostTag) error
	DeleteTag(id uint64) error
}

type hostTagService struct {
	tagRepo repository.HostTagRepository
}

func NewHostTagService(tagRepo repository.HostTagRepository) HostTagService {
	return &hostTagService{tagRepo: tagRepo}
}

func (s *hostTagService) CreateTag(tag *models.HostTag) error {
	return s.tagRepo.Create(tag)
}

func (s *hostTagService) GetTag(id uint64) (*models.HostTag, error) {
	return s.tagRepo.GetByID(id)
}

func (s *hostTagService) ListTags(page, pageSize int, filters map[string]interface{}) ([]models.HostTag, int64, error) {
	offset := (page - 1) * pageSize
	return s.tagRepo.List(offset, pageSize, filters)
}

func (s *hostTagService) UpdateTag(id uint64, tag *models.HostTag) error {
	existingTag, err := s.tagRepo.GetByID(id)
	if err != nil {
		return err
	}

	if tag.Name != "" {
		existingTag.Name = tag.Name
	}
	if tag.Color != "" {
		existingTag.Color = tag.Color
	}
	if tag.Description != "" {
		existingTag.Description = tag.Description
	}

	return s.tagRepo.Update(existingTag)
}

func (s *hostTagService) DeleteTag(id uint64) error {
	return s.tagRepo.Delete(id)
}

