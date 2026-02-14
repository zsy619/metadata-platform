package api

import (
	"context"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"

	"github.com/cloudwego/hertz/pkg/app"
)

// SsoMenuHandler 菜单处理器结构体
type SsoMenuHandler struct {
	menuService service.SsoMenuService
}

// NewSsoMenuHandler 创建菜单处理器实例
func NewSsoMenuHandler(menuService service.SsoMenuService) *SsoMenuHandler {
	return &SsoMenuHandler{menuService: menuService}
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
	Visible   int    `json:"visible" form:"visible"`
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
	AppCode   string `json:"application_code" form:"application_code"`
	MenuName  string `json:"menu_name" form:"menu_name"`
	MenuCode  string `json:"menu_code" form:"menu_code"`
	Status    int    `json:"state" form:"state"`
	DataRange string `json:"data_range" form:"data_range"`
	DataScope string `json:"data_scope" form:"data_scope"`
	Visible   int    `json:"visible" form:"visible"`
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

	// 创建菜单模型
	menu := &model.SsoMenu{
		ParentID:  req.ParentID,
		AppCode:   req.AppCode,
		MenuName:  req.MenuName,
		MenuCode:  req.MenuCode,
		Status:    req.Status,
		DataRange: req.DataRange,
		DataScope: req.DataScope,
		Visible:   req.Visible,
		MenuType:  req.MenuType,
		Icon:      req.Icon,
		URL:       req.URL,
		Method:    req.Method,
		Target:    req.Target,
		Remark:    req.Remark,
		Sort:      req.Sort,
		Tier:      req.Tier,
	}

	// 调用服务层创建菜单
	if err := h.menuService.CreateMenu(menu); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

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

	// 调用服务层获取菜单
	menu, err := h.menuService.GetMenuByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "菜单不存在"})
		return
	}

	// 更新菜单字段
	if req.ParentID != "" {
		menu.ParentID = req.ParentID
	}
	if req.AppCode != "" {
		menu.AppCode = req.AppCode
	}
	if req.MenuName != "" {
		menu.MenuName = req.MenuName
	}
	if req.MenuCode != "" {
		menu.MenuCode = req.MenuCode
	}
	if req.Status != 0 {
		menu.Status = req.Status
	}
	if req.DataRange != "" {
		menu.DataRange = req.DataRange
	}
	if req.DataScope != "" {
		menu.DataScope = req.DataScope
	}
	if req.Visible != 0 {
		menu.Visible = req.Visible
	}
	if req.MenuType != "" {
		menu.MenuType = req.MenuType
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	if req.URL != "" {
		menu.URL = req.URL
	}
	if req.Method != "" {
		menu.Method = req.Method
	}
	if req.Target != "" {
		menu.Target = req.Target
	}
	if req.Remark != "" {
		menu.Remark = req.Remark
	}
	if req.Sort != 0 {
		menu.Sort = req.Sort
	}
	if req.Tier != 0 {
		menu.Tier = req.Tier
	}

	// 调用服务层更新菜单
	if err := h.menuService.UpdateMenu(menu); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, menu)
}

// DeleteMenu 删除菜单
func (h *SsoMenuHandler) DeleteMenu(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层删除菜单
	if err := h.menuService.DeleteMenu(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

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
