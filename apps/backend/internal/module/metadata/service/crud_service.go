package service

import (
	"context"
	"encoding/json"
	"fmt"
	"metadata-platform/internal/module/metadata/engine"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CRUDService 通用数据操作服务接口
type CRUDService interface {
	Create(ctx context.Context, modelID string, data map[string]any) (map[string]any, error)
	Get(modelID string, id string) (map[string]any, error)
	Update(ctx context.Context, modelID string, id string, data map[string]any) error
	Delete(ctx context.Context, modelID string, id string) error
	List(modelID string, queryParams map[string]any) ([]map[string]any, int64, error)
	BatchCreate(modelID string, dataList []map[string]any) ([]map[string]any, error)
	BatchDelete(modelID string, ids []string) error
	Statistics(modelID string, queryParams map[string]any) (map[string]int64, error)
	Aggregate(modelID string, queryParams map[string]any) ([]map[string]any, error)

	// 事务支持方法
	CreateWithTx(ctx context.Context, modelID string, data map[string]any, tx *gorm.DB) (map[string]any, error)
	UpdateWithTx(ctx context.Context, modelID string, id string, data map[string]any, tx *gorm.DB) error
	BatchCreateWithTx(ctx context.Context, modelID string, dataList []map[string]any, tx *gorm.DB) ([]map[string]any, error)
	BatchDeleteWithTx(ctx context.Context, modelID string, ids []string, tx *gorm.DB) error
}

type crudService struct {
	builder         *engine.SQLBuilder
	executor        *engine.SQLExecutor
	validator       DataValidator
	templateService QueryTemplateService
	auditService    AuditService
}

// NewCRUDService 创建通用数据操作服务实例
func NewCRUDService(
	builder *engine.SQLBuilder,
	executor *engine.SQLExecutor,
	validator DataValidator,
	templateService QueryTemplateService,
	auditService AuditService,
) CRUDService {
	return &crudService{
		builder:         builder,
		executor:        executor,
		validator:       validator,
		templateService: templateService,
		auditService:    auditService,
	}
}

// Create 插入数据
func (s *crudService) Create(ctx context.Context, modelID string, data map[string]any) (map[string]any, error) {
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
	tableName := s.getTableName(mainTable)

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

	// 3. 记录审计日志
	s.recordAudit(ctx, modelID, "", "CREATE", nil, data)

	return data, nil
}

// Get 查询单条数据
func (s *crudService) Get(modelID string, id string) (map[string]any, error) {
	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return nil, err
	}

	pkField := s.getPrimaryKey(modelData)
	sqlStr, args, err := s.builder.BuildSQL(modelID, nil)
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
func (s *crudService) Update(ctx context.Context, modelID string, id string, data map[string]any) error {
	// 获取变更前数据 (Audit)
	beforeData, _ := s.Get(modelID, id)

	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return err
	}

	// 0. 执行校验
	if err := s.validator.Validate(modelID, modelData.Fields, data); err != nil {
		return err
	}

	pkField := s.getPrimaryKey(modelData)
	mainTable := s.getMainTable(modelData)
	if mainTable == nil {
		return fmt.Errorf("no table defined for model %s", modelID)
	}

	tableName := s.getTableName(mainTable)

	setClauses := make([]string, 0)
	args := make([]any, 0)

	for k, v := range data {
		if k == pkField {
			continue
		} // 不更新主键
		setClauses = append(setClauses, "`"+k+"` = ?")
		args = append(args, v)
	}

	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE `%s` = ?",
		tableName,
		joinString(setClauses, ", "),
		pkField)
	args = append(args, id)

	_, err = s.executor.Execute(mainTable.ConnID, sqlStr, args...)
	if err != nil {
		return err
	}

	// 记录审计日志
	s.recordAudit(ctx, modelID, id, "UPDATE", beforeData, data)

	return nil
}

// Delete 删除数据
func (s *crudService) Delete(ctx context.Context, modelID string, id string) error {
	// 获取变更前数据 (Audit)
	beforeData, _ := s.Get(modelID, id)

	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return err
	}

	pkField := s.getPrimaryKey(modelData)
	mainTable := s.getMainTable(modelData)
	if mainTable == nil {
		return fmt.Errorf("no table defined for model %s", modelID)
	}

	tableName := s.getTableName(mainTable)

	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE `%s` = ?", tableName, pkField)
	_, err = s.executor.Execute(mainTable.ConnID, sqlStr, id)
	if err != nil {
		return err
	}

	// 记录审计日志
	s.recordAudit(ctx, modelID, id, "DELETE", beforeData, nil)

	return nil
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

	mainTable := s.getMainTable(modelData)
	if mainTable == nil {
		return nil, fmt.Errorf("no table defined for model %s", modelID)
	}
	tableName := s.getTableName(mainTable)

	allFields := make(map[string]bool)
	for _, data := range dataList {
		if err := s.validator.Validate(modelID, modelData.Fields, data); err != nil {
			return nil, err
		}
		for k := range data {
			allFields[k] = true
		}
	}

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
	mainTable := s.getMainTable(modelData)
	if mainTable == nil {
		return fmt.Errorf("no table defined for model %s", modelID)
	}
	tableName := s.getTableName(mainTable)

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

	// 应用过滤逻辑
	s.applyListFilters(modelData, queryParams)

	// 现在根据合并后的 metadata 生成 SQL
	sqlStr, args, err := s.builder.BuildFromMetadata(modelData, queryParams)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
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

