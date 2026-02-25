package service

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

// BusinessStats 业务统计数据
type BusinessStats struct {
	TotalUsers      int64 `json:"total_users"`       // 总用户数
	ActiveUsers     int64 `json:"active_users"`      // 活跃用户数（7天内登录）
	OnlineUsers     int64 `json:"online_users"`      // 在线用户数（30分钟内活跃）
	TodayNewUsers   int64 `json:"today_new_users"`   // 今日新增用户
	TotalRoles      int64 `json:"total_roles"`       // 总角色数
	TotalOrgs       int64 `json:"total_orgs"`        // 总组织数
	TotalApps       int64 `json:"total_apps"`        // 总应用数
	TotalTenants    int64 `json:"total_tenants"`     // 总租户数
	TotalMenus      int64 `json:"total_menus"`       // 总菜单数
	TotalPositions  int64 `json:"total_positions"`   // 总职位数
	TotalRoleGroups int64 `json:"total_role_groups"` // 总角色组数
	TotalUserGroups int64 `json:"total_user_groups"` // 总用户组数
}

// AccessStats 访问统计数据
type AccessStats struct {
	TotalRequests int64   `json:"total_requests"` // 总请求数
	TodayRequests int64   `json:"today_requests"` // 今日请求数
	HourRequests  int64   `json:"hour_requests"`  // 最近1小时请求数
	AvgLatency    float64 `json:"avg_latency"`    // 平均响应时间(ms)
	ErrorCount    int64   `json:"error_count"`    // 错误数量
	ErrorRate     float64 `json:"error_rate"`     // 错误率
	SuccessRate   float64 `json:"success_rate"`   // 成功率
	TotalDataIn   int64   `json:"total_data_in"`  // 总输入数据量
	TotalDataOut  int64   `json:"total_data_out"` // 总输出数据量
	QPS           float64 `json:"qps"`            // 当前QPS（最近1分钟）
	PeakQPS       float64 `json:"peak_qps"`       // 峰值QPS
}

// TrendData 趋势数据
type TrendData struct {
	Time  string `json:"time"`
	Value int64  `json:"value"`
}

// HourlyStats 按小时统计
type HourlyStats struct {
	Hour  int   `json:"hour"`
	Count int64 `json:"count"`
}

// PathStats 路径统计
type PathStats struct {
	Path       string `json:"path"`
	Count      int64  `json:"count"`
	AvgLatency int64  `json:"avg_latency,omitempty"`
}

// StatusStats 状态码统计
type StatusStats struct {
	Status int   `json:"status"`
	Count  int64 `json:"count"`
}

// SlowQuery 慢查询
type SlowQuery struct {
	Query     string `json:"query"`
	Duration  int64  `json:"duration"`
	Timestamp string `json:"timestamp"`
	User      string `json:"user"`
}

// ErrorInterface 错误接口
type ErrorInterface struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	ErrorCount  int64  `json:"error_count"`
	LatestError string `json:"latest_error"`
	Timestamp   string `json:"timestamp"`
}

// LatencyDistribution 响应时间分布
type LatencyDistribution struct {
	Range string `json:"range"`
	Count int64  `json:"count"`
}

// LoginStats 登录统计
type LoginStats struct {
	TodayLogins  int64       `json:"today_logins"`  // 今日登录次数
	TodaySuccess int64       `json:"today_success"` // 今日成功登录次数
	TodayFailed  int64       `json:"today_failed"`  // 今日失败登录次数
	OnlineCount  int64       `json:"online_count"`  // 当前在线人数
	LoginTrend   []TrendData `json:"login_trend"`   // 登录趋势
}

// cachedData 缓存数据结构
type cachedData struct {
	data      interface{}
	timestamp time.Time
}

// StatsService 统计服务
type StatsService struct {
	userDB  *gorm.DB
	auditDB *gorm.DB

	// 缓存
	cache    map[string]*cachedData
	cacheMu  sync.RWMutex
	cacheTTL time.Duration
}

// NewStatsService 创建统计服务实例
func NewStatsService(userDB, auditDB *gorm.DB) *StatsService {
	return &StatsService{
		userDB:   userDB,
		auditDB:  auditDB,
		cache:    make(map[string]*cachedData),
		cacheTTL: 30 * time.Second, // 默认缓存30秒
	}
}

