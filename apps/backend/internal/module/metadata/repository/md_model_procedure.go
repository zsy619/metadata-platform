package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/metadata/model"
)

// MdModelProcedureRepository 模型存储过程/函数仓库接口
type MdModelProcedureRepository interface {
	CreateProcedure(proc *model.MdModelProcedure) error
	GetProcedureByID(id string) (*model.MdModelProcedure, error)
	GetProcedureByName(connID, procSchema, procName string) (*model.MdModelProcedure, error)
	UpdateProcedure(proc *model.MdModelProcedure) error
	DeleteProcedure(id string) error
	GetProceduresByConnID(connID string) ([]model.MdModelProcedure, error)
	GetAllProcedures(tenantID string) ([]model.MdModelProcedure, error)
	CreateProcedureParam(param *model.MdModelProcedureParam) error
	GetParamsByProcID(procID string) ([]model.MdModelProcedureParam, error)
	DeleteParamsByProcID(procID string) error
}

// mdProcedureRepository 存储过程/函数仓库实现
type mdProcedureRepository struct {
	db *gorm.DB
}

// NewMdModelProcedureRepository 创建模型存储过程/函数仓库实例
func NewMdModelProcedureRepository(db *gorm.DB) MdModelProcedureRepository {
	return &mdProcedureRepository{db: db}
}

// CreateProcedure 创建存储过程/函数
func (r *mdProcedureRepository) CreateProcedure(proc *model.MdModelProcedure) error {
	return r.db.Create(proc).Error
}

// GetProcedureByID 根据ID获取存储过程/函数
func (r *mdProcedureRepository) GetProcedureByID(id string) (*model.MdModelProcedure, error) {
	var proc model.MdModelProcedure
	result := r.db.Where("id = ?", id).First(&proc)
	if result.Error != nil {
		return nil, result.Error
	}
	return &proc, nil
}

// GetProcedureByName 根据连接ID、模式和名称获取存储过程/函数
func (r *mdProcedureRepository) GetProcedureByName(connID, procSchema, procName string) (*model.MdModelProcedure, error) {
	var proc model.MdModelProcedure
	result := r.db.Where("conn_id = ? AND proc_schema = ? AND proc_name = ?", connID, procSchema, procName).First(&proc)
	if result.Error != nil {
		return nil, result.Error
	}
	return &proc, nil
}

// UpdateProcedure 更新存储过程/函数
func (r *mdProcedureRepository) UpdateProcedure(proc *model.MdModelProcedure) error {
	return r.db.Save(proc).Error
}

// DeleteProcedure 删除存储过程/函数
func (r *mdProcedureRepository) DeleteProcedure(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.MdModelProcedure{}).Error
}

// GetProceduresByConnID 根据连接ID获取所有存储过程/函数
func (r *mdProcedureRepository) GetProceduresByConnID(connID string) ([]model.MdModelProcedure, error) {
	var procs []model.MdModelProcedure
	result := r.db.Where("conn_id = ?", connID).Order("sort, proc_name").Find(&procs)
	if result.Error != nil {
		return nil, result.Error
	}
	return procs, nil
}

// GetAllProcedures 获取所有存储过程/函数
func (r *mdProcedureRepository) GetAllProcedures(tenantID string) ([]model.MdModelProcedure, error) {
	var procs []model.MdModelProcedure
	result := r.db.Find(&procs)
	if result.Error != nil {
		return nil, result.Error
	}
	return procs, nil
}

// CreateProcedureParam 创建存储过程/函数参数
func (r *mdProcedureRepository) CreateProcedureParam(param *model.MdModelProcedureParam) error {
	return r.db.Create(param).Error
}

// GetParamsByProcID 根据存储过程ID获取所有参数
func (r *mdProcedureRepository) GetParamsByProcID(procID string) ([]model.MdModelProcedureParam, error) {
	var params []model.MdModelProcedureParam
	result := r.db.Where("proc_id = ?", procID).Order("sort").Find(&params)
	if result.Error != nil {
		return nil, result.Error
	}
	return params, nil
}

// DeleteParamsByProcID 根据存储过程ID删除所有参数
func (r *mdProcedureRepository) DeleteParamsByProcID(procID string) error {
	return r.db.Where("proc_id = ?", procID).Delete(&model.MdModelProcedureParam{}).Error
}
