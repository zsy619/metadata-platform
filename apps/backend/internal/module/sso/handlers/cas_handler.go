package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/module/sso/service"
	"time"
)

// CASProtocolHandler CAS 协议处理器
type CASProtocolHandler struct {
	sessionService service.SsoSessionService
}

// NewCASProtocolHandler 创建 CAS 协议处理器
func NewCASProtocolHandler(sessionService service.SsoSessionService) *CASProtocolHandler {
	return &CASProtocolHandler{
		sessionService: sessionService,
	}
}

// GetProtocolType 获取协议类型
func (h *CASProtocolHandler) GetProtocolType() model.ProtocolType {
	return model.ProtocolTypeCAS
}

// HandleAuthorization 处理 CAS 登录请求
func (h *CASProtocolHandler) HandleAuthorization(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	params map[string]any,
) (any, error) {
	casConfig, err := parseCASConfig(config)
	if err != nil {
		return nil, err
	}

	service, _ := params["service"].(string)
	if service == "" {
		service = casConfig.ServiceURL
	}

	ticket := generateCASTicket()

	loginURL := casConfig.LoginURL
	if loginURL == "" {
		loginURL = casConfig.ServerURL + "/login"
	}

	u, err := url.Parse(loginURL)
	if err != nil {
		return nil, fmt.Errorf("解析登录 URL 失败：%w", err)
	}

	query := u.Query()
	query.Set("service", service)
	u.RawQuery = query.Encode()

	extraData, _ := marshalExtraData(map[string]any{
		"service": service,
	})

	session := &model.SsoSession{
		SessionID:        ticket,
		UserID:           "",
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
		"redirect_url": u.String(),
		"ticket":       ticket,
		"service":      service,
	}, nil
}

// HandleToken CAS 不支持 Token 端点
func (h *CASProtocolHandler) HandleToken(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	params map[string]any,
) (any, error) {
	return nil, errors.New("CAS 协议不支持 Token 端点")
}

// HandleUserInfo CAS 不支持 UserInfo 端点
func (h *CASProtocolHandler) HandleUserInfo(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	token string,
) (map[string]any, error) {
	return nil, errors.New("CAS 协议不支持 UserInfo 端点")
}

// HandleLogout 处理 CAS 登出请求
func (h *CASProtocolHandler) HandleLogout(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	session *model.SsoSession,
) error {
	if session == nil {
		return errors.New("会话不能为空")
	}

	casConfig, err := parseCASConfig(config)
	if err != nil {
		return err
	}

	logoutURL := casConfig.LogoutURL
	if logoutURL == "" {
		logoutURL = casConfig.ServerURL + "/logout"
	}

	service := getStringFromExtraData(session.ExtraData, "service")
	if service != "" {
		logoutURL += "?service=" + url.QueryEscape(service)
	}

	session.Status = model.SessionStatusLoggedOut
	session.ExpiresAt = time.Now()

	return h.sessionService.UpdateSession(session)
}

// HandleValidate 处理 CAS 验证请求（支持 CAS 1.0/2.0/3.0）
func (h *CASProtocolHandler) HandleValidate(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	ticket string,
	service string,
) (map[string]any, error) {
	// 验证票据格式
	if len(ticket) < 5 {
		return nil, errors.New("无效的票据格式")
	}

	// TODO: 实现真实的 CAS 服务器验证逻辑
	// 这里是完整的实现框架，实际使用时需要：
	// 1. 构建验证 URL（根据 CAS 版本）
	// 2. 发送 HTTP GET 请求到 CAS 服务器
	// 3. 解析响应（CAS 1.0 是纯文本，CAS 2.0/3.0 是 XML）
	// 4. 验证票据有效性
	// 5. 提取用户信息

	// 模拟验证成功
	userInfo := map[string]any{
		"user":  "cas_user",
		"email": "cas@example.com",
		"name":  "CAS User",
	}

	// CAS 2.0/3.0 支持更多属性
	if config.ProtocolType == "cas_2_0" || config.ProtocolType == "cas_3_0" {
		userInfo["attributes"] = map[string]any{
			"displayName": "CAS User",
			"email":       "cas@example.com",
		}
	}

	// CAS 3.0 支持完整的用户属性
	if config.ProtocolType == "cas_3_0" {
		userInfo["attributes"] = map[string]any{
			"displayName": "CAS User",
			"email":       "cas@example.com",
			"groups":      []string{"users", "cas_users"},
		}
	}

	sessionID := fmt.Sprintf("cas-sid-%d", time.Now().UnixNano())
	extraData, _ := marshalExtraData(userInfo)

	session := &model.SsoSession{
		SessionID:        sessionID,
		UserID:           userInfo["user"].(string),
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
		"success":    true,
	}, nil
}

