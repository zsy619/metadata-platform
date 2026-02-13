package api

import (
	"context"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"

	"github.com/cloudwego/hertz/pkg/app"
)

// SsoAppHandler 应用处理器结构体
type SsoAppHandler struct {
	appService service.SsoAppService
}

// NewSsoAppHandler 创建应用处理器实例
func NewSsoAppHandler(appService service.SsoAppService) *SsoAppHandler {
	return &SsoAppHandler{appService: appService}
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
	ParentID string `json:"parent_id" form:"parent_id"`
	AppName  string `json:"app_name" form:"app_name"`
	AppCode  string `json:"app_code" form:"app_code"`
	Status   int    `json:"state" form:"state"`
	Host     string `json:"host" form:"host"`
	Logo     string `json:"logo" form:"logo"`
	Remark   string `json:"remark" form:"remark"`
	Sort     int    `json:"sort" form:"sort"`
}

// CreateApp 创建应用
func (h *SsoAppHandler) CreateApp(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateAppRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 创建应用模型
	application := &model.SsoApp{
		ParentID: req.ParentID,
		AppName:  req.AppName,
		AppCode:  req.AppCode,
		Status:   req.Status,
		Host:     req.Host,
		Logo:     req.Logo,
		Remark:   req.Remark,
		Sort:     req.Sort,
	}

	// 调用服务层创建应用
	if err := h.appService.CreateApp(application); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

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

	// 调用服务层获取应用
	application, err := h.appService.GetAppByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "应用不存在"})
		return
	}

	// 更新应用字段
	if req.ParentID != "" {
		application.ParentID = req.ParentID
	}
	if req.AppName != "" {
		application.AppName = req.AppName
	}
	if req.AppCode != "" {
		application.AppCode = req.AppCode
	}
	if req.Status != 0 {
		application.Status = req.Status
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
	if err := h.appService.UpdateApp(application); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, application)
}

// DeleteApp 删除应用
func (h *SsoAppHandler) DeleteApp(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层删除应用
	if err := h.appService.DeleteApp(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

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
