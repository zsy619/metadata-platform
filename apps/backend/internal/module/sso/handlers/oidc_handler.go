package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/module/sso/service"
)

// OIDCProtocolHandler OIDC 协议处理器
type OIDCProtocolHandler struct {
	sessionService service.SsoSessionService
}

// NewOIDCProtocolHandler 创建 OIDC 协议处理器
func NewOIDCProtocolHandler(sessionService service.SsoSessionService) *OIDCProtocolHandler {
	return &OIDCProtocolHandler{
		sessionService: sessionService,
	}
}

// GetProtocolType 获取协议类型
func (h *OIDCProtocolHandler) GetProtocolType() model.ProtocolType {
	return model.ProtocolTypeOIDC
}

// HandleAuthorization 处理授权请求
func (h *OIDCProtocolHandler) HandleAuthorization(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	params map[string]any,
) (any, error) {
	responseType, _ := params["response_type"].(string)
	redirectURI, _ := params["redirect_uri"].(string)
	scope, _ := params["scope"].(string)
	state, _ := params["state"].(string)
	nonce, _ := params["nonce"].(string)

	if responseType != "code" {
		return nil, errors.New("不支持的 response_type，仅支持 code")
	}

	if redirectURI == "" {
		return nil, errors.New("redirect_uri 不能为空")
	}

	authCode, err := h.generateAuthCode()
	if err != nil {
		return nil, fmt.Errorf("生成授权码失败：%w", err)
	}

	extraData, _ := marshalExtraData(map[string]any{
		"nonce":        nonce,
		"scope":        scope,
		"redirect_uri": redirectURI,
		"state":        state,
	})

	session := &model.SsoSession{
		SessionID:        authCode,
		UserID:           "pending",
		ClientID:         client.ID,
		ProtocolConfigID: config.ID,
		Status:           model.SessionStatusPending,
		ExpiresAt:        time.Now().Add(10 * time.Minute),
		ExtraData:        extraData,
	}

	if err := h.sessionService.CreateSession(session); err != nil {
		return nil, fmt.Errorf("创建会话失败：%w", err)
	}

	redirectURL, err := url.Parse(redirectURI)
	if err != nil {
		return nil, fmt.Errorf("解析 redirect_uri 失败：%w", err)
	}

	query := redirectURL.Query()
	query.Set("code", authCode)
	if state != "" {
		query.Set("state", state)
	}
	redirectURL.RawQuery = query.Encode()

	return map[string]any{
		"redirect_uri": redirectURL.String(),
		"code":         authCode,
		"state":        state,
	}, nil
}

// HandleToken 处理令牌请求
func (h *OIDCProtocolHandler) HandleToken(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	params map[string]any,
) (any, error) {
	grantType, _ := params["grant_type"].(string)
	code, _ := params["code"].(string)
	redirectURI, _ := params["redirect_uri"].(string)

	switch grantType {
	case "authorization_code":
		return h.handleAuthorizationCodeGrant(config, client, code, redirectURI)
	case "refresh_token":
		return h.handleRefreshTokenGrant(config, client, params)
	case "client_credentials":
		return h.handleClientCredentialsGrant(config, client)
	default:
		return nil, fmt.Errorf("不支持的 grant_type: %s", grantType)
	}
}

