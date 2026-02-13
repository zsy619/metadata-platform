package service

import (
	"errors"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// ssoPosService 职位服务实现
type ssoPosService struct {
	posRepo repository.SsoPosRepository
}

// NewSsoPosService 创建职位服务实例
func NewSsoPosService(posRepo repository.SsoPosRepository) SsoPosService {
	return &ssoPosService{posRepo: posRepo}
}

// CreatePos 创建职位
func (s *ssoPosService) CreatePos(pos *model.SsoPos) error {
	// 检查职位编码是否已存在
	existingPos, err := s.posRepo.GetPosByCode(pos.PosCode)
	if err == nil && existingPos != nil {
		return errors.New("职位编码已存在")
	}

	// 检查父职位是否存在（如果有）
	if pos.ParentID != "" {
		_, err := s.posRepo.GetPosByID(pos.ParentID)
		if err != nil {
			return errors.New("父职位不存在")
		}
	}

	// 使用全局雪花算法生成ID
	pos.ID = utils.GetSnowflake().GenerateIDString()

	// 创建职位CreatePos
	return s.posRepo.CreatePos(pos)
}

// GetPosByID 根据ID获取职位
func (s *ssoPosService) GetPosByID(id string) (*model.SsoPos, error) {
	return s.posRepo.GetPosByID(id)
}

// GetPosByCode 根据编码获取职位
func (s *ssoPosService) GetPosByCode(code string) (*model.SsoPos, error) {
	return s.posRepo.GetPosByCode(code)
}

// UpdatePos 更新职位
func (s *ssoPosService) UpdatePos(pos *model.SsoPos) error {
	return s.posRepo.UpdatePos(pos)
}

// DeletePos 删除职位
func (s *ssoPosService) DeletePos(id string) error {
	return s.posRepo.DeletePos(id)
}

// GetAllPos 获取所有职位
func (s *ssoPosService) GetAllPoss() ([]model.SsoPos, error) {
	return s.posRepo.GetAllPoss()
}
