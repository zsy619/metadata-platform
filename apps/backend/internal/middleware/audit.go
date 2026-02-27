package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"

	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/service"
	"metadata-platform/internal/utils"
)

func AuditMiddleware(auditSvc service.AuditService, source string) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		start := time.Now()

		traceID := string(ctx.Request.Header.Get("X-Trace-ID"))
		if traceID == "" {
			traceID = string(ctx.Request.Header.Get("trace_id"))
		}
		if traceID == "" {
			traceID = uuid.New().String()
		}
		ctx.Set("trace_id", traceID)

		ctx.Next(c)

		if storedTraceID, exists := ctx.Get("trace_id"); exists {
			traceID = storedTraceID.(string)
		}

		latency := time.Since(start).Milliseconds()
		statusCode := ctx.Response.StatusCode()

		userID := ""
		if uid, exists := ctx.Get("user_id"); exists {
			userID = uid.(string)
		}

		userAccount := ""
		if username, exists := ctx.Get("username"); exists {
			userAccount = username.(string)
		}

		tenantID := string(ctx.Request.Header.Get("X-Tenant-ID"))
		if tenantID == "" {
			if tid, exists := ctx.Get("tenant_id"); exists {
				tenantID = fmt.Sprintf("%v", tid)
			}
		}
		if tenantID == "" {
			tenantID = "1"
		}

		clientInfo := utils.ParseUserAgent(
			string(ctx.Request.Header.UserAgent()),
			ctx.ClientIP(),
			string(ctx.Request.Header.Get("Accept-Language")),
		)

		log := &model.SysOperationLog{
			TraceID:        traceID,
			UserID:         userID,
			UserAccount:    userAccount,
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
			Source:         source,
			CreateAt:       time.Now(),
		}

		if statusCode >= 400 {
			log.ErrorMessage = fmtResponseBody(ctx)
		}

		auditSvc.RecordOperation(c, log)
	}
}

func fmtResponseBody(ctx *app.RequestContext) string {
	return string(ctx.Response.Body())
}
