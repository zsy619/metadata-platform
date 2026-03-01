# SSO 协议端点实现文档

## 概述
本文档详细说明了系统支持的单点登录（SSO）协议及其端点实现。

---

## 1. OIDC (OpenID Connect) 协议

### 1.1 协议配置参数
| 参数名 | 说明 | 示例值 | 必填 |
|--------|------|--------|------|
| client_id | 客户端 ID | `my-app-client` | ✅ |
| client_secret | 客户端密钥 | `secret-key-123456` | ✅ |
| issuer | 发行方 URL | `https://sso.example.com` | ✅ |
| authorization_endpoint | 授权端点 | `https://sso.example.com/oauth2/authorize` | ✅ |
| token_endpoint | 令牌端点 | `https://sso.example.com/oauth2/token` | ✅ |
| userinfo_endpoint | 用户信息端点 | `https://sso.example.com/oauth2/userinfo` | ✅ |
| redirect_uri | 重定向 URI | `https://app.example.com/callback` | ✅ |
| scope | 作用域 | `openid profile email` | ✅ |

### 1.2 系统端点

#### 授权端点
```
GET /api/sso/auth/oidc/authorize
```

**请求参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| response_type | string | ✅ | 响应类型，固定为 `code` |
| client_id | string | ✅ | 客户端 ID |
| redirect_uri | string | ✅ | 重定向 URI |
| scope | string | ✅ | 请求的作用域 |
| state | string | ✅ | 状态值，用于防止 CSRF |
| nonce | string | ❌ | 随机值，用于防止重放攻击 |

**响应示例（重定向）：**
```
HTTP/1.1 302 Found
Location: https://app.example.com/callback?code=AUTH_CODE&state=STATE_VALUE
```

#### 令牌端点
```
POST /api/sso/auth/oidc/token
Content-Type: application/x-www-form-urlencoded
```

**请求参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| grant_type | string | ✅ | 授权类型 (`authorization_code`, `refresh_token`, `client_credentials`) |
| code | string | ✅ | 授权码（当 grant_type=authorization_code 时） |
| redirect_uri | string | ✅ | 重定向 URI（必须与授权时一致） |
| client_id | string | ✅ | 客户端 ID |
| client_secret | string | ✅ | 客户端密钥 |
| refresh_token | string | ❌ | 刷新令牌（当 grant_type=refresh_token 时） |

**响应示例：**
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4...",
  "id_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "scope": "openid profile email"
}
```

#### 用户信息端点
```
GET /api/sso/auth/oidc/userinfo
Authorization: Bearer <access_token>
```

**请求头：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Authorization | string | ✅ | Bearer Token |

**响应示例：**
```json
{
  "sub": "1234567890",
  "name": "张三",
  "given_name": "三",
  "family_name": "张",
  "email": "zhangsan@example.com",
  "email_verified": true,
  "picture": "https://example.com/avatar.jpg",
  "locale": "zh-CN"
}
```

#### 登出端点
```
POST /api/sso/auth/oidc/logout
```

**请求参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id_token_hint | string | ❌ | ID Token |
| post_logout_redirect_uri | string | ❌ | 登出后重定向 URI |
| state | string | ❌ | 状态值 |

**响应示例：**
```json
{
  "message": "登出成功",
  "redirect_uri": "https://app.example.com/logged-out"
}
```

---

## 2. SAML 2.0 协议

### 2.1 协议配置参数
| 参数名 | 说明 | 示例值 | 必填 |
|--------|------|--------|------|
| entity_id | 实体 ID | `https://sso.example.com/saml/metadata` | ✅ |
| single_sign_on_url | 单点登录 URL | `https://idp.example.com/sso` | ✅ |
| single_logout_url | 单点登出 URL | `https://idp.example.com/slo` | ❌ |
| x509_certificate | X509 证书 | `MIICpDCCAYwCCQD...` | ✅ |
| acs_url | 断言消费服务 URL | `https://app.example.com/saml/acs` | ✅ |

### 2.2 系统端点

#### SSO 端点（SP 发起）
```
GET /api/sso/auth/saml/sso
```

**请求参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| RelayState | string | ❌ | 中继状态 |

**响应：**
- 重定向到 IdP 的 SSO URL

#### SSO 端点（IdP 响应）
```
POST /api/sso/auth/saml/acs
```

**请求参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| SAMLResponse | string | ✅ | SAML 断言 |
| RelayState | string | ❌ | 中继状态 |

**响应：**
- 验证 SAML 断言
- 创建本地会话
- 重定向到应用

#### 单点登出端点
```
GET|POST /api/sso/auth/saml/slo
```

**请求参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| SAMLRequest | string | ✅ | SAML 登出请求 |
| RelayState | string | ❌ | 中继状态 |

#### Metadata 端点
```
GET /api/sso/auth/saml/metadata
```

**响应：**
```xml
<?xml version="1.0"?>
<EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" 
                  entityID="https://sso.example.com/saml/metadata">
  <SPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" 
                              Location="https://app.example.com/saml/acs" 
                              index="0" isDefault="true"/>
  </SPSSODescriptor>
</EntityDescriptor>
```

---

## 3. CAS 协议

### 3.1 协议配置参数
| 参数名 | 说明 | 示例值 | 必填 |
|--------|------|--------|------|
| server_url | CAS 服务器 URL | `https://cas.example.com` | ✅ |
| service_url | 服务 URL | `https://app.example.com` | ✅ |
| login_url | 登录 URL | `https://cas.example.com/login` | ✅ |
| logout_url | 登出 URL | `https://cas.example.com/logout` | ✅ |
| validate_url | 验证 URL | `https://cas.example.com/serviceValidate` | ✅ |

### 3.2 系统端点

