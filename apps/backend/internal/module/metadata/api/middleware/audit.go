package middleware

import (
	"context"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/service"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
)

// AuditMiddleware 审计日志中间件
func AuditMiddleware(auditSvc service.AuditService) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		start := time.Now()
		
		// 生成 TraceID
		traceID := uuid.New().String()
		ctx.Set("trace_id", traceID)

		// 执行下一个 handler
		ctx.Next(c)

		latency := time.Since(start).Milliseconds()
		statusCode := ctx.Response.StatusCode()
		
		// 暂时简单的 UserID 获取 (后续应结合 JWT Middleware)
		userID := string(ctx.Request.Header.Get("X-User-ID"))
		if userID == "" {
			userID = "anonymous"
		}
		
		tenantID := string(ctx.Request.Header.Get("X-Tenant-ID"))

		log := &model.SysOperationLog{
			TraceID:   traceID,
			UserID:    userID,
			TenantID:  tenantID,
			Method:    string(ctx.Request.Method()),
			Path:      string(ctx.Request.URI().Path()),
			Status:    statusCode,
			Latency:   latency,
			ClientIP:  ctx.ClientIP(),
			UserAgent: string(ctx.Request.Header.UserAgent()),
			CreateAt:  time.Now(),
		}

		// 记录错误信息 (如果有)
		// Hertz 的 Error 机制可能不在 ctx.Errors 中直接体现业务错误，
		// 这里假设如果有 panic recovery 或者 utils.ErrorResponse 可能会留下痕迹
		// 暂时只记录 HTTP 状态码非 200 的情况
		if statusCode >= 400 {
			// 尝试获取 body 作为 error message? 
			// 或者是从 ctx 的 keys 中获取 error?
			// 这里简单记录
			log.ErrorMessage = fmtResponseBody(ctx)
		}

		auditSvc.RecordOperation(c, log)
	}
}

func fmtResponseBody(ctx *app.RequestContext) string {
	// 注意：读取 Response Body 可能会影响性能，仅在错误时读取
	return string(ctx.Response.Body())
}
