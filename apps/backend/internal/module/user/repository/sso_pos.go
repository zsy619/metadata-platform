package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

// ssoPosRepository 职位仓库实现
type ssoPosRepository struct {
	db *gorm.DB
}

// NewSsoPosRepository 创建职位仓库实例
func NewSsoPosRepository(db *gorm.DB) SsoPosRepository {
	return &ssoPosRepository{db: db}
}

// CreatePos 创建职位
func (r *ssoPosRepository) CreatePos(pos *model.SsoPos) error {
	return r.db.Create(pos).Error
}

// GetPosByID 根据ID获取职位
func (r *ssoPosRepository) GetPosByID(id string) (*model.SsoPos, error) {
	var pos model.SsoPos
	result := r.db.Where("id = ?", id).First(&pos)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pos, nil
}

// GetPosByCode 根据编码获取职位
func (r *ssoPosRepository) GetPosByCode(code string) (*model.SsoPos, error) {
	var pos model.SsoPos
	result := r.db.Where("pos_code = ?", code).First(&pos)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pos, nil
}

// UpdatePos 更新职位
func (r *ssoPosRepository) UpdatePos(pos *model.SsoPos) error {
	return r.db.Save(pos).Error
}

// UpdatePosFields 更新职位指定字段
// 使用 map 方式只更新指定的字段，避免全量更新
func (r *ssoPosRepository) UpdatePosFields(id string, fields map[string]any) error {
	return r.db.Model(&model.SsoPos{}).Where("id = ?", id).Updates(fields).Error
}

// DeletePosition 删除职位
func (r *ssoPosRepository) DeletePos(id string) error {
	return r.db.Model(&model.SsoPos{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// GetAllPoss 获取所有职位
func (r *ssoPosRepository) GetAllPoss() ([]model.SsoPos, error) {
	var positions []model.SsoPos
	result := r.db.Find(&positions)
	if result.Error != nil {
		return nil, result.Error
	}
	return positions, nil
}

// GetMaxSort 获取最大排序值
func (r *ssoPosRepository) GetMaxSort() (int, error) {
	var maxSort int
	result := r.db.Model(&model.SsoPos{}).Select("COALESCE(MAX(sort), 0)").Scan(&maxSort)
	if result.Error != nil {
		return 0, result.Error
	}
	return maxSort, nil
}
