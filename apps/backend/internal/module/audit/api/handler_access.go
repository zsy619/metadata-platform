package api

import (
	"context"
	"metadata-platform/internal/module/audit/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

// AccessLogHandler 访问日志处理器
type AccessLogHandler struct {
	service *service.AccessLogService
}

// NewAccessLogHandler 创建访问日志处理器
func NewAccessLogHandler(svc *service.AccessLogService) *AccessLogHandler {
	return &AccessLogHandler{
		service: svc,
	}
}

// GetAccessLogs 获取访问日志
func (h *AccessLogHandler) GetAccessLogs(ctx context.Context, c *app.RequestContext) {
	page, pageSize := utils.GetPaginationParams(c)
	filters := getFilters(c)
	logs, total, err := h.service.GetAccessLogs(page, pageSize, filters)
	if err != nil {
		utils.InternalServerErrorResponse(c, "获取访问日志失败")
		return
	}
	utils.SuccessWithPagination(c, logs, total, page, pageSize)
}

// ExportAccessLogs 导出访问日志
func (h *AccessLogHandler) ExportAccessLogs(ctx context.Context, c *app.RequestContext) {
	filters := getFilters(c)
	logs, err := h.service.GetAllAccessLogs(filters)
	if err != nil {
		utils.InternalServerErrorResponse(c, "导出访问日志失败")
		return
	}
	utils.SuccessResponse(c, logs)
}

// GetAccessStatistics 获取访问统计
func (h *AccessLogHandler) GetAccessStatistics(ctx context.Context, c *app.RequestContext) {
	filters := getFilters(c)
	stats, err := h.service.GetStatistics(filters)
	if err != nil {
		utils.InternalServerErrorResponse(c, "获取访问统计失败")
		return
	}
	utils.SuccessResponse(c, stats)
}

// GetAbnormalAccess 获取异常访问记录
func (h *AccessLogHandler) GetAbnormalAccess(ctx context.Context, c *app.RequestContext) {
	page, pageSize := utils.GetPaginationParams(c)
	filters := getFilters(c)
	logs, total, err := h.service.GetAbnormalAccess(filters, page, pageSize)
	if err != nil {
		utils.InternalServerErrorResponse(c, "获取异常访问记录失败")
		return
	}
	utils.SuccessWithPagination(c, logs, total, page, pageSize)
}
