package service

import (
	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/module/sso/repository"
)

// ssoProtocolConfigService 协议配置服务实现
type ssoProtocolConfigService struct {
	repo repository.SsoProtocolConfigRepository
}

// NewSsoProtocolConfigService 创建协议配置服务实例
func NewSsoProtocolConfigService(repo repository.SsoProtocolConfigRepository) SsoProtocolConfigService {
	return &ssoProtocolConfigService{repo: repo}
}

func (s *ssoProtocolConfigService) CreateConfig(config *model.SsoProtocolConfig) error {
	return s.repo.CreateConfig(config)
}

func (s *ssoProtocolConfigService) GetConfigByID(id string) (*model.SsoProtocolConfig, error) {
	return s.repo.GetConfigByID(id)
}

func (s *ssoProtocolConfigService) GetConfigByProtocolType(protocolType model.ProtocolType) ([]model.SsoProtocolConfig, error) {
	return s.repo.GetConfigByProtocolType(protocolType)
}

func (s *ssoProtocolConfigService) GetEnabledConfigs() ([]model.SsoProtocolConfig, error) {
	return s.repo.GetEnabledConfigs()
}

func (s *ssoProtocolConfigService) UpdateConfig(config *model.SsoProtocolConfig) error {
	return s.repo.UpdateConfig(config)
}

func (s *ssoProtocolConfigService) UpdateConfigFields(id string, fields map[string]any) error {
	return s.repo.UpdateConfigFields(id, fields)
}

func (s *ssoProtocolConfigService) DeleteConfig(id string) error {
	return s.repo.DeleteConfig(id)
}

func (s *ssoProtocolConfigService) GetAllConfigs() ([]model.SsoProtocolConfig, error) {
	return s.repo.GetAllConfigs()
}
