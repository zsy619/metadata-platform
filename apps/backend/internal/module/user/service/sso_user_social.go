package service

import (
	"errors"

	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// ssoUserSocialService 用户第三方账号服务实现
type ssoUserSocialService struct {
	socialRepo repository.SsoUserSocialRepository
}

// NewSsoUserSocialService 创建用户第三方账号服务实例
func NewSsoUserSocialService(socialRepo repository.SsoUserSocialRepository) SsoUserSocialService {
	return &ssoUserSocialService{socialRepo: socialRepo}
}

// GetByUserID 获取用户第三方账号绑定列表
func (s *ssoUserSocialService) GetByUserID(userID string) ([]model.SsoUserSocial, error) {
	return s.socialRepo.GetByUserID(userID)
}

// Bind 绑定第三方账号
func (s *ssoUserSocialService) Bind(social *model.SsoUserSocial) error {
	if social.UserID == "" || social.Provider == "" || social.OpenID == "" {
		return errors.New("用户ID、平台和OpenID不能为空")
	}
	// 检查 openid 是否已被其他账号绑定
	existing, err := s.socialRepo.GetByProviderAndOpenID(social.Provider, social.OpenID)
	if err == nil && existing != nil && existing.UserID != social.UserID {
		return errors.New("该第三方账号已被其他用户绑定")
	}
	social.ID = utils.GetSnowflake().GenerateIDString()
	return s.socialRepo.Create(social)
}

// Unbind 解绑第三方账号
func (s *ssoUserSocialService) Unbind(userID, id string) error {
	social, err := s.socialRepo.GetByID(id)
	if err != nil {
		return errors.New("绑定记录不存在")
	}
	if social.UserID != userID {
		return errors.New("无权解绑此账号")
	}
	return s.socialRepo.Delete(id)
}