// Statistics 简单统计
func (s *crudService) Statistics(modelID string, queryParams map[string]any) (map[string]int64, error) {
	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return nil, err
	}

	s.applyListFilters(modelData, queryParams)

	sqlStr, args, err := s.builder.BuildFromMetadata(modelData, queryParams)
	if err != nil {
		return nil, err
	}

	count, err := s.executor.ExecuteCount(modelData.Model.ConnID, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	return map[string]int64{"total": count}, nil
}

// Aggregate 聚合查询
func (s *crudService) Aggregate(modelID string, queryParams map[string]any) ([]map[string]any, error) {
	return nil, nil
}

// applyListFilters 提取 List 中重复的过滤逻辑
func (s *crudService) applyListFilters(modelData *engine.ModelData, queryParams map[string]any) {
	// 1. 应用查询模板
	templateID := ""
	if queryParams != nil {
		if tid, ok := queryParams["query_template_id"].(string); ok {
			templateID = tid
		}
	}

	if templateID != "" {
		s.templateService.ApplyTemplate(templateID, modelData)
	} else {
		defTemplate, _ := s.templateService.GetDefaultTemplate(modelData.Model.ID)
		if defTemplate != nil {
			s.templateService.ApplyTemplate(defTemplate.ID, modelData)
		}
	}

	// 2. 合并动态过滤条件
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
					if where.Operator1 == "" {
						where.Operator1 = "AND"
					}
					modelData.Wheres = append(modelData.Wheres, where)
				}
			}
		}
	}
}

// recordAudit 辅助记录审计日志
func (s *crudService) recordAudit(ctx context.Context, modelID, recordID, action string, before, after map[string]any) {
	if s.auditService == nil {
		return
	}

	// 从 context 获取 trace info
	var traceID string
	if v := ctx.Value("trace_id"); v != nil {
		traceID = utils.ToString(v)
	}

	// 序列化 json
	var beforeJSON, afterJSON string
	if before != nil {
		if b, err := json.Marshal(before); err == nil {
			beforeJSON = string(b)
		}
	}
	if after != nil {
		if b, err := json.Marshal(after); err == nil {
			afterJSON = string(b)
		}
	}

	s.auditService.RecordDataChange(ctx, &model.SysDataChangeLog{
		ID:         uuid.New().String(), // Ensure uuid import
		TraceID:    traceID,
		ModelID:    modelID,
		RecordID:   recordID,
		Action:     action,
		BeforeData: beforeJSON,
		AfterData:  afterJSON,
		CreateBy:   "system", // Improve: get user from context
		CreateAt:   time.Now(),
	})
}

// ------------------------------------------------------------------------------------------------
// Transactional Methods
// ------------------------------------------------------------------------------------------------

func (s *crudService) CreateWithTx(ctx context.Context, modelID string, data map[string]any, tx *gorm.DB) (map[string]any, error) {
	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return nil, err
	}

	if err := s.validator.Validate(modelID, modelData.Fields, data); err != nil {
		return nil, err
	}

	mainTable := s.getMainTable(modelData)
	if mainTable == nil {
		return nil, fmt.Errorf("no table defined for model %s", modelID)
	}

	tableName := s.getTableName(mainTable)

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

	_, err = s.executor.ExecuteWithTx(tx, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	// Record audit
	s.recordAudit(ctx, modelID, "", "CREATE_TX", nil, data)

	return data, nil
}

