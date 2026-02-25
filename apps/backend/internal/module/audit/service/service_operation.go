package service

import (
	"context"
	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/utils"
)

// RecordOperation 记录操作日志（异步）
func (s *auditService) RecordOperation(ctx context.Context, log *model.SysOperationLog) {
	if s.queue != nil {
		s.queue.PushOperation(log)
	}
}

// GetOperationLogs 获取操作日志列表
func (s *auditService) GetOperationLogs(page, pageSize int, filters map[string]interface{}) ([]model.SysOperationLog, int64, error) {
	var logs []model.SysOperationLog
	var total int64
	db := s.db.Model(&model.SysOperationLog{})
	db = applyFilters(db, filters)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Scopes(utils.Paginate(page, pageSize)).Order("create_at DESC").Find(&logs).Error
	return logs, total, err
}

// GetAllOperationLogs 获取所有操作日志（用于导出）
func (s *auditService) GetAllOperationLogs(filters map[string]interface{}) ([]model.SysOperationLog, error) {
	var logs []model.SysOperationLog
	db := s.db.Model(&model.SysOperationLog{})
	db = applyFilters(db, filters)
	err := db.Order("create_at DESC").Find(&logs).Error
	return logs, err
}

// GetRecentOperationLogs 获取最近的操作日志
func (s *auditService) GetRecentOperationLogs(limit int) ([]model.SysOperationLog, error) {
	var logs []model.SysOperationLog
	err := s.db.Model(&model.SysOperationLog{}).Order("create_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}