#### 登录端点
```
GET /api/sso/auth/cas/login
```

**请求参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| service | string | ✅ | 服务 URL |

**响应：**
- 重定向到 CAS 服务器登录页面

#### 验证端点
```
GET /api/sso/auth/cas/validate
```

**请求参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ticket | string | ✅ | CAS Ticket |
| service | string | ✅ | 服务 URL |

**响应示例（CAS 2.0 XML）：**
```xml
<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
  <cas:authenticationSuccess>
    <cas:user>zhangsan</cas:user>
    <cas:attributes>
      <cas:email>zhangsan@example.com</cas:email>
      <cas:name>张三</cas:name>
    </cas:attributes>
  </cas:authenticationSuccess>
</cas:serviceResponse>
```

#### 登出端点
```
GET|POST /api/sso/auth/cas/logout
```

**请求参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| service | string | ❌ | 服务 URL |

**响应：**
- 清除本地会话
- 重定向到 CAS 服务器登出页面

---

## 4. LDAP 协议

### 4.1 协议配置参数
| 参数名 | 说明 | 示例值 | 必填 |
|--------|------|--------|------|
| server_url | LDAP 服务器 URL | `ldap://ldap.example.com` | ✅ |
| bind_dn | 绑定 DN | `cn=admin,dc=example,dc=com` | ✅ |
| bind_password | 绑定密码 | `admin-password` | ✅ |
| base_dn | 基础 DN | `dc=example,dc=com` | ✅ |
| user_filter | 用户过滤器 | `(uid={username})` | ✅ |
| group_filter | 组过滤器 | `(memberUid={username})` | ❌ |
| port | 端口 | `389` 或 `636`(SSL) | ✅ |
| use_ssl | 使用 SSL | `true` 或 `false` | ❌ |

### 4.2 系统端点

#### 绑定端点
```
POST /api/sso/auth/ldap/bind
```

**请求参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | ✅ | 用户名 |
| password | string | ✅ | 密码 |

**处理流程：**
1. 使用 bind_dn 连接到 LDAP 服务器
2. 使用 user_filter 搜索用户
3. 使用用户 DN 尝试绑定
4. 绑定成功则返回用户信息

**响应示例：**
```json
{
  "success": true,
  "user": {
    "dn": "uid=zhangsan,ou=users,dc=example,dc=com",
    "uid": "zhangsan",
    "mail": "zhangsan@example.com",
    "cn": "张三",
    "memberOf": [
      "cn=developers,ou=groups,dc=example,dc=com"
    ]
  }
}
```

#### 搜索端点
```
POST /api/sso/auth/ldap/search
```

**请求参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| base_dn | string | ✅ | 搜索基础 DN |
| filter | string | ✅ | 搜索过滤器 |
| attributes | array | ❌ | 返回属性列表 |

**响应示例：**
```json
{
  "success": true,
  "entries": [
    {
      "dn": "uid=zhangsan,ou=users,dc=example,dc=com",
      "attributes": {
        "uid": ["zhangsan"],
        "mail": ["zhangsan@example.com"],
        "cn": ["张三"]
      }
    }
  ]
}
```

#### 认证端点
```
POST /api/sso/auth/ldap/authenticate
```

**请求参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | ✅ | 用户名 |
| password | string | ✅ | 密码 |

**响应示例：**
```json
{
  "success": true,
  "message": "认证成功",
  "user": {
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "name": "张三",
    "groups": ["developers", "users"]
  }
}
```

---

## 5. 实现状态

| 协议 | 授权端点 | Token 端点 | 用户信息端点 | 登出端点 | 状态 |
|------|---------|----------|------------|---------|------|
| OIDC | ✅ | ✅ | ✅ | ✅ | 已实现 |
| SAML | ✅ | N/A | N/A | ✅ | 框架已实现 |
| CAS | ✅ | N/A | N/A | ✅ | 框架已实现 |
| LDAP | ✅ | N/A | N/A | N/A | 框架已实现 |

---

## 6. 快速开始

### 6.1 配置 OIDC 协议
1. 访问 `/sso/config/protocol`
2. 点击"新增协议"
3. 选择协议类型为 "OIDC"
4. 填写配置参数
5. 保存配置

### 6.2 测试授权流程
```bash
# 1. 获取授权码
curl -v "http://localhost:8080/api/sso/auth/oidc/authorize?response_type=code&client_id=YOUR_CLIENT_ID&redirect_uri=YOUR_REDIRECT_URI&scope=openid+profile+email&state=RANDOM_STATE"

# 2. 使用授权码获取 Token
curl -X POST "http://localhost:8080/api/sso/auth/oidc/token" \
  -d "grant_type=authorization_code" \
  -d "code=AUTH_CODE" \
  -d "redirect_uri=YOUR_REDIRECT_URI" \
  -d "client_id=YOUR_CLIENT_ID" \
  -d "client_secret=YOUR_CLIENT_SECRET"

# 3. 使用 Token 获取用户信息
curl "http://localhost:8080/api/sso/auth/oidc/userinfo" \
  -H "Authorization: Bearer ACCESS_TOKEN"
```

---

## 7. 错误码说明

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 501 | 功能未实现 |
| 503 | 服务不可用 |

---

## 8. 安全建议

1. **HTTPS**: 所有生产环境的端点都应该使用 HTTPS
2. **Client Secret**: 妥善保管客户端密钥，不要在前端暴露
3. **State 参数**: 始终使用 state 参数防止 CSRF 攻击
4. **Nonce**: OIDC 建议使用 nonce 防止重放攻击
5. **Token 有效期**: 合理设置 access_token 和 refresh_token 的有效期
6. **CORS**: 正确配置跨域策略
7. **速率限制**: 对认证端点实施速率限制防止暴力破解
