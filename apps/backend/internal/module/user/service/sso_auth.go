package service

import (
	"errors"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// ssoAuthService 实现
type ssoAuthService struct {
	userRepo repository.SsoUserRepository
}

// NewSsoAuthService 创建认证服务实例
func NewSsoAuthService(userRepo repository.SsoUserRepository) SsoAuthService {
	return &ssoAuthService{userRepo: userRepo}
}

// Login 验证账号密码并返回 token
func (s *ssoAuthService) Login(account string, password string, tenantID uint) (string, string, error) {
	user, err := s.userRepo.GetUserByAccount(account)
	if err != nil {
		return "", "", err
	}
	if user == nil || !utils.CheckPasswordHash(password, user.Password, user.Salt) {
		return "", "", errors.New("invalid credentials")
	}
	// 生成访问令牌和刷新令牌
	access, err := utils.GenerateToken(user.ID, user.Account, false)
	if err != nil {
		return "", "", err
	}
	refresh, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}
	// 这里可以将 refresh token 保存到数据库或缓存，略过实现
	_ = tenantID // 暂未使用
	return access, refresh, nil
}

// Refresh 使用刷新令牌生成新的访问令牌
func (s *ssoAuthService) Refresh(refreshToken string) (string, error) {
	claims, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}
	// 根据 claims.Subject（userID）重新生成访问令牌
	// 为简化，这里直接使用 userID 生成 token，假设用户名和管理员标识可通过查询获取
	user, err := s.userRepo.GetUserByID(claims.Subject)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}
	return utils.GenerateToken(user.ID, user.Account, false)
}

// GetUserInfo 根据用户ID获取完整信息
func (s *ssoAuthService) GetUserInfo(userID string) (*model.SsoUser, error) {
	return s.userRepo.GetUserByID(userID)
}
