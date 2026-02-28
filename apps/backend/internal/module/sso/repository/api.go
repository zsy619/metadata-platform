package repository

import (
	"metadata-platform/internal/module/sso/model"

	"gorm.io/gorm"
)

// SsoProtocolConfigRepository 协议配置仓库接口
type SsoProtocolConfigRepository interface {
	CreateConfig(config *model.SsoProtocolConfig) error
	GetConfigByID(id string) (*model.SsoProtocolConfig, error)
	GetConfigByProtocolType(protocolType model.ProtocolType) ([]model.SsoProtocolConfig, error)
	GetEnabledConfigs() ([]model.SsoProtocolConfig, error)
	UpdateConfig(config *model.SsoProtocolConfig) error
	UpdateConfigFields(id string, fields map[string]any) error
	DeleteConfig(id string) error
	GetAllConfigs() ([]model.SsoProtocolConfig, error)
}

// SsoClientRepository 客户端配置仓库接口
type SsoClientRepository interface {
	CreateClient(client *model.SsoClient) error
	GetClientByID(id string) (*model.SsoClient, error)
	GetClientByClientID(clientID string) (*model.SsoClient, error)
	GetClientsByProtocolConfigID(protocolConfigID string) ([]model.SsoClient, error)
	UpdateClient(client *model.SsoClient) error
	UpdateClientFields(id string, fields map[string]any) error
	DeleteClient(id string) error
	GetAllClients() ([]model.SsoClient, error)
}

// SsoKeyRepository 密钥管理仓库接口
type SsoKeyRepository interface {
	CreateKey(key *model.SsoKey) error
	GetKeyByID(id string) (*model.SsoKey, error)
	GetKeyByKeyID(keyID string) (*model.SsoKey, error)
	GetKeysByProtocolConfigID(protocolConfigID string) ([]model.SsoKey, error)
	GetPrimaryKey(protocolConfigID string) (*model.SsoKey, error)
	GetValidKeys(protocolConfigID string) ([]model.SsoKey, error)
	UpdateKey(key *model.SsoKey) error
	UpdateKeyFields(id string, fields map[string]any) error
	DeleteKey(id string) error
	GetAllKeys() ([]model.SsoKey, error)
}

// SsoFieldMappingRepository 字段映射仓库接口
type SsoFieldMappingRepository interface {
	CreateMapping(mapping *model.SsoFieldMapping) error
	GetMappingByID(id string) (*model.SsoFieldMapping, error)
	GetMappingsByProtocolConfigID(protocolConfigID string) ([]model.SsoFieldMapping, error)
	GetMappingsByClientID(clientID string) ([]model.SsoFieldMapping, error)
	GetMappingBySourceAndTarget(protocolConfigID, sourceField, targetField string) (*model.SsoFieldMapping, error)
	UpdateMapping(mapping *model.SsoFieldMapping) error
	UpdateMappingFields(id string, fields map[string]any) error
	DeleteMapping(id string) error
	GetAllMappings() ([]model.SsoFieldMapping, error)
}

// SsoSessionRepository 会话仓库接口
type SsoSessionRepository interface {
	CreateSession(session *model.SsoSession) error
	GetSessionByID(id string) (*model.SsoSession, error)
	GetSessionBySessionID(sessionID string) (*model.SsoSession, error)
	GetSessionsByUserID(userID string) ([]model.SsoSession, error)
	GetSessionsByClientID(clientID string) ([]model.SsoSession, error)
	GetActiveSessionsByUserID(userID string) ([]model.SsoSession, error)
	UpdateSession(session *model.SsoSession) error
	UpdateSessionFields(id string, fields map[string]any) error
	UpdateSessionStatus(sessionID string, status model.SessionStatus) error
	DeleteSession(id string) error
	DeleteSessionsByUserID(userID string) error
	RevokeSession(sessionID string) error
	GetAllSessions() ([]model.SsoSession, error)
}

// Repositories SSO模块仓库集合
type Repositories struct {
	ProtocolConfig SsoProtocolConfigRepository
	Client         SsoClientRepository
	Key            SsoKeyRepository
	FieldMapping   SsoFieldMappingRepository
	Session        SsoSessionRepository
}

// NewRepositories 创建SSO模块仓库集合
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		ProtocolConfig: NewSsoProtocolConfigRepository(db),
		Client:         NewSsoClientRepository(db),
		Key:            NewSsoKeyRepository(db),
		FieldMapping:   NewSsoFieldMappingRepository(db),
		Session:        NewSsoSessionRepository(db),
	}
}
