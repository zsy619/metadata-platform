package api

import (
	"context"
	"encoding/json"
	"fmt"
	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"

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
	ParentID  string `json:"parent_id" form:"parent_id"`
	AppCode   string `json:"app_code" form:"app_code"`
	OrgID     string `json:"org_id" form:"org_id"`
	KindCode  string `json:"kind_code" form:"kind_code" binding:"required"`
	PosName   string `json:"pos_name" form:"pos_name" binding:"required"`
	PosCode   string `json:"pos_code" form:"pos_code" binding:"required"`
	DataRange string `json:"data_range" form:"data_range"`
	DataScope string `json:"data_scope" form:"data_scope"`
	Status    int    `json:"status" form:"status"`
	Remark    string `json:"remark" form:"remark"`
	Sort      int    `json:"sort" form:"sort"`
}

// SsoUpdatePosRequest 更新职位请求结构
type SsoUpdatePosRequest struct {
	ParentID  string `json:"parent_id" form:"parent_id"`
	AppCode   string `json:"app_code" form:"app_code"`
	OrgID     string `json:"org_id" form:"org_id"`
	DataRange string `json:"data_range" form:"data_range"`
	DataScope string `json:"data_scope" form:"data_scope"`
	KindCode  string `json:"kind_code" form:"kind_code"`
	PosName   string `json:"pos_name" form:"pos_name"`
	PosCode   string `json:"pos_code" form:"pos_code"`
	Status    int    `json:"status" form:"status"`
	Remark    string `json:"remark" form:"remark"`
	Sort      int    `json:"sort" form:"sort"`
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
		ParentID:  req.ParentID,
		AppCode:   req.AppCode,
		OrgID:     req.OrgID,
		KindCode:  req.KindCode,
		DataRange: req.DataRange,
		DataScope: req.DataScope,
		PosName:   req.PosName,
		PosCode:   req.PosCode,
		Status:    req.Status,
		Remark:    req.Remark,
		Sort:      req.Sort,
		CreateID:  headerUser.UserID,
		CreateBy:  headerUser.UserAccount,
		TenantID:  headerUser.TenantID,
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

	// 构建更新字段（只更新传递的字段）
	// 使用 map 方式实现部分字段更新，避免全量更新
	fields := map[string]any{
		"update_id": headerUser.UserID,
		"update_by": headerUser.UserAccount,
	}

	// 根据请求中的字段构建更新 map
	fields["parent_id"] = req.ParentID
	fields["app_code"] = req.AppCode
	fields["org_id"] = req.OrgID
	fields["kind_code"] = req.KindCode
	fields["pos_name"] = req.PosName
	fields["pos_code"] = req.PosCode
	fields["status"] = req.Status
	fields["data_range"] = req.DataRange
	fields["data_scope"] = req.DataScope
	fields["remark"] = req.Remark
	fields["sort"] = req.Sort

	// 调用服务层更新职位（使用 map 方式）
	if err := h.posService.UpdatePosFields(id, fields); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取更新后的职位数据
	afterData, err := h.posService.GetPosByID(id)
	if err != nil {
		ctx.JSON(200, map[string]string{"message": "更新成功"})
		return
	}

	// 记录数据变更日志
	afterJSON, _ := json.Marshal(afterData)
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

	ctx.JSON(200, afterData)
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

// GetPosRoles 获取职位的角色列表
func (h *SsoPosHandler) GetPosRoles(c context.Context, ctx *app.RequestContext) {
	posID := ctx.Param("id")

	// 调用服务层获取职位的角色ID列表
	roleIDs, err := h.posService.GetPosRoles(posID)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, map[string]any{"role_ids": roleIDs})
}

// UpdatePosRolesRequest 更新职位角色请求结构
type UpdatePosRolesRequest struct {
	RoleIDs []string `json:"role_ids"`
}

// UpdatePosRoles 更新职位的角色关联
func (h *SsoPosHandler) UpdatePosRoles(c context.Context, ctx *app.RequestContext) {
	posID := ctx.Param("id")
	var req UpdatePosRolesRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取当前用户信息
	headerUser := h.GetHeaderUserStruct(c, ctx)

	// 获取更新前的角色列表（用于审计日志）
	beforeRoleIDs, _ := h.posService.GetPosRoles(posID)
	beforeData, _ := json.Marshal(map[string]any{"role_ids": beforeRoleIDs})

	// 调用服务层更新职位角色关联
	if err := h.posService.UpdatePosRoles(posID, req.RoleIDs, headerUser.UserAccount); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取更新后的角色列表（用于审计日志）
	afterRoleIDs, _ := h.posService.GetPosRoles(posID)
	afterData, _ := json.Marshal(map[string]any{"role_ids": afterRoleIDs})

	// 记录数据变更日志
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "pos_role",
		RecordID:   posID,
		Action:     "UPDATE_ROLES",
		BeforeData: string(beforeData),
		AfterData:  string(afterData),
		CreateBy:   headerUser.UserAccount,
		Source:     "pos_service",
	})

	ctx.JSON(200, map[string]string{"message": "更新成功"})
}
