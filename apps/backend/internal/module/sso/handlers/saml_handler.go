package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/module/sso/service"
)

// SAMLProtocolHandler SAML 2.0 协议处理器
type SAMLProtocolHandler struct {
	sessionService service.SsoSessionService
}

// NewSAMLProtocolHandler 创建 SAML 协议处理器
func NewSAMLProtocolHandler(sessionService service.SsoSessionService) *SAMLProtocolHandler {
	return &SAMLProtocolHandler{
		sessionService: sessionService,
	}
}

// GetProtocolType 获取协议类型
func (h *SAMLProtocolHandler) GetProtocolType() model.ProtocolType {
	return model.ProtocolTypeSAML
}

// HandleAuthorization 处理 SAML SSO 请求
func (h *SAMLProtocolHandler) HandleAuthorization(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	params map[string]any,
) (any, error) {
	samlConfig, err := parseSAMLConfig(config)
	if err != nil {
		return nil, err
	}

	requestID := fmt.Sprintf("id-%d", time.Now().UnixNano())
	relayState, _ := params["RelayState"].(string)

	extraData, _ := marshalExtraData(map[string]any{
		"relay_state": relayState,
		"request_id":  requestID,
	})

	session := &model.SsoSession{
		SessionID:        requestID,
		ClientID:         client.ID,
		ProtocolConfigID: config.ID,
		Status:           model.SessionStatusPending,
		ExpiresAt:        time.Now().Add(5 * time.Minute),
		ExtraData:        extraData,
	}

	if err := h.sessionService.CreateSession(session); err != nil {
		return nil, fmt.Errorf("创建会话失败：%w", err)
	}

	return map[string]any{
		"redirect_url": samlConfig.SingleSignOnURL,
		"request_id":   requestID,
		"relay_state":  relayState,
	}, nil
}

// HandleToken SAML 不支持 Token 端点
func (h *SAMLProtocolHandler) HandleToken(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	params map[string]any,
) (any, error) {
	return nil, errors.New("SAML 协议不支持 Token 端点")
}

// HandleUserInfo SAML 不支持 UserInfo 端点
func (h *SAMLProtocolHandler) HandleUserInfo(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	token string,
) (map[string]any, error) {
	return nil, errors.New("SAML 协议不支持 UserInfo 端点")
}

// HandleLogout 处理 SAML SLO 请求
func (h *SAMLProtocolHandler) HandleLogout(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	session *model.SsoSession,
) error {
	if session == nil {
		return errors.New("会话不能为空")
	}

	session.Status = model.SessionStatusLoggedOut
	session.ExpiresAt = time.Now()

	return h.sessionService.UpdateSession(session)
}

// HandleACS 处理 Assertion Consumer Service
func (h *SAMLProtocolHandler) HandleACS(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	samlResponse string,
) (map[string]any, error) {
	// TODO: 实现 SAML Response 解码和验证
	userInfo := map[string]any{
		"sub":   "saml_user",
		"email": "user@example.com",
		"name":  "SAML User",
	}

	sessionID := fmt.Sprintf("sid-%d", time.Now().UnixNano())
	extraData, _ := marshalExtraData(userInfo)

	session := &model.SsoSession{
		SessionID:        sessionID,
		UserID:           userInfo["sub"].(string),
		ClientID:         client.ID,
		ProtocolConfigID: config.ID,
		Status:           model.SessionStatusActive,
		ExpiresAt:        time.Now().Add(2 * time.Hour),
		ExtraData:        extraData,
	}

	if err := h.sessionService.CreateSession(session); err != nil {
		return nil, fmt.Errorf("创建会话失败：%w", err)
	}

	return map[string]any{
		"session_id": sessionID,
		"user_info":  userInfo,
	}, nil
}

func parseSAMLConfig(config *model.SsoProtocolConfig) (*SAMLConfig, error) {
	if config.ConfigData == "" {
		return nil, errors.New("配置数据为空")
	}

	var samlConfig SAMLConfig
	if err := json.Unmarshal([]byte(config.ConfigData), &samlConfig); err != nil {
		return nil, fmt.Errorf("解析配置数据失败：%w", err)
	}

	return &samlConfig, nil
}

// SAMLConfig SAML 配置结构
type SAMLConfig struct {
	EntityID        string `json:"entity_id"`
	SingleSignOnURL string `json:"single_sign_on_url"`
	SingleLogoutURL string `json:"single_logout_url"`
	X509Certificate string `json:"x509_certificate"`
	ACSURL          string `json:"acs_url"`
}
