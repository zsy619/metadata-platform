package api

import (
	"context"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// DataQueryHandler 通用数据查询处理器
type DataQueryHandler struct {
	crudService  service.CRUDService
	modelService service.MdModelService
}

// NewDataQueryHandler 创建处理器实例
func NewDataQueryHandler(crudService service.CRUDService, modelService service.MdModelService) *DataQueryHandler {
	return &DataQueryHandler{
		crudService:  crudService,
		modelService: modelService,
	}
}

// HandleUnifiedQuery 处理统一查询
func (h *DataQueryHandler) HandleUnifiedQuery(c context.Context, ctx *app.RequestContext) {
	// 获取模型ID (通常由路由参数或上下文传递，这里假设调用者已解析并传递 model_id)
	// 但 DynamicRouter 的模式是每次重新解析 code。
	// 为了复用，这里只关注核心逻辑，但也需要知道 modelID。
	// 我们可以让 DynamicRouter 解析好 modelID 后传进来，或者这里再次解析。
	// 考虑到 Hertz Handler 签名固定，我们假设 DynamicRouter 会通过某种方式传递 ModelID，
	// 或者我们这里再次解析 context 中的 param/code。
	//
	// 简单起见，我们遵循 DynamicRouter 的现有模式：在闭包中解析，这里只提供核心逻辑函数，
	// 或者这个 Handler 就像 APIHandler 一样，自己从 Path Param 获取 ID。
	//
	// 既然是 DynamicRouter 调用的，路由参数可能不固定（取决于 Path）。
	// 统一查询路径: /api/data/:model_code/query -> model_code 在 param 中?
	// DynamicRouter 注册路径时是写死的: /api/data/user/query
	// 所以 path param 不能直接用 :model_code 除非注册时用了通配符。
	// APIGenerator 生成的是具体路径 /api/data/user/query。
	// 所以这里无法通过 ctx.Param("model_code") 获取，只能通过上下文或闭包。

	// 因此，这个 Handler 最好不做 Hertz 路由绑定，而是提供方法供 DynamicRouter 调用。
	// 或者 DynamicRouter 注册时，将 modelID 注入到闭包中。
}

// HandleUnifiedQueryWithModelID 供 DynamicRouter 调用的带 ModelID 的处理函数
func (h *DataQueryHandler) HandleUnifiedQueryWithModelID(ctx *app.RequestContext, modelID string) {
	// Bind JSON body directly to map 更好，因为结构不固定
	var body map[string]any
	if err := ctx.BindJSON(&body); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, "Invalid JSON payload")
		return
	}

	results, count, err := h.crudService.List(modelID, body)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, map[string]any{
		"list":  results,
		"total": count,
	})
}

// HandleBatchCreateWithModelID 批量创建
func (h *DataQueryHandler) HandleBatchCreateWithModelID(ctx *app.RequestContext, modelID string) {
	var dataList []map[string]any
	if err := ctx.BindJSON(&dataList); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, "Expected JSON array of objects")
		return
	}

	results, err := h.crudService.BatchCreate(modelID, dataList)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, results)
}

// HandleBatchDeleteWithModelID 批量删除
func (h *DataQueryHandler) HandleBatchDeleteWithModelID(ctx *app.RequestContext, modelID string) {
	var ids []string
	if err := ctx.BindJSON(&ids); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, "Expected JSON array of IDs")
		return
	}

	err := h.crudService.BatchDelete(modelID, ids)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "Deleted successfully")
}

// HandleStatisticsWithModelID 统计查询
func (h *DataQueryHandler) HandleStatisticsWithModelID(ctx *app.RequestContext, modelID string) {
	var body map[string]any
	// 允许空 body
	if string(ctx.Request.Body()) != "" {
		if err := ctx.BindJSON(&body); err != nil {
			utils.ErrorResponse(ctx, consts.StatusBadRequest, "Invalid JSON payload")
			return
		}
	}

	result, err := h.crudService.Statistics(modelID, body)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, result)
}

// HandleAggregateWithModelID 聚合查询
func (h *DataQueryHandler) HandleAggregateWithModelID(ctx *app.RequestContext, modelID string) {
	var body map[string]any
	// 允许空 body
	if string(ctx.Request.Body()) != "" {
		if err := ctx.BindJSON(&body); err != nil {
			utils.ErrorResponse(ctx, consts.StatusBadRequest, "Invalid JSON payload")
			return
		}
	}

	results, err := h.crudService.Aggregate(modelID, body)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, results)
}
