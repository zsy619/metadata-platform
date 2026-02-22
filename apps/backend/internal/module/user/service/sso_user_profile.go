package service

import (
	"errors"

	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// ssoUserProfileService 用户档案服务实现
type ssoUserProfileService struct {
	profileRepo repository.SsoUserProfileRepository
}

// NewSsoUserProfileService 创建用户档案服务实例
func NewSsoUserProfileService(profileRepo repository.SsoUserProfileRepository) SsoUserProfileService {
	return &ssoUserProfileService{profileRepo: profileRepo}
}

// GetByUserID 获取用户档案
func (s *ssoUserProfileService) GetByUserID(userID string) (*model.SsoUserProfile, error) {
	profile, err := s.profileRepo.GetByUserID(userID)
	if err != nil {
		// 档案不存在时返回空档案，不报错
		return &model.SsoUserProfile{UserID: userID}, nil
	}
	return profile, nil
}

// Upsert 创建或更新用户档案
func (s *ssoUserProfileService) Upsert(profile *model.SsoUserProfile) error {
	if profile.UserID == "" {
		return errors.New("用户ID不能为空")
	}
	// 如无 ID，使用雪花算法生成
	if profile.ID == "" {
		profile.ID = utils.GetSnowflake().GenerateIDString()
	}
	return s.profileRepo.Upsert(profile)
}
