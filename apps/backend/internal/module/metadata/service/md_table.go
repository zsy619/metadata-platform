package service

import (
	"errors"

	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/repository"
	"metadata-platform/internal/utils"
)

// MdTableService 数据连接表服务接口
type MdTableService interface {
	CreateTable(table *model.MdTable) error
	GetTableByID(id string) (*model.MdTable, error)
	GetTableByName(connID, tableName string) (*model.MdTable, error)
	UpdateTable(table *model.MdTable) error
	DeleteTable(id string) error
	GetTablesByConnID(connID string) ([]model.MdTable, error)
	GetAllTables(tenantID string) ([]model.MdTable, error)
}

// mdTableService 数据连接表服务实现
type mdTableService struct {
	tableRepo  repository.MdTableRepository
	fieldRepo  repository.MdTableFieldRepository
	snowflake  *utils.Snowflake
}

// NewMdTableService 创建数据连接表服务实例
func NewMdTableService(tableRepo repository.MdTableRepository, fieldRepo repository.MdTableFieldRepository) MdTableService {
	// 创建雪花算法生成器实例，使用默认数据中心ID和机器ID
	snowflake := utils.NewSnowflake(1, 1)
	return &mdTableService{
		tableRepo:  tableRepo,
		fieldRepo:  fieldRepo,
		snowflake: snowflake,
	}
}

// CreateTable 创建数据连接表
func (s *mdTableService) CreateTable(table *model.MdTable) error {
	// 使用雪花算法生成唯一ID
	table.ID = s.snowflake.GenerateIDString()

	// 检查表是否已存在
	existingTable, err := s.tableRepo.GetTableByID(table.ID)
	if err == nil && existingTable != nil {
		return errors.New("表已存在")
	}

	// 创建表
	return s.tableRepo.CreateTable(table)
}

// GetTableByID 根据ID获取数据连接表
func (s *mdTableService) GetTableByID(id string) (*model.MdTable, error) {
	return s.tableRepo.GetTableByID(id)
}

// GetTableByName 根据连接ID和表名获取数据连接表
func (s *mdTableService) GetTableByName(connID, tableName string) (*model.MdTable, error) {
	return s.tableRepo.GetTableByName(connID, tableName)
}

// UpdateTable 更新数据连接表
func (s *mdTableService) UpdateTable(table *model.MdTable) error {
	// 检查表是否存在
	existingTable, err := s.tableRepo.GetTableByID(table.ID)
	if err != nil {
		return errors.New("表不存在")
	}

	// 如果表名发生变化，检查新表名是否已存在
	if existingTable.TableNameStr != table.TableNameStr {
		otherTable, err := s.tableRepo.GetTableByName(table.ConnID, table.TableNameStr)
		if err == nil && otherTable != nil {
			return errors.New("表名已存在")
		}
	}

	// 更新表
	return s.tableRepo.UpdateTable(table)
}

// DeleteTable 删除数据连接表
func (s *mdTableService) DeleteTable(id string) error {
	// 检查表是否存在
	_, err := s.tableRepo.GetTableByID(id)
	if err != nil {
		return errors.New("表不存在")
	}

	// 先删除该表的所有字段
	err = s.fieldRepo.DeleteFieldsByTableID(id)
	if err != nil {
		return err
	}

	// 删除表
	return s.tableRepo.DeleteTable(id)
}

// GetTablesByConnID 根据连接ID获取所有表
func (s *mdTableService) GetTablesByConnID(connID string) ([]model.MdTable, error) {
	return s.tableRepo.GetTablesByConnID(connID)
}

// GetAllTables 获取所有数据连接表
func (s *mdTableService) GetAllTables(tenantID string) ([]model.MdTable, error) {
	return s.tableRepo.GetAllTables(tenantID)
}
