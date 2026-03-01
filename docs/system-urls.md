# 系统默认 URL 地址

## 基础信息
- **前端开发服务器**: `http://localhost:3000`
- **后端 API 服务器**: `http://localhost:8080`
- **API 基础路径**: `/api`

---

## 前端路由

### 1. 首页模块 (`/home`)
| 路由 | URL | 组件 | 说明 |
|------|-----|------|------|
| 首页 | `/home/dashboard` | `@/views/home/Dashboard.vue` | 系统首页 |

### 2. 元数据管理模块 (`/metadata`)
| 路由 | URL | 组件 | 说明 |
|------|-----|------|------|
| 数据源列表 | `/metadata/datasource/list` | `@/views/metadata/datasource/List.vue` | 数据源管理 |
| 添加数据源 | `/metadata/datasource/create` | `@/views/metadata/datasource/Create.vue` | 添加数据源 |
| 编辑数据源 | `/metadata/datasource/:id/edit` | `@/views/metadata/datasource/Edit.vue` | 编辑数据源 |
| 表列表 | `/metadata/maintenance/table` | `@/views/metadata/table/List.vue` | 表维护 |
| 视图列表 | `/metadata/maintenance/view` | `@/views/metadata/table/List.vue` | 视图维护 |
| 存储过程 | `/metadata/maintenance/procedure` | `@/views/metadata/procedure/List.vue` | 存储过程 |
| 函数 | `/metadata/maintenance/function` | `@/views/metadata/procedure/List.vue` | 函数 |
| 字段列表 | `/metadata/field/list` | `@/views/metadata/field/List.vue` | 字段管理 |
| 模型列表 | `/metadata/model/list` | `@/views/model/List.vue` | 模型管理 |
| SQL 模型测试 | `/metadata/model/sql-test` | `@/views/model/SqlTest.vue` | SQL 测试 |
| 表视图测试 | `/metadata/model/table-view-test` | `@/views/model/TableViewTest.vue` | 表视图测试 |
| 创建模型 | `/metadata/model/create` | `@/views/model/Create.vue` | 创建模型 |
| 可视化创建模型 | `/metadata/model/visual-create` | `@/views/model/VisualCreate.vue` | 可视化创建 |

### 3. 系统监控模块 (`/monitor`)
| 路由 | URL | 组件 | 说明 |
|------|-----|------|------|
| 监控仪表盘 | `/monitor/dashboard` | `@/views/monitor/Dashboard.vue` | 监控首页 |
| 性能分析 | `/monitor/performance` | `@/views/monitor/Performance.vue` | 性能分析 |

### 4. 接口管理模块 (`/api`)
| 路由 | URL | 组件 | 说明 |
|------|-----|------|------|
| 接口列表 | `/api/list` | `@/views/api/List.vue` | 接口列表 |
| 创建接口 | `/api/create` | `@/views/api/Create.vue` | 创建接口 |
| 编辑接口 | `/api/:id/edit` | `@/views/api/Edit.vue` | 编辑接口 |

### 5. SSO 管理模块 (`/sso`)
#### 配置管理
| 路由 | URL | 组件 | 说明 |
|------|-----|------|------|
| 协议配置 | `/sso/config/protocol` | `@/views/sso/config/ProtocolConfig.vue` | 协议配置 |
| 客户端配置 | `/sso/config/client` | `@/views/sso/config/ClientConfig.vue` | 客户端配置 |
| 密钥管理 | `/sso/config/key` | `@/views/sso/config/KeyManager.vue` | 密钥管理 |
| 字段映射 | `/sso/config/mapping` | `@/views/sso/config/FieldMapping.vue` | 字段映射 |
| 会话管理 | `/sso/config/session` | `@/views/sso/config/SessionManager.vue` | 会话管理 |

#### 租户与应用
| 路由 | URL | 组件 | 说明 |
|------|-----|------|------|
| 租户管理 | `/sso/tenant` | `@/views/sso/tenant/Tenant.vue` | 租户管理 |
| 应用列表 | `/sso/app` | `@/views/sso/app/App.vue` | 应用管理 |

#### 组织与职位
| 路由 | URL | 组件 | 说明 |
|------|-----|------|------|
| 组织管理 | `/sso/org` | `@/views/sso/org/Org.vue` | 组织管理 |
| 组织类型 | `/sso/orgKind` | `@/views/sso/org/OrgKind.vue` | 组织类型 |
| 职位管理 | `/sso/pos` | `@/views/sso/pos/Pos.vue` | 职位管理 |

#### 菜单与角色
| 路由 | URL | 组件 | 说明 |
|------|-----|------|------|
| 菜单管理 | `/sso/menu` | `@/views/sso/menu/Menu.vue` | 菜单管理 |
| 角色管理 | `/sso/role` | `@/views/sso/role/Role.vue` | 角色管理 |
| 角色组 | `/sso/roleGroup` | `@/views/sso/role/RoleGroup.vue` | 角色组 |
| 用户组 | `/sso/userGroup` | `@/views/sso/user/UserGroup.vue` | 用户组 |
| 用户管理 | `/sso/user` | `@/views/sso/user/User.vue` | 用户管理 |

### 6. 系统设置模块 (`/system`)
| 路由 | URL | 组件 | 说明 |
|------|-----|------|------|
| 系统配置 | `/system/settings` | `@/views/system/Settings.vue` | 系统配置 |
| 登录日志 | `/system/audit/login` | `@/views/system/audit/LoginLog.vue` | 登录日志 |
| 操作日志 | `/system/audit/operation` | `@/views/system/audit/OperationLog.vue` | 操作日志 |
| 数据变更 | `/system/audit/data` | `@/views/system/audit/DataChangeLog.vue` | 数据变更日志 |
| 访问日志 | `/system/audit/access` | `@/views/system/audit/AccessLog.vue` | 访问日志 |

