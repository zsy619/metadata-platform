package repository

import (
	"metadata-platform/internal/module/sso/model"

	"gorm.io/gorm"
)

type ssoKeyRepository struct {
	db *gorm.DB
}

func NewSsoKeyRepository(db *gorm.DB) SsoKeyRepository {
	return &ssoKeyRepository{db: db}
}

func (r *ssoKeyRepository) CreateKey(key *model.SsoKey) error {
	return r.db.Create(key).Error
}

func (r *ssoKeyRepository) GetKeyByID(id string) (*model.SsoKey, error) {
	var key model.SsoKey
	result := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&key)
	if result.Error != nil {
		return nil, result.Error
	}
	return &key, nil
}

func (r *ssoKeyRepository) GetKeyByKeyID(keyID string) (*model.SsoKey, error) {
	var key model.SsoKey
	result := r.db.Where("key_id = ? AND is_deleted = ?", keyID, false).First(&key)
	if result.Error != nil {
		return nil, result.Error
	}
	return &key, nil
}

func (r *ssoKeyRepository) GetKeysByProtocolConfigID(protocolConfigID string) ([]model.SsoKey, error) {
	var keys []model.SsoKey
	result := r.db.Where("protocol_config_id = ? AND is_deleted = ?", protocolConfigID, false).
		Order("is_primary DESC, create_at DESC").Find(&keys)
	if result.Error != nil {
		return nil, result.Error
	}
	return keys, nil
}

func (r *ssoKeyRepository) GetPrimaryKey(protocolConfigID string) (*model.SsoKey, error) {
	var key model.SsoKey
	result := r.db.Where("protocol_config_id = ? AND is_primary = ? AND is_enabled = ? AND is_deleted = ?",
		protocolConfigID, true, true, false).First(&key)
	if result.Error != nil {
		return nil, result.Error
	}
	return &key, nil
}

func (r *ssoKeyRepository) GetValidKeys(protocolConfigID string) ([]model.SsoKey, error) {
	var keys []model.SsoKey
	result := r.db.Where("protocol_config_id = ? AND is_enabled = ? AND is_deleted = ?",
		protocolConfigID, true, false).
		Order("is_primary DESC, create_at DESC").Find(&keys)
	if result.Error != nil {
		return nil, result.Error
	}
	var validKeys []model.SsoKey
	for _, key := range keys {
		if key.IsValid() {
			validKeys = append(validKeys, key)
		}
	}
	return validKeys, nil
}

func (r *ssoKeyRepository) UpdateKey(key *model.SsoKey) error {
	return r.db.Save(key).Error
}

func (r *ssoKeyRepository) UpdateKeyFields(id string, fields map[string]any) error {
	return r.db.Model(&model.SsoKey{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ssoKeyRepository) DeleteKey(id string) error {
	return r.db.Model(&model.SsoKey{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoKeyRepository) GetAllKeys() ([]model.SsoKey, error) {
	var keys []model.SsoKey
	result := r.db.Where("is_deleted = ?", false).
		Order("create_at DESC").Find(&keys)
	if result.Error != nil {
		return nil, result.Error
	}
	return keys, nil
}
