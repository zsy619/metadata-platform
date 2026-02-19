package service

import (
	"errors"

	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// ssoUserService 用户服务实现
type ssoUserService struct {
	userRepo          repository.SsoUserRepository
	orgRepo           repository.SsoOrgRepository
	orgUserRepo       repository.SsoOrgUserRepository
	userRoleRepo      repository.SsoUserRoleRepository
	userPosRepo       repository.SsoUserPosRepository
	userGroupUserRepo repository.SsoUserGroupUserRepository
	userRoleGroupRepo repository.SsoUserRoleGroupRepository
}

// NewSsoUserService 创建用户服务实例
func NewSsoUserService(userRepo repository.SsoUserRepository, orgRepo repository.SsoOrgRepository, orgUserRepo repository.SsoOrgUserRepository, userRoleRepo repository.SsoUserRoleRepository, userPosRepo repository.SsoUserPosRepository, userGroupUserRepo repository.SsoUserGroupUserRepository, userRoleGroupRepo repository.SsoUserRoleGroupRepository) SsoUserService {
	return &ssoUserService{
		userRepo:          userRepo,
		orgRepo:           orgRepo,
		orgUserRepo:       orgUserRepo,
		userRoleRepo:      userRoleRepo,
		userPosRepo:       userPosRepo,
		userGroupUserRepo: userGroupUserRepo,
		userRoleGroupRepo: userRoleGroupRepo,
	}
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

	if user.OrgID != "" {
		_, err := s.orgRepo.GetOrgByID(user.OrgID)
		if err != nil {
			return errors.New("组织不存在")
		}
	}

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

	// 删除用户关联的角色
	if err := s.userRoleRepo.DeleteUserRolesByUserID(id); err != nil {
		utils.SugarLogger.Errorw("删除用户角色关联失败", "userID", id, "error", err)
	}

	// 删除用户关联的职位
	if err := s.userPosRepo.DeleteUserPosByUserID(id); err != nil {
		utils.SugarLogger.Errorw("删除用户职位关联失败", "userID", id, "error", err)
	}

	// 删除用户关联的用户组
	if err := s.userGroupUserRepo.DeleteUserGroupUsersByUserID(id); err != nil {
		utils.SugarLogger.Errorw("删除用户组用户关联失败", "userID", id, "error", err)
	}

	// 删除用户关联的角色组
	if err := s.userRoleGroupRepo.DeleteUserRoleGroupsByUserID(id); err != nil {
		utils.SugarLogger.Errorw("删除用户角色组关联失败", "userID", id, "error", err)
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

// GetUserRoles 获取用户的角色ID列表
func (s *ssoUserService) GetUserRoles(userID string) ([]string, error) {
	userRoles, err := s.userRoleRepo.GetUserRolesByUserID(userID)
	if err != nil {
		return nil, err
	}
	roleIDs := make([]string, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}
	return roleIDs, nil
}

// UpdateUserRoles 更新用户的角色关联
func (s *ssoUserService) UpdateUserRoles(userID string, roleIDs []string, createBy string) error {
	if err := s.userRoleRepo.DeleteUserRolesByUserID(userID); err != nil {
		return err
	}
	for _, roleID := range roleIDs {
		userRole := &model.SsoUserRole{
			ID:       utils.GetSnowflake().GenerateIDString(),
			UserID:   userID,
			RoleID:   roleID,
			CreateBy: createBy,
		}
		if err := s.userRoleRepo.CreateUserRole(userRole); err != nil {
			utils.SugarLogger.Errorw("创建用户角色关联失败", "userID", userID, "roleID", roleID, "error", err)
		}
	}
	return nil
}

// GetUserPos 获取用户的职位ID列表
func (s *ssoUserService) GetUserPos(userID string) ([]string, error) {
	userPoss, err := s.userPosRepo.GetUserPosByUserID(userID)
	if err != nil {
		return nil, err
	}
	posIDs := make([]string, 0, len(userPoss))
	for _, up := range userPoss {
		posIDs = append(posIDs, up.PosID)
	}
	return posIDs, nil
}

// UpdateUserPos 更新用户的职位关联
func (s *ssoUserService) UpdateUserPos(userID string, posIDs []string, createBy string) error {
	if err := s.userPosRepo.DeleteUserPosByUserID(userID); err != nil {
		return err
	}
	for _, posID := range posIDs {
		userPos := &model.SsoUserPos{
			ID:       utils.GetSnowflake().GenerateIDString(),
			UserID:   userID,
			PosID:    posID,
			CreateBy: createBy,
		}
		if err := s.userPosRepo.CreateUserPos(userPos); err != nil {
			utils.SugarLogger.Errorw("创建用户职位关联失败", "userID", userID, "posID", posID, "error", err)
		}
	}
	return nil
}

// GetUserGroups 获取用户的用户组ID列表
func (s *ssoUserService) GetUserGroups(userID string) ([]string, error) {
	userGroups, err := s.userGroupUserRepo.GetUserGroupUsersByUserID(userID)
	if err != nil {
		return nil, err
	}
	groupIDs := make([]string, 0, len(userGroups))
	for _, ug := range userGroups {
		groupIDs = append(groupIDs, ug.GroupID)
	}
	return groupIDs, nil
}

// UpdateUserGroups 更新用户的用户组关联
func (s *ssoUserService) UpdateUserGroups(userID string, groupIDs []string, createBy string) error {
	if err := s.userGroupUserRepo.DeleteUserGroupUsersByUserID(userID); err != nil {
		return err
	}
	for _, groupID := range groupIDs {
		userGroup := &model.SsoUserGroupUser{
			ID:       utils.GetSnowflake().GenerateIDString(),
			UserID:   userID,
			GroupID:  groupID,
			CreateBy: createBy,
		}
		if err := s.userGroupUserRepo.CreateUserGroupUser(userGroup); err != nil {
			utils.SugarLogger.Errorw("创建用户组用户关联失败", "userID", userID, "groupID", groupID, "error", err)
		}
	}
	return nil
}

// GetUserRoleGroups 获取用户的角色组ID列表
func (s *ssoUserService) GetUserRoleGroups(userID string) ([]string, error) {
	userRoleGroups, err := s.userRoleGroupRepo.GetUserRoleGroupsByUserID(userID)
	if err != nil {
		return nil, err
	}
	roleGroupIDs := make([]string, 0, len(userRoleGroups))
	for _, urg := range userRoleGroups {
		roleGroupIDs = append(roleGroupIDs, urg.GroupID)
	}
	return roleGroupIDs, nil
}

// UpdateUserRoleGroups 更新用户的角色组关联
func (s *ssoUserService) UpdateUserRoleGroups(userID string, roleGroupIDs []string, createBy string) error {
	if err := s.userRoleGroupRepo.DeleteUserRoleGroupsByUserID(userID); err != nil {
		return err
	}
	for _, roleGroupID := range roleGroupIDs {
		userRoleGroup := &model.SsoUserRoleGroup{
			ID:       utils.GetSnowflake().GenerateIDString(),
			UserID:   userID,
			GroupID:  roleGroupID,
			CreateBy: createBy,
		}
		if err := s.userRoleGroupRepo.CreateUserRoleGroup(userRoleGroup); err != nil {
			utils.SugarLogger.Errorw("创建用户角色组关联失败", "userID", userID, "roleGroupID", roleGroupID, "error", err)
		}
	}
	return nil
}

// GetUserOrgs 获取用户的组织ID列表
func (s *ssoUserService) GetUserOrgs(userID string) ([]string, error) {
	orgUsers, err := s.orgUserRepo.GetOrgUsersByUserID(userID)
	if err != nil {
		return nil, err
	}
	orgIDs := make([]string, 0, len(orgUsers))
	for _, ou := range orgUsers {
		orgIDs = append(orgIDs, ou.OrgID)
	}
	return orgIDs, nil
}

// UpdateUserOrgs 更新用户的组织关联
func (s *ssoUserService) UpdateUserOrgs(userID string, orgIDs []string, createBy string) error {
	if err := s.orgUserRepo.DeleteOrgUsersByUserID(userID); err != nil {
		return err
	}
	for _, orgID := range orgIDs {
		orgUser := &model.SsoOrgUser{
			ID:       utils.GetSnowflake().GenerateIDString(),
			UserID:   userID,
			OrgID:    orgID,
			CreateBy: createBy,
		}
		if err := s.orgUserRepo.CreateOrgUser(orgUser); err != nil {
			utils.SugarLogger.Errorw("创建组织用户关联失败", "userID", userID, "orgID", orgID, "error", err)
		}
	}
	return nil
}
