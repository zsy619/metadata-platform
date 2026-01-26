package api

import (
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

// SuccessResponse 返回成功响应
func SuccessResponse(c *app.RequestContext, data any) {
	utils.SuccessResponse(c, data)
}

// SuccessWithPagination 返回带分页的成功响应
func SuccessWithPagination(c *app.RequestContext, list any, total int64, page, pageSize int) {
	utils.SuccessWithPagination(c, list, total, page, pageSize)
}

// ErrorResponse 返回错误响应
func ErrorResponse(c *app.RequestContext, code int, message string) {
	utils.ErrorResponse(c, code, message)
}

// BadRequestResponse 返回400错误响应
func BadRequestResponse(c *app.RequestContext, message string) {
	utils.BadRequestResponse(c, message)
}

// UnauthorizedResponse 返回401错误响应
func UnauthorizedResponse(c *app.RequestContext, message string) {
	utils.UnauthorizedResponse(c, message)
}

// ForbiddenResponse 返回403错误响应
func ForbiddenResponse(c *app.RequestContext, message string) {
	utils.ForbiddenResponse(c, message)
}

// NotFoundResponse 返回404错误响应
func NotFoundResponse(c *app.RequestContext, message string) {
	utils.NotFoundResponse(c, message)
}

// InternalServerErrorResponse 返回500错误响应
func InternalServerErrorResponse(c *app.RequestContext, message string) {
	utils.InternalServerErrorResponse(c, message)
}

// GetPaginationParams 从请求中获取分页参数
func GetPaginationParams(c *app.RequestContext) (page, pageSize int) {
	return utils.GetPaginationParams(c)
}

// BindAndValidate 绑定请求体并验证
func BindAndValidate(c *app.RequestContext, obj any) bool {
	return utils.BindAndValidate(c, obj)
}
