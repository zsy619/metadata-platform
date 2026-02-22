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

// SsoUserAddressHandler 用户地址处理器
type SsoUserAddressHandler struct {
	*utils.BaseHandler
	addrSvc service.SsoUserAddressService
	audit   AuditService
}

// NewSsoUserAddressHandler 创建用户地址处理器实例
func NewSsoUserAddressHandler(addrSvc service.SsoUserAddressService, auditQueue *queue.AuditLogQueue) *SsoUserAddressHandler {
	return &SsoUserAddressHandler{
		BaseHandler: utils.NewBaseHandler(),
		addrSvc:     addrSvc,
		audit:       &auditServiceImpl{queue: auditQueue},
	}
}

// GetAddresses 获取用户地址列表
func (h *SsoUserAddressHandler) GetAddresses(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	list, err := h.addrSvc.GetByUserID(userID)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取地址列表失败"})
		return
	}
	ctx.JSON(200, list)
}

// CreateAddressRequest 新增地址请求
type CreateAddressRequest struct {
	Label        string `json:"label" form:"label"`
	ReceiverName string `json:"receiver_name" form:"receiver_name" binding:"required"`
	Phone        string `json:"phone" form:"phone" binding:"required"`
	Province     string `json:"province" form:"province"`
	City         string `json:"city" form:"city"`
	District     string `json:"district" form:"district"`
	Detail       string `json:"detail" form:"detail" binding:"required"`
	Remark       string `json:"remark" form:"remark"`
}

// CreateAddress 新增地址
func (h *SsoUserAddressHandler) CreateAddress(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	headerUser := h.GetHeaderUserStruct(c, ctx)

	var req CreateAddressRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	addr := &userModel.SsoUserAddress{
		UserID:       userID,
		Label:        req.Label,
		ReceiverName: req.ReceiverName,
		Phone:        req.Phone,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Detail:       req.Detail,
		Remark:       req.Remark,
		TenantID:     headerUser.TenantID,
		CreateID:     headerUser.UserID,
		CreateBy:     headerUser.UserAccount,
		UpdateID:     headerUser.UserID,
		UpdateBy:     headerUser.UserAccount,
	}

	if err := h.addrSvc.Create(addr); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	afterJSON, _ := json.Marshal(addr)
	h.audit.RecordDataChange(c, &auditModel.SysDataChangeLog{
		ID:        utils.GetSnowflake().GenerateIDString(),
		TraceID:   headerUser.TraceID,
		ModelID:   "user_address",
		RecordID:  addr.ID,
		Action:    "CREATE",
		AfterData: string(afterJSON),
		CreateBy:  headerUser.UserAccount,
		Source:    "user_address_service",
	})

	ctx.JSON(201, addr)
}

// UpdateAddressRequest 更新地址请求
type UpdateAddressRequest struct {
	Label        string `json:"label" form:"label"`
	ReceiverName string `json:"receiver_name" form:"receiver_name"`
	Phone        string `json:"phone" form:"phone"`
	Province     string `json:"province" form:"province"`
	City         string `json:"city" form:"city"`
	District     string `json:"district" form:"district"`
	Detail       string `json:"detail" form:"detail"`
	Remark       string `json:"remark" form:"remark"`
}

// UpdateAddress 更新地址
func (h *SsoUserAddressHandler) UpdateAddress(c context.Context, ctx *app.RequestContext) {
	addrID := ctx.Param("aid")
	headerUser := h.GetHeaderUserStruct(c, ctx)

	var req UpdateAddressRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	fields := map[string]any{
		"label":         req.Label,
		"receiver_name": req.ReceiverName,
		"phone":         req.Phone,
		"province":      req.Province,
		"city":          req.City,
		"district":      req.District,
		"detail":        req.Detail,
		"remark":        req.Remark,
		"update_id":     headerUser.UserID,
		"update_by":     headerUser.UserAccount,
	}

	if err := h.addrSvc.UpdateFields(addrID, fields); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	ctx.JSON(200, map[string]string{"message": "更新成功"})
}

// SetDefaultAddress 设置默认地址
func (h *SsoUserAddressHandler) SetDefaultAddress(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	addrID := ctx.Param("aid")
	if err := h.addrSvc.SetDefault(userID, addrID); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	ctx.JSON(200, map[string]string{"message": "设置默认地址成功"})
}

// DeleteAddress 删除地址
func (h *SsoUserAddressHandler) DeleteAddress(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	addrID := ctx.Param("aid")
	if err := h.addrSvc.Delete(userID, addrID); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	ctx.JSON(200, map[string]string{"message": "删除成功"})
}
