package api

import (
	"context"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

// GetDataChangeLogs 获取数据变更日志
func (h *AuditHandler) GetDataChangeLogs(ctx context.Context, c *app.RequestContext) {
	page, pageSize := utils.GetPaginationParams(c)
	filters := getFilters(c)
	logs, total, err := h.service.GetDataChangeLogs(page, pageSize, filters)
	if err != nil {
		utils.InternalServerErrorResponse(c, "获取数据变更日志失败")
		return
	}
	utils.SuccessWithPagination(c, logs, total, page, pageSize)
}

// ExportDataChangeLogs 导出数据变更日志
func (h *AuditHandler) ExportDataChangeLogs(ctx context.Context, c *app.RequestContext) {
	filters := getFilters(c)
	logs, err := h.service.GetAllDataChangeLogs(filters)
	if err != nil {
		utils.InternalServerErrorResponse(c, "导出数据变更日志失败")
		return
	}
	utils.SuccessResponse(c, logs)
}
