package service

import (
	"errors"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// ssoRoleService 角色服务实现
type ssoRoleService struct {
	roleRepo repository.SsoRoleRepository
}

// NewSsoRoleService 创建角色服务实例
func NewSsoRoleService(roleRepo repository.SsoRoleRepository) SsoRoleService {
	return &ssoRoleService{roleRepo: roleRepo}
}

// CreateRole 创建角色
func (s *ssoRoleService) CreateRole(role *model.SsoRole) error {
	// 检查角色编码是否已存在
	existingRole, err := s.roleRepo.GetRoleByCode(role.RoleCode)
	if err == nil && existingRole != nil {
		return errors.New("角色编码已存在")
	}

	// 检查父角色是否存在（如果有）
	if role.ParentID != "" {
		_, err := s.roleRepo.GetRoleByID(role.ParentID)
		if err != nil {
			return errors.New("父角色不存在")
		}
	}

	// 使用全局雪花算法生成ID
	role.ID = utils.GetSnowflake().GenerateIDString()

	// 创建角色
	return s.roleRepo.CreateRole(role)
}

// GetRoleByID 根据ID获取角色
func (s *ssoRoleService) GetRoleByID(id string) (*model.SsoRole, error) {
	return s.roleRepo.GetRoleByID(id)
}

// GetRoleByCode 根据编码获取角色
func (s *ssoRoleService) GetRoleByCode(code string) (*model.SsoRole, error) {
	return s.roleRepo.GetRoleByCode(code)
}

// UpdateRole 更新角色
func (s *ssoRoleService) UpdateRole(role *model.SsoRole) error {
	// 检查角色是否存在
	existingRole, err := s.roleRepo.GetRoleByID(role.ID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 如果角色编码发生变化，检查新编码是否已存在
	if existingRole.RoleCode != role.RoleCode {
		anotherRole, err := s.roleRepo.GetRoleByCode(role.RoleCode)
		if err == nil && anotherRole != nil {
			return errors.New("角色编码已存在")
		}
	}

	// 更新角色
	return s.roleRepo.UpdateRole(role)
}

// DeleteRole 删除角色
func (s *ssoRoleService) DeleteRole(id string) error {
	// 检查角色是否存在
	_, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 删除角色
	return s.roleRepo.DeleteRole(id)
}

// GetAllRoles 获取所有角色
func (s *ssoRoleService) GetAllRoles() ([]model.SsoRole, error) {
	return s.roleRepo.GetAllRoles()
}
