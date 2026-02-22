package api

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"

	auditModel "metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/queue"
	userModel "metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"
)

// SsoUserContactHandler 用户联系方式处理器
type SsoUserContactHandler struct {
	*utils.BaseHandler
	contactSvc service.SsoUserContactService
	audit      AuditService
}

// NewSsoUserContactHandler 创建用户联系方式处理器实例
func NewSsoUserContactHandler(contactSvc service.SsoUserContactService, auditQueue *queue.AuditLogQueue) *SsoUserContactHandler {
	return &SsoUserContactHandler{
		BaseHandler: utils.NewBaseHandler(),
		contactSvc:  contactSvc,
		audit:       &auditServiceImpl{queue: auditQueue},
	}
}

// GetContacts 获取联系方式列表
func (h *SsoUserContactHandler) GetContacts(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	list, err := h.contactSvc.GetByUserID(userID)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取联系方式失败"})
		return
	}
	ctx.JSON(200, list)
}

// CreateContactRequest 新增联系方式请求
type CreateContactRequest struct {
	Type   string `json:"type" form:"type" binding:"required"`
	Value  string `json:"value" form:"value" binding:"required"`
	Remark string `json:"remark" form:"remark"`
}

// CreateContact 新增联系方式
func (h *SsoUserContactHandler) CreateContact(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	headerUser := h.GetHeaderUserStruct(c, ctx)

	var req CreateContactRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	contact := &userModel.SsoUserContact{
		UserID:   userID,
		Type:     req.Type,
		Value:    req.Value,
		Remark:   req.Remark,
		TenantID: headerUser.TenantID,
		CreateID: headerUser.UserID,
		CreateBy: headerUser.UserAccount,
		UpdateID: headerUser.UserID,
		UpdateBy: headerUser.UserAccount,
	}

	if err := h.contactSvc.Create(contact); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	afterJSON, _ := json.Marshal(contact)
	h.audit.RecordDataChange(c, &auditModel.SysDataChangeLog{
		ID:        utils.GetSnowflake().GenerateIDString(),
		TraceID:   headerUser.TraceID,
		ModelID:   "user_contact",
		RecordID:  contact.ID,
		Action:    "CREATE",
		AfterData: string(afterJSON),
		CreateBy:  headerUser.UserAccount,
		Source:    "user_contact_service",
	})

	ctx.JSON(201, contact)
}

// UpdateContactRequest 更新联系方式请求
type UpdateContactRequest struct {
	Type   string `json:"type" form:"type"`
	Value  string `json:"value" form:"value"`
	Remark string `json:"remark" form:"remark"`
}

// UpdateContact 更新联系方式
func (h *SsoUserContactHandler) UpdateContact(c context.Context, ctx *app.RequestContext) {
	contactID := ctx.Param("cid")
	headerUser := h.GetHeaderUserStruct(c, ctx)

	var req UpdateContactRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	fields := map[string]any{
		"type":      req.Type,
		"value":     req.Value,
		"remark":    req.Remark,
		"update_id": headerUser.UserID,
		"update_by": headerUser.UserAccount,
	}

	if err := h.contactSvc.UpdateFields(contactID, fields); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	ctx.JSON(200, map[string]string{"message": "更新成功"})
}

// DeleteContact 删除联系方式
func (h *SsoUserContactHandler) DeleteContact(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	contactID := ctx.Param("cid")
	if err := h.contactSvc.Delete(userID, contactID); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	ctx.JSON(200, map[string]string{"message": "删除成功"})
}
