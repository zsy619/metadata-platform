package middleware

import (
	"context"
	"fmt"
	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/service"
	"metadata-platform/internal/utils"
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
		userAccount := string(ctx.Request.Header.Get("X-User-Account"))
		fmt.Println(userID, userAccount)

		tenantID := string(ctx.Request.Header.Get("X-Tenant-ID"))
		if tenantID == "" {
			tenantID = "1"
		}

		// 解析客户端信息
		clientInfo := utils.ParseUserAgent(
			string(ctx.Request.Header.UserAgent()),
			ctx.ClientIP(),
			string(ctx.Request.Header.Get("Accept-Language")),
		)

		log := &model.SysOperationLog{
			TraceID:        traceID,
			UserID:         userID,
			TenantID:       tenantID,
			Method:         string(ctx.Request.Method()),
			Path:           string(ctx.Request.URI().Path()),
			Status:         statusCode,
			Latency:        latency,
			ClientIP:       clientInfo.IP,
			UserAgent:      clientInfo.UserAgent,
			Browser:        clientInfo.Browser,
			BrowserVersion: clientInfo.BrowserVersion,
			OS:             clientInfo.OS,
			OSVersion:      clientInfo.OSVersion,
			DeviceType:     clientInfo.DeviceType,
			Language:       clientInfo.Language,
			Platform:       clientInfo.Platform,
			Source:         "metadata", // 明确标识来源为 metadata 模块
			CreateAt:       time.Now(),
		}

		// 记录错误信息 (如果有)
		if statusCode >= 400 {
			log.ErrorMessage = fmtResponseBody(ctx)
		}

		auditSvc.RecordOperation(c, log)
	}
}

func fmtResponseBody(ctx *app.RequestContext) string {
	// 注意：读取 Response Body 可能会影响性能，仅在错误时读取
	return string(ctx.Response.Body())
}