// getFromCache 从缓存获取数据
func (s *StatsService) getFromCache(key string) (interface{}, bool) {
	s.cacheMu.RLock()
	defer s.cacheMu.RUnlock()

	if cached, ok := s.cache[key]; ok {
		if time.Since(cached.timestamp) < s.cacheTTL {
			return cached.data, true
		}
	}
	return nil, false
}

// setCache 设置缓存
func (s *StatsService) setCache(key string, data interface{}) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	s.cache[key] = &cachedData{
		data:      data,
		timestamp: time.Now(),
	}
}

// GetBusinessStats 获取业务统计数据
func (s *StatsService) GetBusinessStats() (*BusinessStats, error) {
	// 检查缓存
	if cached, ok := s.getFromCache("business_stats"); ok {
		if stats, ok := cached.(*BusinessStats); ok {
			return stats, nil
		}
	}

	stats := &BusinessStats{}
	today := time.Now().Format("2006-01-02")
	thirtyMinutesAgo := time.Now().Add(-30 * time.Minute).Format("2006-01-02 15:04:05")

	// 使用并发查询提高性能
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 用户统计
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.userDB.Table("sso_user").Where("is_deleted = ?", false).Count(&count)
		mu.Lock()
		stats.TotalUsers = count
		mu.Unlock()
	}()

	// 活跃用户（最近7天登录）
	wg.Add(1)
	go func() {
		defer wg.Done()
		sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		var count int64
		s.auditDB.Table("sys_login_log").
			Where("DATE(create_at) >= ? AND login_status = 1", sevenDaysAgo).
			Distinct("user_id").
			Count(&count)
		mu.Lock()
		stats.ActiveUsers = count
		mu.Unlock()
	}()

	// 在线用户（30分钟内有活动）
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.auditDB.Table("sys_access_log").
			Where("create_at >= ?", thirtyMinutesAgo).
			Distinct("user_id").
			Count(&count)
		mu.Lock()
		stats.OnlineUsers = count
		mu.Unlock()
	}()

	// 今日新增用户
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.userDB.Table("sso_user").
			Where("DATE(create_at) = ? AND is_deleted = ?", today, false).
			Count(&count)
		mu.Lock()
		stats.TodayNewUsers = count
		mu.Unlock()
	}()

	// 角色统计
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.userDB.Table("sso_role").Count(&count)
		mu.Lock()
		stats.TotalRoles = count
		mu.Unlock()
	}()

	// 组织统计
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.userDB.Table("sso_org").Count(&count)
		mu.Lock()
		stats.TotalOrgs = count
		mu.Unlock()
	}()

	// 应用统计
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.userDB.Table("sso_app").Count(&count)
		mu.Lock()
		stats.TotalApps = count
		mu.Unlock()
	}()

	// 租户统计
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.userDB.Table("sso_tenant").Count(&count)
		mu.Lock()
		stats.TotalTenants = count
		mu.Unlock()
	}()

	// 菜单统计
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.userDB.Table("sso_menu").Count(&count)
		mu.Lock()
		stats.TotalMenus = count
		mu.Unlock()
	}()

	// 职位统计
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.userDB.Table("sso_pos").Count(&count)
		mu.Lock()
		stats.TotalPositions = count
		mu.Unlock()
	}()

	// 角色组统计
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.userDB.Table("sso_role_group").Count(&count)
		mu.Lock()
		stats.TotalRoleGroups = count
		mu.Unlock()
	}()

	// 用户组统计
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.userDB.Table("sso_user_group").Count(&count)
		mu.Lock()
		stats.TotalUserGroups = count
		mu.Unlock()
	}()

	wg.Wait()

	// 设置缓存
	s.setCache("business_stats", stats)

	return stats, nil
}

