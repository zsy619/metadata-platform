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
	AppCode   string `json:"app_code" form:"app_code"`
	MenuName  string `json:"menu_name" form:"menu_name" binding:"required"`
	MenuCode  string `json:"menu_code" form:"menu_code" binding:"required"`
	Status    int    `json:"status" form:"status"`
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
	ParentID  string `json:"parent_id" form:"parent_id"`
	AppCode   string `json:"app_code" form:"app_code"`
	MenuName  string `json:"menu_name" form:"menu_name"`
	MenuCode  string `json:"menu_code" form:"menu_code"`
	Status    int    `json:"status" form:"status"`
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
		TraceID:   headerUser.TraceID,
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

	// 构建更新字段（只更新传递的字段）
	// 使用 map 方式实现部分字段更新，避免全量更新
	fields := map[string]any{
		"update_id": headerUser.UserID,
		"update_by": headerUser.UserAccount,
	}

	// 根据请求中的字段构建更新 map
	fields["parent_id"] = req.ParentID
	fields["app_code"] = req.AppCode
	fields["menu_name"] = req.MenuName
	fields["menu_code"] = req.MenuCode
	fields["status"] = req.Status
	fields["data_range"] = req.DataRange
	fields["data_scope"] = req.DataScope
	fields["is_visible"] = req.IsVisible
	fields["menu_type"] = req.MenuType
	fields["icon"] = req.Icon
	fields["url"] = req.URL
	fields["method"] = req.Method
	fields["target"] = req.Target
	fields["remark"] = req.Remark
	fields["sort"] = req.Sort
	fields["tier"] = req.Tier

	// 调用服务层更新菜单（使用 map 方式）
	if err := h.menuService.UpdateMenuFields(id, fields); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取更新后的菜单数据
	afterData, err := h.menuService.GetMenuByID(id)
	if err != nil {
		ctx.JSON(200, map[string]string{"message": "更新成功"})
		return
	}

	// 记录数据变更日志
	afterJSON, _ := json.Marshal(afterData)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "menu",
		RecordID:   id,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "menu_service",
	})

	ctx.JSON(200, afterData)
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
		TraceID:    headerUser.TraceID,
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
