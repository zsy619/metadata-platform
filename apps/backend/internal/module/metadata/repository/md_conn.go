package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/metadata/model"
)

// MdConnRepository 数据连接仓库接口
type MdConnRepository interface {
	CreateConn(conn *model.MdConn) error
	GetConnByID(id string) (*model.MdConn, error)
	GetConnByName(name string) (*model.MdConn, error)
	UpdateConn(conn *model.MdConn) error
	DeleteConn(id string) error
	GetAllConns(tenantID string) ([]model.MdConn, error)
	GetConnsByParentID(parentID string) ([]model.MdConn, error)
	GetMDConnByID(id string) (*model.MdConn, error)
}

// mdConnRepository 数据连接仓库实现
type mdConnRepository struct {
	db *gorm.DB
}

// NewMdConnRepository 创建数据连接仓库实例
func NewMdConnRepository(db *gorm.DB) MdConnRepository {
	return &mdConnRepository{db: db}
}

// CreateConn 创建数据连接
func (r *mdConnRepository) CreateConn(conn *model.MdConn) error {
	return r.db.Create(conn).Error
}

// GetConnByID 根据ID获取数据连接
func (r *mdConnRepository) GetConnByID(id string) (*model.MdConn, error) {
	var conn model.MdConn
	result := r.db.Where("id = ?", id).First(&conn)
	if result.Error != nil {
		return nil, result.Error
	}
	return &conn, nil
}

// GetConnByName 根据名称获取数据连接
func (r *mdConnRepository) GetConnByName(name string) (*model.MdConn, error) {
	var conn model.MdConn
	result := r.db.Where("conn_name = ?", name).First(&conn)
	if result.Error != nil {
		return nil, result.Error
	}
	return &conn, nil
}

// UpdateConn 更新数据连接
func (r *mdConnRepository) UpdateConn(conn *model.MdConn) error {
	return r.db.Save(conn).Error
}

// DeleteConn 删除数据连接
func (r *mdConnRepository) DeleteConn(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.MdConn{}).Error
}

// GetAllConns 获取所有数据连接
func (r *mdConnRepository) GetAllConns(tenantID string) ([]model.MdConn, error) {
	var conns []model.MdConn
	result := r.db.Find(&conns)
	if result.Error != nil {
		return nil, result.Error
	}
	return conns, nil
}

// GetConnsByParentID 根据父ID获取数据连接
func (r *mdConnRepository) GetConnsByParentID(parentID string) ([]model.MdConn, error) {
	var conns []model.MdConn
	result := r.db.Where("parent_id = ?", parentID).Find(&conns)
	if result.Error != nil {
		return nil, result.Error
	}
	return conns, nil
}

// GetMDConnByID 根据ID获取数据连接 (Engine专用名)
func (r *mdConnRepository) GetMDConnByID(id string) (*model.MdConn, error) {
	return r.GetConnByID(id)
}
