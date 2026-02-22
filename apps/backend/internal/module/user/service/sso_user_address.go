package service

import (
	"errors"

	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// ssoUserAddressService 用户地址服务实现
type ssoUserAddressService struct {
	addrRepo repository.SsoUserAddressRepository
}

// NewSsoUserAddressService 创建用户地址服务实例
func NewSsoUserAddressService(addrRepo repository.SsoUserAddressRepository) SsoUserAddressService {
	return &ssoUserAddressService{addrRepo: addrRepo}
}

// GetByUserID 获取用户地址列表
func (s *ssoUserAddressService) GetByUserID(userID string) ([]model.SsoUserAddress, error) {
	return s.addrRepo.GetByUserID(userID)
}

// Create 新增地址
func (s *ssoUserAddressService) Create(addr *model.SsoUserAddress) error {
	if addr.UserID == "" {
		return errors.New("用户ID不能为空")
	}
	addr.ID = utils.GetSnowflake().GenerateIDString()
	// 如首条地址，设为默认
	existing, _ := s.addrRepo.GetByUserID(addr.UserID)
	if len(existing) == 0 {
		addr.IsDefault = true
	}
	return s.addrRepo.Create(addr)
}

// UpdateFields 更新地址
func (s *ssoUserAddressService) UpdateFields(id string, fields map[string]any) error {
	_, err := s.addrRepo.GetByID(id)
	if err != nil {
		return errors.New("地址不存在")
	}
	return s.addrRepo.UpdateFields(id, fields)
}

// SetDefault 将指定地址设为默认
func (s *ssoUserAddressService) SetDefault(userID, id string) error {
	addr, err := s.addrRepo.GetByID(id)
	if err != nil {
		return errors.New("地址不存在")
	}
	if addr.UserID != userID {
		return errors.New("无权操作此地址")
	}
	// 先清除该用户所有默认标记
	if err := s.addrRepo.ClearDefault(userID); err != nil {
		return err
	}
	// 设置新的默认地址
	return s.addrRepo.UpdateFields(id, map[string]any{"is_default": true})
}

// Delete 删除地址
func (s *ssoUserAddressService) Delete(userID, id string) error {
	addr, err := s.addrRepo.GetByID(id)
	if err != nil {
		return errors.New("地址不存在")
	}
	if addr.UserID != userID {
		return errors.New("无权删除此地址")
	}
	return s.addrRepo.Delete(id)
}