func (h *OIDCProtocolHandler) handleAuthorizationCodeGrant(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	code string,
	redirectURI string,
) (any, error) {
	if code == "" {
		return nil, errors.New("授权码不能为空")
	}

	session, err := h.sessionService.GetSessionBySessionID(code)
	if err != nil {
		return nil, errors.New("无效的授权码")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("授权码已过期")
	}

	if session.Status != model.SessionStatusPending {
		return nil, errors.New("授权码已被使用")
	}

	storedRedirectURI := getStringFromExtraData(session.ExtraData, "redirect_uri")
	if storedRedirectURI != "" && storedRedirectURI != redirectURI {
		return nil, errors.New("redirect_uri 不匹配")
	}

	accessToken, err := h.generateAccessToken()
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败：%w", err)
	}

	refreshToken, err := h.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败：%w", err)
	}

	idToken, err := h.generateIDToken(config, client, session)
	if err != nil {
		return nil, fmt.Errorf("生成 ID Token 失败：%w", err)
	}

	session.SessionID = accessToken
	session.Status = model.SessionStatusActive
	
	newExtraData, _ := marshalExtraData(map[string]any{
		"refresh_token": refreshToken,
		"id_token":      idToken,
		"client_id":     client.ClientID,
	})
	session.ExtraData = newExtraData
	
	session.ExpiresAt = time.Now().Add(2 * time.Hour)

	if err := h.sessionService.UpdateSession(session); err != nil {
		return nil, fmt.Errorf("更新会话失败：%w", err)
	}

	scope := "openid"
	storedScope := getStringFromExtraData(session.ExtraData, "scope")
	if storedScope != "" {
		scope = storedScope
	}

	return map[string]any{
		"access_token":  accessToken,
		"token_type":    "Bearer",
		"expires_in":    7200,
		"refresh_token": refreshToken,
		"id_token":      idToken,
		"scope":         scope,
	}, nil
}

func (h *OIDCProtocolHandler) handleRefreshTokenGrant(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	params map[string]any,
) (any, error) {
	refreshToken, _ := params["refresh_token"].(string)
	if refreshToken == "" {
		return nil, errors.New("刷新令牌不能为空")
	}

	accessToken, err := h.generateAccessToken()
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败：%w", err)
	}

	return map[string]any{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   7200,
	}, nil
}

func (h *OIDCProtocolHandler) handleClientCredentialsGrant(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
) (any, error) {
	accessToken, err := h.generateAccessToken()
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败：%w", err)
	}

	return map[string]any{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   3600,
		"scope":        "client",
	}, nil
}

// HandleUserInfo 处理用户信息请求
func (h *OIDCProtocolHandler) HandleUserInfo(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	token string,
) (map[string]any, error) {
	if token == "" {
		return nil, errors.New("访问令牌不能为空")
	}

	session, err := h.sessionService.GetSessionBySessionID(token)
	if err != nil {
		return nil, errors.New("无效的访问令牌")
	}

	if session.Status != model.SessionStatusActive {
		return nil, errors.New("会话未激活")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("会话已过期")
	}

	userInfo := map[string]any{
		"sub":            session.UserID,
		"name":           "示例用户",
		"given_name":     "示例",
		"family_name":    "用户",
		"email":          "user@example.com",
		"email_verified": true,
		"locale":         "zh-CN",
	}

	return userInfo, nil
}

// HandleLogout 处理登出请求
func (h *OIDCProtocolHandler) HandleLogout(
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

func (h *OIDCProtocolHandler) generateAuthCode() (string, error) {
	return h.generateRandomString(32)
}

func (h *OIDCProtocolHandler) generateAccessToken() (string, error) {
	return h.generateRandomString(64)
}

func (h *OIDCProtocolHandler) generateRefreshToken() (string, error) {
	return h.generateRandomString(128)
}

func (h *OIDCProtocolHandler) generateIDToken(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	session *model.SsoSession,
) (string, error) {
	nonce := getStringFromExtraData(session.ExtraData, "nonce")
	
	idTokenData := map[string]any{
		"iss":       "metadata-platform",
		"sub":       session.UserID,
		"aud":       client.ClientID,
		"exp":       session.ExpiresAt.Unix(),
		"iat":       time.Now().Unix(),
		"auth_time": time.Now().Unix(),
		"nonce":     nonce,
	}

	data, err := json.Marshal(idTokenData)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(data), nil
}

func (h *OIDCProtocolHandler) generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// OIDCConfig OIDC 配置结构
type OIDCConfig struct {
	ClientID            string `json:"client_id"`
	ClientSecret        string `json:"client_secret"`
	Issuer              string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint       string `json:"token_endpoint"`
	UserInfoEndpoint    string `json:"userinfo_endpoint"`
	RedirectURI         string `json:"redirect_uri"`
	Scope               string `json:"scope"`
}
