# 元数据管理平台 (Metadata Management Platform)

这是一个基于现代化技术栈构建的企业级元数据管理平台，旨在提供高效的数据治理、元数据采集与管理功能。采用前后端分离架构，后端使用 Go 语言的高性能框架 Hertz，前端使用 Vue 3 和 TypeScript 开发。

## ✨ 核心特性

- **数据源管理**：支持多种数据库类型的元数据采集与连接管理。
- **模型管理**：灵活的元数据建模与血缘分析。
- **单点登录与权限 (SSO)**：基于 Casbin 的精细化权限控制，支持多租户、多应用管理。
- **现代化 UI**：
  - 侧边栏可折叠设计，具有平滑的过渡动画。
  - 响应式布局，适配多种屏幕尺寸。
  - 集成了 Element Plus 组件库，提供极致交互体验。

## 🛠 技术栈

### 后端 (Backend)
*   **语言**: Go 1.25+
*   **Web 框架**: [CloudWeGo Hertz](https://www.cloudwego.io/zh/docs/hertz/) - 高性能 HTTP 框架
*   **ORM**: [GORM](https://gorm.io/zh_CN/) - 强大的 ORM 库
*   **权限引擎**: [Casbin](https://casbin.org/) - 灵活的角色访问控制 (RBAC)
*   **数据库**: MySQL 8.0+
*   **配置管理**: Viper
*   **日志**: Zap + Lumberjack
*   **认证**: JWT

### 前端 (Frontend)
*   **框架**: Vue 3 (Composition API)
*   **构建工具**: Vite
*   **语言**: TypeScript
*   **UI 组件库**: Element Plus
*   **图标库**: @element-plus/icons-vue
*   **状态管理**: Pinia
*   **路由**: Vue Router

## 📂 项目结构

```
.
├── apps
│   ├── backend      # 后端项目代码
│   │   ├── cmd      # 入口文件 (main.go)
│   │   ├── configs  # 配置文件与 RBAC 定义
│   │   ├── internal # 内部业务逻辑 (Module-based API, Service, Model, Repository)
│   │   └── ...
│   └── frontend     # 前端项目代码
│       ├── src      # 源代码 (Components, Views, Assets)
│       ├── public   # 静态资源
│       └── ...
├── docs             # 项目文档 (设计文档、接口定义等)
└── README.md        # 项目说明
```

## 🚀 快速开始

### 环境要求
*   Go 1.25+
*   Node.js 18+
*   MySQL 8.0+

### 后端启动

1.  **进入后端目录**：
    ```bash
    cd apps/backend
    ```
2.  **下载依赖**：
    ```bash
    go mod tidy
    ```
3.  **配置数据库**：
    修改 `configs` 目录下的配置文件（如 `config.yaml` 或其示例文件），设置正确的 MySQL 连接信息。
4.  **初始化与迁移**：
    启动时会自动执行数据库迁移 (`Migrate`) 和种子数据初始化 (`SeedData`)。
5.  **运行服务**：
    ```bash
    go run cmd/main.go
    ```

### 前端启动

1.  **进入前端目录**：
    ```bash
    cd apps/frontend
    ```
2.  **安装依赖**：
    ```bash
    npm install
    # 或者使用 pnpm
    pnpm install
    ```
3.  **启动开发服务器**：
    ```bash
    npm run dev
    ```
    访问 http://localhost:5173 (或终端显示的端口)。

## 📝 文档与规范

更多详细文档请查看 `docs` 目录。
- **重构说明**：用户模块已统一采用 `Sso` 前缀命名，并补全了 `form` 标签以支持多种参数绑定。
- **代码规范**：遵循 Effective Go 和 Vue 3 最佳实践。

---

Copyright &copy; 2026 元数据管理平台. All rights reserved.
