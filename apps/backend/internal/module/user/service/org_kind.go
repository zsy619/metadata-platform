package service

import (
	"context"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

// OrgKindService 组织类型服务
type OrgKindService struct {
	db *gorm.DB
}

// NewOrgKindService 创建组织类型服务实例
func NewOrgKindService(db *gorm.DB) *OrgKindService {
	return &OrgKindService{db: db}
}

// CreateOrgKind 创建组织类型
func (s *OrgKindService) CreateOrgKind(ctx context.Context, data *model.SsoOrgKind) error {
	return s.db.WithContext(ctx).Create(data).Error
}

// GetOrgKindByID 根据ID获取组织类型
func (s *OrgKindService) GetOrgKindByID(ctx context.Context, id string) (*model.SsoOrgKind, error) {
	var orgKind model.SsoOrgKind
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&orgKind).Error
	if err != nil {
		return nil, err
	}
	return &orgKind, nil
}

// UpdateOrgKind 更新组织类型
func (s *OrgKindService) UpdateOrgKind(ctx context.Context, id string, data map[string]interface{}) error {
	return s.db.WithContext(ctx).Model(&model.SsoOrgKind{}).Where("id = ?", id).Updates(data).Error
}

// DeleteOrgKind 删除组织类型
func (s *OrgKindService) DeleteOrgKind(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&model.SsoOrgKind{}).Error
}

// GetAllOrgKinds 获取所有组织类型
func (s *OrgKindService) GetAllOrgKinds(ctx context.Context, tenantID string) ([]model.SsoOrgKind, error) {
	var orgKinds []model.SsoOrgKind
	query := s.db.WithContext(ctx).Where("is_deleted = ?", false).Order("sort ASC, create_at DESC")
	if tenantID != "" && tenantID != "0" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	err := query.Find(&orgKinds).Error
	return orgKinds, err
}

// CheckOrgKindExists 检查组织类型是否存在
func (s *OrgKindService) CheckOrgKindExists(ctx context.Context, name, code string, tenantID string) (bool, error) {
	var count int64
	query := s.db.WithContext(ctx).Model(&model.SsoOrgKind{}).Where("is_deleted = ?", false)
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
type OrgKindHandler struct {
	service *OrgKindService
}

// NewOrgKindHandler 创建组织类型处理器
func NewOrgKindHandler(service *OrgKindService) *OrgKindHandler {
	return &OrgKindHandler{service: service}
}

// CreateOrgKind 创建组织类型
func (h *OrgKindHandler) CreateOrgKind(ctx context.Context, c *app.RequestContext) {
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
		tenantID = "0"
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

	// 获取当前用户ID
	userID := string(c.Request.Header.Get("X-User-ID"))
	if userID == "" {
		userID = "0"
	}

	orgKind := &model.SsoOrgKind{
		ID:       utils.GetSnowflake().GenerateIDString(),
		ParentID: req.ParentID,
		KindName: req.Name,
		KindCode: req.Code,
		Sort:     req.Sort,
		Status:   req.Status,
		Remark:   req.Remark,
		TenantID: tenantID,
		CreateID: userID,
		CreateBy: userID,
	}

	if err := h.service.CreateOrgKind(ctx, orgKind); err != nil {
		hlog.Error("Failed to create org kind: ", err)
		utils.InternalServerErrorResponse(c, "创建组织类型失败")
		return
	}

	utils.SuccessResponse(c, orgKind)
}

// GetOrgKindByID 根据ID获取组织类型
func (h *OrgKindHandler) GetOrgKindByID(ctx context.Context, c *app.RequestContext) {
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
func (h *OrgKindHandler) UpdateOrgKind(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequestResponse(c, "ID不能为空")
		return
	}

	var req struct {
		ParentID string `json:"parent_id"`
		Name     string `json:"name"`
		Code     string `json:"code"`
		Sort     int    `json:"sort"`
		Status   int    `json:"status"`
		Remark   string `json:"remark"`
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

	// 获取当前用户ID
	userID := string(c.Request.Header.Get("X-User-ID"))
	if userID == "" {
		userID = "0"
	}

	data := map[string]any{
		"parent_id": req.ParentID,
		"kind_name": req.Name,
		"kind_code": req.Code,
		"sort":      req.Sort,
		"status":    req.Status,
		"remark":    req.Remark,
		"update_id": userID,
		"update_by": userID,
	}

	if err := h.service.UpdateOrgKind(ctx, id, data); err != nil {
		hlog.Error("Failed to update org kind: ", err)
		utils.InternalServerErrorResponse(c, "更新组织类型失败")
		return
	}

	utils.SuccessResponse(c, nil)
}

// DeleteOrgKind 删除组织类型
func (h *OrgKindHandler) DeleteOrgKind(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequestResponse(c, "ID不能为空")
		return
	}

	if err := h.service.DeleteOrgKind(ctx, id); err != nil {
		hlog.Error("Failed to delete org kind: ", err)
		utils.InternalServerErrorResponse(c, "删除组织类型失败")
		return
	}

	utils.SuccessResponse(c, nil)
}

// GetAllOrgKinds 获取所有组织类型
func (h *OrgKindHandler) GetAllOrgKinds(ctx context.Context, c *app.RequestContext) {
	tenantID := string(c.Request.Header.Get("X-Tenant-ID"))

	orgKinds, err := h.service.GetAllOrgKinds(ctx, tenantID)
	if err != nil {
		hlog.Error("Failed to get all org kinds: ", err)
		utils.InternalServerErrorResponse(c, "获取组织类型列表失败")
		return
	}

	utils.SuccessResponse(c, orgKinds)
}
