package utils

import (
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Response 统一响应结构
type Response struct {
	Code    int    `json:"code"`               // 状态码
	Message string `json:"message"`            // 响应信息
	Data    any    `json:"data"`               // 响应数据
	TraceID string `json:"trace_id,omitempty"` // 链路追踪ID
}

// Pagination 分页响应结构
type Pagination struct {
	Total    int64 `json:"total"`     // 总条数
	Page     int   `json:"page"`      // 当前页码
	PageSize int   `json:"page_size"` // 每页条数
	List     any   `json:"list"`      // 数据列表
}

// SuccessResponse 返回成功响应
func SuccessResponse(c *app.RequestContext, data any) {
	c.JSON(consts.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithPagination 返回带分页的成功响应
func SuccessWithPagination(c *app.RequestContext, list any, total int64, page, pageSize int) {
	c.JSON(consts.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data: Pagination{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			List:     list,
		},
	})
}

// ErrorResponse 返回错误响应
func ErrorResponse(c *app.RequestContext, code int, message string) {
	c.JSON(consts.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequestResponse 返回400错误响应
func BadRequestResponse(c *app.RequestContext, message string) {
	c.JSON(consts.StatusBadRequest, Response{
		Code:    400,
		Message: message,
		Data:    nil,
	})
}

// UnauthorizedResponse 返回401错误响应
func UnauthorizedResponse(c *app.RequestContext, message string) {
	c.JSON(consts.StatusUnauthorized, Response{
		Code:    401,
		Message: message,
		Data:    nil,
	})
}

// ForbiddenResponse 返回403错误响应
func ForbiddenResponse(c *app.RequestContext, message string) {
	c.JSON(consts.StatusForbidden, Response{
		Code:    403,
		Message: message,
		Data:    nil,
	})
}

// NotFoundResponse 返回404错误响应
func NotFoundResponse(c *app.RequestContext, message string) {
	c.JSON(consts.StatusNotFound, Response{
		Code:    404,
		Message: message,
		Data:    nil,
	})
}

// InternalServerErrorResponse 返回500错误响应
func InternalServerErrorResponse(c *app.RequestContext, message string) {
	c.JSON(consts.StatusInternalServerError, Response{
		Code:    500,
		Message: message,
		Data:    nil,
	})
}

// GetPaginationParams 从请求中获取分页参数
func GetPaginationParams(c *app.RequestContext) (page, pageSize int) {
	// 默认值
	page = 1
	pageSize = 20

	// 从查询参数中获取
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	return
}

// BindAndValidate 绑定请求体并验证
func BindAndValidate(c *app.RequestContext, obj any) bool {
	if err := c.Bind(obj); err != nil {
		BadRequestResponse(c, "Invalid request body: "+err.Error())
		return false
	}
	return true
}
