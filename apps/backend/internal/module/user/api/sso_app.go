package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"

	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/queue"
	userModel "metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"
)

// SsoAppHandler 应用处理器结构体
type SsoAppHandler struct {
	*utils.BaseHandler
	appService service.SsoAppService
	audit      AuditService
}

// NewSsoAppHandler 创建应用处理器实例
func NewSsoAppHandler(appService service.SsoAppService, auditQueue *queue.AuditLogQueue) *SsoAppHandler {
	return &SsoAppHandler{
		appService: appService,
		audit:      &auditServiceImpl{queue: auditQueue},
	}
}

// SsoCreateAppRequest 创建应用请求结构
type SsoCreateAppRequest struct {
	ParentID string `json:"parent_id" form:"parent_id"`
	AppName  string `json:"app_name" form:"app_name" binding:"required"`
	AppCode  string `json:"app_code" form:"app_code" binding:"required"`
	Status   int    `json:"state" form:"state"`
	Host     string `json:"host" form:"host"`
	Logo     string `json:"logo" form:"logo"`
	Remark   string `json:"remark" form:"remark"`
	Sort     int    `json:"sort" form:"sort"`
}

// SsoUpdateAppRequest 更新应用请求结构
type SsoUpdateAppRequest struct {
	ParentID *string `json:"parent_id" form:"parent_id"`
	AppName  *string `json:"app_name" form:"app_name"`
	AppCode  *string `json:"app_code" form:"app_code"`
	Status   *int    `json:"state" form:"state"`
	Host     *string `json:"host" form:"host"`
	Logo     *string `json:"logo" form:"logo"`
	Remark   *string `json:"remark" form:"remark"`
	Sort     *int    `json:"sort" form:"sort"`
}

// CreateApp 创建应用
func (h *SsoAppHandler) CreateApp(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateAppRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	userID, userAccount, tenantID := h.GetHeaderUser(c, ctx)
	fmt.Println(userID, userAccount, tenantID)

	// 创建应用模型
	application := &userModel.SsoApp{
		ParentID: req.ParentID,
		AppName:  req.AppName,
		AppCode:  req.AppCode,
		Status:   req.Status,
		Host:     req.Host,
		Logo:     req.Logo,
		Remark:   req.Remark,
		Sort:     req.Sort,
		CreateID: userID,
		CreateBy: userAccount,
		TenantID: tenantID,
	}

	// 调用服务层创建应用
	if err := h.appService.CreateApp(application); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	afterData, _ := json.Marshal(application)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:        utils.GetSnowflake().GenerateIDString(),
		TraceID:   ctx.GetString("trace_id"),
		ModelID:   "app",
		RecordID:  application.ID,
		Action:    "CREATE",
		AfterData: string(afterData),
		CreateBy:  userAccount,
		Source:    "app_service",
	})

	ctx.JSON(201, application)
}

// GetAppByID 根据ID获取应用
func (h *SsoAppHandler) GetAppByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取应用
	application, err := h.appService.GetAppByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "应用不存在"})
		return
	}

	ctx.JSON(200, application)
}

// UpdateApp 更新应用
func (h *SsoAppHandler) UpdateApp(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req SsoUpdateAppRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取变更前数据
	beforeData, err := h.appService.GetAppByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "应用不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	// 构建更新字段（只更新传递的字段）
	fields := map[string]any{
		"update_id": headerUser.UserID,
		"update_by": headerUser.UserAccount,
	}

	if req.ParentID != nil {
		fields["parent_id"] = *req.ParentID
	}
	if req.AppName != nil {
		fields["app_name"] = *req.AppName
	}
	if req.AppCode != nil {
		fields["app_code"] = *req.AppCode
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}
	if req.Host != nil {
		fields["host"] = *req.Host
	}
	if req.Logo != nil {
		fields["logo"] = *req.Logo
	}
	if req.Remark != nil {
		fields["remark"] = *req.Remark
	}
	if req.Sort != nil {
		fields["sort"] = *req.Sort
	}

	// 调用服务层更新应用
	if err := h.appService.UpdateAppFields(id, fields); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取更新后的应用
	application, err := h.appService.GetAppByID(id)
	if err != nil {
		ctx.JSON(200, map[string]string{"message": "更新成功"})
		return
	}

	// 记录数据变更日志
	afterJSON, _ := json.Marshal(application)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "app",
		RecordID:   id,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "app_service",
	})

	ctx.JSON(200, application)
}

// DeleteApp 删除应用
func (h *SsoAppHandler) DeleteApp(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取删除前数据
	beforeData, err := h.appService.GetAppByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "应用不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	// 调用服务层删除应用
	if err := h.appService.DeleteApp(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "app",
		RecordID:   id,
		Action:     "DELETE",
		BeforeData: string(beforeJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "app_service",
	})

	ctx.JSON(200, map[string]string{"message": "应用删除成功"})
}

// GetAllApps 获取所有应用
func (h *SsoAppHandler) GetAllApps(c context.Context, ctx *app.RequestContext) {
	// 调用服务层获取所有应用
	applications, err := h.appService.GetAllApps()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取应用列表失败"})
		return
	}

	ctx.JSON(200, applications)
}
