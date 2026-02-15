package api

import (
	"context"
	"fmt"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// DataIOHandler 数据导入导出处理器
type DataIOHandler struct {
	*utils.BaseHandler
	ioService service.DataIOService
}

// NewDataIOHandler 创建数据导入导出处理器实例
func NewDataIOHandler(ioService service.DataIOService) *DataIOHandler {
	return &DataIOHandler{
		BaseHandler: utils.NewBaseHandler(),
		ioService:   ioService,
	}
}

// ExportData 导出数据
func (h *DataIOHandler) ExportData(c context.Context, ctx *app.RequestContext) {
	modelID := ctx.Param("model_id")
	format := ctx.Query("format")

	queryParams := make(map[string]any)
	ctx.BindQuery(&queryParams)

	// Hertz Response Writer adapter
	ctx.SetStatusCode(consts.StatusOK)

	var err error
	if format == "json" {
		ctx.Header("Content-Type", "application/json")
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", modelID))

		// Use stream writer
		writer := ctx.Response.BodyWriter()
		err = h.ioService.ExportToJSON(modelID, queryParams, writer)
	} else {
		// Default to Excel
		ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", modelID))

		writer := ctx.Response.BodyWriter()
		err = h.ioService.ExportToExcel(modelID, queryParams, writer)
	}

	if err != nil {
		// Since we might have already started writing headers/body, standard error response logic is tricky.
		// Detailed error handling for streaming response usually implies logging and potentially appending error to stream or closing connection abruptly.
		utils.SugarLogger.Errorf("Export failed: %v", err)
	}
}

// ImportTemplate 下载导入模板
func (h *DataIOHandler) ImportTemplate(c context.Context, ctx *app.RequestContext) {
	modelID := ctx.Param("model_id")

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s_template.xlsx", modelID))

	writer := ctx.Response.BodyWriter()
	err := h.ioService.GenerateExcelTemplate(modelID, writer)
	if err != nil {
		utils.SugarLogger.Errorf("Template generation failed: %v", err)
	}
}

// ImportData 导入数据
func (h *DataIOHandler) ImportData(c context.Context, ctx *app.RequestContext) {
	modelID := ctx.Param("model_id")

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, "File is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, "Failed to open uploaded file")
		return
	}
	defer file.Close()

	// Check extension or content type to decide format
	filename := fileHeader.Filename
	var success int
	var errors []string

	if len(filename) > 5 && filename[len(filename)-5:] == ".json" {
		success, errors, err = h.ioService.ImportFromJSON(modelID, file)
	} else {
		success, errors, err = h.ioService.ImportFromExcel(modelID, file)
	}

	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, map[string]any{
		"success_count": success,
		"errors":        errors,
	})
}
