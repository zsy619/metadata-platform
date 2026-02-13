package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

// ssoUserPosRepository 用户职位仓库实现
type ssoUserPosRepository struct {
	db *gorm.DB
}

// NewSsoUserPosRepository 创建用户职位仓库实例
func NewSsoUserPosRepository(db *gorm.DB) SsoUserPosRepository {
	return &ssoUserPosRepository{db: db}
}

// CreateUserPos 创建用户职位
func (r *ssoUserPosRepository) CreateUserPos(userPos *model.SsoUserPos) error {
	return r.db.Create(userPos).Error
}

// GetUserPosByID 根据ID获取用户职位
func (r *ssoUserPosRepository) GetUserPosByID(id string) (*model.SsoUserPos, error) {
	var userPos model.SsoUserPos
	result := r.db.Where("id = ?", id).First(&userPos)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userPos, nil
}

// GetUserPosByUserID 根据用户ID获取所有职位关联
func (r *ssoUserPosRepository) GetUserPosByUserID(userID string) ([]model.SsoUserPos, error) {
	var userPositions []model.SsoUserPos
	result := r.db.Where("user_id = ?", userID).Find(&userPositions)
	if result.Error != nil {
		return nil, result.Error
	}
	return userPositions, nil
}

// GetUserPosByPosID 根据职位ID获取所有关联
func (r *ssoUserPosRepository) GetUserPosByPosID(posID string) ([]model.SsoUserPos, error) {
	var userPositions []model.SsoUserPos
	result := r.db.Where("pos_id = ?", posID).Find(&userPositions)
	if result.Error != nil {
		return nil, result.Error
	}
	return userPositions, nil
}

// DeleteUserPosition 删除关联
func (r *ssoUserPosRepository) DeleteUserPos(id string) error {
	return r.db.Delete(&model.SsoUserPos{}, "id = ?", id).Error
}

// DeleteUserPositionsByUserID 根据用户ID删除其所有职位关联
func (r *ssoUserPosRepository) DeleteUserPosByUserID(userID string) error {
	return r.db.Delete(&model.SsoUserPos{}, "user_id = ?", userID).Error
}

// DeleteUserPositionsByPosID 根据职位ID删除其所有关联
func (r *ssoUserPosRepository) DeleteUserPosByPosID(posID string) error {
	return r.db.Delete(&model.SsoUserPos{}, "pos_id = ?", posID).Error
}
