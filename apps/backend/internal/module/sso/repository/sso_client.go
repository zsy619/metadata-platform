package repository

import (
	"metadata-platform/internal/module/sso/model"

	"gorm.io/gorm"
)

type ssoClientRepository struct {
	db *gorm.DB
}

func NewSsoClientRepository(db *gorm.DB) SsoClientRepository {
	return &ssoClientRepository{db: db}
}

func (r *ssoClientRepository) CreateClient(client *model.SsoClient) error {
	return r.db.Create(client).Error
}

func (r *ssoClientRepository) GetClientByID(id string) (*model.SsoClient, error) {
	var client model.SsoClient
	result := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&client)
	if result.Error != nil {
		return nil, result.Error
	}
	return &client, nil
}

func (r *ssoClientRepository) GetClientByClientID(clientID string) (*model.SsoClient, error) {
	var client model.SsoClient
	result := r.db.Where("client_id = ? AND is_deleted = ?", clientID, false).First(&client)
	if result.Error != nil {
		return nil, result.Error
	}
	return &client, nil
}

func (r *ssoClientRepository) GetClientsByProtocolConfigID(protocolConfigID string) ([]model.SsoClient, error) {
	var clients []model.SsoClient
	result := r.db.Where("protocol_config_id = ? AND is_deleted = ?", protocolConfigID, false).
		Order("sort ASC, create_at DESC").Find(&clients)
	if result.Error != nil {
		return nil, result.Error
	}
	return clients, nil
}

func (r *ssoClientRepository) UpdateClient(client *model.SsoClient) error {
	return r.db.Save(client).Error
}

func (r *ssoClientRepository) UpdateClientFields(id string, fields map[string]any) error {
	return r.db.Model(&model.SsoClient{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ssoClientRepository) DeleteClient(id string) error {
	return r.db.Model(&model.SsoClient{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoClientRepository) GetAllClients() ([]model.SsoClient, error) {
	var clients []model.SsoClient
	result := r.db.Where("is_deleted = ?", false).
		Order("sort ASC, create_at DESC").Find(&clients)
	if result.Error != nil {
		return nil, result.Error
	}
	return clients, nil
}
