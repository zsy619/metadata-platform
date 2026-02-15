package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoTenantRepository 租户仓库实现
type ssoTenantRepository struct {
	db *gorm.DB
}

// NewSsoTenantRepository 创建租户仓库实例
func NewSsoTenantRepository(db *gorm.DB) SsoTenantRepository {
	return &ssoTenantRepository{db: db}
}

// CreateTenant 创建租户
func (r *ssoTenantRepository) CreateTenant(tenant *model.SsoTenant) error {
	return r.db.Create(tenant).Error
}

// GetTenantByID 根据ID获取租户
func (r *ssoTenantRepository) GetTenantByID(id string) (*model.SsoTenant, error) {
	var tenant model.SsoTenant
	result := r.db.Where("id = ?", id).First(&tenant)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tenant, nil
}

// GetTenantByCode 根据编码获取租户
func (r *ssoTenantRepository) GetTenantByCode(code string) (*model.SsoTenant, error) {
	var tenant model.SsoTenant
	result := r.db.Where("tenant_code = ?", code).First(&tenant)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tenant, nil
}

// UpdateTenant 更新租户
func (r *ssoTenantRepository) UpdateTenant(tenant *model.SsoTenant) error {
	return r.db.Save(tenant).Error
}

// DeleteTenant 删除租户
func (r *ssoTenantRepository) DeleteTenant(id string) error {
	return r.db.Model(&model.SsoTenant{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// GetAllTenants 获取所有租户
func (r *ssoTenantRepository) GetAllTenants() ([]model.SsoTenant, error) {
	var tenants []model.SsoTenant
	result := r.db.Unscoped().Where("is_deleted = ?", false).Find(&tenants)
	if result.Error != nil {
		return nil, result.Error
	}
	return tenants, nil
}

// GetMaxSort 获取最大排序值
func (r *ssoOrgRepository) GetMaxSort() (int, error) {
	var maxSort int
	result := r.db.Model(&model.SsoOrg{}).Select("COALESCE(MAX(sort), 0)").Scan(&maxSort)
	if result.Error != nil {
		return 0, result.Error
	}
	return maxSort, nil
}
