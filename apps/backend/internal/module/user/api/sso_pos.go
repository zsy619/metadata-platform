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

// SsoPosHandler 职位处理器结构体
type SsoPosHandler struct {
	*utils.BaseHandler
	posService service.SsoPosService
	audit      AuditService
}

// NewSsoPosHandler 创建职位处理器实例
func NewSsoPosHandler(posService service.SsoPosService, auditQueue *queue.AuditLogQueue) *SsoPosHandler {
	return &SsoPosHandler{
		BaseHandler: utils.NewBaseHandler(),
		posService:  posService,
		audit:       &auditServiceImpl{queue: auditQueue},
	}
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
	ParentID *string `json:"parent_id" form:"parent_id"`
	AppCode  *string `json:"app_code" form:"app_code"`
	OrgID    *string `json:"org_id" form:"org_id"`
	KindCode *string `json:"kind_code" form:"kind_code"`
	PosName  *string `json:"pos_name" form:"pos_name"`
	PosCode  *string `json:"pos_code" form:"pos_code"`
	Status   *int    `json:"state" form:"state"`
	Remark   *string `json:"remark" form:"remark"`
	Sort     *int    `json:"sort" form:"sort"`
}

// CreatePosition 创建职位
func (h *SsoPosHandler) CreatePos(c context.Context, ctx *app.RequestContext) {
	var req SsoCreatePosRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 创建职位模型
	position := &userModel.SsoPos{
		ParentID: req.ParentID,
		AppCode:  req.AppCode,
		OrgID:    req.OrgID,
		KindCode: req.KindCode,
		PosName:  req.PosName,
		PosCode:  req.PosCode,
		Status:   req.Status,
		Remark:   req.Remark,
		Sort:     req.Sort,
		CreateID: headerUser.UserID,
		CreateBy: headerUser.UserAccount,
		TenantID: headerUser.TenantID,
	}

	// 生成职位ID
	position.ID = utils.GetSnowflake().GenerateIDString()

	// 调用服务层创建职位
	if err := h.posService.CreatePos(position); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	afterData, _ := json.Marshal(position)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:        utils.GetSnowflake().GenerateIDString(),
		TraceID:   headerUser.TraceID,
		ModelID:   "pos",
		RecordID:  position.ID,
		Action:    "CREATE",
		AfterData: string(afterData),
		CreateBy:  headerUser.UserAccount,
		Source:    "pos_service",
	})

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

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取变更前数据
	beforeData, err := h.posService.GetPosByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "职位不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	// 更新职位字段
	if req.ParentID != nil {
		beforeData.ParentID = *req.ParentID
	}
	if req.AppCode != nil {
		beforeData.AppCode = *req.AppCode
	}
	if req.OrgID != nil {
		beforeData.OrgID = *req.OrgID
	}
	if req.KindCode != nil {
		beforeData.KindCode = *req.KindCode
	}
	if req.PosName != nil {
		beforeData.PosName = *req.PosName
	}
	if req.PosCode != nil {
		beforeData.PosCode = *req.PosCode
	}
	if req.Status != nil {
		beforeData.Status = *req.Status
	}
	if req.Remark != nil {
		beforeData.Remark = *req.Remark
	}
	if req.Sort != nil {
		beforeData.Sort = *req.Sort
	}

	// 设置更新人信息
	beforeData.UpdateID = headerUser.UserID
	beforeData.UpdateBy = headerUser.UserAccount

	// 调用服务层更新职位
	if err := h.posService.UpdatePos(beforeData); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	afterJSON, _ := json.Marshal(beforeData)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "pos",
		RecordID:   id,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "pos_service",
	})

	ctx.JSON(200, beforeData)
}

// DeletePos 删除职位
func (h *SsoPosHandler) DeletePos(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取删除前数据
	beforeData, err := h.posService.GetPosByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "职位不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	// 调用服务层删除职位
	if err := h.posService.DeletePos(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "pos",
		RecordID:   id,
		Action:     "DELETE",
		BeforeData: string(beforeJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "pos_service",
	})

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
