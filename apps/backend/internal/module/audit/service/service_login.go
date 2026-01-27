package service

import (
	"context"
	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/utils"
)

// RecordLogin 记录登录日志（异步）
func (s *auditService) RecordLogin(ctx context.Context, log *model.SysLoginLog) {
	if s.queue != nil {
		s.queue.PushLogin(log)
	}
}

// GetLoginLogs 获取登录日志列表
func (s *auditService) GetLoginLogs(page, pageSize int, filters map[string]interface{}) ([]model.SysLoginLog, int64, error) {
	var logs []model.SysLoginLog
	var total int64
	db := s.db.Model(&model.SysLoginLog{})
	db = applyFilters(db, filters)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Scopes(utils.Paginate(page, pageSize)).Order("create_at DESC").Find(&logs).Error
	return logs, total, err
}

// GetAllLoginLogs 获取所有登录日志（用于导出）
func (s *auditService) GetAllLoginLogs(filters map[string]interface{}) ([]model.SysLoginLog, error) {
	var logs []model.SysLoginLog
	db := s.db.Model(&model.SysLoginLog{})
	db = applyFilters(db, filters)
	err := db.Order("create_at DESC").Find(&logs).Error
	return logs, err
}
