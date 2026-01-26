package service

import (
	"fmt"
	"metadata-platform/internal/module/metadata/engine"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/utils"
	"strings"
)

// CRUDService 通用数据操作服务接口
type CRUDService interface {
	Create(modelID string, data map[string]any) (map[string]any, error)
	Get(modelID string, id string) (map[string]any, error)
	Update(modelID string, id string, data map[string]any) error
	Delete(modelID string, id string) error
	List(modelID string, queryParams map[string]any) ([]map[string]any, int64, error)
	BatchCreate(modelID string, dataList []map[string]any) ([]map[string]any, error)
	BatchDelete(modelID string, ids []string) error
}

type crudService struct {
	builder         *engine.SQLBuilder
	executor        *engine.SQLExecutor
	validator       DataValidator
	templateService QueryTemplateService
}

// NewCRUDService 创建通用数据操作服务实例
func NewCRUDService(
	builder *engine.SQLBuilder,
	executor *engine.SQLExecutor,
	validator DataValidator,
	templateService QueryTemplateService,
) CRUDService {
	return &crudService{
		builder:         builder,
		executor:        executor,
		validator:       validator,
		templateService: templateService,
	}
}

// Create 插入数据
func (s *crudService) Create(modelID string, data map[string]any) (map[string]any, error) {
	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return nil, err
	}

	// 0. 执行校验
	if err := s.validator.Validate(modelID, modelData.Fields, data); err != nil {
		return nil, err
	}

	// 1. 找到主表
	var mainTable *model.MdModelTable
	for _, t := range modelData.Tables {
		if t.IsMain {
			mainTable = t
			break
		}
	}
	if mainTable == nil && len(modelData.Tables) > 0 {
		mainTable = modelData.Tables[0]
	}
	if mainTable == nil {
		return nil, fmt.Errorf("no table defined for model %s", modelID)
	}

	// 2. 构造 INSERT 语句
	// 注意：这里需要一个专门的子句构建器用于 CUD 操作，或者直接在 service 中实现
	// 为简化，目前先在这里实现基础逻辑
	tableName := "`" + mainTable.TableNameStr + "`"
	if mainTable.TableSchema != "" {
		tableName = "`" + mainTable.TableSchema + "`.`" + mainTable.TableNameStr + "`"
	}

	fields := make([]string, 0)
	placeholders := make([]string, 0)
	args := make([]any, 0)

	for k, v := range data {
		fields = append(fields, "`"+k+"`")
		placeholders = append(placeholders, "?")
		args = append(args, v)
	}

	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", 
		tableName, 
		joinString(fields, ", "), 
		joinString(placeholders, ", "))

	_, err = s.executor.Execute(mainTable.ConnID, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Get 查询单条数据
func (s *crudService) Get(modelID string, id string) (map[string]any, error) {
	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return nil, err
	}

	pkField := s.getPrimaryKey(modelData)
	sqlStr, args, err := s.builder.BuildSQL(modelID)
	if err != nil {
		return nil, err
	}

	// 添加针对 ID 的过滤
	if strings.Contains(strings.ToUpper(sqlStr), " WHERE ") {
		sqlStr += " AND `" + pkField + "` = ?"
	} else {
		sqlStr += " WHERE `" + pkField + "` = ?"
	}
	args = append(args, id)

	results, err := s.executor.Execute(modelData.Model.ConnID, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results[0], nil
}

// Update 更新数据
func (s *crudService) Update(modelID string, id string, data map[string]any) error {
	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return err
	}

	// 0. 执行校验
	if err := s.validator.Validate(modelID, modelData.Fields, data); err != nil {
		return err
	}

	pkField := s.getPrimaryKey(modelData)
	var mainTable *model.MdModelTable
	for _, t := range modelData.Tables {
		if t.IsMain {
			mainTable = t
			break
		}
	}
	if mainTable == nil && len(modelData.Tables) > 0 {
		mainTable = modelData.Tables[0]
	}

	tableName := "`" + mainTable.TableNameStr + "`"
	if mainTable.TableSchema != "" {
		tableName = "`" + mainTable.TableSchema + "`.`" + mainTable.TableNameStr + "`"
	}

	setClauses := make([]string, 0)
	args := make([]any, 0)

	for k, v := range data {
		if k == pkField { continue } // 不更新主键
		setClauses = append(setClauses, "`"+k+"` = ?")
		args = append(args, v)
	}

	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE `%s` = ?", 
		tableName, 
		joinString(setClauses, ", "),
		pkField)
	args = append(args, id)

	_, err = s.executor.Execute(mainTable.ConnID, sqlStr, args...)
	return err
}

// Delete 删除数据
func (s *crudService) Delete(modelID string, id string) error {
	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return err
	}

	pkField := s.getPrimaryKey(modelData)
	var mainTable *model.MdModelTable
	for _, t := range modelData.Tables {
		if t.IsMain {
			mainTable = t
			break
		}
	}
	if mainTable == nil && len(modelData.Tables) > 0 {
		mainTable = modelData.Tables[0]
	}

	tableName := "`" + mainTable.TableNameStr + "`"
	if mainTable.TableSchema != "" {
		tableName = "`" + mainTable.TableSchema + "`.`" + mainTable.TableNameStr + "`"
	}

	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE `%s` = ?", tableName, pkField)
	_, err = s.executor.Execute(mainTable.ConnID, sqlStr, id)
	return err
}

