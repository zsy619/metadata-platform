package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoUserPositionRepository 用户职位仓库实现
type ssoUserPositionRepository struct {
	db *gorm.DB
}

// NewSsoUserPositionRepository 创建用户职位仓库实例
func NewSsoUserPositionRepository(db *gorm.DB) SsoUserPositionRepository {
	return &ssoUserPositionRepository{db: db}
}

// CreateUserPosition 创建用户职位
func (r *ssoUserPositionRepository) CreateUserPosition(userPos *model.SsoUserPosition) error {
	return r.db.Create(userPos).Error
}

// GetUserPositionByID 根据ID获取用户职位
func (r *ssoUserPositionRepository) GetUserPositionByID(id string) (*model.SsoUserPosition, error) {
	var userPos model.SsoUserPosition
	result := r.db.Where("id = ?", id).First(&userPos)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userPos, nil
}

// GetUserPositionsByUserID 根据用户ID获取所有职位关联
func (r *ssoUserPositionRepository) GetUserPositionsByUserID(userID string) ([]model.SsoUserPosition, error) {
	var userPositions []model.SsoUserPosition
	result := r.db.Where("user_id = ?", userID).Find(&userPositions)
	if result.Error != nil {
		return nil, result.Error
	}
	return userPositions, nil
}

// GetUserPositionsByPosID 根据职位ID获取所有关联
func (r *ssoUserPositionRepository) GetUserPositionsByPosID(posID string) ([]model.SsoUserPosition, error) {
	var userPositions []model.SsoUserPosition
	result := r.db.Where("position_id = ?", posID).Find(&userPositions)
	if result.Error != nil {
		return nil, result.Error
	}
	return userPositions, nil
}

// DeleteUserPosition 删除关联
func (r *ssoUserPositionRepository) DeleteUserPosition(id string) error {
	return r.db.Delete(&model.SsoUserPosition{}, "id = ?", id).Error
}

// DeleteUserPositionsByUserID 根据用户ID删除其所有职位关联
func (r *ssoUserPositionRepository) DeleteUserPositionsByUserID(userID string) error {
	return r.db.Delete(&model.SsoUserPosition{}, "user_id = ?", userID).Error
}

// DeleteUserPositionsByPosID 根据职位ID删除其所有关联
func (r *ssoUserPositionRepository) DeleteUserPositionsByPosID(posID string) error {
	return r.db.Delete(&model.SsoUserPosition{}, "position_id = ?", posID).Error
}