// GetAccessStats 获取访问统计数据
func (s *StatsService) GetAccessStats() (*AccessStats, error) {
	// 检查缓存
	if cached, ok := s.getFromCache("access_stats"); ok {
		if stats, ok := cached.(*AccessStats); ok {
			return stats, nil
		}
	}

	stats := &AccessStats{}
	today := time.Now().Format("2006-01-02")
	oneHourAgo := time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05")
	oneMinuteAgo := time.Now().Add(-1 * time.Minute).Format("2006-01-02 15:04:05")

	// 使用并发查询
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 总请求数
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.auditDB.Table("sys_access_log").Count(&count)
		mu.Lock()
		stats.TotalRequests = count
		mu.Unlock()
	}()

	// 今日请求数
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.auditDB.Table("sys_access_log").
			Where("DATE(create_at) = ?", today).
			Count(&count)
		mu.Lock()
		stats.TodayRequests = count
		mu.Unlock()
	}()

	// 最近1小时请求数
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.auditDB.Table("sys_access_log").
			Where("create_at >= ?", oneHourAgo).
			Count(&count)
		mu.Lock()
		stats.HourRequests = count
		mu.Unlock()
	}()

	// 平均响应时间
	wg.Add(1)
	go func() {
		defer wg.Done()
		var avgLatency sql.NullFloat64
		s.auditDB.Table("sys_access_log").
			Select("AVG(latency)").
			Scan(&avgLatency)
		if avgLatency.Valid {
			mu.Lock()
			stats.AvgLatency = avgLatency.Float64
			mu.Unlock()
		}
	}()

	// 错误数量（状态码 >= 400）
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.auditDB.Table("sys_access_log").
			Where("status >= 400").
			Count(&count)
		mu.Lock()
		stats.ErrorCount = count
		mu.Unlock()
	}()

	// 总输入数据量
	wg.Add(1)
	go func() {
		defer wg.Done()
		var totalIn sql.NullInt64
		s.auditDB.Table("sys_access_log").
			Select("SUM(request_size)").
			Scan(&totalIn)
		if totalIn.Valid {
			mu.Lock()
			stats.TotalDataIn = totalIn.Int64
			mu.Unlock()
		}
	}()

	// 总输出数据量
	wg.Add(1)
	go func() {
		defer wg.Done()
		var totalOut sql.NullInt64
		s.auditDB.Table("sys_access_log").
			Select("SUM(response_size)").
			Scan(&totalOut)
		if totalOut.Valid {
			mu.Lock()
			stats.TotalDataOut = totalOut.Int64
			mu.Unlock()
		}
	}()

	// 最近1分钟请求数（用于计算QPS）
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.auditDB.Table("sys_access_log").
			Where("create_at >= ?", oneMinuteAgo).
			Count(&count)
		mu.Lock()
		stats.QPS = float64(count) / 60.0
		mu.Unlock()
	}()

	// 峰值QPS（按分钟统计，取最大值）
	wg.Add(1)
	go func() {
		defer wg.Done()
		var maxCount int64
		dialect := s.auditDB.Dialector.Name()
		var minuteExpr string
		switch dialect {
		case "postgres":
			minuteExpr = "DATE_TRUNC('minute', create_at)"
		case "mysql":
			minuteExpr = "DATE_FORMAT(create_at, '%Y-%m-%d %H:%i')"
		case "sqlite":
			minuteExpr = "strftime('%Y-%m-%d %H:%M', create_at)"
		default:
			minuteExpr = "DATE_FORMAT(create_at, '%Y-%m-%d %H:%i')"
		}
		s.auditDB.Table("sys_access_log").
			Select("COUNT(*) as cnt").
			Where("DATE(create_at) = ?", today).
			Group(minuteExpr).
			Order("cnt DESC").
			Limit(1).
			Scan(&maxCount)
		mu.Lock()
		stats.PeakQPS = float64(maxCount) / 60.0
		mu.Unlock()
	}()

	wg.Wait()

	// 计算错误率和成功率
	if stats.TotalRequests > 0 {
		stats.ErrorRate = float64(stats.ErrorCount) / float64(stats.TotalRequests) * 100
		stats.SuccessRate = 100 - stats.ErrorRate
	}

	// 设置缓存
	s.setCache("access_stats", stats)

	return stats, nil
}

// GetHourlyTrend 获取按小时趋势数据
func (s *StatsService) GetHourlyTrend(date string) ([]HourlyStats, error) {
	cacheKey := fmt.Sprintf("hourly_trend_%s", date)
	if cached, ok := s.getFromCache(cacheKey); ok {
		if stats, ok := cached.([]HourlyStats); ok {
			return stats, nil
		}
	}

	var stats []HourlyStats

	query := s.auditDB.Table("sys_access_log")

	if date != "" {
		query = query.Where("DATE(create_at) = ?", date)
	} else {
		query = query.Where("DATE(create_at) = ?", time.Now().Format("2006-01-02"))
	}

	// 根据数据库类型选择日期格式化函数
	dialect := s.auditDB.Dialector.Name()
	var hourExpr string
	switch dialect {
	case "postgres":
		hourExpr = "EXTRACT(HOUR FROM create_at)::int"
	case "mysql":
		hourExpr = "HOUR(create_at)"
	case "sqlite":
		hourExpr = "strftime('%H', create_at)"
	default:
		hourExpr = "HOUR(create_at)"
	}

	query.Select(hourExpr + " as hour, COUNT(*) as count").
		Group(hourExpr).
		Order("hour").
		Scan(&stats)

	s.setCache(cacheKey, stats)

	return stats, nil
}

