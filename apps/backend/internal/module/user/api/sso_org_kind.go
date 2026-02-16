package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"

	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"
	userModel "metadata-platform/internal/module/user/model"
)

type SsoOrgKindHandler struct {
	*utils.BaseHandler
	orgKindService service.SsoOrgKindService
	audit          AuditService
}

func NewSsoOrgKindHandler(orgKindService service.SsoOrgKindService, auditQueue *queue.AuditLogQueue) *SsoOrgKindHandler {
	return &SsoOrgKindHandler{
		BaseHandler:    utils.NewBaseHandler(),
		orgKindService: orgKindService,
		audit:          &auditServiceImpl{queue: auditQueue},
	}
}

type SsoCreateOrgKindRequest struct {
	ParentID string `json:"parent_id" form:"parent_id"`
	KindName string `json:"kind_name" form:"kind_name" binding:"required"`
	KindCode string `json:"kind_code" form:"kind_code" binding:"required"`
	KindTag  string `json:"kind_tag" form:"kind_tag"`
	Status   int    `json:"status" form:"status"`
	Remark   string `json:"remark" form:"remark"`
	Sort     int    `json:"sort" form:"sort"`
}

type SsoUpdateOrgKindRequest struct {
	ParentID *string `json:"parent_id" form:"parent_id"`
	KindName *string `json:"kind_name" form:"kind_name"`
	KindCode *string `json:"kind_code" form:"kind_code"`
	KindTag  *string `json:"kind_tag" form:"kind_tag"`
	Status   *int    `json:"status" form:"status"`
	Remark   *string `json:"remark" form:"remark"`
	Sort     *int    `json:"sort" form:"sort"`
}

func (h *SsoOrgKindHandler) CreateOrgKind(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateOrgKindRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		utils.BadRequestResponse(ctx, err.Error())
		return
	}

	headerUser := h.GetHeaderUserStruct(c, ctx)

	orgKind := &userModel.SsoOrgKind{
		ParentID: req.ParentID,
		KindName: req.KindName,
		KindCode: req.KindCode,
		KindTag:  req.KindTag,
		Status:   req.Status,
		Remark:   req.Remark,
		Sort:     req.Sort,
		TenantID: headerUser.TenantID,
		CreateID: headerUser.UserID,
		CreateBy: headerUser.UserAccount,
	}

	if err := h.orgKindService.CreateOrgKind(orgKind); err != nil {
		utils.BadRequestResponse(ctx, err.Error())
		return
	}

	afterData, _ := json.Marshal(orgKind)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:        utils.GetSnowflake().GenerateIDString(),
		TraceID:   ctx.GetString("trace_id"),
		ModelID:   "org_kind",
		RecordID:  orgKind.ID,
		Action:    "CREATE",
		AfterData: string(afterData),
		CreateBy:  headerUser.UserAccount,
		Source:    "sso",
	})

	utils.SuccessResponse(ctx, orgKind)
}

func (h *SsoOrgKindHandler) GetOrgKindByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	if id == "" {
		utils.BadRequestResponse(ctx, "ID不能为空")
		return
	}

	orgKind, err := h.orgKindService.GetOrgKindByID(id)
	if err != nil {
		utils.BadRequestResponse(ctx, "组织类型不存在")
		return
	}

	utils.SuccessResponse(ctx, orgKind)
}

func (h *SsoOrgKindHandler) UpdateOrgKind(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	if id == "" {
		utils.BadRequestResponse(ctx, "ID不能为空")
		return
	}

	var req SsoUpdateOrgKindRequest
	if err := ctx.Bind(&req); err != nil {
		utils.BadRequestResponse(ctx, err.Error())
		return
	}

	headerUser := h.GetHeaderUserStruct(c, ctx)

	fields := map[string]any{
		"update_id": headerUser.UserID,
		"update_by": headerUser.UserAccount,
	}

	if req.ParentID != nil {
		fields["parent_id"] = *req.ParentID
	}
	if req.KindName != nil {
		fields["kind_name"] = *req.KindName
	}
	if req.KindCode != nil {
		fields["kind_code"] = *req.KindCode
	}
	if req.KindTag != nil {
		fields["kind_tag"] = *req.KindTag
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

	beforeData, err := h.orgKindService.GetOrgKindByID(id)
	if err != nil {
		utils.BadRequestResponse(ctx, "组织类型不存在")
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	if err := h.orgKindService.UpdateOrgKindFields(id, fields); err != nil {
		utils.BadRequestResponse(ctx, err.Error())
		return
	}

	afterData, _ := h.orgKindService.GetOrgKindByID(id)
	afterJSON, _ := json.Marshal(afterData)

	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    ctx.GetString("trace_id"),
		ModelID:    "org_kind",
		RecordID:   id,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "sso",
	})

	utils.SuccessResponse(ctx, afterData)
}

func (h *SsoOrgKindHandler) DeleteOrgKind(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	if id == "" {
		utils.BadRequestResponse(ctx, "ID不能为空")
		return
	}

	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserAccount)

	beforeData, err := h.orgKindService.GetOrgKindByID(id)
	if err != nil {
		utils.BadRequestResponse(ctx, "组织类型不存在")
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	if err := h.orgKindService.DeleteOrgKind(id); err != nil {
		utils.BadRequestResponse(ctx, err.Error())
		return
	}

	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    ctx.GetString("trace_id"),
		ModelID:    "org_kind",
		RecordID:   id,
		Action:     "DELETE",
		BeforeData: string(beforeJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "sso",
	})

	utils.SuccessResponse(ctx, nil)
}

func (h *SsoOrgKindHandler) GetAllOrgKinds(c context.Context, ctx *app.RequestContext) {
	orgKinds, err := h.orgKindService.GetAllOrgKinds()
	if err != nil {
		utils.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, orgKinds)
}
