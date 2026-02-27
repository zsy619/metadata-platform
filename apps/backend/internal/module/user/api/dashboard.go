package api

import (
	"context"
	"time"

	auditService "metadata-platform/internal/module/audit/service"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	db           *gorm.DB
	repos        *repository.Repositories
	auditService auditService.AuditService
}

func NewDashboardHandler(db *gorm.DB, repos *repository.Repositories, auditService auditService.AuditService) *DashboardHandler {
	return &DashboardHandler{
		db:           db,
		repos:        repos,
		auditService: auditService,
	}
}

// StatsResponse 统计数据响应
type StatsResponse struct {
	UserCount      int64 `json:"user_count"`
	RoleCount      int64 `json:"role_count"`
	OrgCount       int64 `json:"org_count"`
	MenuCount      int64 `json:"menu_count"`
	PosCount       int64 `json:"pos_count"`
	UserGroupCount int64 `json:"user_group_count"`
	RoleGroupCount int64 `json:"role_group_count"`
}

// LoginTrendItem 登录趋势项
type LoginTrendItem struct {
	Date    string `json:"date"`
	Success int64  `json:"success"`
	Fail    int64  `json:"fail"`
}

// UserStatusDistribution 用户状态分布
type UserStatusDistribution struct {
	Active   int64 `json:"active"`
	Inactive int64 `json:"inactive"`
	Locked   int64 `json:"locked"`
	Pending  int64 `json:"pending"`
}

// OperationStats 操作统计
type OperationStats struct {
	Create int64 `json:"create"`
	Update int64 `json:"update"`
	Delete int64 `json:"delete"`
	Query  int64 `json:"query"`
	Export int64 `json:"export"`
}

// OrgDistribution 组织分布项
type OrgDistribution struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// GetStats 获取统计数据
func (h *DashboardHandler) GetStats(c context.Context, ctx *app.RequestContext) {
	userCount, _ := h.repos.User.Count()
	roleCount, _ := h.repos.Role.Count()
	orgCount, _ := h.repos.Org.Count()
	menuCount, _ := h.repos.Menu.Count()
	posCount, _ := h.repos.Pos.Count()
	userGroupCount, _ := h.repos.UserGroup.Count()
	roleGroupCount, _ := h.repos.RoleGroup.Count()

	utils.SuccessResponse(ctx, StatsResponse{
		UserCount:      userCount,
		RoleCount:      roleCount,
		OrgCount:       orgCount,
		MenuCount:      menuCount,
		PosCount:       posCount,
		UserGroupCount: userGroupCount,
		RoleGroupCount: roleGroupCount,
	})
}

// GetRecentLoginLogs 获取最近登录日志
func (h *DashboardHandler) GetRecentLoginLogs(c context.Context, ctx *app.RequestContext) {
	logs, _ := h.auditService.GetRecentLoginLogs(10)
	utils.SuccessResponse(ctx, logs)
}

// GetRecentOperationLogs 获取最近操作日志
func (h *DashboardHandler) GetRecentOperationLogs(c context.Context, ctx *app.RequestContext) {
	logs, _ := h.auditService.GetRecentOperationLogs(10)
	utils.SuccessResponse(ctx, logs)
}

// GetLoginTrend 获取登录趋势（最近7天）
func (h *DashboardHandler) GetLoginTrend(c context.Context, ctx *app.RequestContext) {
	trends := make([]LoginTrendItem, 0)

	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("01-02")

		startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endTime := startTime.Add(24 * time.Hour)

		var successCount, failCount int64

		h.db.Table("sys_login_log").
			Where("create_at >= ? AND create_at < ? AND login_status = ?", startTime, endTime, 1).
			Count(&successCount)

		h.db.Table("sys_login_log").
			Where("create_at >= ? AND create_at < ? AND login_status = ?", startTime, endTime, 0).
			Count(&failCount)

		trends = append(trends, LoginTrendItem{
			Date:    dateStr,
			Success: successCount,
			Fail:    failCount,
		})
	}

	utils.SuccessResponse(ctx, trends)
}

// GetUserStatusDistribution 获取用户状态分布
func (h *DashboardHandler) GetUserStatusDistribution(c context.Context, ctx *app.RequestContext) {
	var active, inactive, locked, pending int64

	h.db.Table("sso_user").Where("status = ?", 1).Count(&active)
	h.db.Table("sso_user").Where("status = ?", 0).Count(&inactive)
	h.db.Table("sso_user").Where("status = ?", 2).Count(&locked)
	h.db.Table("sso_user").Where("status = ?", 3).Count(&pending)

	utils.SuccessResponse(ctx, UserStatusDistribution{
		Active:   active,
		Inactive: inactive,
		Locked:   locked,
		Pending:  pending,
	})
}

// GetOperationStats 获取操作统计
func (h *DashboardHandler) GetOperationStats(c context.Context, ctx *app.RequestContext) {
	var create, update, delete, query, export int64

	h.db.Table("sys_operation_log").Where("action = ?", "create").Count(&create)
	h.db.Table("sys_operation_log").Where("action = ?", "update").Count(&update)
	h.db.Table("sys_operation_log").Where("action = ?", "delete").Count(&delete)
	h.db.Table("sys_operation_log").Where("action = ?", "query").Count(&query)
	h.db.Table("sys_operation_log").Where("action = ?", "export").Count(&export)

	utils.SuccessResponse(ctx, OperationStats{
		Create: create,
		Update: update,
		Delete: delete,
		Query:  query,
		Export: export,
	})
}

// GetOrgDistribution 获取组织分布
func (h *DashboardHandler) GetOrgDistribution(c context.Context, ctx *app.RequestContext) {
	distributions := make([]OrgDistribution, 0)

	h.db.Table("sso_org as o").
		Select("ok.name as name, count(*) as value").
		Joins("LEFT JOIN sso_org_kind ok ON o.org_kind_id = ok.id").
		Where("o.is_deleted = ?", false).
		Group("ok.id, ok.name").
		Scan(&distributions)

	if len(distributions) == 0 {
		var orgs []struct {
			OrgName string `gorm:"column:org_name"`
		}
		h.db.Table("sso_org").
			Select("org_name").
			Where("is_deleted = ? AND parent_id = ''", false).
			Limit(5).
			Scan(&orgs)

		for _, org := range orgs {
			distributions = append(distributions, OrgDistribution{
				Name:  org.OrgName,
				Value: 1,
			})
		}
	}

	utils.SuccessResponse(ctx, distributions)
}
