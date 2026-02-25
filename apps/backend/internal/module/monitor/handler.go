package monitor

import (
	"context"
	"math/rand"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/websocket"
	"gorm.io/gorm"

	"metadata-platform/internal/middleware"
	"metadata-platform/internal/module/monitor/api"
	"metadata-platform/internal/module/monitor/service"
)

var upgrader = websocket.HertzUpgrader{
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
	},
}

// RegisterRoutes 注册监控模块路由
func RegisterRoutes(h *server.Hertz) {
	// 初始化服务
	monitorSvc := service.NewMonitorService()
	statsSvc := service.NewStatsService(nil, nil) // 暂时为空，需要通过 SetDB 方法设置

	// 初始化处理器
	handler := api.NewMonitorHandler(monitorSvc, statsSvc)

	// 公开路由组（WebSocket）
	wsGroup := h.Group("/api/monitor")
	wsGroup.GET("/ws", WSHandler)

	// 需要认证的路由组
	group := h.Group("/api/monitor")
	group.Use(middleware.AuthMiddleware())

	// 仪表盘汇总
	group.GET("/dashboard", handler.GetDashboardSummary)

	// 系统监控
	group.GET("/system", handler.GetSystemStats)

	// 业务统计
	group.GET("/business", handler.GetBusinessStats)

	// 访问统计
	group.GET("/access", handler.GetAccessStats)

	// 趋势数据
	group.GET("/trend/hourly", handler.GetHourlyTrend)
	group.GET("/trend/daily", handler.GetDailyTrend)

	// TOP数据
	group.GET("/top-paths", handler.GetTopPaths)

	// 分布数据
	group.GET("/status-distribution", handler.GetStatusDistribution)
	group.GET("/latency-distribution", handler.GetLatencyDistribution)

	// 性能分析
	group.GET("/slow-queries", handler.GetSlowQueries)
	group.GET("/error-interfaces", handler.GetErrorInterfaces)

	// 登录统计
	group.GET("/login-stats", handler.GetLoginStats)
}

// RegisterRoutesWithDB 注册监控模块路由（带数据库连接）
func RegisterRoutesWithDB(h *server.Hertz, userDB, auditDB *gorm.DB) {
	// 初始化服务
	monitorSvc := service.NewMonitorService()
	statsSvc := service.NewStatsService(userDB, auditDB)

	// 初始化处理器
	handler := api.NewMonitorHandler(monitorSvc, statsSvc)

	// 公开路由组（WebSocket）
	wsGroup := h.Group("/api/monitor")
	wsGroup.GET("/ws", WSHandler)

	// 需要认证的路由组
	group := h.Group("/api/monitor")
	group.Use(middleware.AuthMiddleware())

	// 仪表盘汇总
	group.GET("/dashboard", handler.GetDashboardSummary)

	// 系统监控
	group.GET("/system", handler.GetSystemStats)

	// 业务统计
	group.GET("/business", handler.GetBusinessStats)

	// 访问统计
	group.GET("/access", handler.GetAccessStats)

	// 趋势数据
	group.GET("/trend/hourly", handler.GetHourlyTrend)
	group.GET("/trend/daily", handler.GetDailyTrend)

	// TOP数据
	group.GET("/top-paths", handler.GetTopPaths)

	// 分布数据
	group.GET("/status-distribution", handler.GetStatusDistribution)
	group.GET("/latency-distribution", handler.GetLatencyDistribution)

	// 性能分析
	group.GET("/slow-queries", handler.GetSlowQueries)
	group.GET("/error-interfaces", handler.GetErrorInterfaces)

	// 登录统计
	group.GET("/login-stats", handler.GetLoginStats)
}

// WSHandler WebSocket处理函数
func WSHandler(ctx context.Context, c *app.RequestContext) {
	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		defer conn.Close()
		for {
			// 模拟实时数据
			data := map[string]interface{}{
				"timestamp":  time.Now().Format("15:04:05"),
				"requests":   1000 + rand.Intn(500),
				"qps":        300 + rand.Intn(200),
				"error_rate": float64(rand.Intn(10)) / 100.0,
				"latency":    50 + rand.Intn(100),
				"type":       "realtime",
			}

			if err := conn.WriteJSON(data); err != nil {
				break
			}
			time.Sleep(2 * time.Second)
		}
	})
	if err != nil {
		return
	}
}
