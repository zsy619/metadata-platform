package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"
)

// SsoApplicationHandler 应用处理器结构体
type SsoApplicationHandler struct {
	appService service.SsoApplicationService
}

// NewSsoApplicationHandler 创建应用处理器实例
func NewSsoApplicationHandler(appService service.SsoApplicationService) *SsoApplicationHandler {
	return &SsoApplicationHandler{appService: appService}
}

// SsoCreateAppRequest 创建应用请求结构
type SsoCreateAppRequest struct {
	ParentID        string `json:"parent_id" form:"parent_id"`
	ApplicationName string `json:"application_name" form:"application_name" binding:"required"`
	ApplicationCode string `json:"application_code" form:"application_code" binding:"required"`
	State           int    `json:"state" form:"state"`
	Host            string `json:"host" form:"host"`
	Logo            string `json:"logo" form:"logo"`
	Remark          string `json:"remark" form:"remark"`
	Sort            int    `json:"sort" form:"sort"`
}

// SsoUpdateAppRequest 更新应用请求结构
type SsoUpdateAppRequest struct {
	ParentID        string `json:"parent_id" form:"parent_id"`
	ApplicationName string `json:"application_name" form:"application_name"`
	ApplicationCode string `json:"application_code" form:"application_code"`
	State           int    `json:"state" form:"state"`
	Host            string `json:"host" form:"host"`
	Logo            string `json:"logo" form:"logo"`
	Remark          string `json:"remark" form:"remark"`
	Sort            int    `json:"sort" form:"sort"`
}

// CreateApplication 创建应用
func (h *SsoApplicationHandler) CreateApplication(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateAppRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 创建应用模型
	application := &model.SsoApplication{
		ParentID:        req.ParentID,
		ApplicationName: req.ApplicationName,
		ApplicationCode: req.ApplicationCode,
		State:           req.State,
		Host:            req.Host,
		Logo:            req.Logo,
		Remark:          req.Remark,
		Sort:            req.Sort,
	}

	// 调用服务层创建应用
	if err := h.appService.CreateApplication(application); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(201, application)
}

// GetApplicationByID 根据ID获取应用
func (h *SsoApplicationHandler) GetApplicationByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取应用
	application, err := h.appService.GetApplicationByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "应用不存在"})
		return
	}

	ctx.JSON(200, application)
}

// UpdateApplication 更新应用
func (h *SsoApplicationHandler) UpdateApplication(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req SsoUpdateAppRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 调用服务层获取应用
	application, err := h.appService.GetApplicationByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "应用不存在"})
		return
	}

	// 更新应用字段
	if req.ParentID != "" {
		application.ParentID = req.ParentID
	}
	if req.ApplicationName != "" {
		application.ApplicationName = req.ApplicationName
	}
	if req.ApplicationCode != "" {
		application.ApplicationCode = req.ApplicationCode
	}
	if req.State != 0 {
		application.State = req.State
	}
	if req.Host != "" {
		application.Host = req.Host
	}
	if req.Logo != "" {
		application.Logo = req.Logo
	}
	if req.Remark != "" {
		application.Remark = req.Remark
	}
	if req.Sort != 0 {
		application.Sort = req.Sort
	}

	// 调用服务层更新应用
	if err := h.appService.UpdateApplication(application); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, application)
}

// DeleteApplication 删除应用
func (h *SsoApplicationHandler) DeleteApplication(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层删除应用
	if err := h.appService.DeleteApplication(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, map[string]string{"message": "应用删除成功"})
}

// GetAllApplications 获取所有应用
func (h *SsoApplicationHandler) GetAllApplications(c context.Context, ctx *app.RequestContext) {
	// 调用服务层获取所有应用
	applications, err := h.appService.GetAllApplications()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取应用列表失败"})
		return
	}

	ctx.JSON(200, applications)
}
