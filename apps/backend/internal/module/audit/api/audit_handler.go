package api

import (
	"metadata-platform/internal/module/audit/service"

	"github.com/cloudwego/hertz/pkg/app"
)

type AuditHandler struct {
	service service.AuditService
}

func NewAuditHandler(s service.AuditService) *AuditHandler {
	return &AuditHandler{service: s}
}

func getFilters(c *app.RequestContext) map[string]interface{} {
	filters := make(map[string]interface{})
	// Generic filter extraction
	if v := c.Query("user_id"); v != "" {
		filters["user_id"] = v
	}
	if v := c.Query("account"); v != "" {
		filters["account"] = v
	}
	if v := c.Query("module"); v != "" {
		// Frontend uses module, backend model uses source
		filters["source"] = v
	}
	if v := c.Query("status"); v != "" {
		filters["status"] = v
	}
	if v := c.Query("start_time"); v != "" {
		filters["start_time"] = v
	}
	if v := c.Query("end_time"); v != "" {
		filters["end_time"] = v
	}
	// Data Change specific
	if v := c.Query("table_name"); v != "" {
		filters["model_id"] = v
	}
	if v := c.Query("data_type"); v != "" {
		filters["action"] = v
	}
	// Operation specific
	if v := c.Query("type"); v != "" {
		filters["method"] = v
	}
	return filters
}
