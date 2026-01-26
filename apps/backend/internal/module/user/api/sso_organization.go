package api

import (
	"context"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"

	"github.com/cloudwego/hertz/pkg/app"
)

// SsoOrganizationHandler 组织处理器结构体
type SsoOrganizationHandler struct {
	orgService service.SsoOrganizationService
}

// NewSsoOrganizationHandler 创建组织处理器实例
func NewSsoOrganizationHandler(orgService service.SsoOrganizationService) *SsoOrganizationHandler {
	return &SsoOrganizationHandler{orgService: orgService}
}

// SsoCreateOrganizationRequest 创建组织请求结构
type SsoCreateOrganizationRequest struct {
	ParentID        string `json:"parent_id" form:"parent_id"`
	FromID          string `json:"from_id" form:"from_id"`
	ApplicationCode string `json:"application_code" form:"application_code"`
	UnitName        string `json:"unit_name" form:"unit_name" binding:"required"`
	UnitShort       string `json:"unit_short" form:"unit_short"`
	UnitEn          string `json:"unit_en" form:"unit_en"`
	UnitEnShort     string `json:"unit_en_short" form:"unit_en_short"`
	UnitCode        string `json:"unit_code" form:"unit_code" binding:"required"`
	KindCode        string `json:"kind_code" form:"kind_code" binding:"required"`
	Logo            string `json:"logo" form:"logo"`
	Host            string `json:"host" form:"host"`
	Contact         string `json:"contact" form:"contact"`
	Phone           string `json:"phone" form:"phone"`
	Address         string `json:"address" form:"address"`
	Postcode        string `json:"postcode" form:"postcode"`
	State           int    `json:"state" form:"state"`
	Remark          string `json:"remark" form:"remark"`
	Sort            int    `json:"sort" form:"sort"`
}

// SsoUpdateOrganizationRequest 更新组织请求结构
type SsoUpdateOrganizationRequest struct {
	ParentID        string `json:"parent_id" form:"parent_id"`
	FromID          string `json:"from_id" form:"from_id"`
	ApplicationCode string `json:"application_code" form:"application_code"`
	UnitName        string `json:"unit_name" form:"unit_name"`
	UnitShort       string `json:"unit_short" form:"unit_short"`
	UnitEn          string `json:"unit_en" form:"unit_en"`
	UnitEnShort     string `json:"unit_en_short" form:"unit_en_short"`
	UnitCode        string `json:"unit_code" form:"unit_code"`
	KindCode        string `json:"kind_code" form:"kind_code"`
	Logo            string `json:"logo" form:"logo"`
	Host            string `json:"host" form:"host"`
	Contact         string `json:"contact" form:"contact"`
	Phone           string `json:"phone" form:"phone"`
	Address         string `json:"address" form:"address"`
	Postcode        string `json:"postcode" form:"postcode"`
	State           int    `json:"state" form:"state"`
	Remark          string `json:"remark" form:"remark"`
	Sort            int    `json:"sort" form:"sort"`
}

// CreateOrganization 创建组织
func (h *SsoOrganizationHandler) CreateOrganization(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateOrganizationRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 创建组织模型
	unit := &model.SsoOrganization{
		ParentID:        req.ParentID,
		FromID:          req.FromID,
		ApplicationCode: req.ApplicationCode,
		UnitName:        req.UnitName,
		UnitShort:       req.UnitShort,
		UnitEn:          req.UnitEn,
		UnitEnShort:     req.UnitEnShort,
		UnitCode:        req.UnitCode,
		KindCode:        req.KindCode,
		Logo:            req.Logo,
		Host:            req.Host,
		Contact:         req.Contact,
		Phone:           req.Phone,
		Address:         req.Address,
		Postcode:        req.Postcode,
		State:           req.State,
		Remark:          req.Remark,
		Sort:            req.Sort,
	}

	// 调用服务层创建组织
	if err := h.orgService.CreateOrganization(unit); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(201, unit)
}

// GetOrganizationByID 根据ID获取组织
func (h *SsoOrganizationHandler) GetOrganizationByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取组织
	unit, err := h.orgService.GetOrganizationByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "组织不存在"})
		return
	}

	ctx.JSON(200, unit)
}

// UpdateOrganization 更新组织
func (h *SsoOrganizationHandler) UpdateOrganization(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req SsoUpdateOrganizationRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 调用服务层获取组织
	unit, err := h.orgService.GetOrganizationByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "组织不存在"})
		return
	}

	// 更新组织字段
	if req.ParentID != "" {
		unit.ParentID = req.ParentID
	}
	if req.FromID != "" {
		unit.FromID = req.FromID
	}
	if req.ApplicationCode != "" {
		unit.ApplicationCode = req.ApplicationCode
	}
	if req.UnitName != "" {
		unit.UnitName = req.UnitName
	}
	if req.UnitShort != "" {
		unit.UnitShort = req.UnitShort
	}
	if req.UnitEn != "" {
		unit.UnitEn = req.UnitEn
	}
	if req.UnitEnShort != "" {
		unit.UnitEnShort = req.UnitEnShort
	}
	if req.UnitCode != "" {
		unit.UnitCode = req.UnitCode
	}
	if req.KindCode != "" {
		unit.KindCode = req.KindCode
	}
	if req.Logo != "" {
		unit.Logo = req.Logo
	}
	if req.Host != "" {
		unit.Host = req.Host
	}
	if req.Contact != "" {
		unit.Contact = req.Contact
	}
	if req.Phone != "" {
		unit.Phone = req.Phone
	}
	if req.Address != "" {
		unit.Address = req.Address
	}
	if req.Postcode != "" {
		unit.Postcode = req.Postcode
	}
	if req.State != 0 {
		unit.State = req.State
	}
	if req.Remark != "" {
		unit.Remark = req.Remark
	}
	if req.Sort != 0 {
		unit.Sort = req.Sort
	}

	// 调用服务层更新组织
	if err := h.orgService.UpdateOrganization(unit); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, unit)
}

// DeleteOrganization 删除组织
func (h *SsoOrganizationHandler) DeleteOrganization(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层删除组织
	if err := h.orgService.DeleteOrganization(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, map[string]string{"message": "组织删除成功"})
}

// GetAllOrganizations 获取所有组织
func (h *SsoOrganizationHandler) GetAllOrganizations(c context.Context, ctx *app.RequestContext) {
	// 调用服务层获取所有组织
	units, err := h.orgService.GetAllOrganizations()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取组织列表失败"})
		return
	}

	ctx.JSON(200, units)
}
