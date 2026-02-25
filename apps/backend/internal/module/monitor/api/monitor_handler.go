package api

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"metadata-platform/internal/module/monitor/service"
)

// MonitorHandler 监控处理器
type MonitorHandler struct {
	monitorSvc *service.MonitorService
	statsSvc   *service.StatsService
}

// NewMonitorHandler 创建监控处理器实例
func NewMonitorHandler(monitorSvc *service.MonitorService, statsSvc *service.StatsService) *MonitorHandler {
	return &MonitorHandler{
		monitorSvc: monitorSvc,
		statsSvc:   statsSvc,
	}
}

// parseLimit 解析 limit 参数，默认值为 10，最大值为 100
func parseLimit(c *app.RequestContext, defaultLimit int) int {
	if defaultLimit <= 0 {
		defaultLimit = 10
	}

	limit := defaultLimit
	if limitStr := string(c.Query("limit")); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			if v > 100 {
				limit = 100
			} else {
				limit = v
			}
		}
	}
	return limit
}

// GetSystemStats 获取系统统计数据
// @Summary 获取系统统计数据
// @Description 获取CPU、内存、磁盘等系统资源使用情况
// @Tags 监控
// @Accept json
// @Produce json
// @Success 200 {object} service.SystemStats
// @Router /api/monitor/system [get]
func (h *MonitorHandler) GetSystemStats(ctx context.Context, c *app.RequestContext) {
	stats, err := h.monitorSvc.GetSystemStats()
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get system stats: %v", err)
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取系统统计数据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}

// GetBusinessStats 获取业务统计数据
// @Summary 获取业务统计数据
// @Description 获取用户、角色、组织等业务数据统计
// @Tags 监控
// @Accept json
// @Produce json
// @Success 200 {object} service.BusinessStats
// @Router /api/monitor/business [get]
func (h *MonitorHandler) GetBusinessStats(ctx context.Context, c *app.RequestContext) {
	stats, err := h.statsSvc.GetBusinessStats()
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get business stats: %v", err)
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取业务统计数据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}

// GetAccessStats 获取访问统计数据
// @Summary 获取访问统计数据
// @Description 获取请求量、响应时间、错误率等访问统计
// @Tags 监控
// @Accept json
// @Produce json
// @Success 200 {object} service.AccessStats
// @Router /api/monitor/access [get]
func (h *MonitorHandler) GetAccessStats(ctx context.Context, c *app.RequestContext) {
	stats, err := h.statsSvc.GetAccessStats()
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get access stats: %v", err)
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取访问统计数据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}

// GetHourlyTrend 获取按小时趋势数据
// @Summary 获取按小时趋势数据
// @Description 获取指定日期的按小时访问趋势
// @Tags 监控
// @Accept json
// @Produce json
// @Param date query string false "日期 (YYYY-MM-DD)"
// @Success 200 {array} service.HourlyStats
// @Router /api/monitor/trend/hourly [get]
func (h *MonitorHandler) GetHourlyTrend(ctx context.Context, c *app.RequestContext) {
	date := string(c.Query("date"))

	stats, err := h.statsSvc.GetHourlyTrend(date)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get hourly trend: %v", err)
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取小时趋势数据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}

// GetDailyTrend 获取按日趋势数据
// @Summary 获取按日趋势数据
// @Description 获取最近7天的访问趋势
// @Tags 监控
// @Accept json
// @Produce json
// @Success 200 {array} service.TrendData
// @Router /api/monitor/trend/daily [get]
func (h *MonitorHandler) GetDailyTrend(ctx context.Context, c *app.RequestContext) {
	stats, err := h.statsSvc.GetDailyTrend()
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get daily trend: %v", err)
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取日趋势数据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}

// GetTopPaths 获取访问量TOP路径
// @Summary 获取访问量TOP路径
// @Description 获取访问量最高的接口路径
// @Tags 监控
// @Accept json
// @Produce json
// @Param limit query int false "返回数量限制" default(10)
// @Success 200 {array} service.PathStats
// @Router /api/monitor/top-paths [get]
func (h *MonitorHandler) GetTopPaths(ctx context.Context, c *app.RequestContext) {
	limit := parseLimit(c, 10)

	stats, err := h.statsSvc.GetTopPaths(limit)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get top paths: %v", err)
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取TOP路径数据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}

// GetStatusDistribution 获取状态码分布
// @Summary 获取状态码分布
// @Description 获取HTTP状态码分布情况
// @Tags 监控
// @Accept json
// @Produce json
// @Success 200 {array} service.StatusStats
// @Router /api/monitor/status-distribution [get]
func (h *MonitorHandler) GetStatusDistribution(ctx context.Context, c *app.RequestContext) {
	stats, err := h.statsSvc.GetStatusDistribution()
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get status distribution: %v", err)
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取状态码分布数据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}

