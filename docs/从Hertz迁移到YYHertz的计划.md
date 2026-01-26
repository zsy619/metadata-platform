# 从Hertz迁移到YYHertz的计划

## 1. YYHertz核心组件分析

YYHertz是一个基于CloudWeGo-Hertz的企业级Go Web框架，提供了完整的MVC架构和Beego风格的开发体验。其核心组件包括：

### 1.1 架构设计
- **MVC架构**：标准的Model-View-Controller设计模式
- **模块化设计**：领域驱动的包结构
- **Beego兼容**：100%兼容Beego命名空间路由系统
- **智能路由**：自动路由注册 + 手动路由映射

### 1.2 核心包结构
- `github.com/zsy619/yyhertz/framework/mvc`：核心MVC框架
- `github.com/zsy619/yyhertz/framework/mvc/middleware`：内置中间件
- `github.com/zsy619/yyhertz/framework/orm`：统一ORM解决方案

## 2. 迁移步骤

### 2.1 项目结构调整
1. **创建控制器目录**：`internal/controllers/`
2. **创建模型目录**：`internal/models/`（已有model目录，可复用）
3. **创建视图目录**：`views/`（可选，根据需要）
4. **调整配置文件**：适配yyhertz的配置系统

### 2.2 依赖更新
1. **添加yyhertz依赖**：
   ```go
   require github.com/zsy619/yyhertz v0.0.0-20251209140455-f7ee7cf87897
   ```
2. **保留核心依赖**：
   - jwt：用户认证
   - gorm：数据库操作
   - viper：配置管理
   - zap：日志系统

### 2.3 核心代码迁移

#### 2.3.1 入口文件迁移
将`cmd/main.go`迁移为yyhertz风格：
```go
package main

import (
    "github.com/zsy619/yyhertz/framework/mvc"
    "github.com/zsy619/yyhertz/framework/mvc/middleware"
    "metadata-platform/internal/controllers"
)

func main() {
    // 创建应用实例
    app := mvc.HertzApp
    
    // 添加中间件
    app.Use(
        middleware.RecoveryMiddleware(),
        middleware.LoggerMiddleware(),
        middleware.CORSMiddleware(),
    )
    
    // 初始化控制器
    dataSourceController := &controllers.DataSourceController{}
    modelController := &controllers.ModelController{}
    
    // 注册路由
    app.RouterAuto(dataSourceController, modelController)
    
    // 启动服务器
    app.Run()
}
```

#### 2.3.2 控制器迁移
将API处理器迁移为yyhertz控制器：
```go
package controllers

import (
    "github.com/zsy619/yyhertz/framework/mvc"
    "metadata-platform/internal/service"
)

// DataSourceController 数据源控制器
type DataSourceController struct {
    mvc.BaseController
    dataSourceService service.DataSourceService
}

// GetIndex 获取数据源列表
func (c *DataSourceController) GetIndex() {
    // 实现获取数据源列表逻辑
}

// PostCreate 创建数据源
func (c *DataSourceController) PostCreate() {
    // 实现创建数据源逻辑
}
```

#### 2.3.3 中间件迁移
将现有中间件适配为yyhertz中间件格式：
```go
package middleware

import (
    "github.com/zsy619/yyhertz/framework/mvc"
)

// TenantMiddleware 租户中间件
func TenantMiddleware() mvc.HandlerFunc {
    return func(ctx *mvc.Context) {
        // 实现租户中间件逻辑
        ctx.Next()
    }
}
```

#### 2.3.4 路由迁移
将现有路由配置迁移为yyhertz命名空间路由：
```go
package main

import (
    "github.com/zsy619/yyhertz/framework/mvc"
    "metadata-platform/internal/controllers"
)

func main() {
    // ...
    
    // 使用Beego风格的Namespace功能
    nsApi := mvc.NewNamespace("/api",
        // 数据源管理路由
        mvc.NSNamespace("/data-sources",
            mvc.NSAutoRouter(&controllers.DataSourceController{}),
        ),
        
        // 模型管理路由
        mvc.NSNamespace("/models",
            mvc.NSAutoRouter(&controllers.ModelController{}),
        ),
    )
    
    // 添加命名空间到全局应用
    mvc.AddNamespace(nsApi)
    
    // ...
}
```

## 3. 迁移注意事项

### 3.1 API兼容性
- YYHertz的API与原始Hertz有较大差异，需要完全重写控制器和路由
- 中间件签名不同，需要重新实现
- 上下文对象不同，需要调整请求处理逻辑

### 3.2 配置系统
- YYHertz有自己的配置系统，需要调整配置加载逻辑
- 支持多种配置格式：yaml、toml、json等

### 3.3 数据库操作
- YYHertz提供了统一ORM解决方案，支持GORM和MyBatis双引擎
- 可以复用现有GORM模型，但需要适配yyhertz的ORM配置

### 3.4 日志系统
- YYHertz有内置日志系统，也可以集成现有zap日志
- 需要调整日志初始化逻辑

## 4. 测试计划

1. **单元测试**：测试核心业务逻辑
2. **集成测试**：测试控制器和路由
3. **端到端测试**：测试完整业务流程
4. **性能测试**：验证迁移后的性能表现

## 5. 风险评估

### 5.1 风险点
- **API兼容性**：YYHertz的API与原始Hertz差异较大，迁移工作量大
- **项目结构调整**：需要重新组织项目结构，可能影响现有代码
- **配置系统迁移**：需要调整配置加载逻辑，可能引入配置错误
- **中间件迁移**：需要重新实现中间件，可能影响现有功能

### 5.2 缓解措施
- **逐步迁移**：分模块逐步迁移，降低风险
- **充分测试**：迁移前编写测试用例，确保功能正常
- **保留备份**：迁移前备份现有代码，以便回滚
- **文档更新**：更新项目文档，说明迁移后的架构和使用方式

## 6. 迁移时间估算

| 阶段 | 时间 | 主要工作 |
|------|------|----------|
| 需求分析 | 1天 | 分析YYHertz架构和API |
| 项目结构调整 | 1天 | 调整项目目录结构 |
| 依赖更新 | 0.5天 | 更新go.mod和go.sum |
| 核心代码迁移 | 3天 | 迁移入口文件、控制器、中间件和路由 |
| 测试 | 2天 | 编写和执行测试用例 |
| 文档更新 | 0.5天 | 更新项目文档 |
| 总计 | 8天 | 完成从Hertz到YYHertz的迁移 |

## 7. 结论

从Hertz迁移到YYHertz是一个较大的重构项目，需要调整项目结构、更新依赖、重写核心代码和测试。但是，YYHertz提供了更完整的MVC框架和丰富的内置功能，可以提高开发效率和代码质量。

建议采用**逐步迁移**的方式，先迁移核心功能，再迁移次要功能，确保迁移过程平稳进行。