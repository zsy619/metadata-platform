package service

import (
	"errors"

	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// ssoUserContactService 用户联系方式服务实现
type ssoUserContactService struct {
	contactRepo repository.SsoUserContactRepository
}

// NewSsoUserContactService 创建用户联系方式服务实例
func NewSsoUserContactService(contactRepo repository.SsoUserContactRepository) SsoUserContactService {
	return &ssoUserContactService{contactRepo: contactRepo}
}

// GetByUserID 获取用户联系方式列表
func (s *ssoUserContactService) GetByUserID(userID string) ([]model.SsoUserContact, error) {
	return s.contactRepo.GetByUserID(userID)
}

// Create 新增联系方式
func (s *ssoUserContactService) Create(contact *model.SsoUserContact) error {
	if contact.UserID == "" {
		return errors.New("用户ID不能为空")
	}
	if contact.Type == "" || contact.Value == "" {
		return errors.New("联系方式类型和值不能为空")
	}
	contact.ID = utils.GetSnowflake().GenerateIDString()
	return s.contactRepo.Create(contact)
}

// UpdateFields 更新联系方式
func (s *ssoUserContactService) UpdateFields(id string, fields map[string]any) error {
	_, err := s.contactRepo.GetByID(id)
	if err != nil {
		return errors.New("联系方式不存在")
	}
	return s.contactRepo.UpdateFields(id, fields)
}

// Delete 删除联系方式
func (s *ssoUserContactService) Delete(userID, id string) error {
	contact, err := s.contactRepo.GetByID(id)
	if err != nil {
		return errors.New("联系方式不存在")
	}
	if contact.UserID != userID {
		return errors.New("无权删除此联系方式")
	}
	return s.contactRepo.Delete(id)
}
