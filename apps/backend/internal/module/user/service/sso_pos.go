package service

import (
	"errors"

	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// ssoPosService 职位服务实现
type ssoPosService struct {
	posRepo     repository.SsoPosRepository
	posRoleRepo repository.SsoPosRoleRepository
	userPosRepo repository.SsoUserPosRepository
	casbinSync  SsoCasbinSyncService
}

// NewSsoPosService 创建职位服务实例
func NewSsoPosService(posRepo repository.SsoPosRepository, posRoleRepo repository.SsoPosRoleRepository, userPosRepo repository.SsoUserPosRepository, casbinSync SsoCasbinSyncService) SsoPosService {
	return &ssoPosService{
		posRepo:     posRepo,
		posRoleRepo: posRoleRepo,
		userPosRepo: userPosRepo,
		casbinSync:  casbinSync,
	}
}

// CreatePos 创建职位
func (s *ssoPosService) CreatePos(pos *model.SsoPos) error {
	// 检查职位编码是否已存在
	existingPos, err := s.posRepo.GetPosByCode(pos.PosCode)
	if err == nil && existingPos != nil {
		return errors.New("职位编码已存在")
	}

	// 检查父职位是否存在（如果有）
	if pos.ParentID != "" {
		_, err := s.posRepo.GetPosByID(pos.ParentID)
		if err != nil {
			return errors.New("父职位不存在")
		}
	}

	// 使用全局雪花算法生成ID
	pos.ID = utils.GetSnowflake().GenerateIDString()

	// 自动获取 Sort 最大值并加1
	if pos.Sort == 0 {
		maxSort, err := s.posRepo.GetMaxSort()
		if err == nil {
			pos.Sort = maxSort + 1
		}
	}

	// 创建职位CreatePos
	return s.posRepo.CreatePos(pos)
}

// GetPosByID 根据ID获取职位
func (s *ssoPosService) GetPosByID(id string) (*model.SsoPos, error) {
	return s.posRepo.GetPosByID(id)
}

// GetPosByCode 根据编码获取职位
func (s *ssoPosService) GetPosByCode(code string) (*model.SsoPos, error) {
	return s.posRepo.GetPosByCode(code)
}

// UpdatePos 更新职位
func (s *ssoPosService) UpdatePos(pos *model.SsoPos) error {
	return s.posRepo.UpdatePos(pos)
}

// UpdatePosFields 更新职位指定字段
// 使用 map 方式只更新指定的字段，避免全量更新
// 会检查职位是否存在，以及如果更新了职位编码，会检查新编码是否已存在
func (s *ssoPosService) UpdatePosFields(id string, fields map[string]any) error {
	// 检查职位是否存在
	_, err := s.posRepo.GetPosByID(id)
	if err != nil {
		return errors.New("职位不存在")
	}

	// 如果更新了职位编码，检查新编码是否已存在
	if posCode, ok := fields["pos_code"]; ok && posCode != "" {
		anotherPos, err := s.posRepo.GetPosByCode(posCode.(string))
		if err == nil && anotherPos != nil && anotherPos.ID != id {
			return errors.New("职位编码已存在")
		}
	}

	return s.posRepo.UpdatePosFields(id, fields)
}

// DeletePos 删除职位
func (s *ssoPosService) DeletePos(id string) error {
	// 检查职位是否存在
	pos, err := s.posRepo.GetPosByID(id)
	if err != nil {
		return errors.New("职位不存在")
	}

	// 检查是否为系统内置职位
	if pos.IsSystem {
		return errors.New("系统内置职位不允许删除")
	}

	// 删除职位关联的角色
	if err := s.posRoleRepo.DeletePosRolesByPosID(id); err != nil {
		utils.SugarLogger.Errorw("删除职位角色关联失败", "posID", id, "error", err)
	}

	// 删除职位关联的用户
	if err := s.userPosRepo.DeleteUserPosByPosID(id); err != nil {
		utils.SugarLogger.Errorw("删除用户职位关联失败", "posID", id, "error", err)
	}

	// 删除职位
	err = s.posRepo.DeletePos(id)
	if err == nil {
		_ = s.casbinSync.SyncPos(id) // SyncPos 里会清空关联 p
	}
	return err
}

// GetAllPos 获取所有职位
func (s *ssoPosService) GetAllPoss() ([]model.SsoPos, error) {
	return s.posRepo.GetAllPoss()
}

// GetPosRoles 获取职位的角色ID列表
func (s *ssoPosService) GetPosRoles(posID string) ([]string, error) {
	// 检查职位是否存在
	_, err := s.posRepo.GetPosByID(posID)
	if err != nil {
		return nil, errors.New("职位不存在")
	}

	// 获取职位角色关联
	posRoles, err := s.posRoleRepo.GetPosRolesByPosID(posID)
	if err != nil {
		return nil, err
	}

	// 提取角色ID列表
	roleIDs := make([]string, 0, len(posRoles))
	for _, pr := range posRoles {
		roleIDs = append(roleIDs, pr.RoleID)
	}

	return roleIDs, nil
}

// UpdatePosRoles 更新职位的角色关联
func (s *ssoPosService) UpdatePosRoles(posID string, roleIDs []string, createBy string) error {
	// 检查职位是否存在
	_, err := s.posRepo.GetPosByID(posID)
	if err != nil {
		return errors.New("职位不存在")
	}

	// 删除原有的职位角色关联
	if err := s.posRoleRepo.DeletePosRolesByPosID(posID); err != nil {
		return err
	}

	// 创建新的职位角色关联
	for _, roleID := range roleIDs {
		posRole := &model.SsoPosRole{
			ID:       utils.GetSnowflake().GenerateIDString(),
			PosID:    posID,
			RoleID:   roleID,
			CreateBy: createBy,
		}
		if err := s.posRoleRepo.CreatePosRole(posRole); err != nil {
			utils.SugarLogger.Errorw("创建职位角色关联失败", "posID", posID, "roleID", roleID, "error", err)
		}
	}

	_ = s.casbinSync.SyncPos(posID)
	return nil
}
