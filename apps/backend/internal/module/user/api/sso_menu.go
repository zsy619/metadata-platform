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

// SsoMenuHandler 菜单处理器结构体
type SsoMenuHandler struct {
	*utils.BaseHandler
	menuService service.SsoMenuService
	audit       AuditService
}

// NewSsoMenuHandler 创建菜单处理器实例
func NewSsoMenuHandler(menuService service.SsoMenuService, auditQueue *queue.AuditLogQueue) *SsoMenuHandler {
	return &SsoMenuHandler{
		BaseHandler: utils.NewBaseHandler(),
		menuService: menuService,
		audit:       &auditServiceImpl{queue: auditQueue},
	}
}

// SsoCreateMenuRequest 创建菜单请求结构
type SsoCreateMenuRequest struct {
	ParentID  string `json:"parent_id" form:"parent_id"`
	AppCode   string `json:"application_code" form:"application_code"`
	MenuName  string `json:"menu_name" form:"menu_name" binding:"required"`
	MenuCode  string `json:"menu_code" form:"menu_code" binding:"required"`
	Status    int    `json:"state" form:"state"`
	DataRange string `json:"data_range" form:"data_range"`
	DataScope string `json:"data_scope" form:"data_scope"`
	IsVisible bool   `json:"is_visible" form:"is_visible"`
	MenuType  string `json:"menu_type" form:"menu_type"`
	Icon      string `json:"icon" form:"icon"`
	URL       string `json:"url" form:"url"`
	Method    string `json:"method" form:"method"`
	Target    string `json:"target" form:"target"`
	Remark    string `json:"remark" form:"remark"`
	Sort      int    `json:"sort" form:"sort"`
	Tier      int    `json:"tier" form:"tier"`
}

// SsoUpdateMenuRequest 更新菜单请求结构
type SsoUpdateMenuRequest struct {
	ParentID  *string `json:"parent_id" form:"parent_id"`
	AppCode   *string `json:"application_code" form:"application_code"`
	MenuName  *string `json:"menu_name" form:"menu_name"`
	MenuCode  *string `json:"menu_code" form:"menu_code"`
	Status    *int    `json:"state" form:"state"`
	DataRange *string `json:"data_range" form:"data_range"`
	DataScope *string `json:"data_scope" form:"data_scope"`
	IsVisible *bool   `json:"is_visible" form:"is_visible"`
	MenuType  *string `json:"menu_type" form:"menu_type"`
	Icon      *string `json:"icon" form:"icon"`
	URL       *string `json:"url" form:"url"`
	Method    *string `json:"method" form:"method"`
	Target    *string `json:"target" form:"target"`
	Remark    *string `json:"remark" form:"remark"`
	Sort      *int    `json:"sort" form:"sort"`
	Tier      *int    `json:"tier" form:"tier"`
}

// CreateMenu 创建菜单
func (h *SsoMenuHandler) CreateMenu(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateMenuRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 创建菜单模型
	menu := &userModel.SsoMenu{
		ParentID:  req.ParentID,
		AppCode:   req.AppCode,
		MenuName:  req.MenuName,
		MenuCode:  req.MenuCode,
		Status:    req.Status,
		DataRange: req.DataRange,
		DataScope: req.DataScope,
		IsVisible: req.IsVisible,
		MenuType:  req.MenuType,
		Icon:      req.Icon,
		URL:       req.URL,
		Method:    req.Method,
		Target:    req.Target,
		Remark:    req.Remark,
		Sort:      req.Sort,
		Tier:      req.Tier,
		CreateID:  headerUser.UserID,
		CreateBy:  headerUser.UserAccount,
		TenantID:  headerUser.TenantID,
	}

	// 调用服务层创建菜单
	if err := h.menuService.CreateMenu(menu); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	afterData, _ := json.Marshal(menu)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:        utils.GetSnowflake().GenerateIDString(),
		TraceID:   ctx.GetString("trace_id"),
		ModelID:   "menu",
		RecordID:  menu.ID,
		Action:    "CREATE",
		AfterData: string(afterData),
		CreateBy:  headerUser.UserAccount,
		Source:    "menu_service",
	})

	ctx.JSON(201, menu)
}

// GetMenuByID 根据ID获取菜单
func (h *SsoMenuHandler) GetMenuByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取菜单
	menu, err := h.menuService.GetMenuByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "菜单不存在"})
		return
	}

	ctx.JSON(200, menu)
}

// UpdateMenu 更新菜单
func (h *SsoMenuHandler) UpdateMenu(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req SsoUpdateMenuRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取变更前数据
	beforeData, err := h.menuService.GetMenuByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "菜单不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	// 更新菜单字段
	if req.ParentID != nil {
		beforeData.ParentID = *req.ParentID
	}
	if req.AppCode != nil {
		beforeData.AppCode = *req.AppCode
	}
	if req.MenuName != nil {
		beforeData.MenuName = *req.MenuName
	}
	if req.MenuCode != nil {
		beforeData.MenuCode = *req.MenuCode
	}
	if req.Status != nil {
		beforeData.Status = *req.Status
	}
	if req.DataRange != nil {
		beforeData.DataRange = *req.DataRange
	}
	if req.DataScope != nil {
		beforeData.DataScope = *req.DataScope
	}
	if req.IsVisible != nil {
		beforeData.IsVisible = *req.IsVisible
	}
	if req.MenuType != nil {
		beforeData.MenuType = *req.MenuType
	}
	if req.Icon != nil {
		beforeData.Icon = *req.Icon
	}
	if req.URL != nil {
		beforeData.URL = *req.URL
	}
	if req.Method != nil {
		beforeData.Method = *req.Method
	}
	if req.Target != nil {
		beforeData.Target = *req.Target
	}
	if req.Remark != nil {
		beforeData.Remark = *req.Remark
	}
	if req.Sort != nil {
		beforeData.Sort = *req.Sort
	}
	if req.Tier != nil {
		beforeData.Tier = *req.Tier
	}

	// 设置更新人信息
	beforeData.UpdateID = headerUser.UserID
	beforeData.UpdateBy = headerUser.UserAccount

	// 调用服务层更新菜单
	if err := h.menuService.UpdateMenu(beforeData); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	afterJSON, _ := json.Marshal(beforeData)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    ctx.GetString("trace_id"),
		ModelID:    "menu",
		RecordID:   id,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "menu_service",
	})

	ctx.JSON(200, beforeData)
}

// DeleteMenu 删除菜单
func (h *SsoMenuHandler) DeleteMenu(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取删除前数据
	beforeData, err := h.menuService.GetMenuByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "菜单不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	// 调用服务层删除菜单
	if err := h.menuService.DeleteMenu(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    ctx.GetString("trace_id"),
		ModelID:    "menu",
		RecordID:   id,
		Action:     "DELETE",
		BeforeData: string(beforeJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "menu_service",
	})

	ctx.JSON(200, map[string]string{"message": "菜单删除成功"})
}

// GetAllMenus 获取所有菜单
func (h *SsoMenuHandler) GetAllMenus(c context.Context, ctx *app.RequestContext) {
	// 调用服务层获取所有菜单
	menus, err := h.menuService.GetAllMenus()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取菜单列表失败"})
		return
	}

	ctx.JSON(200, menus)
}
