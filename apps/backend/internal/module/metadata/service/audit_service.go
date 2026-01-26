package service

import (
	"context"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"
)

// AuditService 审计日志服务接口
type AuditService interface {
	RecordOperation(ctx context.Context, log *model.SysOperationLog)
	RecordDataChange(ctx context.Context, log *model.SysDataChangeLog)
}

type auditService struct {
	db *gorm.DB
}

// NewAuditService 创建审计日志服务实例
func NewAuditService(db *gorm.DB) AuditService {
	return &auditService{db: db}
}

// RecordOperation 记录操作日志（异步）
func (s *auditService) RecordOperation(ctx context.Context, log *model.SysOperationLog) {
	// 使用 goroutine 异步记录，避免阻塞主流程
	go func(l *model.SysOperationLog) {
		// 注意：这里的 ctx 如果是 request context，可能在 goroutine 执行时已经 cancel
		// 所以最好只用 context 传参，或者创建一个 detached context
		// 简单起见，这里忽略 context timeout
		if err := s.db.Create(l).Error; err != nil {
			utils.SugarLogger.Errorf("Failed to record operation log: %v", err)
		}
	}(log)
}

// RecordDataChange 记录数据变更日志（异步）
func (s *auditService) RecordDataChange(ctx context.Context, log *model.SysDataChangeLog) {
	go func(l *model.SysDataChangeLog) {
		if err := s.db.Create(l).Error; err != nil {
			utils.SugarLogger.Errorf("Failed to record data change log: %v", err)
		}
	}(log)
}
