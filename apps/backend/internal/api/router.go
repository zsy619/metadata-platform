package api

import (
	"fmt"
	"strings"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"

	"metadata-platform/internal/utils"
)

// Route 路由定义结构
type Route struct {
	Method      string            `json:"method"`      // HTTP方法
	Path        string            `json:"path"`        // 路由路径
	Handler     any       `json:"handler"`     // 处理函数
	Middlewares []string          `json:"middlewares"` // 中间件列表
	Priority    int               `json:"priority"`    // 路由优先级
	Group       string            `json:"group"`       // 路由分组
	Metadata    map[string]string `json:"metadata"`    // 路由元数据
}

// Router 路由管理器
type Router struct {
	hz           *server.Hertz              // Hertz引擎
	Routes       map[string][]*Route        // 已注册路由
	routeLock    sync.RWMutex               // 路由锁
	middlewareMap map[string]app.HandlerFunc // 中间件映射
}

// NewRouter 创建路由管理器
func NewRouter(hz *server.Hertz) *Router {
	return &Router{
		hz:           hz,
		Routes:       make(map[string][]*Route),
		middlewareMap: make(map[string]app.HandlerFunc),
	}
}

// RegisterMiddleware 注册中间件
func (r *Router) RegisterMiddleware(name string, mw app.HandlerFunc) {
	r.middlewareMap[name] = mw
}

// RegisterRoute 注册单个路由
func (r *Router) RegisterRoute(route *Route) error {
	// 路由冲突检测
	if err := r.checkRouteConflict(route); err != nil {
		return err
	}

	// 查找或创建路由组
	group := r.hz.Group(route.Group)

	// 添加中间件
	for _, mwName := range route.Middlewares {
		if mw, ok := r.middlewareMap[mwName]; ok {
			group.Use(mw)
		} else {
			return fmt.Errorf("middleware %s not found", mwName)
		}
	}

	// 注册路由
	switch route.Method {
	case "GET":
		group.GET(route.Path, route.Handler.(app.HandlerFunc))
	case "POST":
		group.POST(route.Path, route.Handler.(app.HandlerFunc))
	case "PUT":
		group.PUT(route.Path, route.Handler.(app.HandlerFunc))
	case "DELETE":
		group.DELETE(route.Path, route.Handler.(app.HandlerFunc))
	case "PATCH":
		group.PATCH(route.Path, route.Handler.(app.HandlerFunc))
	case "OPTIONS":
		group.OPTIONS(route.Path, route.Handler.(app.HandlerFunc))
	default:
		return fmt.Errorf("unsupported HTTP method: %s", route.Method)
	}

	// 保存路由信息
	r.routeLock.Lock()
	defer r.routeLock.Unlock()
	r.Routes[route.Group] = append(r.Routes[route.Group], route)

	utils.SugarLogger.Infof("Route registered: %s %s%s", route.Method, route.Group, route.Path)

	return nil
}

// BatchRegisterRoutes 批量注册路由
func (r *Router) BatchRegisterRoutes(routes []*Route) []error {
	var errors []error

	for _, route := range routes {
		if err := r.RegisterRoute(route); err != nil {
			errors = append(errors, fmt.Errorf("failed to register route %s %s: %w", route.Method, route.Path, err))
		}
	}

	return errors
}

// checkRouteConflict 检查路由冲突
func (r *Router) checkRouteConflict(route *Route) error {
	r.routeLock.RLock()
	defer r.routeLock.RUnlock()

	if groupRoutes, ok := r.Routes[route.Group]; ok {
		for _, existingRoute := range groupRoutes {
			if existingRoute.Method == route.Method && existingRoute.Path == route.Path {
				return fmt.Errorf("route conflict: %s %s already exists", route.Method, route.Group+route.Path)
			}
		}
	}

	return nil
}

// GetRoutes 获取所有已注册路由
func (r *Router) GetRoutes() map[string][]*Route {
	r.routeLock.RLock()
	defer r.routeLock.RUnlock()

	// 返回路由副本
	routesCopy := make(map[string][]*Route)
	for group, routes := range r.Routes {
		routesCopy[group] = make([]*Route, len(routes))
		copy(routesCopy[group], routes)
	}

	return routesCopy
}

// GetRoutesByGroup 获取指定分组的路由
func (r *Router) GetRoutesByGroup(group string) []*Route {
	r.routeLock.RLock()
	defer r.routeLock.RUnlock()

	if routes, ok := r.Routes[group]; ok {
		// 返回路由副本
		routesCopy := make([]*Route, len(routes))
		copy(routesCopy, routes)
		return routesCopy
	}

	return nil
}

// DeleteRoute 删除路由
func (r *Router) DeleteRoute(method, path, group string) error {
	// 注意：Hertz目前不支持运行时删除路由，这里只从记录中移除
	r.routeLock.Lock()
	defer r.routeLock.Unlock()

	if groupRoutes, ok := r.Routes[group]; ok {
		for i, route := range groupRoutes {
			if route.Method == method && route.Path == path {
				// 从切片中移除路由
				r.Routes[group] = append(groupRoutes[:i], groupRoutes[i+1:]...)
				utils.SugarLogger.Infof("Route deleted: %s %s%s", method, group, path)
				return nil
			}
		}
	}

	return fmt.Errorf("route not found: %s %s%s", method, group, path)
}

// UpdateRoute 更新路由
func (r *Router) UpdateRoute(route *Route) error {
	// 注意：Hertz目前不支持运行时更新路由，这里先删除再重新注册
	if err := r.DeleteRoute(route.Method, route.Path, route.Group); err != nil {
		// 如果路由不存在，直接注册
		if !strings.Contains(err.Error(), "route not found") {
			return err
		}
	}

	return r.RegisterRoute(route)
}

// PrintRoutes 打印所有已注册路由
func (r *Router) PrintRoutes() {
	r.routeLock.RLock()
	defer r.routeLock.RUnlock()

	utils.SugarLogger.Info("All registered routes:")
	for group, routes := range r.Routes {
		for _, route := range routes {
			utils.SugarLogger.Infof("%s %s%s (Priority: %d)", route.Method, group, route.Path, route.Priority)
		}
	}
}
