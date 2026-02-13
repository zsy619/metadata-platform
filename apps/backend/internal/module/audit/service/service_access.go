package service

import (
	"context"
	"database/sql"
	"fmt"
	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/utils"
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

// AccessLogService 访问日志服务
type AccessLogService struct {
	db       *gorm.DB
	auditLog *queue.AuditLogQueue
}

// NewAccessLogService 创建访问日志服务实例
func NewAccessLogService(db *gorm.DB, auditLog *queue.AuditLogQueue) *AccessLogService {
	return &AccessLogService{
		db:       db,
		auditLog: auditLog,
	}
}

// GetAccessLogs 获取访问日志列表（分页）
func (s *AccessLogService) GetAccessLogs(page, pageSize int, filters map[string]interface{}) ([]model.SysAccessLog, int64, error) {
	var logs []model.SysAccessLog
	var total int64

	query := s.db.Model(&model.SysAccessLog{})

	// 应用过滤条件
	if userID, ok := filters["user_id"]; ok && userID != "" {
		query = query.Where("user_id LIKE ?", "%"+fmt.Sprintf("%v", userID)+"%")
	}
	if path, ok := filters["path"]; ok && path != "" {
		query = query.Where("path LIKE ?", "%"+fmt.Sprintf("%v", path)+"%")
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if method, ok := filters["method"]; ok && method != "" {
		query = query.Where("method = ?", method)
	}
	if startTime, ok := filters["start_time"]; ok && startTime != "" {
		query = query.Where("create_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok && endTime != "" {
		query = query.Where("create_at <= ?", endTime)
	}
	if ip, ok := filters["client_ip"]; ok && ip != "" {
		query = query.Where("client_ip LIKE ?", "%"+fmt.Sprintf("%v", ip)+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		hlog.Error("Failed to count access logs: ", err)
		return nil, 0, err
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		query = query.Scopes(utils.Paginate(page, pageSize))
	}

	// 按创建时间倒序
	query = query.Order("create_at DESC")

	if err := query.Find(&logs).Error; err != nil {
		hlog.Error("Failed to get access logs: ", err)
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAllAccessLogs 获取所有访问日志（用于导出）
func (s *AccessLogService) GetAllAccessLogs(filters map[string]interface{}) ([]model.SysAccessLog, error) {
	var logs []model.SysAccessLog

	query := s.db.Model(&model.SysAccessLog{})

	// 应用过滤条件
	if userID, ok := filters["user_id"]; ok && userID != "" {
		query = query.Where("user_id LIKE ?", "%"+fmt.Sprintf("%v", userID)+"%")
	}
	if path, ok := filters["path"]; ok && path != "" {
		query = query.Where("path LIKE ?", "%"+fmt.Sprintf("%v", path)+"%")
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if method, ok := filters["method"]; ok && method != "" {
		query = query.Where("method = ?", method)
	}
	if startTime, ok := filters["start_time"]; ok && startTime != "" {
		query = query.Where("create_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok && endTime != "" {
		query = query.Where("create_at <= ?", endTime)
	}
	if ip, ok := filters["client_ip"]; ok && ip != "" {
		query = query.Where("client_ip LIKE ?", "%"+fmt.Sprintf("%v", ip)+"%")
	}

	// 按创建时间倒序
	query = query.Order("create_at DESC")

	if err := query.Find(&logs).Error; err != nil {
		hlog.Error("Failed to get all access logs: ", err)
		return nil, err
	}

	return logs, nil
}

// RecordAccess 记录访问日志
func (s *AccessLogService) RecordAccess(ctx context.Context, log *model.SysAccessLog) error {
	// 异步写入队列
	if s.auditLog != nil {
		s.auditLog.PushAccess(log)
		return nil
	}

	// 如果队列不可用，直接写入数据库
	return s.db.WithContext(ctx).Create(log).Error
}

// GetStatistics 获取访问统计数据
func (s *AccessLogService) GetStatistics(filters map[string]interface{}) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	baseQuery := s.db.Model(&model.SysAccessLog{})

	// 应用时间过滤
	if startTime, ok := filters["start_time"]; ok && startTime != "" {
		baseQuery = baseQuery.Where("create_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok && endTime != "" {
		baseQuery = baseQuery.Where("create_at <= ?", endTime)
	}

	// 总访问量
	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// 平均响应时间 - 单独查询，处理空表情况
	var avgLatency sql.NullFloat64
	if err := s.db.Model(&model.SysAccessLog{}).Select("AVG(latency)").Scan(&avgLatency).Error; err != nil {
		return nil, err
	}
	if avgLatency.Valid {
		stats["avg_latency"] = avgLatency.Float64
	} else {
		stats["avg_latency"] = float64(0)
	}

	// 状态码统计 - 使用新的query
	statusQuery := s.db.Model(&model.SysAccessLog{})
	if startTime, ok := filters["start_time"]; ok && startTime != "" {
		statusQuery = statusQuery.Where("create_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok && endTime != "" {
		statusQuery = statusQuery.Where("create_at <= ?", endTime)
	}
	var statusCounts []struct {
		Status int
		Count  int64
	}
	if err := statusQuery.Select("status, COUNT(*) as count").Group("status").Scan(&statusCounts).Error; err != nil {
		return nil, err
	}
	statusMap := make(map[int]int64)
	for _, sc := range statusCounts {
		statusMap[sc.Status] = sc.Count
	}
	stats["status_counts"] = statusMap

	// TOP 10 访问路径 - 使用新的query
	pathQuery := s.db.Model(&model.SysAccessLog{})
	if startTime, ok := filters["start_time"]; ok && startTime != "" {
		pathQuery = pathQuery.Where("create_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok && endTime != "" {
		pathQuery = pathQuery.Where("create_at <= ?", endTime)
	}
	var topPaths []struct {
		Path  string
		Count int64
	}
	if err := pathQuery.Select("path, COUNT(*) as count").Group("path").Order("count DESC").Limit(10).Scan(&topPaths).Error; err != nil {
		return nil, err
	}
	stats["top_paths"] = topPaths

	// 按小时统计访问量 - 使用新的query
	hourQuery := s.db.Model(&model.SysAccessLog{})
	if startTime, ok := filters["start_time"]; ok && startTime != "" {
		hourQuery = hourQuery.Where("create_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok && endTime != "" {
		hourQuery = hourQuery.Where("create_at <= ?", endTime)
	}
	var hourlyCounts []struct {
		Hour  int
		Count int64
	}
	dateFormat := "DATE_FORMAT(create_at, '%H')"
	if s.db.Dialector.Name() == "postgres" {
		dateFormat = "TO_CHAR(create_at, 'HH24')"
	}
	if err := hourQuery.Select(dateFormat+" as hour, COUNT(*) as count").Group(dateFormat).Order("hour").Scan(&hourlyCounts).Error; err != nil {
		if !strings.Contains(err.Error(), "invalid") {
			return nil, err
		}
	}
	stats["hourly_counts"] = hourlyCounts

	return stats, nil
}

// GetAbnormalAccess 获取异常访问记录
func (s *AccessLogService) GetAbnormalAccess(filters map[string]interface{}, page, pageSize int) ([]model.SysAccessLog, int64, error) {
	var logs []model.SysAccessLog
	var total int64

	query := s.db.Model(&model.SysAccessLog{})

	// 异常访问：状态码 >= 400 或 响应时间 > 5秒
	query = query.Where("status >= 400 OR latency > 5000")

	// 应用时间过滤
	if startTime, ok := filters["start_time"]; ok && startTime != "" {
		query = query.Where("create_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok && endTime != "" {
		query = query.Where("create_at <= ?", endTime)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		query = query.Scopes(utils.Paginate(page, pageSize))
	}

	query = query.Order("create_at DESC")

	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
