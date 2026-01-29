package service

import (
	"errors"

	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/repository"
	"metadata-platform/internal/utils"
)

// MdTableFieldService 数据连接表字段服务接口
type MdTableFieldService interface {
	CreateField(field *model.MdTableField) error
	GetFieldByID(id string) (*model.MdTableField, error)
	GetFieldsByTableID(tableID string) ([]model.MdTableField, error)
	UpdateField(field *model.MdTableField) error
	DeleteField(id string) error
	DeleteFieldsByTableID(tableID string) error
	GetAllFields(connID string, tableID string) ([]model.MdTableField, error)
}

// mdTableFieldService 数据连接表字段服务实现
type mdTableFieldService struct {
	fieldRepo  repository.MdTableFieldRepository
	snowflake  *utils.Snowflake
}

// NewMdTableFieldService 创建数据连接表字段服务实例
func NewMdTableFieldService(fieldRepo repository.MdTableFieldRepository) MdTableFieldService {
	// 创建雪花算法生成器实例，使用默认数据中心ID和机器ID
	snowflake := utils.NewSnowflake(1, 1)
	return &mdTableFieldService{
		fieldRepo:  fieldRepo,
		snowflake: snowflake,
	}
}

// CreateField 创建数据连接表字段
func (s *mdTableFieldService) CreateField(field *model.MdTableField) error {
	// 使用雪花算法生成唯一ID
	field.ID = s.snowflake.GenerateIDString()

	// 兜底逻辑：如果排序字段为0，自动计算下一个排序值
	if field.Sort == 0 {
		fields, err := s.fieldRepo.GetFieldsByTableID(field.TableID)
		if err == nil {
			field.Sort = len(fields) + 1
		} else {
			field.Sort = 1
		}
	}

	// 创建字段
	return s.fieldRepo.CreateField(field)
}

// GetFieldByID 根据ID获取数据连接表字段
func (s *mdTableFieldService) GetFieldByID(id string) (*model.MdTableField, error) {
	return s.fieldRepo.GetFieldByID(id)
}

// GetFieldsByTableID 根据表ID获取所有字段
func (s *mdTableFieldService) GetFieldsByTableID(tableID string) ([]model.MdTableField, error) {
	return s.fieldRepo.GetFieldsByTableID(tableID)
}

// UpdateField 更新数据连接表字段
func (s *mdTableFieldService) UpdateField(field *model.MdTableField) error {
	// 检查字段是否存在
	existingField, err := s.fieldRepo.GetFieldByID(field.ID)
	if err != nil {
		return errors.New("字段不存在")
	}

	// 更新字段
	existingField.ColumnName = field.ColumnName
	existingField.ColumnTitle = field.ColumnTitle
	existingField.ColumnType = field.ColumnType
	existingField.ColumnLength = field.ColumnLength
	existingField.ColumnComment = field.ColumnComment
	existingField.IsNullable = field.IsNullable
	existingField.IsPrimaryKey = field.IsPrimaryKey
	existingField.IsAutoIncrement = field.IsAutoIncrement
	existingField.DefaultValue = field.DefaultValue
	existingField.ExtraInfo = field.ExtraInfo
	existingField.Sort = field.Sort

	return s.fieldRepo.UpdateField(existingField)
}

// DeleteField 删除数据连接表字段
func (s *mdTableFieldService) DeleteField(id string) error {
	// 检查字段是否存在
	_, err := s.fieldRepo.GetFieldByID(id)
	if err != nil {
		return errors.New("字段不存在")
	}

	// 删除字段
	return s.fieldRepo.DeleteField(id)
}

// DeleteFieldsByTableID 根据表ID删除所有字段
func (s *mdTableFieldService) DeleteFieldsByTableID(tableID string) error {
	// 删除表下所有字段
	return s.fieldRepo.DeleteFieldsByTableID(tableID)
}

// GetAllFields 获取所有数据连接表字段
func (s *mdTableFieldService) GetAllFields(connID string, tableID string) ([]model.MdTableField, error) {
	return s.fieldRepo.GetAllFields(connID, tableID)
}
