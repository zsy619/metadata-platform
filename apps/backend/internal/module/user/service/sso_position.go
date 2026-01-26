package service

import (
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
)

// ssoPositionService 职位服务实现
type ssoPositionService struct {
	posRepo repository.SsoPositionRepository
}

// NewSsoPositionService 创建职位服务实例
func NewSsoPositionService(posRepo repository.SsoPositionRepository) SsoPositionService {
	return &ssoPositionService{posRepo: posRepo}
}

// CreatePosition 创建职位
func (s *ssoPositionService) CreatePosition(pos *model.SsoPosition) error {
	return s.posRepo.CreatePosition(pos)
}

// GetPositionByID 根据ID获取职位
func (s *ssoPositionService) GetPositionByID(id string) (*model.SsoPosition, error) {
	return s.posRepo.GetPositionByID(id)
}

// GetPositionByCode 根据编码获取职位
func (s *ssoPositionService) GetPositionByCode(code string) (*model.SsoPosition, error) {
	return s.posRepo.GetPositionByCode(code)
}

// UpdatePosition 更新职位
func (s *ssoPositionService) UpdatePosition(pos *model.SsoPosition) error {
	return s.posRepo.UpdatePosition(pos)
}

// DeletePosition 删除职位
func (s *ssoPositionService) DeletePosition(id string) error {
	return s.posRepo.DeletePosition(id)
}

// GetAllPositions 获取所有职位
func (s *ssoPositionService) GetAllPositions() ([]model.SsoPosition, error) {
	return s.posRepo.GetAllPositions()
}
