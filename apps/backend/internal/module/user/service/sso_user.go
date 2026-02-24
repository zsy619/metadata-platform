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
	// 用户扩展信息仓库
	userProfileRepo repository.SsoUserProfileRepository
	userAddressRepo repository.SsoUserAddressRepository
	userContactRepo repository.SsoUserContactRepository
	userSocialRepo  repository.SsoUserSocialRepository
	casbinSync      SsoCasbinSyncService
}

// NewSsoUserService 创建用户服务实例
func NewSsoUserService(
	userRepo repository.SsoUserRepository,
	orgRepo repository.SsoOrgRepository,
	orgUserRepo repository.SsoOrgUserRepository,
	userRoleRepo repository.SsoUserRoleRepository,
	userPosRepo repository.SsoUserPosRepository,
	userGroupUserRepo repository.SsoUserGroupUserRepository,
	userRoleGroupRepo repository.SsoUserRoleGroupRepository,
	userProfileRepo repository.SsoUserProfileRepository,
	userAddressRepo repository.SsoUserAddressRepository,
	userContactRepo repository.SsoUserContactRepository,
	userSocialRepo repository.SsoUserSocialRepository,
	casbinSync SsoCasbinSyncService,
) SsoUserService {
	return &ssoUserService{
		userRepo:          userRepo,
		orgRepo:           orgRepo,
		orgUserRepo:       orgUserRepo,
		userRoleRepo:      userRoleRepo,
		userPosRepo:       userPosRepo,
		userGroupUserRepo: userGroupUserRepo,
		userRoleGroupRepo: userRoleGroupRepo,
		userProfileRepo:   userProfileRepo,
		userAddressRepo:   userAddressRepo,
		userContactRepo:   userContactRepo,
		userSocialRepo:    userSocialRepo,
		casbinSync:        casbinSync,
	}
}

// CreateUser 创建用户
func (s *ssoUserService) CreateUser(user *model.SsoUser) error {
	// 检查账号是否已存在
	existingUser, err := s.userRepo.GetUserByAccount(user.Account)
	if err == nil && existingUser != nil {
		return errors.New("账号已存在")
	}
	// 检查手机号唯一性
	if user.Mobile != "" {
		if m, err := s.userRepo.GetUserByMobile(user.Mobile); err == nil && m != nil {
			return errors.New("手机号已被其他用户使用")
		}
	}
	// 检查邮箱唯一性
	if user.Email != "" {
		if m, err := s.userRepo.GetUserByEmail(user.Email); err == nil && m != nil {
			return errors.New("邮箱已被其他用户使用")
		}
	}

	// 生成雪花 ID
	if user.ID == "" {
		user.ID = utils.GetSnowflake().GenerateIDString()
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

	updates := make(map[string]any)

	// 如果账号发生变化，检查新账号是否已存在
	if user.Account != "" && existingUser.Account != user.Account {
		otherUser, err := s.userRepo.GetUserByAccount(user.Account)
		if err == nil && otherUser != nil {
			return errors.New("账号已存在")
		}
		updates["account"] = user.Account
	}
	
	// 检查手机号唯一性（排除自身）
	if user.Mobile != "" && user.Mobile != existingUser.Mobile {
		if m, err := s.userRepo.GetUserByMobile(user.Mobile); err == nil && m != nil && m.ID != user.ID {
			return errors.New("手机号已被其他用户使用")
		}
	}
	// 手机号允许清空，所以即便是空也可能需要更新（此处由传入的 user 决定，如果用 map 逻辑，前端传了就可以更新。如果 user 结构体未区分未传和清空，这会导致空值覆盖，这里需要处理：通常业务层是全量覆盖可选项）
	// 因为传进来的是组装好的 user，我们按现有逻辑直接写入
	updates["mobile"] = user.Mobile

	// 检查邮箱唯一性（排除自身）
	if user.Email != "" && user.Email != existingUser.Email {
		if m, err := s.userRepo.GetUserByEmail(user.Email); err == nil && m != nil && m.ID != user.ID {
			return errors.New("邮箱已被其他用户使用")
		}
	}
	updates["email"] = user.Email

	// 如果密码不为空，则对密码进行哈希加密
	if user.Password != "" {
		hashedPassword, salt, err := utils.HashPassword(user.Password)
		if err != nil {
			return errors.New("密码加密失败")
		}
		updates["password"] = hashedPassword
		updates["salt"] = salt
	}

	if user.Name != "" {
		updates["name"] = user.Name
	}
	if user.Kind != 0 {
		updates["kind"] = user.Kind
	}
	// Status 是 int，0 也有意义，前面 handler 已确保传入的是修改后的值
	updates["status"] = user.Status
	// 备注可以直接覆盖
	updates["remark"] = user.Remark
	// 过期时间覆盖
	updates["end_time"] = user.EndTime

	// 更新用户
	return s.userRepo.UpdateUser(user.ID, updates)
}

// DeleteUser 删除用户
func (s *ssoUserService) DeleteUser(id string) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 系统内置用户不允许删除
	if user.IsSystem {
		return errors.New("系统内置用户不允许删除")
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
	// 删除用户档案
	if err := s.userProfileRepo.Delete(id); err != nil {
		utils.SugarLogger.Warnw("删除用户档案失败", "userID", id, "error", err)
	}
	// 删除用户地址簿
	addrs, _ := s.userAddressRepo.GetByUserID(id)
	for _, addr := range addrs {
		_ = s.userAddressRepo.Delete(addr.ID)
	}
	// 删除用户联系方式
	contacts, _ := s.userContactRepo.GetByUserID(id)
	for _, c := range contacts {
		_ = s.userContactRepo.Delete(c.ID)
	}
	// 删除用户第三方账号
	socials, _ := s.userSocialRepo.GetByUserID(id)
	for _, s2 := range socials {
		_ = s.userSocialRepo.Delete(s2.ID)
	}

	// 删除用户
	if err := s.userRepo.DeleteUser(id); err != nil {
		return err
	}
	// 同步清空 Casbin 权限
	_ = s.casbinSync.SyncUser(id)
	return nil
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
	// 同步增量 Casbin 规则
	_ = s.casbinSync.SyncUser(userID)
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
	// 同步增量 Casbin 规则
	_ = s.casbinSync.SyncUser(userID)
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
	// 同步增量 Casbin 规则
	_ = s.casbinSync.SyncUser(userID)
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
	// 同步增量 Casbin 规则
	_ = s.casbinSync.SyncUser(userID)
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
	// 同步增量 Casbin 规则
	_ = s.casbinSync.SyncUser(userID)
	return nil
}
