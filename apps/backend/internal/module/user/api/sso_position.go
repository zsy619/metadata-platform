package api

import (
	"context"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"

	"github.com/cloudwego/hertz/pkg/app"
)

// SsoPositionHandler 职位处理器结构体
type SsoPositionHandler struct {
	posService service.SsoPositionService
}

// NewSsoPositionHandler 创建职位处理器实例
func NewSsoPositionHandler(posService service.SsoPositionService) *SsoPositionHandler {
	return &SsoPositionHandler{posService: posService}
}

// SsoCreatePosRequest 创建职位请求结构
type SsoCreatePosRequest struct {
	ParentID        string `json:"parent_id" form:"parent_id"`
	ApplicationCode string `json:"application_code" form:"application_code"`
	OrganizationID  string `json:"organization_id" form:"organization_id"`
	KindCode        string `json:"kind_code" form:"kind_code" binding:"required"`
	PosName         string `json:"pos_name" form:"pos_name" binding:"required"`
	PosCode         string `json:"pos_code" form:"pos_code" binding:"required"`
	State           int    `json:"state" form:"state"`
	Remark          string `json:"remark" form:"remark"`
	Sort            int    `json:"sort" form:"sort"`
}

// SsoUpdatePosRequest 更新职位请求结构
type SsoUpdatePosRequest struct {
	ParentID        string `json:"parent_id" form:"parent_id"`
	ApplicationCode string `json:"application_code" form:"application_code"`
	OrganizationID  string `json:"organization_id" form:"organization_id"`
	KindCode        string `json:"kind_code" form:"kind_code"`
	PosName         string `json:"pos_name" form:"pos_name"`
	PosCode         string `json:"pos_code" form:"pos_code"`
	State           int    `json:"state" form:"state"`
	Remark          string `json:"remark" form:"remark"`
	Sort            int    `json:"sort" form:"sort"`
}

// CreatePosition 创建职位
func (h *SsoPositionHandler) CreatePosition(c context.Context, ctx *app.RequestContext) {
	var req SsoCreatePosRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 创建职位模型
	position := &model.SsoPosition{
		ParentID:        req.ParentID,
		ApplicationCode: req.ApplicationCode,
		OrganizationID:  req.OrganizationID,
		KindCode:        req.KindCode,
		PosName:         req.PosName,
		PosCode:         req.PosCode,
		State:           req.State,
		Remark:          req.Remark,
		Sort:            req.Sort,
	}

	// 调用服务层创建职位
	if err := h.posService.CreatePosition(position); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(201, position)
}

// GetPositionByID 根据ID获取职位
func (h *SsoPositionHandler) GetPositionByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取职位
	position, err := h.posService.GetPositionByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "职位不存在"})
		return
	}

	ctx.JSON(200, position)
}

// UpdatePosition 更新职位
func (h *SsoPositionHandler) UpdatePosition(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req SsoUpdatePosRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 调用服务层获取职位
	position, err := h.posService.GetPositionByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "职位不存在"})
		return
	}

	// 更新职位字段
	if req.ParentID != "" {
		position.ParentID = req.ParentID
	}
	if req.ApplicationCode != "" {
		position.ApplicationCode = req.ApplicationCode
	}
	if req.OrganizationID != "" {
		position.OrganizationID = req.OrganizationID
	}
	if req.KindCode != "" {
		position.KindCode = req.KindCode
	}
	if req.PosName != "" {
		position.PosName = req.PosName
	}
	if req.PosCode != "" {
		position.PosCode = req.PosCode
	}
	if req.State != 0 {
		position.State = req.State
	}
	if req.Remark != "" {
		position.Remark = req.Remark
	}
	if req.Sort != 0 {
		position.Sort = req.Sort
	}

	// 调用服务层更新职位
	if err := h.posService.UpdatePosition(position); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, position)
}

// DeletePosition 删除职位
func (h *SsoPositionHandler) DeletePosition(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层删除职位
	if err := h.posService.DeletePosition(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, map[string]string{"message": "职位删除成功"})
}

// GetAllPositions 获取所有职位
func (h *SsoPositionHandler) GetAllPositions(c context.Context, ctx *app.RequestContext) {
	// 调用服务层获取所有职位
	positions, err := h.posService.GetAllPositions()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取职位列表失败"})
		return
	}

	ctx.JSON(200, positions)
}
