package api

import (
	"context"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

// SsoOrgHandler 组织处理器结构体
type SsoOrgHandler struct {
	orgService service.SsoOrgService
}

// NewSsoOrgHandler 创建组织处理器实例
func NewSsoOrgHandler(orgService service.SsoOrgService) *SsoOrgHandler {
	return &SsoOrgHandler{orgService: orgService}
}

// SsoCreateOrgRequest 创建组织请求结构
type SsoCreateOrgRequest struct {
	ParentID   string `json:"parent_id" form:"parent_id"`
	FromID     string `json:"from_id" form:"from_id"`
	AppCode    string `json:"app_code" form:"app_code"`
	OrgName    string `json:"org_name" form:"org_name" binding:"required"`
	OrgShort   string `json:"org_short" form:"org_short"`
	OrgEn      string `json:"org_en" form:"org_en"`
	OrgEnShort string `json:"org_en_short" form:"org_en_short"`
	OrgCode    string `json:"org_code" form:"org_code" binding:"required"`
	KindCode   string `json:"kind_code" form:"kind_code" binding:"required"`
	Logo       string `json:"logo" form:"logo"`
	Host       string `json:"host" form:"host"`
	Contact    string `json:"contact" form:"contact"`
	Phone      string `json:"phone" form:"phone"`
	Address    string `json:"address" form:"address"`
	Postcode   string `json:"postcode" form:"postcode"`
	Status     int    `json:"state" form:"state"`
	Remark     string `json:"remark" form:"remark"`
	Sort       int    `json:"sort" form:"sort"`
}

// SsoUpdateOrgRequest 更新组织请求结构
type SsoUpdateOrgRequest struct {
	ParentID   string `json:"parent_id" form:"parent_id"`
	FromID     string `json:"from_id" form:"from_id"`
	AppCode    string `json:"app_code" form:"app_code"`
	OrgName    string `json:"org_name" form:"org_name"`
	OrgShort   string `json:"org_short" form:"org_short"`
	OrgEn      string `json:"org_en" form:"org_en"`
	OrgEnShort string `json:"org_en_short" form:"org_en_short"`
	OrgCode    string `json:"org_code" form:"org_code"`
	KindCode   string `json:"kind_code" form:"kind_code"`
	Logo       string `json:"logo" form:"logo"`
	Host       string `json:"host" form:"host"`
	Contact    string `json:"contact" form:"contact"`
	Phone      string `json:"phone" form:"phone"`
	Address    string `json:"address" form:"address"`
	Postcode   string `json:"postcode" form:"postcode"`
	Status     int    `json:"state" form:"state"`
	Remark     string `json:"remark" form:"remark"`
	Sort       int    `json:"sort" form:"sort"`
}

// CreateOrg 创建组织
func (h *SsoOrgHandler) CreateOrg(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateOrgRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 创建组织模型
	unit := &model.SsoOrg{
		ParentID:   req.ParentID,
		FromID:     req.FromID,
		AppCode:    req.AppCode,
		OrgName:    req.OrgName,
		OrgShort:   req.OrgShort,
		OrgEn:      req.OrgEn,
		OrgEnShort: req.OrgEnShort,
		OrgCode:    req.OrgCode,
		KindCode:   req.KindCode,
		Logo:       req.Logo,
		Host:       req.Host,
		Contact:    req.Contact,
		Phone:      req.Phone,
		Address:    req.Address,
		Postcode:   req.Postcode,
		Status:     req.Status,
		Remark:     req.Remark,
		Sort:       req.Sort,
	}

	// 生成组织ID
	unit.ID = utils.GetSnowflake().GenerateIDString()

	// 调用服务层创建组织
	if err := h.orgService.CreateOrg(unit); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(201, unit)
}

// GetOrgByID 根据ID获取组织
func (h *SsoOrgHandler) GetOrgByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取组织
	unit, err := h.orgService.GetOrgByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "组织不存在"})
		return
	}

	ctx.JSON(200, unit)
}

// UpdateOrg 更新组织
func (h *SsoOrgHandler) UpdateOrg(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req SsoUpdateOrgRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 调用服务层获取组织
	unit, err := h.orgService.GetOrgByID(id)
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
	if req.AppCode != "" {
		unit.AppCode = req.AppCode
	}
	if req.OrgName != "" {
		unit.OrgName = req.OrgName
	}
	if req.OrgShort != "" {
		unit.OrgShort = req.OrgShort
	}
	if req.OrgEn != "" {
		unit.OrgEn = req.OrgEn
	}
	if req.OrgEnShort != "" {
		unit.OrgEnShort = req.OrgEnShort
	}
	if req.OrgCode != "" {
		unit.OrgCode = req.OrgCode
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
	if req.Status != 0 {
		unit.Status = req.Status
	}
	if req.Remark != "" {
		unit.Remark = req.Remark
	}
	if req.Sort != 0 {
		unit.Sort = req.Sort
	}

	// 调用服务层更新组织
	if err := h.orgService.UpdateOrg(unit); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, unit)
}

// DeleteOrg 删除组织
func (h *SsoOrgHandler) DeleteOrg(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层删除组织
	if err := h.orgService.DeleteOrg(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, map[string]string{"message": "组织删除成功"})
}

// GetAllOrgs 获取所有组织
func (h *SsoOrgHandler) GetAllOrgs(c context.Context, ctx *app.RequestContext) {
	// 调用服务层获取所有组织
	units, err := h.orgService.GetAllOrgs()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取组织列表失败"})
		return
	}

	ctx.JSON(200, units)
}
