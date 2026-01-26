package middleware

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// CORSMiddleware 跨域中间件
func CORSMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		// 处理OPTIONS请求
		if string(ctx.Method()) == "OPTIONS" {
			ctx.AbortWithStatus(consts.StatusNoContent)
			return
		}

		ctx.Next(c)
	}
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		ctx.Next(c)

		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := string(ctx.Method())
		// 请求路由
		reqURI := string(ctx.Request.URI().Path())
		// 状态码
		statusCode := ctx.Response.StatusCode()
		// 请求IP
		clientIP := ctx.ClientIP()

		// 日志格式
		log.Printf("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqURI,
		)
	}
}

// AuthMiddleware 认证中间件
func AuthMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 从请求头中获取Authorization
		authHeader := string(ctx.GetHeader("Authorization"))
		if authHeader == "" {
			ctx.JSON(consts.StatusUnauthorized, map[string]any{
				"code":    401,
				"message": "未提供Authorization头",
			})
			ctx.Abort()
			return
		}

		// 检查Authorization格式，应为"Bearer {token}"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(consts.StatusUnauthorized, map[string]any{
				"code":    401,
				"message": "Authorization格式不正确",
			})
			ctx.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]any{
				"code":    401,
				"message": "无效的token: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		// 将用户信息存入上下文
		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Set("is_admin", claims.IsAdmin)

		// 自动续约逻辑：如果 token 快过期了，生成个新的在 Header 里带回去
		if utils.ShouldRefresh(claims) {
			newToken, err := utils.GenerateToken(claims.UserID, claims.Username, claims.IsAdmin)
			if err == nil {
				ctx.Header("New-Token", newToken)
				ctx.Header("Access-Control-Expose-Headers", "New-Token")
			}
		}

		ctx.Next(c)
	}
}

// PermissionMiddleware 权限中间件
func PermissionMiddleware(permission string) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 从上下文中获取用户信息
		isAdmin, exists := ctx.Get("is_admin")
		if !exists {
			ctx.JSON(consts.StatusUnauthorized, map[string]any{
				"code":    401,
				"message": "未认证",
			})
			ctx.Abort()
			return
		}

		// 管理员拥有所有权限
		if isAdmin.(bool) {
			ctx.Next(c)
			return
		}

		// TODO: 实现权限检查逻辑
		// 这里简化处理，实际需要从数据库中查询用户权限
		// 目前直接返回403，后续需要完善
		ctx.JSON(consts.StatusForbidden, map[string]any{
			"code":    403,
			"message": "没有权限访问该资源",
		})
		ctx.Abort()
	}
}

// TenantMiddleware 租户中间件
func TenantMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 从请求头中获取租户ID
		tenantIDStr := string(ctx.GetHeader("X-Tenant-ID"))
		if tenantIDStr == "" {
			// 如果请求头中没有租户ID，默认使用1作为租户ID
			tenantIDStr = "1"
		}

		// 解析租户ID
		tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
		if err != nil {
			ctx.JSON(consts.StatusBadRequest, map[string]any{
				"code":    400,
				"message": "无效的租户ID",
			})
			ctx.Abort()
			return
		}

		// 将租户ID存入上下文
		ctx.Set("tenant_id", uint(tenantID))

		ctx.Next(c)
	}
}