// BatchCreate 批量插入数据
func (s *crudService) BatchCreate(modelID string, dataList []map[string]any) ([]map[string]any, error) {
	if len(dataList) == 0 {
		return nil, nil
	}

	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return nil, err
	}

	// 1. 找到主表
	var mainTable *model.MdModelTable
	for _, t := range modelData.Tables {
		if t.IsMain {
			mainTable = t
			break
		}
	}
	if mainTable == nil && len(modelData.Tables) > 0 {
		mainTable = modelData.Tables[0]
	}
	if mainTable == nil {
		return nil, fmt.Errorf("no table defined for model %s", modelID)
	}

	tableName := "`" + mainTable.TableNameStr + "`"
	if mainTable.TableSchema != "" {
		tableName = "`" + mainTable.TableSchema + "`.`" + mainTable.TableNameStr + "`"
	}

	// 2. 逐条校验并收集字段
	allFields := make(map[string]bool)
	for _, data := range dataList {
		if err := s.validator.Validate(modelID, modelData.Fields, data); err != nil {
			return nil, err
		}
		for k := range data {
			allFields[k] = true
		}
	}

	// 3. 构造批量 INSERT
	fieldList := make([]string, 0, len(allFields))
	for f := range allFields {
		fieldList = append(fieldList, "`"+f+"`")
	}

	placeholders := make([]string, 0, len(dataList))
	args := make([]any, 0, len(dataList)*len(fieldList))

	rowPlaceholder := "(" + strings.Repeat("?,", len(fieldList))
	rowPlaceholder = rowPlaceholder[:len(rowPlaceholder)-1] + ")"

	for _, data := range dataList {
		placeholders = append(placeholders, rowPlaceholder)
		for _, f := range fieldList {
			// 将带引号的字段还原
			cleanField := f[1 : len(f)-1]
			args = append(args, data[cleanField])
		}
	}

	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		tableName,
		strings.Join(fieldList, ", "),
		strings.Join(placeholders, ", "))

	_, err = s.executor.Execute(mainTable.ConnID, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	return dataList, nil
}

// BatchDelete 批量删除数据
func (s *crudService) BatchDelete(modelID string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return err
	}

	pkField := s.getPrimaryKey(modelData)
	var mainTable *model.MdModelTable
	for _, t := range modelData.Tables {
		if t.IsMain {
			mainTable = t
			break
		}
	}
	if mainTable == nil && len(modelData.Tables) > 0 {
		mainTable = modelData.Tables[0]
	}

	tableName := "`" + mainTable.TableNameStr + "`"
	if mainTable.TableSchema != "" {
		tableName = "`" + mainTable.TableSchema + "`.`" + mainTable.TableNameStr + "`"
	}

	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE `%s` IN (%s)", tableName, pkField, placeholders)
	
	anyIds := make([]any, len(ids))
	for i, id := range ids {
		anyIds[i] = id
	}

	_, err = s.executor.Execute(mainTable.ConnID, sqlStr, anyIds...)
	return err
}

func (s *crudService) getPrimaryKey(data *engine.ModelData) string {
	for _, f := range data.Fields {
		if f.IsPrimaryKey {
			return f.ColumnName
		}
	}
	return "id" // 默认退回到 id
}

func joinString(slice []string, sep string) string {
	res := ""
	for i, s := range slice {
		if i > 0 {
			res += sep
		}
		res += s
	}
	return res
}

// List 分页查询
func (s *crudService) List(modelID string, queryParams map[string]any) ([]map[string]any, int64, error) {
	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return nil, 0, err
	}

	// 1. 应用查询模板 (如果指定了 template_id 或者有默认模板)
	templateID := ""
	if queryParams != nil {
		if tid, ok := queryParams["query_template_id"].(string); ok {
			templateID = tid
		}
	}

	if templateID != "" {
		if err := s.templateService.ApplyTemplate(templateID, modelData); err != nil {
			return nil, 0, fmt.Errorf("failed to apply template: %v", err)
		}
	} else {
		// 检查默认模板
		defTemplate, _ := s.templateService.GetDefaultTemplate(modelID)
		if defTemplate != nil {
			if err := s.templateService.ApplyTemplate(defTemplate.ID, modelData); err != nil {
				// 默认模板加载失败可作为 warning，或根据业务需求报错
			}
		}
	}

	// 2. 合并动态过滤条件 (如果是统一查询入口)
	if queryParams != nil {
		if filters, ok := queryParams["filters"].([]any); ok {
			for _, item := range filters {
				if f, ok := item.(map[string]any); ok {
					where := &model.MdModelWhere{
						TableNameStr: utils.ToString(f["table_name"]),
						ColumnName:   utils.ToString(f["column_name"]),
						Operator2:    utils.ToString(f["operator"]),
						Value1:       utils.ToString(f["value"]),
					}
					// 默认 AND
					if where.Operator1 == "" {
						where.Operator1 = "AND"
					}
					modelData.Wheres = append(modelData.Wheres, where)
				}
			}
		}
	}

	// 现在根据合并后的 metadata 生成 SQL
	sqlStr, args, err := s.builder.BuildFromMetadata(modelData)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数 (注意：ExecuteCount 内部会处理 LIMIT 的剥离)
	count, err := s.executor.ExecuteCount(modelData.Model.ConnID, sqlStr, args...)
	if err != nil {
		return nil, 0, err
	}

	// 执行查询
	results, err := s.executor.Execute(modelData.Model.ConnID, sqlStr, args...)
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}
