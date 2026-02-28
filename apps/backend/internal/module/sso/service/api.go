package service

import (
	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/module/sso/repository"

	"gorm.io/gorm"
)

// SsoProtocolConfigService 协议配置服务接口
type SsoProtocolConfigService interface {
	CreateConfig(config *model.SsoProtocolConfig) error
	GetConfigByID(id string) (*model.SsoProtocolConfig, error)
	GetConfigByProtocolType(protocolType model.ProtocolType) ([]model.SsoProtocolConfig, error)
	GetEnabledConfigs() ([]model.SsoProtocolConfig, error)
	UpdateConfig(config *model.SsoProtocolConfig) error
	UpdateConfigFields(id string, fields map[string]any) error
	DeleteConfig(id string) error
	GetAllConfigs() ([]model.SsoProtocolConfig, error)
}

// SsoClientService 客户端配置服务接口
type SsoClientService interface {
	CreateClient(client *model.SsoClient) error
	GetClientByID(id string) (*model.SsoClient, error)
	GetClientByClientID(clientID string) (*model.SsoClient, error)
	GetClientsByProtocolConfigID(protocolConfigID string) ([]model.SsoClient, error)
	UpdateClient(client *model.SsoClient) error
	UpdateClientFields(id string, fields map[string]any) error
	DeleteClient(id string) error
	GetAllClients() ([]model.SsoClient, error)
}

// SsoKeyService 密钥管理服务接口
type SsoKeyService interface {
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
	GenerateKeyPair(keyType model.KeyType, algorithm string) (*model.SsoKey, error)
}

// SsoFieldMappingService 字段映射服务接口
type SsoFieldMappingService interface {
	CreateMapping(mapping *model.SsoFieldMapping) error
	GetMappingByID(id string) (*model.SsoFieldMapping, error)
	GetMappingsByProtocolConfigID(protocolConfigID string) ([]model.SsoFieldMapping, error)
	GetMappingsByClientID(clientID string) ([]model.SsoFieldMapping, error)
	UpdateMapping(mapping *model.SsoFieldMapping) error
	UpdateMappingFields(id string, fields map[string]any) error
	DeleteMapping(id string) error
	GetAllMappings() ([]model.SsoFieldMapping, error)
	MapUserFields(sourceData map[string]any, mappings []model.SsoFieldMapping) (map[string]any, error)
}

// SsoSessionService 会话服务接口
type SsoSessionService interface {
	CreateSession(session *model.SsoSession) error
	GetSessionByID(id string) (*model.SsoSession, error)
	GetSessionBySessionID(sessionID string) (*model.SsoSession, error)
	GetSessionsByUserID(userID string) ([]model.SsoSession, error)
	GetActiveSessionsByUserID(userID string) ([]model.SsoSession, error)
	UpdateSession(session *model.SsoSession) error
	UpdateSessionFields(id string, fields map[string]any) error
	UpdateSessionStatus(sessionID string, status model.SessionStatus) error
	DeleteSession(id string) error
	DeleteSessionsByUserID(userID string) error
	RevokeSession(sessionID string) error
	GetAllSessions() ([]model.SsoSession, error)
}

// ProtocolHandler 协议处理器接口
type ProtocolHandler interface {
	GetProtocolType() model.ProtocolType
	HandleAuthorization(config *model.SsoProtocolConfig, client *model.SsoClient, params map[string]any) (any, error)
	HandleToken(config *model.SsoProtocolConfig, client *model.SsoClient, params map[string]any) (any, error)
	HandleUserInfo(config *model.SsoProtocolConfig, client *model.SsoClient, token string) (map[string]any, error)
	HandleLogout(config *model.SsoProtocolConfig, client *model.SsoClient, session *model.SsoSession) error
}

// Services SSO模块服务集合
type Services struct {
	ProtocolConfig SsoProtocolConfigService
	Client         SsoClientService
	Key            SsoKeyService
	FieldMapping   SsoFieldMappingService
	Session        SsoSessionService
	ProtocolHandlers map[model.ProtocolType]ProtocolHandler
}

// NewServices 创建SSO模块服务集合
func NewServices(repos *repository.Repositories, db *gorm.DB) *Services {
	return &Services{
		ProtocolConfig: NewSsoProtocolConfigService(repos.ProtocolConfig),
		Client:         NewSsoClientService(repos.Client),
		Key:            NewSsoKeyService(repos.Key),
		FieldMapping:   NewSsoFieldMappingService(repos.FieldMapping),
		Session:        NewSsoSessionService(repos.Session),
		ProtocolHandlers: make(map[model.ProtocolType]ProtocolHandler),
	}
}
