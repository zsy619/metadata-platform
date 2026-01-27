package service

import (
	"context"
	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/utils"
)

// RecordDataChange 记录数据变更日志（异步）
func (s *auditService) RecordDataChange(ctx context.Context, log *model.SysDataChangeLog) {
	if s.queue != nil {
		s.queue.PushDataChange(log)
	}
}

// GetDataChangeLogs 获取数据变更日志列表
func (s *auditService) GetDataChangeLogs(page, pageSize int, filters map[string]interface{}) ([]model.SysDataChangeLog, int64, error) {
	var logs []model.SysDataChangeLog
	var total int64
	db := s.db.Model(&model.SysDataChangeLog{})
	db = applyFilters(db, filters)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Scopes(utils.Paginate(page, pageSize)).Order("create_time DESC").Find(&logs).Error
	return logs, total, err
}

// GetAllDataChangeLogs 获取所有数据变更日志（用于导出）
func (s *auditService) GetAllDataChangeLogs(filters map[string]interface{}) ([]model.SysDataChangeLog, error) {
	var logs []model.SysDataChangeLog
	db := s.db.Model(&model.SysDataChangeLog{})
	db = applyFilters(db, filters)
	err := db.Order("create_time DESC").Find(&logs).Error
	return logs, err
}
