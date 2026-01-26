package service

import (
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/repository"
)

// MdModelFieldEnhancementService 字段增强规范服务接口
type MdModelFieldEnhancementService interface {
	CreateEnhancement(enh *model.MdModelFieldEnhancement) error
	GetEnhancementByFieldID(fieldID string) (*model.MdModelFieldEnhancement, error)
	GetEnhancementsByModelID(modelID string) ([]model.MdModelFieldEnhancement, error)
	UpdateEnhancement(enh *model.MdModelFieldEnhancement) error
	DeleteEnhancement(id string) error
	BatchUpdateEnhancements(enhancements []model.MdModelFieldEnhancement) error
}

type mdModelFieldEnhancementService struct {
	enhRepo repository.MdModelFieldEnhancementRepository
}

// NewMdModelFieldEnhancementService 创建字段增强规范服务实例
func NewMdModelFieldEnhancementService(enhRepo repository.MdModelFieldEnhancementRepository) MdModelFieldEnhancementService {
	return &mdModelFieldEnhancementService{enhRepo: enhRepo}
}

func (s *mdModelFieldEnhancementService) CreateEnhancement(enh *model.MdModelFieldEnhancement) error {
	return s.enhRepo.CreateEnhancement(enh)
}

func (s *mdModelFieldEnhancementService) GetEnhancementByFieldID(fieldID string) (*model.MdModelFieldEnhancement, error) {
	return s.enhRepo.GetEnhancementByFieldID(fieldID)
}

func (s *mdModelFieldEnhancementService) GetEnhancementsByModelID(modelID string) ([]model.MdModelFieldEnhancement, error) {
	return s.enhRepo.GetEnhancementsByModelID(modelID)
}

func (s *mdModelFieldEnhancementService) UpdateEnhancement(enh *model.MdModelFieldEnhancement) error {
	return s.enhRepo.UpdateEnhancement(enh)
}

func (s *mdModelFieldEnhancementService) DeleteEnhancement(id string) error {
	return s.enhRepo.DeleteEnhancement(id)
}

func (s *mdModelFieldEnhancementService) BatchUpdateEnhancements(enhs []model.MdModelFieldEnhancement) error {
	return s.enhRepo.BatchUpdateEnhancements(enhs)
}
