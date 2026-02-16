package middleware

import (
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"

	"metadata-platform/internal/utils"
)

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	Name     string         `json:"name"`     // 中间件名称
	Enabled  bool           `json:"enabled"`  // 是否启用
	Priority int            `json:"priority"` // 中间件优先级
	Config   map[string]any `json:"config"`   // 中间件配置
}

// MiddlewareFactory 中间件工厂函数类型
type MiddlewareFactory func(config map[string]any) app.HandlerFunc

// MiddlewareLoader 中间件加载器
type MiddlewareLoader struct {
	middlewareFactories map[string]MiddlewareFactory // 中间件工厂映射
	middlewareConfigs   []*MiddlewareConfig          // 中间件配置列表
	middlewareChain     []app.HandlerFunc            // 中间件链
	middlewareLock      sync.RWMutex                 // 中间件锁
}

// NewMiddlewareLoader 创建中间件加载器
func NewMiddlewareLoader() *MiddlewareLoader {
	return &MiddlewareLoader{
		middlewareFactories: make(map[string]MiddlewareFactory),
		middlewareConfigs:   make([]*MiddlewareConfig, 0),
		middlewareChain:     make([]app.HandlerFunc, 0),
	}
}

// RegisterMiddleware 注册中间件工厂
func (ml *MiddlewareLoader) RegisterMiddleware(name string, factory MiddlewareFactory) {
	ml.middlewareLock.Lock()
	defer ml.middlewareLock.Unlock()
	ml.middlewareFactories[name] = factory
	utils.SugarLogger.Infof("Middleware factory registered: %s", name)
}

// LoadMiddlewareConfig 加载中间件配置
func (ml *MiddlewareLoader) LoadMiddlewareConfig(configs []*MiddlewareConfig) {
	ml.middlewareLock.Lock()
	defer ml.middlewareLock.Unlock()
	ml.middlewareConfigs = configs
	utils.SugarLogger.Infof("Loaded %d middleware configs", len(configs))
}

// BuildMiddlewareChain 构建中间件链
func (ml *MiddlewareLoader) BuildMiddlewareChain() []app.HandlerFunc {
	ml.middlewareLock.Lock()
	defer ml.middlewareLock.Unlock()

	// 按优先级排序中间件配置
	ml.sortMiddlewareConfigs()

	// 构建中间件链
	chain := make([]app.HandlerFunc, 0)

	for _, config := range ml.middlewareConfigs {
		if !config.Enabled {
			continue
		}

		factory, ok := ml.middlewareFactories[config.Name]
		if !ok {
			utils.SugarLogger.Warnf("Middleware factory not found: %s", config.Name)
			continue
		}

		mw := factory(config.Config)
		chain = append(chain, mw)
		utils.SugarLogger.Infof("Middleware added to chain: %s", config.Name)
	}

	ml.middlewareChain = chain
	return chain
}

// UseMiddlewareChain 将中间件链应用到Hertz引擎
func (ml *MiddlewareLoader) UseMiddlewareChain(r *server.Hertz) {
	chain := ml.BuildMiddlewareChain()
	for _, mw := range chain {
		r.Use(mw)
	}
	utils.SugarLogger.Infof("Applied middleware chain with %d middlewares", len(chain))
}

// GetMiddlewareChain 获取中间件链
func (ml *MiddlewareLoader) GetMiddlewareChain() []app.HandlerFunc {
	ml.middlewareLock.RLock()
	defer ml.middlewareLock.RUnlock()
	// 返回中间件链副本
	chain := make([]app.HandlerFunc, len(ml.middlewareChain))
	copy(chain, ml.middlewareChain)
	return chain
}

// GetMiddlewareFactories 获取中间件工厂映射
func (ml *MiddlewareLoader) GetMiddlewareFactories() map[string]MiddlewareFactory {
	ml.middlewareLock.RLock()
	defer ml.middlewareLock.RUnlock()
	// 返回中间件工厂映射副本
	factories := make(map[string]MiddlewareFactory)
	for name, factory := range ml.middlewareFactories {
		factories[name] = factory
	}
	return factories
}

// sortMiddlewareConfigs 按优先级排序中间件配置
func (ml *MiddlewareLoader) sortMiddlewareConfigs() {
	// 使用冒泡排序按优先级从高到低排序
	for i := 0; i < len(ml.middlewareConfigs)-1; i++ {
		for j := 0; j < len(ml.middlewareConfigs)-i-1; j++ {
			if ml.middlewareConfigs[j].Priority < ml.middlewareConfigs[j+1].Priority {
				ml.middlewareConfigs[j], ml.middlewareConfigs[j+1] = ml.middlewareConfigs[j+1], ml.middlewareConfigs[j]
			}
		}
	}
}

// RegisterDefaultMiddlewares 注册默认中间件
func (ml *MiddlewareLoader) RegisterDefaultMiddlewares() {
	// 注册CORS中间件
	ml.RegisterMiddleware("cors", func(config map[string]any) app.HandlerFunc {
		return CORSMiddleware()
	})

	// 注册日志中间件
	ml.RegisterMiddleware("logger", func(config map[string]any) app.HandlerFunc {
		return LoggerMiddleware()
	})

	// 注册认证中间件
	ml.RegisterMiddleware("auth", func(config map[string]any) app.HandlerFunc {
		return AuthMiddleware()
	})

	// 注册恢复中间件
	ml.RegisterMiddleware("recovery", func(config map[string]any) app.HandlerFunc {
		return RecoveryMiddleware()
	})

	// 注册权限中间件 (Casbin)
	ml.RegisterMiddleware("casbin", func(config map[string]any) app.HandlerFunc {
		return CasbinMiddleware()
	})

	utils.SugarLogger.Info("Registered all default middlewares")
}

// GetDefaultMiddlewareConfig 获取默认中间件配置
func GetDefaultMiddlewareConfig() []*MiddlewareConfig {
	return []*MiddlewareConfig{
		{
			Name:     "cors",
			Enabled:  true,
			Priority: 100,
			Config:   make(map[string]any),
		},
		{
			Name:     "logger",
			Enabled:  true,
			Priority: 90,
			Config:   make(map[string]any),
		},
		{
			Name:     "recovery",
			Enabled:  true,
			Priority: 80,
			Config:   make(map[string]any),
		},
		{
			Name:     "auth",
			Enabled:  false, // 认证中间件默认不全局启用，而是在路由组中使用
			Priority: 70,
			Config:   make(map[string]any),
		},
		{
			Name:     "casbin",
			Enabled:  false,
			Priority: 60,
			Config:   make(map[string]any),
		},
	}
}
