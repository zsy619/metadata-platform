package api

import (
	"context"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// QueryTemplateHandler 查询模板处理器
type QueryTemplateHandler struct {
	*utils.BaseHandler
	templateService service.QueryTemplateService
}

// NewQueryTemplateHandler 创建查询模板处理器实例
func NewQueryTemplateHandler(templateService service.QueryTemplateService) *QueryTemplateHandler {
	return &QueryTemplateHandler{
		BaseHandler:     utils.NewBaseHandler(),
		templateService: templateService,
	}
}

// CreateTemplate 创建查询模板
func (h *QueryTemplateHandler) CreateTemplate(c context.Context, ctx *app.RequestContext) {
	var template model.MdQueryTemplate
	if err := ctx.BindJSON(&template); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	// 租户/用户上下文提取
	tenantID, _ := ctx.Get("tenant_id")
	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")

	template.TenantID = utils.ToString(tenantID)
	template.CreateID = utils.ToString(userID)
	template.CreateBy = utils.ToString(username)

	if err := h.templateService.CreateTemplate(&template); err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, template)
}

// GetTemplateByID 获取详情
func (h *QueryTemplateHandler) GetTemplateByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	template, err := h.templateService.GetTemplateByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "模板不存在")
		return
	}
	utils.SuccessResponse(ctx, template)
}

// GetTemplatesByModelID 获取模型关联的模板列表
func (h *QueryTemplateHandler) GetTemplatesByModelID(c context.Context, ctx *app.RequestContext) {
	modelID := ctx.Param("id")
	templates, err := h.templateService.GetTemplatesByModelID(modelID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, templates)
}

// UpdateTemplate 更新模板
func (h *QueryTemplateHandler) UpdateTemplate(c context.Context, ctx *app.RequestContext) {
	var template model.MdQueryTemplate
	if err := ctx.BindJSON(&template); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	if err := h.templateService.UpdateTemplate(&template); err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, template)
}

// SetDefault 设置默认模板
func (h *QueryTemplateHandler) SetDefault(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	modelID := ctx.Query("model_id")
	if modelID == "" {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, "missing model_id")
		return
	}

	if err := h.templateService.SetDefault(modelID, id); err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, "set successfully")
}

// DeleteTemplate 删除模板
func (h *QueryTemplateHandler) DeleteTemplate(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	if err := h.templateService.DeleteTemplate(id); err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, "deleted")
}

// DuplicateTemplate 复制模板
func (h *QueryTemplateHandler) DuplicateTemplate(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("templateId")

	newTemplate, err := h.templateService.DuplicateTemplate(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, newTemplate)
}

// PreviewTemplate 预览模板SQL和结果
func (h *QueryTemplateHandler) PreviewTemplate(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("templateId")

	// 获取模板
	template, err := h.templateService.GetTemplateByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	if template == nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "模板不存在")
		return
	}

	// 返回模板配置和条件信息
	// 注意: 实际的SQL预览需要结合SQLBuilder和具体的Model配置
	// 这里返回模板的条件配置,前端可以据此展示预览信息
	utils.SuccessResponse(ctx, map[string]interface{}{
		"template": template,
		"preview_info": map[string]interface{}{
			"model_id":         template.ModelID,
			"conditions_count": len(template.Conditions),
			"conditions":       template.Conditions,
			"message":          "模板条件预览,实际SQL生成需要结合模型配置",
		},
	})
}
