package service

import (
	"context"
	"errors"
	auditModel "metadata-platform/internal/module/audit/model"
	auditService "metadata-platform/internal/module/audit/service"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
	"time"

	"github.com/google/uuid"
)

// ssoAuthService 实现
type ssoAuthService struct {
	userRepo     repository.SsoUserRepository
	auditService auditService.AuditService
}

// NewSsoAuthService 创建认证服务实例
func NewSsoAuthService(userRepo repository.SsoUserRepository, auditService auditService.AuditService) SsoAuthService {
	return &ssoAuthService{
		userRepo:     userRepo,
		auditService: auditService,
	}
}

// Login 验证账号密码并返回 token
func (s *ssoAuthService) Login(account string, password string, tenantID uint, clientInfo utils.ClientInfo) (string, string, error) {
	loginTime := time.Now()
	var loginStatus int = 1 // default success
	var errMsg string

	// 记录日志闭包
	defer func() {
		// Log login attempt
		userID := ""
		// If user found, use ID
		user, _ := s.userRepo.GetUserByAccount(account)
		if user != nil {
			userID = user.ID
		}

		s.auditService.RecordLogin(context.Background(), &auditModel.SysLoginLog{
			ID:               uuid.New().String(),
			UserID:           userID,
			Account:          account,
			LoginStatus:      loginStatus,
			ClientIP:         clientInfo.IP,
			Browser:          clientInfo.Browser,
			BrowserVersion:   clientInfo.BrowserVersion,
			BrowserEngine:    clientInfo.BrowserEngine,
			Language:         clientInfo.Language,
			UserAgent:        clientInfo.UserAgent,
			OS:               clientInfo.OS,
			OSVersion:        clientInfo.OSVersion,
			OSArch:           clientInfo.OSArch,
			DeviceType:       clientInfo.DeviceType,
			DeviceModel:      clientInfo.DeviceModel,
			ScreenResolution: clientInfo.ScreenResolution,
			Timezone:         clientInfo.Timezone,
			IPLocation:       clientInfo.IPLocation,
			Platform:         clientInfo.Platform,
			// Device field if needed for backward compact or composite string
			CreateAt:     loginTime,
			ErrorMessage: errMsg,
		})
	}()

	user, err := s.userRepo.GetUserByAccount(account)
	if err != nil {
		loginStatus = 0
		errMsg = err.Error()
		return "", "", err
	}
	// 校验密码
	if user == nil || !utils.CheckPasswordHash(password, user.Password, user.Salt) {
		loginStatus = 0
		errMsg = "用户名或密码错误"
		if user != nil {
			_ = s.userRepo.IncrementLoginError(user.ID)
		}
		return "", "", errors.New("用户名或密码错误")
	}

	// 校验状态
	if user.State != 1 {
		loginStatus = 0
		errMsg = "账号已被禁用"
		return "", "", errors.New("账号已被禁用")
	}

	// 校验过期时间
	if user.EndTime != nil && user.EndTime.Before(time.Now()) {
		loginStatus = 0
		errMsg = "账号已过期"
		return "", "", errors.New("账号已过期")
	}

	// 登录成功，更新审计信息
	_ = s.userRepo.UpdateLoginInfo(user.ID, clientInfo.IP)

	// 生成访问令牌和刷新令牌
	access, err := utils.GenerateToken(user.ID, user.Account, false)
	if err != nil {
		loginStatus = 0
		errMsg = "生成令牌失败: " + err.Error()
		return "", "", err
	}
	refresh, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		loginStatus = 0
		errMsg = "生成刷新令牌失败: " + err.Error()
		return "", "", err
	}
	// 这里可以将 refresh token 保存到数据库或缓存，略过实现
	_ = tenantID // 暂未使用
	return access, refresh, nil
}

// Logout 退出登录
func (s *ssoAuthService) Logout(ctx context.Context, userID string, clientInfo utils.ClientInfo) error {
	// Record logout log
	account := ""
	if userID != "" {
		user, _ := s.userRepo.GetUserByID(userID)
		if user != nil {
			account = user.Account
		}
	}

	s.auditService.RecordLogin(ctx, &auditModel.SysLoginLog{
		ID:               uuid.New().String(),
		UserID:           userID,
		Account:          account,
		LoginStatus:      2, // 2: Logout
		ClientIP:         clientInfo.IP,
		Browser:          clientInfo.Browser,
		BrowserVersion:   clientInfo.BrowserVersion,
		BrowserEngine:    clientInfo.BrowserEngine,
		Language:         clientInfo.Language,
		UserAgent:        clientInfo.UserAgent,
		OS:               clientInfo.OS,
		OSVersion:        clientInfo.OSVersion,
		OSArch:           clientInfo.OSArch,
		DeviceType:       clientInfo.DeviceType,
		DeviceModel:      clientInfo.DeviceModel,
		ScreenResolution: clientInfo.ScreenResolution,
		Timezone:         clientInfo.Timezone,
		IPLocation:       clientInfo.IPLocation,
		Platform:         clientInfo.Platform,
		CreateAt:         time.Now(),
		ErrorMessage:     "",
	})

	return nil
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
	return s.userRepo.GetUserWithDetails(userID)
}

// ChangePassword 修改密码
func (s *ssoAuthService) ChangePassword(userID string, oldPassword string, newPassword string) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 校验旧密码
	if !utils.CheckPasswordHash(oldPassword, user.Password, user.Salt) {
		return errors.New("原密码错误")
	}

	// 生成新盐值和哈希
	hash, salt, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hash
	user.Salt = salt

	// 更新用户
	return s.userRepo.UpdateUser(user)
}
