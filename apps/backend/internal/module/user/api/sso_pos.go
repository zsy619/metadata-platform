package api

import (
	"context"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

// SsoPosHandler 职位处理器结构体
type SsoPosHandler struct {
	posService service.SsoPosService
}

// NewSsoPosHandler 创建职位处理器实例
func NewSsoPosHandler(posService service.SsoPosService) *SsoPosHandler {
	return &SsoPosHandler{posService: posService}
}

// SsoCreatePosRequest 创建职位请求结构
type SsoCreatePosRequest struct {
	ParentID string `json:"parent_id" form:"parent_id"`
	AppCode  string `json:"app_code" form:"app_code"`
	OrgID    string `json:"org_id" form:"org_id"`
	KindCode string `json:"kind_code" form:"kind_code" binding:"required"`
	PosName  string `json:"pos_name" form:"pos_name" binding:"required"`
	PosCode  string `json:"pos_code" form:"pos_code" binding:"required"`
	Status   int    `json:"state" form:"state"`
	Remark   string `json:"remark" form:"remark"`
	Sort     int    `json:"sort" form:"sort"`
}

// SsoUpdatePosRequest 更新职位请求结构
type SsoUpdatePosRequest struct {
	ParentID string `json:"parent_id" form:"parent_id"`
	AppCode  string `json:"app_code" form:"app_code"`
	OrgID    string `json:"org_id" form:"org_id"`
	KindCode string `json:"kind_code" form:"kind_code"`
	PosName  string `json:"pos_name" form:"pos_name"`
	PosCode  string `json:"pos_code" form:"pos_code"`
	Status   int    `json:"state" form:"state"`
	Remark   string `json:"remark" form:"remark"`
	Sort     int    `json:"sort" form:"sort"`
}

// CreatePosition 创建职位
func (h *SsoPosHandler) CreatePos(c context.Context, ctx *app.RequestContext) {
	var req SsoCreatePosRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 创建职位模型
	position := &model.SsoPos{
		ParentID: req.ParentID,
		AppCode:  req.AppCode,
		OrgID:    req.OrgID,
		KindCode: req.KindCode,
		PosName:  req.PosName,
		PosCode:  req.PosCode,
		Status:   req.Status,
		Remark:   req.Remark,
		Sort:     req.Sort,
	}

	// 生成职位ID
	position.ID = utils.GetSnowflake().GenerateIDString()

	// 调用服务层创建职位
	if err := h.posService.CreatePos(position); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(201, position)
}

// GetPositionByID 根据ID获取职位
func (h *SsoPosHandler) GetPosByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取职位
	position, err := h.posService.GetPosByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "职位不存在"})
		return
	}

	ctx.JSON(200, position)
}

// UpdatePosition 更新职位
func (h *SsoPosHandler) UpdatePos(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req SsoUpdatePosRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 调用服务层获取职位
	position, err := h.posService.GetPosByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "职位不存在"})
		return
	}

	// 更新职位字段
	if req.ParentID != "" {
		position.ParentID = req.ParentID
	}
	if req.AppCode != "" {
		position.AppCode = req.AppCode
	}
	if req.OrgID != "" {
		position.OrgID = req.OrgID
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
	if req.Status != 0 {
		position.Status = req.Status
	}
	if req.Remark != "" {
		position.Remark = req.Remark
	}
	if req.Sort != 0 {
		position.Sort = req.Sort
	}

	// 调用服务层更新职位
	if err := h.posService.UpdatePos(position); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, position)
}

// DeletePos 删除职位
func (h *SsoPosHandler) DeletePos(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层删除职位
	if err := h.posService.DeletePos(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, map[string]string{"message": "职位删除成功"})
}

// GetAllPoss 获取所有职位
func (h *SsoPosHandler) GetAllPoss(c context.Context, ctx *app.RequestContext) {
	// 调用服务层获取所有职位
	positions, err := h.posService.GetAllPoss()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取职位列表失败"})
		return
	}

	ctx.JSON(200, positions)
}
