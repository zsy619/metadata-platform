package service

import (
	"errors"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// ssoUserService 用户服务实现
type ssoUserService struct {
	userRepo repository.SsoUserRepository
}

// NewSsoUserService 创建用户服务实例
func NewSsoUserService(userRepo repository.SsoUserRepository) SsoUserService {
	return &ssoUserService{userRepo: userRepo}
}

// CreateUser 创建用户
func (s *ssoUserService) CreateUser(user *model.SsoUser) error {
	// 检查账号是否已存在
	existingUser, err := s.userRepo.GetUserByAccount(user.Account)
	if err == nil && existingUser != nil {
		return errors.New("账号已存在")
	}

	// 对密码进行哈希加密
	hashedPassword, salt, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New("密码加密失败")
	}
	user.Password = hashedPassword
	user.Salt = salt

	// 创建用户
	return s.userRepo.CreateUser(user)
}

// GetUserByID 根据ID获取用户
func (s *ssoUserService) GetUserByID(id string) (*model.SsoUser, error) {
	return s.userRepo.GetUserByID(id)
}

// GetUserByAccount 根据账号获取用户
func (s *ssoUserService) GetUserByAccount(account string) (*model.SsoUser, error) {
	return s.userRepo.GetUserByAccount(account)
}

// UpdateUser 更新用户
func (s *ssoUserService) UpdateUser(user *model.SsoUser) error {
	// 检查用户是否存在
	existingUser, err := s.userRepo.GetUserByID(user.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 如果账号发生变化，检查新账号是否已存在
	if existingUser.Account != user.Account {
		otherUser, err := s.userRepo.GetUserByAccount(user.Account)
		if err == nil && otherUser != nil {
			return errors.New("账号已存在")
		}
	}

	// 如果密码不为空，则对密码进行哈希加密
	if user.Password != "" {
		hashedPassword, salt, err := utils.HashPassword(user.Password)
		if err != nil {
			return errors.New("密码加密失败")
		}
		user.Password = hashedPassword
		user.Salt = salt
	} else {
		// 保持原密码不变
		user.Password = existingUser.Password
	}

	// 更新用户
	return s.userRepo.UpdateUser(user)
}

// DeleteUser 删除用户
func (s *ssoUserService) DeleteUser(id string) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 删除用户
	return s.userRepo.DeleteUser(id)
}

// GetAllUsers 获取所有用户
func (s *ssoUserService) GetAllUsers() ([]model.SsoUser, error) {
	return s.userRepo.GetAllUsers()
}

// Login 用户登录
func (s *ssoUserService) Login(account, password string) (string, error) {
	// 根据账号获取用户
	user, err := s.userRepo.GetUserByAccount(account)
	if err != nil {
		return "", errors.New("账号或密码错误")
	}

	// 验证密码
	if !utils.CheckPasswordHash(password, user.Password, user.Salt) {
		return "", errors.New("账号或密码错误")
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Account, user.Kind == 1)
	if err != nil {
		return "", errors.New("生成令牌失败")
	}

	return token, nil
}
