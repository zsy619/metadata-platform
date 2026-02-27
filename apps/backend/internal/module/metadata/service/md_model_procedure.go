package service

import (
	"errors"

	"metadata-platform/internal/module/metadata/adapter"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/repository"
	"metadata-platform/internal/utils"
)

// MdProcedureService 存储过程/函数服务接口
type MdModelProcedureService interface {
	CreateProcedure(proc *model.MdModelProcedure) error
	GetProcedureByID(id string) (*model.MdModelProcedure, error)
	GetProcedureByName(connID, procSchema, procName string) (*model.MdModelProcedure, error)
	UpdateProcedure(proc *model.MdModelProcedure) error
	DeleteProcedure(id string) error
	GetProceduresByConnID(connID string) ([]model.MdModelProcedure, error)
	GetAllProcedures(tenantID string) ([]model.MdModelProcedure, error)
	SaveSelectedProcedures(connID, tenantID string, procSchema string, procedures []adapter.ProcedureInfo, connName string) error
	GetParamsByProcID(procID string) ([]model.MdModelProcedureParam, error)
}

// mdProcedureService 存储过程/函数服务实现
type mdModelProcedureService struct {
	procRepo  repository.MdModelProcedureRepository
	connRepo  repository.MdConnRepository
	snowflake *utils.Snowflake
}

// NewMdModelProcedureService 创建存储过程/函数服务实例
func NewMdModelProcedureService(procRepo repository.MdModelProcedureRepository, connRepo repository.MdConnRepository) MdModelProcedureService {
	snowflake := utils.NewSnowflake(1, 1)
	return &mdModelProcedureService{
		procRepo:  procRepo,
		connRepo:  connRepo,
		snowflake: snowflake,
	}
}

// CreateProcedure 创建存储过程/函数
func (s *mdModelProcedureService) CreateProcedure(proc *model.MdModelProcedure) error {
	proc.ID = s.snowflake.GenerateIDString()

	existingProc, err := s.procRepo.GetProcedureByID(proc.ID)
	if err == nil && existingProc != nil {
		return errors.New("存储过程/函数已存在")
	}

	return s.procRepo.CreateProcedure(proc)
}

// GetProcedureByID 根据ID获取存储过程/函数
func (s *mdModelProcedureService) GetProcedureByID(id string) (*model.MdModelProcedure, error) {
	return s.procRepo.GetProcedureByID(id)
}

// GetProcedureByName 根据连接ID、模式和名称获取存储过程/函数
func (s *mdModelProcedureService) GetProcedureByName(connID, procSchema, procName string) (*model.MdModelProcedure, error) {
	return s.procRepo.GetProcedureByName(connID, procSchema, procName)
}

// UpdateProcedure 更新存储过程/函数
func (s *mdModelProcedureService) UpdateProcedure(proc *model.MdModelProcedure) error {
	existingProc, err := s.procRepo.GetProcedureByID(proc.ID)
	if err != nil {
		return errors.New("存储过程/函数不存在")
	}

	if existingProc.ProcName != proc.ProcName || existingProc.ProcSchema != proc.ProcSchema {
		otherProc, err := s.procRepo.GetProcedureByName(proc.ConnID, proc.ProcSchema, proc.ProcName)
		if err == nil && otherProc != nil && otherProc.ID != proc.ID {
			return errors.New("存储过程/函数名称已存在")
		}
	}

	return s.procRepo.UpdateProcedure(proc)
}

// DeleteProcedure 删除存储过程/函数
func (s *mdModelProcedureService) DeleteProcedure(id string) error {
	_, err := s.procRepo.GetProcedureByID(id)
	if err != nil {
		return errors.New("存储过程/函数不存在")
	}

	err = s.procRepo.DeleteParamsByProcID(id)
	if err != nil {
		return err
	}

	return s.procRepo.DeleteProcedure(id)
}

// GetProceduresByConnID 根据连接ID获取所有存储过程/函数
func (s *mdModelProcedureService) GetProceduresByConnID(connID string) ([]model.MdModelProcedure, error) {
	return s.procRepo.GetProceduresByConnID(connID)
}

// GetAllProcedures 获取所有存储过程/函数
func (s *mdModelProcedureService) GetAllProcedures(tenantID string) ([]model.MdModelProcedure, error) {
	return s.procRepo.GetAllProcedures(tenantID)
}

// SaveSelectedProcedures 保存选中的存储过程/函数
func (s *mdModelProcedureService) SaveSelectedProcedures(connID, tenantID string, procSchema string, procedures []adapter.ProcedureInfo, connName string) error {
	for _, procInfo := range procedures {
		existingProc, err := s.procRepo.GetProcedureByName(connID, procSchema, procInfo.Name)
		if err == nil && existingProc != nil {
			existingProc.ProcTitle = procInfo.Name
			existingProc.ProcComment = procInfo.Comment
			existingProc.Definition = procInfo.Definition
			existingProc.ReturnType = procInfo.ReturnType
			existingProc.Language = procInfo.Language
			err = s.procRepo.UpdateProcedure(existingProc)
			if err != nil {
				return err
			}
			continue
		}

		proc := &model.MdModelProcedure{
			ID:          s.snowflake.GenerateIDString(),
			TenantID:    tenantID,
			ConnID:      connID,
			ConnName:    connName,
			ProcSchema:  procSchema,
			ProcName:    procInfo.Name,
			ProcTitle:   procInfo.Name,
			ProcType:    procInfo.Type,
			ProcComment: procInfo.Comment,
			Definition:  procInfo.Definition,
			ReturnType:  procInfo.ReturnType,
			Language:    procInfo.Language,
		}
		err = s.CreateProcedure(proc)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetParamsByProcID 根据存储过程ID获取所有参数
func (s *mdModelProcedureService) GetParamsByProcID(procID string) ([]model.MdModelProcedureParam, error) {
	return s.procRepo.GetParamsByProcID(procID)
}
