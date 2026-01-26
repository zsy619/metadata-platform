package api

import (
	"context"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// MasterDetailHandler 主子表处理器
type MasterDetailHandler struct {
	mdService service.MasterDetailService
}

// NewMasterDetailHandler 创建主子表处理器实例
func NewMasterDetailHandler(mdService service.MasterDetailService) *MasterDetailHandler {
	return &MasterDetailHandler{mdService: mdService}
}

// CreateMasterDetail 创建主子表数据
func (h *MasterDetailHandler) CreateMasterDetail(c context.Context, ctx *app.RequestContext) {
	masterModelID := ctx.Param("master")
	detailModelID := ctx.Param("detail")

	var payload map[string]any
	if err := ctx.BindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := h.mdService.CreateMasterDetail(c, masterModelID, detailModelID, payload); err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, "Created successfully")
}
