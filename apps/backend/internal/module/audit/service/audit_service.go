package service

import (
	"context"
	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/queue"

	"gorm.io/gorm"
)

// AuditService 审计日志服务接口
type AuditService interface {
	RecordOperation(ctx context.Context, log *model.SysOperationLog)
	RecordDataChange(ctx context.Context, log *model.SysDataChangeLog)
	RecordLogin(ctx context.Context, log *model.SysLoginLog)

	GetLoginLogs(page, pageSize int, filters map[string]interface{}) ([]model.SysLoginLog, int64, error)
	GetOperationLogs(page, pageSize int, filters map[string]interface{}) ([]model.SysOperationLog, int64, error)
	GetDataChangeLogs(page, pageSize int, filters map[string]interface{}) ([]model.SysDataChangeLog, int64, error)

	GetAllLoginLogs(filters map[string]interface{}) ([]model.SysLoginLog, error)
	GetAllOperationLogs(filters map[string]interface{}) ([]model.SysOperationLog, error)
	GetAllDataChangeLogs(filters map[string]interface{}) ([]model.SysDataChangeLog, error)
}

type auditService struct {
	db    *gorm.DB
	queue *queue.AuditLogQueue
}

// NewAuditService 创建审计日志服务实例
func NewAuditService(db *gorm.DB, q *queue.AuditLogQueue) AuditService {
	return &auditService{
		db:    db,
		queue: q,
	}
}

func applyFilters(db *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for k, v := range filters {
		if v != nil && v != "" {
			switch k {
			case "start_time":
				db = db.Where("create_at >= ?", v)
			case "end_time":
				db = db.Where("create_at <= ?", v)
			case "source", "user_id", "trace_id", "module", "status", "action", "model_id", "record_id", "account", "method", "path", "client_ip", "login_status":
				// Only allow known safe columns to be passed directly to Where
				db = db.Where(k+" = ?", v)
			}
		}
	}
	return db
}
