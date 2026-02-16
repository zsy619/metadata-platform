package middleware

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"metadata-platform/internal/utils"
)

// RecoveryMiddleware 全局异常捕获中间件
func RecoveryMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stack := string(debug.Stack())
				utils.SugarLogger.Errorf("Panic recovered: %v\nStack: %s", err, stack)

				// 根据错误类型返回不同的响应
				switch e := err.(type) {
				case *utils.AppError:
					// 自定义应用错误
					statusCode := getStatusCodeFromErrorCode(e.Code)
					ctx.JSON(statusCode, map[string]any{
						"code":    int(e.Code),
						"message": e.Message,
						"data":    nil,
					})
				default:
					// 未知错误
					ctx.JSON(consts.StatusInternalServerError, map[string]any{
						"code":    500,
						"message": "Internal server error",
						"data":    nil,
					})
				}
			}
		}()

		ctx.Next(c)
	}
}

// getStatusCodeFromErrorCode 根据错误码获取HTTP状态码
func getStatusCodeFromErrorCode(code utils.ErrorCode) int {
	switch {
	case code >= 400000 && code < 500000:
		// 业务错误
		return http.StatusBadRequest
	case code >= 401000 && code < 402000:
		// 认证错误
		return http.StatusUnauthorized
	case code >= 500000 && code < 600000:
		// 系统错误
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