### 7. 用户模块 (`/user`)
| 路由 | URL | 组件 | 说明 |
|------|-----|------|------|
| 个人设置 | `/user/profile` | `@/views/user/Profile.vue` | 个人设置 |

### 8. 认证与错误页面
| 路由 | URL | 组件 | 说明 |
|------|-----|------|------|
| 登录 | `/login` | `@/views/login/Index.vue` | 登录页面 |
| 403 错误 | `/403` | `@/views/error/403.vue` | 无权限 |
| 500 错误 | `/500` | `@/views/error/500.vue` | 服务器错误 |
| 503 错误 | `/503` | `@/views/error/503.vue` | 服务维护 |
| 网络错误 | `/network-error` | `@/views/error/NetworkError.vue` | 网络错误 |
| 404 错误 | `/*` | `@/views/error/404.vue` | 页面不存在 |

---

## 后端 API

### SSO 模块 API (`/api/sso`)

#### 协议配置
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/sso/config` | 获取协议配置列表 |
| GET | `/api/sso/config/:id` | 获取单个协议配置 |
| POST | `/api/sso/config` | 创建协议配置 |
| PUT | `/api/sso/config/:id` | 更新协议配置 |
| DELETE | `/api/sso/config/:id` | 删除协议配置 |
| GET | `/api/sso/config/type/:type` | 按类型获取配置 |

#### 客户端管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/sso/client` | 获取客户端列表 |
| GET | `/api/sso/client/:id` | 获取单个客户端 |
| POST | `/api/sso/client` | 创建客户端 |
| PUT | `/api/sso/client/:id` | 更新客户端 |
| DELETE | `/api/sso/client/:id` | 删除客户端 |

#### 密钥管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/sso/key` | 获取密钥列表 |
| GET | `/api/sso/key/:id` | 获取单个密钥 |
| POST | `/api/sso/key` | 创建密钥 |
| POST | `/api/sso/key/generate` | 生成密钥对 |
| PUT | `/api/sso/key/:id` | 更新密钥 |
| DELETE | `/api/sso/key/:id` | 删除密钥 |

#### 字段映射
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/sso/mapping` | 获取字段映射列表 |
| GET | `/api/sso/mapping/:id` | 获取单个映射 |
| POST | `/api/sso/mapping` | 创建字段映射 |
| PUT | `/api/sso/mapping/:id` | 更新字段映射 |
| DELETE | `/api/sso/mapping/:id` | 删除字段映射 |

#### 会话管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/sso/session` | 获取会话列表 |
| GET | `/api/sso/session/:id` | 获取单个会话 |
| DELETE | `/api/sso/session/:id` | 删除会话 |
| POST | `/api/sso/session/:id/revoke` | 撤销会话 |

#### 认证端点 (⚠️ 待实现)
| 方法 | 路径 | 说明 | 状态 |
|------|------|------|------|
| GET | `/api/sso/auth/:protocol/authorize` | 授权端点 | ❌ 未实现 |
| POST | `/api/sso/auth/:protocol/token` | Token 端点 | ❌ 未实现 |
| GET | `/api/sso/auth/:protocol/userinfo` | 用户信息端点 | ❌ 未实现 |
| POST | `/api/sso/auth/:protocol/logout` | 登出端点 | ❌ 未实现 |

---

## 认证协议端点实现状态

### OIDC (OpenID Connect)
- ❌ 授权端点：`/api/sso/auth/oidc/authorize`
- ❌ Token 端点：`/api/sso/auth/oidc/token`
- ❌ 用户信息端点：`/api/sso/auth/oidc/userinfo`
- ❌ 登出端点：`/api/sso/auth/oidc/logout`

### SAML 2.0
- ❌ SSO 端点：`/api/sso/auth/saml/sso`
- ❌ SLO 端点：`/api/sso/auth/saml/slo`
- ❌ Metadata 端点：`/api/sso/auth/saml/metadata`

### CAS
- ❌ 登录端点：`/api/sso/auth/cas/login`
- ❌ 验证端点：`/api/sso/auth/cas/validate`
- ❌ 登出端点：`/api/sso/auth/cas/logout`

### LDAP
- ❌ 绑定端点：`/api/sso/auth/ldap/bind`
- ❌ 搜索端点：`/api/sso/auth/ldap/search`
- ❌ 认证端点：`/api/sso/auth/ldap/authenticate`

---

## 默认登录信息
- **默认用户名**: `admin`
- **默认密码**: 见配置文件或初始化脚本

---

## 快速访问链接

### 开发环境
- 前端：http://localhost:3000
- 后端：http://localhost:8080
- 登录页：http://localhost:3000/login
- 首页：http://localhost:3000/home/dashboard

### 主要功能入口
1. **元数据管理**: http://localhost:3000/metadata/datasource/list
2. **模型管理**: http://localhost:3000/metadata/model/list
3. **SSO 配置**: http://localhost:3000/sso/config/protocol
4. **用户管理**: http://localhost:3000/sso/user
5. **系统设置**: http://localhost:3000/system/settings
6. **监控仪表盘**: http://localhost:3000/monitor/dashboard

---

## 注意事项
1. 所有认证协议端点目前仅有路由注册，具体实现逻辑需要在对应的 Handler 中完成
2. 前端路由使用 Vue Router 的 history 模式，需要服务器配置支持
3. API 路径统一使用 `/api` 前缀，便于代理配置
4. 错误页面支持自定义路由和状态码显示
