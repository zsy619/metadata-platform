package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoPositionRepository 职位仓库实现
type ssoPositionRepository struct {
	db *gorm.DB
}

// NewSsoPositionRepository 创建职位仓库实例
func NewSsoPositionRepository(db *gorm.DB) SsoPositionRepository {
	return &ssoPositionRepository{db: db}
}

// CreatePosition 创建职位
func (r *ssoPositionRepository) CreatePosition(pos *model.SsoPosition) error {
	return r.db.Create(pos).Error
}

// GetPositionByID 根据ID获取职位
func (r *ssoPositionRepository) GetPositionByID(id string) (*model.SsoPosition, error) {
	var pos model.SsoPosition
	result := r.db.Where("id = ?", id).First(&pos)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pos, nil
}

// GetPositionByCode 根据编码获取职位
func (r *ssoPositionRepository) GetPositionByCode(code string) (*model.SsoPosition, error) {
	var pos model.SsoPosition
	result := r.db.Where("pos_code = ?", code).First(&pos)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pos, nil
}

// UpdatePosition 更新职位
func (r *ssoPositionRepository) UpdatePosition(pos *model.SsoPosition) error {
	return r.db.Save(pos).Error
}

// DeletePosition 删除职位
func (r *ssoPositionRepository) DeletePosition(id string) error {
	return r.db.Model(&model.SsoPosition{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// GetAllPositions 获取所有职位
func (r *ssoPositionRepository) GetAllPositions() ([]model.SsoPosition, error) {
	var positions []model.SsoPosition
	result := r.db.Find(&positions)
	if result.Error != nil {
		return nil, result.Error
	}
	return positions, nil
}
