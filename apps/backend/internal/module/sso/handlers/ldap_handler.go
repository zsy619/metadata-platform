package handlers

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/module/sso/service"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// LDAPProtocolHandler LDAP 协议处理器
type LDAPProtocolHandler struct {
	sessionService service.SsoSessionService
}

// NewLDAPProtocolHandler 创建 LDAP 协议处理器
func NewLDAPProtocolHandler(sessionService service.SsoSessionService) *LDAPProtocolHandler {
	return &LDAPProtocolHandler{
		sessionService: sessionService,
	}
}

// GetProtocolType 获取协议类型
func (h *LDAPProtocolHandler) GetProtocolType() model.ProtocolType {
	return model.ProtocolTypeLDAP
}

// HandleAuthorization LDAP 不支持授权端点
func (h *LDAPProtocolHandler) HandleAuthorization(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	params map[string]any,
) (any, error) {
	return nil, errors.New("LDAP 协议不支持授权端点")
}

// HandleToken LDAP 不支持 Token 端点
func (h *LDAPProtocolHandler) HandleToken(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	params map[string]any,
) (any, error) {
	return nil, errors.New("LDAP 协议不支持 Token 端点")
}

// HandleUserInfo LDAP 不支持 UserInfo 端点
func (h *LDAPProtocolHandler) HandleUserInfo(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	token string,
) (map[string]any, error) {
	return nil, errors.New("LDAP 协议不支持 UserInfo 端点")
}

// HandleLogout 处理 LDAP 登出请求
func (h *LDAPProtocolHandler) HandleLogout(
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

// HandleAuthenticate 处理 LDAP 认证请求
func (h *LDAPProtocolHandler) HandleAuthenticate(
	config *model.SsoProtocolConfig,
	client *model.SsoClient,
	username string,
	password string,
) (map[string]any, error) {
	ldapConfig, err := parseLDAPConfig(config)
	if err != nil {
		return nil, err
	}

	// 连接到 LDAP 服务器
	conn, err := h.connectLDAP(ldapConfig)
	if err != nil {
		return nil, fmt.Errorf("连接 LDAP 服务器失败：%w", err)
	}
	defer conn.Close()

	// 使用管理员账号绑定
	err = conn.Bind(ldapConfig.BindDN, ldapConfig.BindPassword)
	if err != nil {
		return nil, fmt.Errorf("管理员绑定失败：%w", err)
	}

	// 搜索用户
	userDN, userInfo, err := h.searchUser(conn, ldapConfig, username)
	if err != nil {
		return nil, fmt.Errorf("搜索用户失败：%w", err)
	}

	// 使用用户凭据尝试绑定验证
	err = conn.Bind(userDN, password)
	if err != nil {
		return nil, fmt.Errorf("用户认证失败：%w", err)
	}

	// 获取用户组信息
	groups, err := h.getGroups(conn, ldapConfig, username)
	if err != nil {
		// 组查询失败不影响认证结果，仅记录警告
		groups = []string{"users"}
	}
	userInfo["groups"] = groups

	// 创建会话
	sessionID := fmt.Sprintf("ldap-sid-%d", time.Now().UnixNano())
	extraData, _ := marshalExtraData(userInfo)

	session := &model.SsoSession{
		SessionID:        sessionID,
		UserID:           username,
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
		"success":    true,
		"session_id": sessionID,
		"user_info":  userInfo,
	}, nil
}

// connectLDAP 连接到 LDAP 服务器
func (h *LDAPProtocolHandler) connectLDAP(config *LDAPConfig) (*ldap.Conn, error) {
	ldapURL := fmt.Sprintf("%s:%d", config.ServerURL, config.Port)
	var conn *ldap.Conn
	var err error

	// 根据配置选择连接方式
	if config.UseSSL {
		// SSL 连接 (ldaps://)
		conn, err = ldap.DialTLS("tcp", ldapURL, &tls.Config{
			InsecureSkipVerify: false,
		})
	} else {
		// 普通连接
		conn, err = ldap.Dial("tcp", ldapURL)
	}

	if err != nil {
		return nil, fmt.Errorf("连接 LDAP 服务器失败：%w", err)
	}

	// 如果使用 TLS 但不是 SSL，启动 StartTLS
	if config.UseTLS && !config.UseSSL {
		err = conn.StartTLS(&tls.Config{
			InsecureSkipVerify: false,
		})
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("启动 TLS 失败：%w", err)
		}
	}

	// 注意：go-ldap/v3 默认使用 LDAP v3，不支持显式设置版本
	// 引用跟随也不直接支持，需要在搜索时处理

	return conn, nil
}

