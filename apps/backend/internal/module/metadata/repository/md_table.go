package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/metadata/model"
)

// MdTableRepository 数据连接表仓库接口
type MdTableRepository interface {
	CreateTable(table *model.MdTable) error
	GetTableByID(id string) (*model.MdTable, error)
	GetTableByName(connID, tableName string) (*model.MdTable, error)
	UpdateTable(table *model.MdTable) error
	DeleteTable(id string) error
	GetTablesByConnID(connID string) ([]model.MdTable, error)
	GetAllTables(tenantID string) ([]model.MdTable, error)
}

// mdTableRepository 数据连接表仓库实现
type mdTableRepository struct {
	db *gorm.DB
}

// NewMdTableRepository 创建数据连接表仓库实例
func NewMdTableRepository(db *gorm.DB) MdTableRepository {
	return &mdTableRepository{db: db}
}

// CreateTable 创建数据连接表
func (r *mdTableRepository) CreateTable(table *model.MdTable) error {
	return r.db.Create(table).Error
}

// GetTableByID 根据ID获取数据连接表
func (r *mdTableRepository) GetTableByID(id string) (*model.MdTable, error) {
	var table model.MdTable
	result := r.db.Where("id = ?", id).First(&table)
	if result.Error != nil {
		return nil, result.Error
	}
	return &table, nil
}

// GetTableByName 根据连接ID和表名获取数据连接表
func (r *mdTableRepository) GetTableByName(connID, tableName string) (*model.MdTable, error) {
	var table model.MdTable
	result := r.db.Where("conn_id = ? AND table_name = ?", connID, tableName).First(&table)
	if result.Error != nil {
		return nil, result.Error
	}
	return &table, nil
}

// UpdateTable 更新数据连接表
func (r *mdTableRepository) UpdateTable(table *model.MdTable) error {
	return r.db.Save(table).Error
}

// DeleteTable 删除数据连接表
func (r *mdTableRepository) DeleteTable(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.MdTable{}).Error
}

// GetTablesByConnID 根据连接ID获取所有表
func (r *mdTableRepository) GetTablesByConnID(connID string) ([]model.MdTable, error) {
	var tables []model.MdTable
	result := r.db.Where("conn_id = ?", connID).Find(&tables)
	if result.Error != nil {
		return nil, result.Error
	}
	return tables, nil
}

// GetAllTables 获取所有数据连接表
func (r *mdTableRepository) GetAllTables(tenantID string) ([]model.MdTable, error) {
	var tables []model.MdTable
	result := r.db.Find(&tables)
	if result.Error != nil {
		return nil, result.Error
	}
	return tables, nil
}
