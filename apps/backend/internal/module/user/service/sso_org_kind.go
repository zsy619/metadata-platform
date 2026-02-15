package service

import (
	"context"
	"encoding/json"
	"fmt"
	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"

	userModel "metadata-platform/internal/module/user/model"
)

// SsoOrgKindService 组织类型服务
type SsoOrgKindService struct {
	db *gorm.DB
}

// NewOrgKindService 创建组织类型服务实例
func NewOrgKindService(db *gorm.DB) *SsoOrgKindService {
	return &SsoOrgKindService{db: db}
}

// CreateOrgKind 创建组织类型
func (s *SsoOrgKindService) CreateOrgKind(ctx context.Context, data *userModel.SsoOrgKind) error {
	return s.db.WithContext(ctx).Create(data).Error
}

// GetOrgKindByID 根据ID获取组织类型
func (s *SsoOrgKindService) GetOrgKindByID(ctx context.Context, id string) (*userModel.SsoOrgKind, error) {
	var orgKind userModel.SsoOrgKind
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&orgKind).Error
	if err != nil {
		return nil, err
	}
	return &orgKind, nil
}

// UpdateOrgKind 更新组织类型
func (s *SsoOrgKindService) UpdateOrgKind(ctx context.Context, id string, data map[string]interface{}) error {
	return s.db.WithContext(ctx).Model(&userModel.SsoOrgKind{}).Where("id = ?", id).Updates(data).Error
}

// DeleteOrgKind 删除组织类型
func (s *SsoOrgKindService) DeleteOrgKind(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&userModel.SsoOrgKind{}).Error
}

// GetAllOrgKinds 获取所有组织类型
func (s *SsoOrgKindService) GetAllOrgKinds(ctx context.Context, tenantID string) ([]userModel.SsoOrgKind, error) {
	var orgKinds []userModel.SsoOrgKind
	query := s.db.WithContext(ctx).Where("is_deleted = ?", false).Order("sort ASC, create_at DESC")
	if tenantID != "" && tenantID != "0" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	err := query.Find(&orgKinds).Error
	return orgKinds, err
}

// CheckOrgKindExists 检查组织类型是否存在
func (s *SsoOrgKindService) CheckOrgKindExists(ctx context.Context, name, code string, tenantID string) (bool, error) {
	var count int64
	query := s.db.WithContext(ctx).Model(&userModel.SsoOrgKind{}).Where("is_deleted = ?", false)
	if tenantID != "" && tenantID != "0" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	err := query.Where("kind_name = ? OR kind_code = ?", name, code).Count(&count).Error
	if err != nil {
		hlog.Error("Failed to check org kind exists: ", err)
		return false, err
	}
	return count > 0, nil
}

// OrgKindHandler 组织类型处理器
type SsoOrgKindHandler struct {
	service *SsoOrgKindService
	audit   AuditService
}

// AuditService 审计服务接口
type AuditService interface {
	RecordDataChange(ctx context.Context, log *model.SysDataChangeLog)
}

// NewOrgKindHandler 创建组织类型处理器
func NewSsoOrgKindHandler(service *SsoOrgKindService, auditQueue *queue.AuditLogQueue) *SsoOrgKindHandler {
	return &SsoOrgKindHandler{
		service: service,
		audit:   &auditServiceImpl{queue: auditQueue},
	}
}

type auditServiceImpl struct {
	queue *queue.AuditLogQueue
}

func (a *auditServiceImpl) RecordDataChange(ctx context.Context, log *model.SysDataChangeLog) {
	if a.queue != nil {
		a.queue.PushDataChange(log)
	}
}

// CreateOrgKind 创建组织类型
func (h *SsoOrgKindHandler) CreateOrgKind(ctx context.Context, c *app.RequestContext) {
	var req struct {
		ParentID string `json:"parent_id"`
		Name     string `json:"name"`
		Code     string `json:"code"`
		Sort     int    `json:"sort"`
		Status   int    `json:"status"`
		Remark   string `json:"remark"`
		TenantID string `json:"tenant_id"`
	}

	if err := c.Bind(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	// 参数验证
	if req.Name == "" {
		utils.BadRequestResponse(c, "组织类型名称不能为空")
		return
	}
	if req.Code == "" {
		utils.BadRequestResponse(c, "组织类型编码不能为空")
		return
	}

	// 获取租户ID
	tenantID := string(c.Request.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		tenantID = "1"
	}

	// 检查是否存在
	exists, err := h.service.CheckOrgKindExists(ctx, req.Name, req.Code, tenantID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "检查组织类型是否存在失败")
		return
	}
	if exists {
		utils.BadRequestResponse(c, "组织类型名称或编码已存在")
		return
	}

	// 获取当前用户ID和账户
	userID := string(c.Request.Header.Get("X-User-ID"))
	userAccount := string(c.Request.Header.Get("X-User-Account"))

	orgKind := &userModel.SsoOrgKind{
		ID:       utils.GetSnowflake().GenerateIDString(),
		ParentID: req.ParentID,
		KindName: req.Name,
		KindCode: req.Code,
		Sort:     req.Sort,
		Status:   req.Status,
		Remark:   req.Remark,
		TenantID: tenantID,
		CreateID: userID,
		CreateBy: userAccount,
	}

	if err := h.service.CreateOrgKind(ctx, orgKind); err != nil {
		hlog.Error("Failed to create org kind: ", err)
		utils.InternalServerErrorResponse(c, "创建组织类型失败")
		return
	}

	// 记录数据变更日志
	afterData, _ := json.Marshal(orgKind)
	h.audit.RecordDataChange(ctx, &model.SysDataChangeLog{
		ID:        utils.GetSnowflake().GenerateIDString(),
		TraceID:   c.GetString("trace_id"),
		ModelID:   "org_kind",
		RecordID:  orgKind.ID,
		Action:    "CREATE",
		AfterData: string(afterData),
		CreateBy:  userAccount,
		Source:    "org_kind_service",
	})

	utils.SuccessResponse(c, orgKind)
}