// HandleProxyValidate 处理 CAS 2.0/3.0 代理验证请求
func (h *CASProtocolHandler) HandleProxyValidate(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	ticket string,
	service string,
) (map[string]any, error) {
	// 只有 CAS 2.0 和 3.0 支持代理验证
	if config.ProtocolType == "cas_1_0" {
		return nil, errors.New("CAS 1.0 不支持代理验证")
	}

	// 重用标准验证逻辑
	return h.HandleValidate(config, client, ticket, service)
}

// HandleProxy 处理 CAS 2.0/3.0 代理票据请求
func (h *CASProtocolHandler) HandleProxy(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	targetService string,
) (map[string]any, error) {
	// 只有 CAS 2.0 和 3.0 支持代理
	if config.ProtocolType == "cas_1_0" {
		return nil, errors.New("CAS 1.0 不支持代理")
	}

	// 生成代理票据
	proxyTicket := "PT-" + generateRandomString(32)

	return map[string]any{
		"proxy_ticket":       proxyTicket,
		"target_service":     targetService,
		"proxy_callback_url": "", // TODO: 从配置中获取
	}, nil
}

// HandleProfile 处理 CAS 3.0 用户信息请求
func (h *CASProtocolHandler) HandleProfile(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	session *model.SsoSession,
) (map[string]any, error) {
	// 只有 CAS 3.0 支持用户信息端点
	if config.ProtocolType != "cas_3_0" {
		return nil, errors.New("只有 CAS 3.0 支持用户信息端点")
	}

	if session == nil {
		return nil, errors.New("会话不能为空")
	}

	// 从会话中提取用户信息
	userInfo, err := unmarshalExtraData(session.ExtraData)
	if err != nil {
		return nil, fmt.Errorf("提取用户信息失败：%w", err)
	}

	return map[string]any{
		"success":    true,
		"attributes": userInfo,
	}, nil
}

func parseCASConfig(config *model.SsoProtocolConfig) (*CASConfig, error) {
	if config.ConfigData == "" {
		return nil, errors.New("配置数据为空")
	}

	var casConfig CASConfig
	if err := json.Unmarshal([]byte(config.ConfigData), &casConfig); err != nil {
		return nil, fmt.Errorf("解析配置数据失败：%w", err)
	}

	return &casConfig, nil
}

// CASConfig CAS 配置结构
type CASConfig struct {
	ServerURL            string `json:"server_url"`
	ServiceURL           string `json:"service_url"`
	LoginURL             string `json:"login_url"`
	LogoutURL            string `json:"logout_url"`
	ValidateURL          string `json:"validate_url"`
	ServiceValidateURL   string `json:"service_validate_url"`
	ProxyValidateURL     string `json:"proxy_validate_url"`
	ProxyURL             string `json:"proxy_url"`
	ProxyCallbackURL     string `json:"proxy_callback_url"`
	ProfileURL           string `json:"profile_url"`
}

func generateCASTicket() string {
	return "ST-" + generateRandomString(32)
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(b)
}
