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

// SsoOrgHandler 组织处理器结构体
type SsoOrgHandler struct {
	*utils.BaseHandler
	orgService service.SsoOrgService
	audit      AuditService
}

// NewSsoOrgHandler 创建组织处理器实例
func NewSsoOrgHandler(orgService service.SsoOrgService, auditQueue *queue.AuditLogQueue) *SsoOrgHandler {
	return &SsoOrgHandler{
		BaseHandler: utils.NewBaseHandler(),
		orgService:  orgService,
		audit:       &auditServiceImpl{queue: auditQueue},
	}
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
	ParentID   *string `json:"parent_id" form:"parent_id"`
	FromID     *string `json:"from_id" form:"from_id"`
	AppCode    *string `json:"app_code" form:"app_code"`
	OrgName    *string `json:"org_name" form:"org_name"`
	OrgShort   *string `json:"org_short" form:"org_short"`
	OrgEn      *string `json:"org_en" form:"org_en"`
	OrgEnShort *string `json:"org_en_short" form:"org_en_short"`
	OrgCode    *string `json:"org_code" form:"org_code"`
	KindCode   *string `json:"kind_code" form:"kind_code"`
	Logo       *string `json:"logo" form:"logo"`
	Host       *string `json:"host" form:"host"`
	Contact    *string `json:"contact" form:"contact"`
	Phone      *string `json:"phone" form:"phone"`
	Address    *string `json:"address" form:"address"`
	Postcode   *string `json:"postcode" form:"postcode"`
	Status     *int    `json:"state" form:"state"`
	Remark     *string `json:"remark" form:"remark"`
	Sort       *int    `json:"sort" form:"sort"`
}

// CreateOrg 创建组织
func (h *SsoOrgHandler) CreateOrg(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateOrgRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 创建组织模型
	unit := &userModel.SsoOrg{
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
		CreateID:   headerUser.UserID,
		CreateBy:   headerUser.UserAccount,
		TenantID:   headerUser.TenantID,
	}

	// 生成组织ID
	unit.ID = utils.GetSnowflake().GenerateIDString()

	// 调用服务层创建组织
	if err := h.orgService.CreateOrg(unit); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	afterData, _ := json.Marshal(unit)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:        utils.GetSnowflake().GenerateIDString(),
		TraceID:   headerUser.TraceID,
		ModelID:   "org",
		RecordID:  unit.ID,
		Action:    "CREATE",
		AfterData: string(afterData),
		CreateBy:  headerUser.UserAccount,
		Source:    "org_service",
	})

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

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取变更前数据
	beforeData, err := h.orgService.GetOrgByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "组织不存在"})
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
	if req.FromID != nil {
		fields["from_id"] = *req.FromID
	}
	if req.AppCode != nil {
		fields["app_code"] = *req.AppCode
	}
	if req.OrgName != nil {
		fields["org_name"] = *req.OrgName
	}
	if req.OrgShort != nil {
		fields["org_short"] = *req.OrgShort
	}
	if req.OrgEn != nil {
		fields["org_en"] = *req.OrgEn
	}
	if req.OrgEnShort != nil {
		fields["org_en_short"] = *req.OrgEnShort
	}
	if req.OrgCode != nil {
		fields["org_code"] = *req.OrgCode
	}
	if req.KindCode != nil {
		fields["kind_code"] = *req.KindCode
	}
	if req.Logo != nil {
		fields["logo"] = *req.Logo
	}
	if req.Host != nil {
		fields["host"] = *req.Host
	}
	if req.Contact != nil {
		fields["contact"] = *req.Contact
	}
	if req.Phone != nil {
		fields["phone"] = *req.Phone
	}
	if req.Address != nil {
		fields["address"] = *req.Address
	}
	if req.Postcode != nil {
		fields["postcode"] = *req.Postcode
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}
	if req.Remark != nil {
		fields["remark"] = *req.Remark
	}
	if req.Sort != nil {
		fields["sort"] = *req.Sort
	}

	// 调用服务层更新组织
	if err := h.orgService.UpdateOrgFields(id, fields); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取更新后的组织
	unit, err := h.orgService.GetOrgByID(id)
	if err != nil {
		ctx.JSON(200, map[string]string{"message": "更新成功"})
		return
	}

	// 记录数据变更日志
	afterJSON, _ := json.Marshal(unit)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "org",
		RecordID:   id,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "org_service",
	})

	ctx.JSON(200, unit)
}

// DeleteOrg 删除组织
func (h *SsoOrgHandler) DeleteOrg(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取删除前数据
	beforeData, err := h.orgService.GetOrgByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "组织不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	// 调用服务层删除组织
	if err := h.orgService.DeleteOrg(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "org",
		RecordID:   id,
		Action:     "DELETE",
		BeforeData: string(beforeJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "org_service",
	})

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