func (s *crudService) UpdateWithTx(ctx context.Context, modelID string, id string, data map[string]any, tx *gorm.DB) error {
	// 获取变更前数据 (Audit) - 注意：事务内查询应该使用 txExecutor，但 Get 没 exposing withTx?
	// 简单实现：使用 id 查，如果事务隔离级别允许读取已提交，或者就是查旧值。
	// 但这里我们没有 GetWithTx.
	// 优化：暂时略过 BeforeData 获取，或者假设事务内不影响旧数据读取（只要没改id）
	// 或者 we assume we update an existing record visible to outside or same txn?
	// 事务内 visible? Get uses `s.executor.Execute` which uses a NEW connection/session unless we propagate tx context.
	// Since Get creates new connection from pool (via ConnID), it won't see changes in current tx unless committed.
	// But `UpdateWithTx` is updating it. Before updating it, the old data is in DB. Get() reads from DB.
	// So calling Get() here is safe to get "Before" state, provided we haven't updated it yet in this txn.
	// Wait, we are *inside* UpdateWithTx, before executing update SQL. So Get() returns "Before" state.
	// But `s.Get` uses `s.executor.Execute`. `s.executor` gets connection by ConnID.
	// `CreateMasterDetail` starts a tx on a connection.
	// If `Get` uses a DIFFERENT connection, it sees committed data.
	// If `CreateMasterDetail` locks the rows (e.g. SelectForUpdate), `Get` might block?
	// Usually Update doesn't block Reader (MVCC).

	beforeData, _ := s.Get(modelID, id)
	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return err
	}

	if err := s.validator.Validate(modelID, modelData.Fields, data); err != nil {
		return err
	}

	pkField := s.getPrimaryKey(modelData)
	mainTable := s.getMainTable(modelData)
	if mainTable == nil {
		return fmt.Errorf("no table defined for model %s", modelID)
	}
	tableName := s.getTableName(mainTable)

	setClauses := make([]string, 0)
	args := make([]any, 0)

	for k, v := range data {
		if k == pkField {
			continue
		}
		setClauses = append(setClauses, "`"+k+"` = ?")
		args = append(args, v)
	}

	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE `%s` = ?",
		tableName,
		joinString(setClauses, ", "),
		pkField)
	args = append(args, id)

	_, err = s.executor.ExecuteWithTx(tx, sqlStr, args...)
	if err == nil {
		s.recordAudit(ctx, modelID, id, "UPDATE_TX", beforeData, data)
	}
	return err
}

func (s *crudService) BatchCreateWithTx(ctx context.Context, modelID string, dataList []map[string]any, tx *gorm.DB) ([]map[string]any, error) {
	if len(dataList) == 0 {
		return nil, nil
	}

	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return nil, err
	}

	mainTable := s.getMainTable(modelData)
	if mainTable == nil {
		return nil, fmt.Errorf("no table defined for model %s", modelID)
	}
	tableName := s.getTableName(mainTable)

	allFields := make(map[string]bool)
	for _, data := range dataList {
		if err := s.validator.Validate(modelID, modelData.Fields, data); err != nil {
			return nil, err
		}
		for k := range data {
			allFields[k] = true
		}
	}

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
			cleanField := f[1 : len(f)-1]
			args = append(args, data[cleanField])
		}
	}

	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		tableName,
		strings.Join(fieldList, ", "),
		strings.Join(placeholders, ", "))

	_, err = s.executor.ExecuteWithTx(tx, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	// 批量审计
	// 记录一次操作，包含 list count? 或者 逐条？
	// 逐条太慢。
	// AuditService.RecordDataChange 只支持单条。
	// 这里这做简化：只记录 Action=BATCH_CREATE, Before=nil, After={"count": len(dataList)}
	// 或者循环调用 recordAudit

	// 暂时只记录最后一条或者 summary
	s.recordAudit(ctx, modelID, "BATCH", "BATCH_CREATE_TX", nil, map[string]any{"count": len(dataList)})

	return dataList, nil
}

func (s *crudService) BatchDeleteWithTx(ctx context.Context, modelID string, ids []string, tx *gorm.DB) error {
	if len(ids) == 0 {
		return nil
	}

	modelData, err := s.builder.LoadModelData(modelID)
	if err != nil {
		return err
	}

	pkField := s.getPrimaryKey(modelData)
	mainTable := s.getMainTable(modelData)
	if mainTable == nil {
		return fmt.Errorf("no table defined for model %s", modelID)
	}
	tableName := s.getTableName(mainTable)

	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE `%s` IN (%s)", tableName, pkField, placeholders)

	anyIds := make([]any, len(ids))
	for i, id := range ids {
		anyIds[i] = id
	}

	_, err = s.executor.ExecuteWithTx(tx, sqlStr, anyIds...)
	if err == nil {
		s.recordAudit(ctx, modelID, "BATCH", "BATCH_DELETE_TX", nil, map[string]any{"ids": ids})
	}
	return err
}

func (s *crudService) getMainTable(modelData *engine.ModelData) *model.MdModelTable {
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
	return mainTable
}

func (s *crudService) getTableName(mainTable *model.MdModelTable) string {
	tableName := "`" + mainTable.TableNameStr + "`"
	if mainTable.TableSchema != "" {
		tableName = "`" + mainTable.TableSchema + "`.`" + mainTable.TableNameStr + "`"
	}
	return tableName
}
