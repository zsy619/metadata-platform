package api

import (
	"context"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

// GetLoginLogs 获取登录日志
func (h *AuditHandler) GetLoginLogs(ctx context.Context, c *app.RequestContext) {
	page, pageSize := utils.GetPaginationParams(c)
	filters := getFilters(c)
	logs, total, err := h.service.GetLoginLogs(page, pageSize, filters)
	if err != nil {
		utils.InternalServerErrorResponse(c, "获取登录日志失败")
		return
	}
	utils.SuccessWithPagination(c, logs, total, page, pageSize)
}

// ExportLoginLogs 导出登录日志
func (h *AuditHandler) ExportLoginLogs(ctx context.Context, c *app.RequestContext) {
	filters := getFilters(c)
	logs, err := h.service.GetAllLoginLogs(filters)
	if err != nil {
		utils.InternalServerErrorResponse(c, "导出登录日志失败")
		return
	}
	// TODO: Generate Excel/CSV
	utils.SuccessResponse(c, logs) // Temporary return JSON
}