// GetDailyTrend 获取按日趋势数据（最近7天）
func (s *StatsService) GetDailyTrend() ([]TrendData, error) {
	if cached, ok := s.getFromCache("daily_trend"); ok {
		if stats, ok := cached.([]TrendData); ok {
			return stats, nil
		}
	}

	var stats []TrendData

	// 获取最近7天的数据
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	dialect := s.auditDB.Dialector.Name()
	var dateExpr string
	switch dialect {
	case "postgres":
		dateExpr = "DATE(create_at)"
	case "mysql":
		dateExpr = "DATE(create_at)"
	case "sqlite":
		dateExpr = "date(create_at)"
	default:
		dateExpr = "DATE(create_at)"
	}

	s.auditDB.Table("sys_access_log").
		Select(dateExpr+" as time, COUNT(*) as value").
		Where("create_at >= ?", sevenDaysAgo).
		Group(dateExpr).
		Order("time").
		Scan(&stats)

	s.setCache("daily_trend", stats)

	return stats, nil
}

// GetTopPaths 获取访问量TOP路径
func (s *StatsService) GetTopPaths(limit int) ([]PathStats, error) {
	cacheKey := fmt.Sprintf("top_paths_%d", limit)
	if cached, ok := s.getFromCache(cacheKey); ok {
		if stats, ok := cached.([]PathStats); ok {
			return stats, nil
		}
	}

	var stats []PathStats

	if limit <= 0 {
		limit = 10
	}

	s.auditDB.Table("sys_access_log").
		Select("path, COUNT(*) as count, AVG(latency) as avg_latency").
		Group("path").
		Order("count DESC").
		Limit(limit).
		Scan(&stats)

	s.setCache(cacheKey, stats)

	return stats, nil
}

// GetStatusDistribution 获取状态码分布
func (s *StatsService) GetStatusDistribution() ([]StatusStats, error) {
	if cached, ok := s.getFromCache("status_distribution"); ok {
		if stats, ok := cached.([]StatusStats); ok {
			return stats, nil
		}
	}

	var stats []StatusStats

	s.auditDB.Table("sys_access_log").
		Select("status, COUNT(*) as count").
		Group("status").
		Order("count DESC").
		Scan(&stats)

	s.setCache("status_distribution", stats)

	return stats, nil
}

// GetSlowQueries 获取慢查询列表
func (s *StatsService) GetSlowQueries(limit int) ([]SlowQuery, error) {
	cacheKey := fmt.Sprintf("slow_queries_%d", limit)
	if cached, ok := s.getFromCache(cacheKey); ok {
		if queries, ok := cached.([]SlowQuery); ok {
			return queries, nil
		}
	}

	var queries []SlowQuery

	if limit <= 0 {
		limit = 10
	}

	// 从访问日志中获取慢请求（响应时间 > 1秒）
	s.auditDB.Table("sys_access_log").
		Select("path as query, latency as duration, create_at as timestamp, user_id as user").
		Where("latency > 1000").
		Order("latency DESC").
		Limit(limit).
		Scan(&queries)

	s.setCache(cacheKey, queries)

	return queries, nil
}

// GetErrorInterfaces 获取错误接口列表
func (s *StatsService) GetErrorInterfaces(limit int) ([]ErrorInterface, error) {
	cacheKey := fmt.Sprintf("error_interfaces_%d", limit)
	if cached, ok := s.getFromCache(cacheKey); ok {
		if interfaces, ok := cached.([]ErrorInterface); ok {
			return interfaces, nil
		}
	}

	var interfaces []ErrorInterface

	if limit <= 0 {
		limit = 10
	}

	// 查询错误次数最多的接口
	dialect := s.auditDB.Dialector.Name()
	var concatExpr string
	switch dialect {
	case "postgres":
		concatExpr = "status || ' Error'"
	case "mysql":
		concatExpr = "CONCAT(status, ' Error')"
	case "sqlite":
		concatExpr = "status || ' Error'"
	default:
		concatExpr = "CONCAT(status, ' Error')"
	}

	s.auditDB.Table("sys_access_log").
		Select(fmt.Sprintf("path, method, COUNT(*) as error_count, MAX(%s) as latest_error, MAX(create_at) as timestamp", concatExpr)).
		Where("status >= 400").
		Group("path, method").
		Order("error_count DESC").
		Limit(limit).
		Scan(&interfaces)

	s.setCache(cacheKey, interfaces)

	return interfaces, nil
}

