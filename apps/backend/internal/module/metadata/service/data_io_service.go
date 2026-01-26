package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/repository"
	// "github.com/xuri/excelize/v2" // Uncomment when available
)

// DataIOService 数据导入导出服务接口
type DataIOService interface {
	ExportToExcel(modelID string, queryParams map[string]any, writer io.Writer) error
	ExportToJSON(modelID string, queryParams map[string]any, writer io.Writer) error
	GenerateExcelTemplate(modelID string, writer io.Writer) error
	ImportFromExcel(modelID string, reader io.Reader) (int, []string, error)
	ImportFromJSON(modelID string, reader io.Reader) (int, []string, error)
}

type dataIOService struct {
	crudSvc        CRUDService
	modelRepo      repository.MdModelRepository
	modelFieldRepo repository.MdModelFieldRepository
	validator      DataValidator
}

// NewDataIOService 创建数据导入导出服务实例
func NewDataIOService(
	crudSvc CRUDService,
	modelRepo repository.MdModelRepository,
	modelFieldRepo repository.MdModelFieldRepository,
	validator DataValidator,
) DataIOService {
	return &dataIOService{
		crudSvc:        crudSvc,
		modelRepo:      modelRepo,
		modelFieldRepo: modelFieldRepo,
		validator:      validator,
	}
}

// ExportToExcel 导出 Excel
func (s *dataIOService) ExportToExcel(modelID string, queryParams map[string]any, writer io.Writer) error {
	return errors.New("Excel support is currently disabled due to missing dependency (excelize)")

	/* Uncomment when excelize is installed
	modelData, err := s.crudSvc.(*crudService).builder.LoadModelData(modelID)
	if err != nil {
		return err
	}

	// Stream query results
	queryParams["page_size"] = 10000 // Large batch
	list, _, err := s.crudSvc.List(modelID, queryParams)
	if err != nil {
		return err
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			utils.SugarLogger.Error(err)
		}
	}()

	// Create header
	headers := []string{}
	for _, field := range modelData.Fields {
		headers = append(headers, field.Title)
	}
	f.SetSheetRow("Sheet1", "A1", &headers)

	// Write data
	for i, row := range list {
		values := []interface{}{}
		for _, field := range modelData.Fields {
			val := row[field.ColumnName]
			values = append(values, val)
		}
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		f.SetSheetRow("Sheet1", cell, &values)
	}

	return f.Write(writer)
	*/
}

// ExportToJSON 导出 JSON (Streaming)
func (s *dataIOService) ExportToJSON(modelID string, queryParams map[string]any, writer io.Writer) error {
	// 获取数据
	// 注意：对于大数据量，CRUDService.List 可能会撑爆内存
	// 理想情况下应该让 CRUDService 返回一个 Iterator 或者 Channel
	// 这里暂时使用分页循环获取
	
	encoder := json.NewEncoder(writer)
	
	// 写入数组开始符
	if _, err := writer.Write([]byte("[")); err != nil {
		return err
	}

	page := 1
	pageSize := 100
	first := true

	for {
		queryParams["page"] = page
		queryParams["page_size"] = pageSize
		list, _, err := s.crudSvc.List(modelID, queryParams)
		if err != nil {
			return err
		}

		if len(list) == 0 {
			break
		}

		for _, item := range list {
			if !first {
				if _, err := writer.Write([]byte(",")); err != nil {
					return err
				}
			}
			if err := encoder.Encode(item); err != nil {
				return err
			}
			first = false
		}

		if len(list) < pageSize {
			break
		}
		page++
	}

	// 写入数组结束符
	if _, err := writer.Write([]byte("]")); err != nil {
		return err
	}

	return nil
}

// GenerateExcelTemplate 生成 Excel 模板
func (s *dataIOService) GenerateExcelTemplate(modelID string, writer io.Writer) error {
	return errors.New("Excel support is currently disabled due to missing dependency (excelize)")
	/*
	md, err := s.modelRepo.GetModelByID(modelID)
	if err != nil {
		return err
	}
	
	fields, err := s.modelRepo.GetModelFields(modelID)
	if err != nil {
		return err
	}

	f := excelize.NewFile()
	headers := []string{}
	for _, field := range fields {
		if field.IsPrimaryKey { continue } // Skip ID for template?
		headers = append(headers, field.Title)
	}
	f.SetSheetRow("Sheet1", "A1", &headers)
	return f.Write(writer)
	*/
}

// ImportFromExcel 导入 Excel
func (s *dataIOService) ImportFromExcel(modelID string, reader io.Reader) (int, []string, error) {
	return 0, nil, errors.New("Excel support is currently disabled due to missing dependency (excelize)")
}

// ImportFromJSON 导入 JSON
func (s *dataIOService) ImportFromJSON(modelID string, reader io.Reader) (int, []string, error) {
	md, err := s.modelRepo.GetModelByID(modelID)
	if err != nil {
		return 0, nil, err
	}
	fields, err := s.modelFieldRepo.GetFieldsByModelID(modelID)
	if err != nil {
		return 0, nil, err
	}

	decoder := json.NewDecoder(reader)
	
	// Expect array
	if token, err := decoder.Token(); err != nil || token != json.Delim('[') {
		return 0, nil, errors.New("expected JSON array start")
	}

	successCount := 0
	errorReport := []string{}
	rowIndex := 0

	batchSize := 100
	batchData := []map[string]any{}

	processBatch := func() error {
		if len(batchData) == 0 { return nil }
		_, err := s.crudSvc.BatchCreate(modelID, batchData)
		if err != nil {
			// 如果批量失败，整个批次都标记失败，或者尝试逐条插入？
			// 简单策略：整个批次失败
			errMsg := fmt.Sprintf("Batch error (rows %d-%d): %v", rowIndex-len(batchData)+1, rowIndex, err)
			errorReport = append(errorReport, errMsg)
		} else {
			successCount += len(batchData)
		}
		batchData = []map[string]any{} // clear
		return nil
	}

	for decoder.More() {
		rowIndex++
		var data map[string]any
		if err := decoder.Decode(&data); err != nil {
			errorReport = append(errorReport, fmt.Sprintf("Row %d: JSON decode error: %v", rowIndex, err))
			continue
		}

		// Convert fields to ptr slice for validator
		fieldPtrs := make([]*model.MdModelField, len(fields))
		for i := range fields {
			fieldPtrs[i] = &fields[i]
		}
		
		// Validate
		if err := s.validator.Validate(md.ID, fieldPtrs, data); err != nil {
			errorReport = append(errorReport, fmt.Sprintf("Row %d: Validation error: %v", rowIndex, err))
			continue
		}

		batchData = append(batchData, data)

		if len(batchData) >= batchSize {
			_ = processBatch()
		}
	}

	// Process remaining
	_ = processBatch()

	if _, err := decoder.Token(); err != nil {
		// End of array
	}

	return successCount, errorReport, nil
}