// searchUser 搜索 LDAP 用户
func (h *LDAPProtocolHandler) searchUser(conn *ldap.Conn, config *LDAPConfig, username string) (string, map[string]any, error) {
	// 确定搜索基础 DN
	searchBase := config.BaseDN
	if config.UserBaseDN != "" {
		searchBase = config.UserBaseDN
	}

	// 构建搜索过滤器
	filter := config.UserFilter
	if filter == "" {
		filter = "(uid={username})"
	}
	filter = strings.ReplaceAll(filter, "{username}", username)

	// 如果用户对象类不为空，添加到过滤器中
	if config.UserObjectClass != "" {
		filter = fmt.Sprintf("(&%s(objectClass=%s))", filter, config.UserObjectClass)
	}

	// 执行搜索
	searchRequest := ldap.NewSearchRequest(
		searchBase,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,     // 时间限制
		0,     // 结果大小限制
		false, // 仅返回类型
		filter,
		[]string{"dn", "uid", "cn", "sn", "givenName", "mail", "email", "displayName", "memberOf"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return "", nil, fmt.Errorf("搜索用户失败：%w", err)
	}

	if len(result.Entries) == 0 {
		return "", nil, errors.New("未找到用户")
	}

	if len(result.Entries) > 1 {
		return "", nil, errors.New("找到多个匹配用户")
	}

	// 获取用户 DN
	userEntry := result.Entries[0]
	userDN := userEntry.DN

	// 提取用户属性
	userInfo := map[string]any{
		"username":    username,
		"uid":         getAttributeValue(userEntry, "uid"),
		"cn":          getAttributeValue(userEntry, "cn"),
		"sn":          getAttributeValue(userEntry, "sn"),
		"givenName":   getAttributeValue(userEntry, "givenName"),
		"mail":        getAttributeValue(userEntry, "mail"),
		"email":       getAttributeValue(userEntry, "email"),
		"displayName": getAttributeValue(userEntry, "displayName"),
	}

	// 获取组成员关系
	var userGroups []string
	memberOf := userEntry.GetAttributeValues("memberOf")
	if len(memberOf) > 0 {
		userGroups = memberOf
	}

	userInfo["groups"] = userGroups

	return userDN, userInfo, nil
}

// getGroups 获取用户所属组
func (h *LDAPProtocolHandler) getGroups(conn *ldap.Conn, config *LDAPConfig, username string) ([]string, error) {
	// 确定搜索基础 DN
	searchBase := config.BaseDN
	if config.GroupBaseDN != "" {
		searchBase = config.GroupBaseDN
	}

	// 构建组搜索过滤器
	filter := config.GroupFilter
	if filter == "" {
		filter = "(memberUid={username})"
	}
	filter = strings.ReplaceAll(filter, "{username}", username)

	// 如果组对象类不为空，添加到过滤器中
	if config.GroupObjectClass != "" {
		filter = fmt.Sprintf("(&%s(objectClass=%s))", filter, config.GroupObjectClass)
	}

	// 执行搜索
	searchRequest := ldap.NewSearchRequest(
		searchBase,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,     // 时间限制
		0,     // 结果大小限制
		false, // 仅返回类型
		filter,
		[]string{"cn", "dn"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		// 组查询失败不影响认证，返回默认组
		return []string{"users"}, nil
	}

	// 提取组名
	var groups []string
	for _, entry := range result.Entries {
		cn := entry.GetAttributeValue("cn")
		if cn != "" {
			groups = append(groups, cn)
		}
	}

	if len(groups) == 0 {
		groups = []string{"users"}
	}

	return groups, nil
}

// getAttributeValue 获取 LDAP 属性值
func getAttributeValue(entry *ldap.Entry, attributeName string) string {
	attr := entry.GetAttributeValue(attributeName)
	if attr == "" {
		// 尝试获取同义词
		switch attributeName {
		case "mail":
			return entry.GetAttributeValue("email")
		case "email":
			return entry.GetAttributeValue("mail")
		case "uid":
			if v := entry.GetAttributeValue("cn"); v != "" {
				return v
			}
		}
	}
	return attr
}

func parseLDAPConfig(config *model.SsoProtocolConfig) (*LDAPConfig, error) {
	if config.ConfigData == "" {
		return nil, errors.New("配置数据为空")
	}

	var ldapConfig LDAPConfig
	if err := json.Unmarshal([]byte(config.ConfigData), &ldapConfig); err != nil {
		return nil, fmt.Errorf("解析配置数据失败：%w", err)
	}

	return &ldapConfig, nil
}

// LDAPConfig LDAP 配置结构
type LDAPConfig struct {
	// 基础配置
	ServerURL    string `json:"server_url"`
	BindDN       string `json:"bind_dn"`
	BindPassword string `json:"bind_password"`
	BaseDN       string `json:"base_dn"`
	Port         int    `json:"port"`
	UseSSL       bool   `json:"use_ssl"`
	UseTLS       bool   `json:"use_tls"`
	Version      string `json:"version"`

	// 高级配置
	UserFilter        string `json:"user_filter"`
	GroupFilter       string `json:"group_filter"`
	UserBaseDN        string `json:"user_base_dn"`
	GroupBaseDN       string `json:"group_base_dn"`
	UserObjectClass   string `json:"user_object_class"`
	GroupObjectClass  string `json:"group_object_class"`
	ReferralFollowing bool   `json:"referral_following"`
	ConnectionTimeout int    `json:"connection_timeout"`
}