// GetSlowQueries 获取慢查询列表
// @Summary 获取慢查询列表
// @Description 获取响应时间最慢的请求列表
// @Tags 监控
// @Accept json
// @Produce json
// @Param limit query int false "返回数量限制" default(10)
// @Success 200 {array} service.SlowQuery
// @Router /api/monitor/slow-queries [get]
func (h *MonitorHandler) GetSlowQueries(ctx context.Context, c *app.RequestContext) {
	limit := parseLimit(c, 10)

	queries, err := h.statsSvc.GetSlowQueries(limit)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get slow queries: %v", err)
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取慢查询数据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    queries,
	})
}

// GetErrorInterfaces 获取错误接口列表
// @Summary 获取错误接口列表
// @Description 获取错误次数最多的接口列表
// @Tags 监控
// @Accept json
// @Produce json
// @Param limit query int false "返回数量限制" default(10)
// @Success 200 {array} service.ErrorInterface
// @Router /api/monitor/error-interfaces [get]
func (h *MonitorHandler) GetErrorInterfaces(ctx context.Context, c *app.RequestContext) {
	limit := parseLimit(c, 10)

	interfaces, err := h.statsSvc.GetErrorInterfaces(limit)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get error interfaces: %v", err)
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取错误接口数据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    interfaces,
	})
}

// GetLatencyDistribution 获取响应时间分布
// @Summary 获取响应时间分布
// @Description 获取API响应时间分布情况
// @Tags 监控
// @Accept json
// @Produce json
// @Success 200 {array} service.LatencyDistribution
// @Router /api/monitor/latency-distribution [get]
func (h *MonitorHandler) GetLatencyDistribution(ctx context.Context, c *app.RequestContext) {
	distribution, err := h.statsSvc.GetLatencyDistribution()
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get latency distribution: %v", err)
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取响应时间分布数据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    distribution,
	})
}

// GetLoginStats 获取登录统计
// @Summary 获取登录统计
// @Description 获取登录相关统计数据
// @Tags 监控
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/monitor/login-stats [get]
func (h *MonitorHandler) GetLoginStats(ctx context.Context, c *app.RequestContext) {
	stats, err := h.statsSvc.GetLoginStats()
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get login stats: %v", err)
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取登录统计数据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}

// GetDashboardSummary 获取仪表盘汇总数据
// @Summary 获取仪表盘汇总数据
// @Description 获取仪表盘所需的所有汇总数据
// @Tags 监控
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/monitor/dashboard [get]
func (h *MonitorHandler) GetDashboardSummary(ctx context.Context, c *app.RequestContext) {
	result := make(map[string]interface{})

	// 并发获取各项数据，提高性能
	type statsResult struct {
		key   string
		value interface{}
	}

	resultChan := make(chan statsResult, 6)

	// 并发获取数据
	go func() {
		systemStats, err := h.monitorSvc.GetSystemStats()
		if err != nil {
			hlog.CtxWarnf(ctx, "Failed to get system stats: %v", err)
		}
		resultChan <- statsResult{key: "system", value: systemStats}
	}()

	go func() {
		businessStats, err := h.statsSvc.GetBusinessStats()
		if err != nil {
			hlog.CtxWarnf(ctx, "Failed to get business stats: %v", err)
		}
		resultChan <- statsResult{key: "business", value: businessStats}
	}()

	go func() {
		accessStats, err := h.statsSvc.GetAccessStats()
		if err != nil {
			hlog.CtxWarnf(ctx, "Failed to get access stats: %v", err)
		}
		resultChan <- statsResult{key: "access", value: accessStats}
	}()

	go func() {
		loginStats, err := h.statsSvc.GetLoginStats()
		if err != nil {
			hlog.CtxWarnf(ctx, "Failed to get login stats: %v", err)
		}
		resultChan <- statsResult{key: "login", value: loginStats}
	}()

	go func() {
		dailyTrend, err := h.statsSvc.GetDailyTrend()
		if err != nil {
			hlog.CtxWarnf(ctx, "Failed to get daily trend: %v", err)
		}
		resultChan <- statsResult{key: "daily_trend", value: dailyTrend}
	}()

	go func() {
		statusDist, err := h.statsSvc.GetStatusDistribution()
		if err != nil {
			hlog.CtxWarnf(ctx, "Failed to get status distribution: %v", err)
		}
		resultChan <- statsResult{key: "status_distribution", value: statusDist}
	}()

	// 收集结果
	for i := 0; i < 6; i++ {
		r := <-resultChan
		result[r.key] = r.value
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    result,
	})
}