// GetOrgKindByID 根据ID获取组织类型
func (h *SsoOrgKindHandler) GetOrgKindByID(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequestResponse(c, "ID不能为空")
		return
	}

	orgKind, err := h.service.GetOrgKindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "组织类型不存在")
			return
		}
		hlog.Error("Failed to get org kind: ", err)
		utils.InternalServerErrorResponse(c, "获取组织类型失败")
		return
	}

	utils.SuccessResponse(c, orgKind)
}

// UpdateOrgKind 更新组织类型
func (h *SsoOrgKindHandler) UpdateOrgKind(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequestResponse(c, "ID不能为空")
		return
	}

	var req struct {
		ParentID *string `json:"parent_id"`
		Name     *string `json:"name"`
		Code     *string `json:"code"`
		Sort     *int    `json:"sort"`
		Status   *int    `json:"status"`
		Remark   *string `json:"remark"`
	}

	if err := c.Bind(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	// 获取当前用户ID和账户
	userID := string(c.Request.Header.Get("X-User-ID"))
	userAccount := string(c.Request.Header.Get("X-User-Account"))

	// 如果用户ID为空，设置为默认值
	if userID == "" {
		userID = "0"
	}

	// 如果用户ID为0（未登录），则update_id和update_by都设置为空字符串
	var updateID, updateBy string
	if userID == "0" {
		updateID = ""
		updateBy = ""
	} else {
		updateID = userID
		updateBy = userAccount
	}

	// 构建更新字段（只更新传递的字段）
	fields := map[string]any{
		"update_id": updateID,
		"update_by": updateBy,
	}

	if req.ParentID != nil {
		fields["parent_id"] = *req.ParentID
	}
	if req.Name != nil {
		fields["kind_name"] = *req.Name
	}
	if req.Code != nil {
		fields["kind_code"] = *req.Code
	}
	if req.Sort != nil {
		fields["sort"] = *req.Sort
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}
	if req.Remark != nil {
		fields["remark"] = *req.Remark
	}

	// 获取变更前数据
	beforeData, err := h.service.GetOrgKindByID(ctx, id)
	if err != nil {
		utils.BadRequestResponse(c, "组织类型不存在")
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	if err := h.service.UpdateOrgKind(ctx, id, fields); err != nil {
		hlog.Error("Failed to update org kind: ", err)
		utils.InternalServerErrorResponse(c, "更新组织类型失败")
		return
	}

	// 获取变更后数据
	afterData, _ := h.service.GetOrgKindByID(ctx, id)
	afterJSON, _ := json.Marshal(afterData)

	// 记录数据变更日志
	h.audit.RecordDataChange(ctx, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    c.GetString("trace_id"),
		ModelID:    "org_kind",
		RecordID:   id,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   updateBy,
		Source:     "org_kind_service",
	})

	utils.SuccessResponse(c, afterData)
}

// DeleteOrgKind 删除组织类型
func (h *SsoOrgKindHandler) DeleteOrgKind(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequestResponse(c, "ID不能为空")
		return
	}

	// 获取当前用户ID和账户
	userID := string(c.Request.Header.Get("X-User-ID"))
	userAccount := string(c.Request.Header.Get("X-User-Account"))
	fmt.Println(userID, userAccount)

	// 获取删除前数据
	beforeData, err := h.service.GetOrgKindByID(ctx, id)
	if err != nil {
		utils.BadRequestResponse(c, "组织类型不存在")
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	if err := h.service.DeleteOrgKind(ctx, id); err != nil {
		hlog.Error("Failed to delete org kind: ", err)
		utils.InternalServerErrorResponse(c, "删除组织类型失败")
		return
	}

	// 记录数据变更日志
	h.audit.RecordDataChange(ctx, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    c.GetString("trace_id"),
		ModelID:    "org_kind",
		RecordID:   id,
		Action:     "DELETE",
		BeforeData: string(beforeJSON),
		CreateBy:   userAccount,
		Source:     "org_kind_service",
	})

	utils.SuccessResponse(c, nil)
}

// GetAllOrgKinds 获取所有组织类型
func (h *SsoOrgKindHandler) GetAllOrgKinds(ctx context.Context, c *app.RequestContext) {
	tenantID := string(c.Request.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		tenantID = "1"
	}

	orgKinds, err := h.service.GetAllOrgKinds(ctx, tenantID)
	if err != nil {
		hlog.Error("Failed to get all org kinds: ", err)
		utils.InternalServerErrorResponse(c, "获取组织类型列表失败")
		return
	}

	utils.SuccessResponse(c, orgKinds)
}