// GetLatencyDistribution 获取响应时间分布
func (s *StatsService) GetLatencyDistribution() ([]LatencyDistribution, error) {
	if cached, ok := s.getFromCache("latency_distribution"); ok {
		if distribution, ok := cached.([]LatencyDistribution); ok {
			return distribution, nil
		}
	}

	// 定义时间范围
	ranges := []struct {
		name  string
		minMs int64
		maxMs int64
	}{
		{"<10ms", 0, 10},
		{"10-50ms", 10, 50},
		{"50-100ms", 50, 100},
		{"100-500ms", 100, 500},
		{"500ms-1s", 500, 1000},
		{">1s", 1000, 999999999},
	}

	// 使用并发查询
	var wg sync.WaitGroup
	var mu sync.Mutex
	distribution := make([]LatencyDistribution, 0, len(ranges))

	for i, r := range ranges {
		wg.Add(1)
		go func(idx int, name string, minMs, maxMs int64) {
			defer wg.Done()
			var count int64
			if maxMs == 999999999 {
				s.auditDB.Table("sys_access_log").
					Where("latency >= ?", minMs).
					Count(&count)
			} else {
				s.auditDB.Table("sys_access_log").
					Where("latency >= ? AND latency < ?", minMs, maxMs).
					Count(&count)
			}
			mu.Lock()
			distribution = append(distribution, LatencyDistribution{
				Range: name,
				Count: count,
			})
			mu.Unlock()
		}(i, r.name, r.minMs, r.maxMs)
	}

	wg.Wait()

	// 按顺序排列
	result := make([]LatencyDistribution, len(ranges))
	for i, r := range ranges {
		for _, d := range distribution {
			if d.Range == r.name {
				result[i] = d
				break
			}
		}
	}

	s.setCache("latency_distribution", result)

	return result, nil
}

// GetLoginStats 获取登录统计
func (s *StatsService) GetLoginStats() (*LoginStats, error) {
	if cached, ok := s.getFromCache("login_stats"); ok {
		if stats, ok := cached.(*LoginStats); ok {
			return stats, nil
		}
	}

	stats := &LoginStats{}
	today := time.Now().Format("2006-01-02")
	thirtyMinutesAgo := time.Now().Add(-30 * time.Minute).Format("2006-01-02 15:04:05")

	// 使用并发查询
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 今日登录次数
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.auditDB.Table("sys_login_log").
			Where("DATE(create_at) = ?", today).
			Count(&count)
		mu.Lock()
		stats.TodayLogins = count
		mu.Unlock()
	}()

	// 今日成功登录次数 (login_status = 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.auditDB.Table("sys_login_log").
			Where("DATE(create_at) = ? AND login_status = 1", today).
			Count(&count)
		mu.Lock()
		stats.TodaySuccess = count
		mu.Unlock()
	}()

	// 今日失败登录次数 (login_status = 0)
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.auditDB.Table("sys_login_log").
			Where("DATE(create_at) = ? AND login_status = 0", today).
			Count(&count)
		mu.Lock()
		stats.TodayFailed = count
		mu.Unlock()
	}()

	// 当前在线人数（30分钟内有访问）
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		s.auditDB.Table("sys_access_log").
			Where("create_at >= ?", thirtyMinutesAgo).
			Distinct("user_id").
			Count(&count)
		mu.Lock()
		stats.OnlineCount = count
		mu.Unlock()
	}()

	wg.Wait()

	// 最近7天登录趋势
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	var loginTrend []TrendData
	s.auditDB.Table("sys_login_log").
		Select("DATE(create_at) as time, COUNT(*) as value").
		Where("create_at >= ? AND login_status = 1", sevenDaysAgo).
		Group("DATE(create_at)").
		Order("time").
		Scan(&loginTrend)
	stats.LoginTrend = loginTrend

	s.setCache("login_stats", stats)

	return stats, nil
}

// ClearCache 清除缓存
func (s *StatsService) ClearCache() {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	s.cache = make(map[string]*cachedData)
}

// SetCacheTTL 设置缓存过期时间
func (s *StatsService) SetCacheTTL(ttl time.Duration) {
	s.cacheTTL = ttl
}
