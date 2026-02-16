package api

import (
	"context"

	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/queue"
)

type AuditService interface {
	RecordDataChange(ctx context.Context, log *model.SysDataChangeLog)
}

type auditServiceImpl struct {
	queue *queue.AuditLogQueue
}

func (a *auditServiceImpl) RecordDataChange(ctx context.Context, log *model.SysDataChangeLog) {
	if a.queue != nil {
		a.queue.PushDataChange(log)
	}
}
