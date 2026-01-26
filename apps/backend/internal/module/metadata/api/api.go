package api

import (
	"context"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// APIHandler API处理器结构体
type APIHandler struct {
	apiService service.APIService
}

// NewAPIHandler 创建API处理器实例
func NewAPIHandler(apiService service.APIService) *APIHandler {
	return &APIHandler{apiService: apiService}
}

// CreateAPIRequest 创建API请求结构
type CreateAPIRequest struct {
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Path     string `json:"path" binding:"required"`
	Method   string `json:"method" binding:"required"`
	IsPublic bool   `json:"is_public"`
	State    int    `json:"state"`
	Remark   string `json:"remark"`
	Sort     int    `json:"sort"`
}

// UpdateAPIRequest 更新API请求结构
type UpdateAPIRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Path     string `json:"path"`
	Method   string `json:"method"`
	IsPublic bool   `json:"is_public"`
	State    int    `json:"state"`
	Remark   string `json:"remark"`
	Sort     int    `json:"sort"`
}

// CreateAPI 创建API
func (h *APIHandler) CreateAPI(c context.Context, ctx *app.RequestContext) {
	var req CreateAPIRequest
	if err := ctx.BindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	// 创建API模型
	api := &model.API{
		Name:     req.Name,
		Code:     req.Code,
		Path:     req.Path,
		Method:   req.Method,
		IsPublic: req.IsPublic,
		State:    req.State,
		Remark:   req.Remark,
		Sort:     req.Sort,
	}

	// 调用服务层创建API
	if err := h.apiService.CreateAPI(api); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, api)
}

// GetAPIByID 根据ID获取API
func (h *APIHandler) GetAPIByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取API
	api, err := h.apiService.GetAPIByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "API不存在")
		return
	}

	utils.SuccessResponse(ctx, api)
}

// UpdateAPI 更新API
func (h *APIHandler) UpdateAPI(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req UpdateAPIRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 调用服务层获取API
	api, err := h.apiService.GetAPIByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "API不存在"})
		return
	}

	// 更新API字段
	if req.Name != "" {
		api.Name = req.Name
	}
	if req.Code != "" {
		api.Code = req.Code
	}
	if req.Path != "" {
		api.Path = req.Path
	}
	if req.Method != "" {
		api.Method = req.Method
	}
	api.IsPublic = req.IsPublic
	if req.State != 0 {
		api.State = req.State
	}
	if req.Remark != "" {
		api.Remark = req.Remark
	}
	if req.Sort != 0 {
		api.Sort = req.Sort
	}

	// 调用服务层更新API
	if err := h.apiService.UpdateAPI(api); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, api)
}

// DeleteAPI 删除API
func (h *APIHandler) DeleteAPI(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层删除API
	if err := h.apiService.DeleteAPI(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, map[string]string{"message": "API删除成功"})
}

// GetAllAPIs 获取所有API
func (h *APIHandler) GetAllAPIs(c context.Context, ctx *app.RequestContext) {
	// 调用服务层获取所有API
	apis, err := h.apiService.GetAllAPIs()
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, "获取API列表失败")
		return
	}

	utils.SuccessResponse(ctx, apis)
}

// EnableAPI 启用/禁用API
func (h *APIHandler) EnableAPI(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 获取当前API
	api, err := h.apiService.GetAPIByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "API不存在")
		return
	}

	// 切换状态: 1=启用, 0=禁用
	if api.State == 1 {
		api.State = 0
	} else {
		api.State = 1
	}

	// 更新API
	if err := h.apiService.UpdateAPI(api); err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, map[string]interface{}{
		"id":    api.ID,
		"state": api.State,
		"message": func() string {
			if api.State == 1 {
				return "API已启用"
			}
			return "API已禁用"
		}(),
	})
}

// TestAPI 测试API执行
func (h *APIHandler) TestAPI(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 获取API配置
	api, err := h.apiService.GetAPIByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "API不存在")
		return
	}

	// 返回API配置信息供测试
	// 注意: 实际的API测试需要根据API类型执行相应逻辑
	// 这里返回基本信息,前端可以据此构造测试请求
	utils.SuccessResponse(ctx, map[string]interface{}{
		"api": api,
		"test_info": map[string]string{
			"method": api.Method,
			"path":   api.Path,
			"state": func() string {
				if api.State == 1 {
					return "enabled"
				} else {
					return "disabled"
				}
			}(),
			"message": "请使用返回的API信息构造测试请求",
		},
	})
}
