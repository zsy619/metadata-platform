package api

import (
	"context"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

// GetOperationLogs 获取操作日志
func (h *AuditHandler) GetOperationLogs(ctx context.Context, c *app.RequestContext) {
	page, pageSize := utils.GetPaginationParams(c)
	filters := getFilters(c)
	logs, total, err := h.service.GetOperationLogs(page, pageSize, filters)
	if err != nil {
		utils.InternalServerErrorResponse(c, "获取操作日志失败")
		return
	}
	utils.SuccessWithPagination(c, logs, total, page, pageSize)
}

// ExportOperationLogs 导出操作日志
func (h *AuditHandler) ExportOperationLogs(ctx context.Context, c *app.RequestContext) {
	filters := getFilters(c)
	logs, err := h.service.GetAllOperationLogs(filters)
	if err != nil {
		utils.InternalServerErrorResponse(c, "导出操作日志失败")
		return
	}
	utils.SuccessResponse(c, logs)
}
