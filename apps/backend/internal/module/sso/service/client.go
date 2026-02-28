package service

import (
	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/module/sso/repository"
)

// ssoClientService 客户端配置服务实现
type ssoClientService struct {
	repo repository.SsoClientRepository
}

// NewSsoClientService 创建客户端配置服务实例
func NewSsoClientService(repo repository.SsoClientRepository) SsoClientService {
	return &ssoClientService{repo: repo}
}

func (s *ssoClientService) CreateClient(client *model.SsoClient) error {
	return s.repo.CreateClient(client)
}

func (s *ssoClientService) GetClientByID(id string) (*model.SsoClient, error) {
	return s.repo.GetClientByID(id)
}

func (s *ssoClientService) GetClientByClientID(clientID string) (*model.SsoClient, error) {
	return s.repo.GetClientByClientID(clientID)
}

func (s *ssoClientService) GetClientsByProtocolConfigID(protocolConfigID string) ([]model.SsoClient, error) {
	return s.repo.GetClientsByProtocolConfigID(protocolConfigID)
}

func (s *ssoClientService) UpdateClient(client *model.SsoClient) error {
	return s.repo.UpdateClient(client)
}

func (s *ssoClientService) UpdateClientFields(id string, fields map[string]any) error {
	return s.repo.UpdateClientFields(id, fields)
}

func (s *ssoClientService) DeleteClient(id string) error {
	return s.repo.DeleteClient(id)
}

func (s *ssoClientService) GetAllClients() ([]model.SsoClient, error) {
	return s.repo.GetAllClients()
}
