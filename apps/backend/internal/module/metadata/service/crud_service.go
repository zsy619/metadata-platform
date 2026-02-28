package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"metadata-platform/internal/module/audit/service"
	"metadata-platform/internal/module/metadata/engine"
)

// CRUDService CRUD服务接口
type CRUDService interface {
	Create(ctx context.Context, modelID string, data map[string]any) (map[string]any, error)
	CreateWithTx(ctx context.Context, modelID string, data map[string]any, tx *gorm.DB) (map[string]any, error)
	Get(modelID, id string) (map[string]any, error)
	Update(ctx context.Context, modelID, id string, data map[string]any) error
	Delete(ctx context.Context, modelID, id string) error
	List(modelID string, params map[string]any) ([]map[string]any, int64, error)
	BatchCreate(modelID string, dataList []map[string]any) ([]map[string]any, error)
	BatchCreateWithTx(ctx context.Context, modelID string, dataList []map[string]any, tx *gorm.DB) ([]map[string]any, error)
	BatchDelete(modelID string, ids []string) error
	Statistics(modelID string, queryParams map[string]any) (map[string]int64, error)
	Aggregate(modelID string, queryParams map[string]any) ([]map[string]any, error)
	BuildSQLFromData(data *engine.ModelData, params map[string]any) (string, []any, error)
	ExecuteModelData(data *engine.ModelData, params map[string]any) ([]map[string]any, int64, error)
}

// crudService CRUD服务实现
type crudService struct {
	sqlBuilder       *engine.SQLBuilder
	sqlExecutor      *engine.SQLExecutor
	validator        DataValidator
	queryTemplateSvc QueryTemplateService
	auditSvc         service.AuditService
}

// NewCRUDService 创建CRUD服务实例
func NewCRUDService(
	sqlBuilder *engine.SQLBuilder,
	sqlExecutor *engine.SQLExecutor,
	validator DataValidator,
	queryTemplateSvc QueryTemplateService,
	auditSvc service.AuditService,
) CRUDService {
	return &crudService{
		sqlBuilder:       sqlBuilder,
		sqlExecutor:      sqlExecutor,
		validator:        validator,
		queryTemplateSvc: queryTemplateSvc,
		auditSvc:         auditSvc,
	}
}

// Create 创建数据
func (s *crudService) Create(ctx context.Context, modelID string, data map[string]any) (map[string]any, error) {
	// 1. 加载模型
	md, err := s.sqlBuilder.LoadModelData(modelID)
	if err != nil {
		return nil, fmt.Errorf("加载模型失败: %w", err)
	}

	// 2. 验证数据
	if err := s.validator.Validate(modelID, md.Fields, data); err != nil {
		return nil, fmt.Errorf("数据验证失败: %w", err)
	}

	// 3. 构建插入SQL
	sql, args, err := s.buildInsertSQL(md, data)
	if err != nil {
		return nil, fmt.Errorf("构建插入SQL失败: %w", err)
	}

	// 4. 执行SQL
	connID := s.getConnID(md)
	result, err := s.sqlExecutor.Execute(connID, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("执行插入失败: %w", err)
	}

	// 5. 获取插入后的ID
	id := ""
	if len(result) > 0 {
		if idValue, ok := result[0]["id"]; ok {
			id = fmt.Sprintf("%v", idValue)
		}
	}
	if id == "" {
		// 尝试从data中获取ID
		if idValue, ok := data[s.getPrimaryKey(md)]; ok {
			id = fmt.Sprintf("%v", idValue)
		} else {
			return nil, errors.New("无法获取插入后的ID")
		}
	}

	// 6. 查询插入后的数据
	return s.Get(modelID, id)
}

// Get 获取数据
func (s *crudService) Get(modelID, id string) (map[string]any, error) {
	// 1. 加载模型
	md, err := s.sqlBuilder.LoadModelData(modelID)
	if err != nil {
		return nil, fmt.Errorf("加载模型失败: %w", err)
	}

	// 2. 构建查询SQL
	sql, args, err := s.buildGetSQL(md, id)
	if err != nil {
		return nil, fmt.Errorf("构建查询SQL失败: %w", err)
	}

	// 3. 执行SQL
	connID := s.getConnID(md)
	result, err := s.sqlExecutor.Execute(connID, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("执行查询失败: %w", err)
	}

	// 4. 处理结果
	if len(result) == 0 {
		return nil, nil
	}

	return result[0], nil
}

// Update 更新数据
func (s *crudService) Update(ctx context.Context, modelID, id string, data map[string]any) error {
	// 1. 加载模型
	md, err := s.sqlBuilder.LoadModelData(modelID)
	if err != nil {
		return fmt.Errorf("加载模型失败: %w", err)
	}

	// 2. 验证数据
	if err := s.validator.Validate(modelID, md.Fields, data); err != nil {
		return fmt.Errorf("数据验证失败: %w", err)
	}

	// 3. 构建更新SQL
	sql, args, err := s.buildUpdateSQL(md, id, data)
	if err != nil {
		return fmt.Errorf("构建更新SQL失败: %w", err)
	}

	// 4. 执行SQL
	connID := s.getConnID(md)
	_, err = s.sqlExecutor.Execute(connID, sql, args...)
	if err != nil {
		return fmt.Errorf("执行更新失败: %w", err)
	}

	return nil
}

// Delete 删除数据
func (s *crudService) Delete(ctx context.Context, modelID, id string) error {
	// 1. 加载模型
	md, err := s.sqlBuilder.LoadModelData(modelID)
	if err != nil {
		return fmt.Errorf("加载模型失败: %w", err)
	}

	// 2. 构建删除SQL
	sql, args, err := s.buildDeleteSQL(md, id)
	if err != nil {
		return fmt.Errorf("构建删除SQL失败: %w", err)
	}

	// 3. 执行SQL
	connID := s.getConnID(md)
	_, err = s.sqlExecutor.Execute(connID, sql, args...)
	if err != nil {
		return fmt.Errorf("执行删除失败: %w", err)
	}

	return nil
}

// List 查询列表
func (s *crudService) List(modelID string, params map[string]any) ([]map[string]any, int64, error) {
	// 1. 加载模型
	md, err := s.sqlBuilder.LoadModelData(modelID)
	if err != nil {
		return nil, 0, fmt.Errorf("加载模型失败: %w", err)
	}

	// 2. 构建查询SQL
	sql, countSql, args, err := s.buildListSQL(md, params)
	if err != nil {
		return nil, 0, fmt.Errorf("构建查询SQL失败: %w", err)
	}

	connID := s.getConnID(md)

	// 3. 执行计数查询
	var total int64
	if countSql != "" {
		countResult, err := s.sqlExecutor.Execute(connID, countSql, args...)
		if err != nil {
			return nil, 0, fmt.Errorf("执行计数查询失败: %w", err)
		}
		if len(countResult) > 0 {
			total = s.toInt64(countResult[0]["count"])
		}
	}

	// 4. 执行列表查询
	result, err := s.sqlExecutor.Execute(connID, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("执行列表查询失败: %w", err)
	}

	return result, total, nil
}

// CreateWithTx 在事务中创建数据
func (s *crudService) CreateWithTx(ctx context.Context, modelID string, data map[string]any, tx *gorm.DB) (map[string]any, error) {
	// 1. 加载模型
	md, err := s.sqlBuilder.LoadModelData(modelID)
	if err != nil {
		return nil, fmt.Errorf("加载模型失败: %w", err)
	}

	// 2. 验证数据
	if err := s.validator.Validate(modelID, md.Fields, data); err != nil {
		return nil, fmt.Errorf("数据验证失败: %w", err)
	}

	// 3. 构建插入SQL
	sql, args, err := s.buildInsertSQL(md, data)
	if err != nil {
		return nil, fmt.Errorf("构建插入SQL失败: %w", err)
	}

	// 4. 在事务中执行SQL
	result, err := s.sqlExecutor.ExecuteWithTx(tx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("执行插入失败: %w", err)
	}

	// 5. 获取插入后的ID
	id := ""
	if len(result) > 0 {
		if idValue, ok := result[0]["id"]; ok {
			id = fmt.Sprintf("%v", idValue)
		}
	}
	if id == "" {
		// 尝试从data中获取ID
		if idValue, ok := data[s.getPrimaryKey(md)]; ok {
			id = fmt.Sprintf("%v", idValue)
		} else {
			return nil, errors.New("无法获取插入后的ID")
		}
	}

	// 6. 查询插入后的数据
	return s.Get(modelID, id)
}

// BatchCreate 批量创建
func (s *crudService) BatchCreate(modelID string, dataList []map[string]any) ([]map[string]any, error) {
	results := make([]map[string]any, 0, len(dataList))
	for _, data := range dataList {
		result, err := s.Create(context.Background(), modelID, data)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// BatchCreateWithTx 在事务中批量创建
func (s *crudService) BatchCreateWithTx(ctx context.Context, modelID string, dataList []map[string]any, tx *gorm.DB) ([]map[string]any, error) {
	results := make([]map[string]any, 0, len(dataList))
	for _, data := range dataList {
		result, err := s.CreateWithTx(ctx, modelID, data, tx)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// BatchDelete 批量删除
func (s *crudService) BatchDelete(modelID string, ids []string) error {
	for _, id := range ids {
		if err := s.Delete(context.Background(), modelID, id); err != nil {
			return err
		}
	}
	return nil
}

// Statistics 统计查询
func (s *crudService) Statistics(modelID string, queryParams map[string]any) (map[string]int64, error) {
	// 1. 加载模型
	md, err := s.sqlBuilder.LoadModelData(modelID)
	if err != nil {
		return nil, fmt.Errorf("加载模型失败: %w", err)
	}

	connID := s.getConnID(md)

	// 2. 构建基础SQL
	sql, _, err := s.sqlBuilder.BuildSQL(modelID, queryParams)
	if err != nil {
		return nil, fmt.Errorf("构建SQL失败: %w", err)
	}

	// 3. 执行统计查询
	stats := make(map[string]int64)

	// 总记录数
	countSQL := fmt.Sprintf("SELECT COUNT(*) as count FROM (%s) AS t", sql)
	countResult, err := s.sqlExecutor.Execute(connID, countSQL)
	if err != nil {
		return nil, err
	}
	if len(countResult) > 0 {
		stats["total"] = s.toInt64(countResult[0]["count"])
	}

	return stats, nil
}

// Aggregate 聚合查询
func (s *crudService) Aggregate(modelID string, queryParams map[string]any) ([]map[string]any, error) {
	// 1. 加载模型
	md, err := s.sqlBuilder.LoadModelData(modelID)
	if err != nil {
		return nil, fmt.Errorf("加载模型失败: %w", err)
	}

	connID := s.getConnID(md)

	// 2. 构建基础SQL
	sql, args, err := s.sqlBuilder.BuildSQL(modelID, queryParams)
	if err != nil {
		return nil, fmt.Errorf("构建SQL失败: %w", err)
	}

	// 3. 执行查询
	return s.sqlExecutor.Execute(connID, sql, args...)
}

// BuildSQLFromData 从ModelData构建SQL
func (s *crudService) BuildSQLFromData(data *engine.ModelData, params map[string]any) (string, []any, error) {
	return s.sqlBuilder.BuildSQL(data.Model.ID, params)
}

// ExecuteModelData 执行ModelData查询
func (s *crudService) ExecuteModelData(data *engine.ModelData, params map[string]any) ([]map[string]any, int64, error) {
	connID := s.getConnID(data)

	// 1. 构建SQL
	sql, args, err := s.sqlBuilder.BuildSQL(data.Model.ID, params)
	if err != nil {
		return nil, 0, err
	}

	// 2. 执行计数查询
	var total int64
	countSQL := fmt.Sprintf("SELECT COUNT(*) as count FROM (%s) AS t", sql)
	countResult, err := s.sqlExecutor.Execute(connID, countSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	if len(countResult) > 0 {
		total = s.toInt64(countResult[0]["count"])
	}

	// 3. 执行列表查询
	results, err := s.sqlExecutor.Execute(connID, sql, args...)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// 辅助方法

func (s *crudService) getConnID(md *engine.ModelData) string {
	if md.Model != nil && md.Model.ConnID != "" {
		return md.Model.ConnID
	}
	// 从主表获取连接ID
	for _, t := range md.Tables {
		if t.IsMain && t.ConnID != "" {
			return t.ConnID
		}
	}
	return ""
}

func (s *crudService) getPrimaryKey(md *engine.ModelData) string {
	for _, f := range md.Fields {
		if f.IsPrimaryKey {
			return f.ColumnName
		}
	}
	return "id"
}

func (s *crudService) toInt64(v any) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case int:
		return int64(val)
	case int32:
		return int64(val)
	case float64:
		return int64(val)
	case float32:
		return int64(val)
	case string:
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

func (s *crudService) buildInsertSQL(md *engine.ModelData, data map[string]any) (string, []any, error) {
	var columns []string
	var placeholders []string
	var args []any

	for _, field := range md.Fields {
		if val, ok := data[field.ColumnName]; ok {
			columns = append(columns, "`"+field.ColumnName+"`")
			placeholders = append(placeholders, "?")
			args = append(args, val)
		}
	}

	if len(columns) == 0 {
		return "", nil, errors.New("no columns to insert")
	}

	// 获取主表名
	tableName := s.getMainTableName(md)
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	return sql, args, nil
}

func (s *crudService) buildGetSQL(md *engine.ModelData, id string) (string, []any, error) {
	tableName := s.getMainTableName(md)
	primaryKey := s.getPrimaryKey(md)

	sql := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = ?", tableName, primaryKey)
	return sql, []any{id}, nil
}

func (s *crudService) buildUpdateSQL(md *engine.ModelData, id string, data map[string]any) (string, []any, error) {
	var setClauses []string
	var args []any

	primaryKey := s.getPrimaryKey(md)

	for _, field := range md.Fields {
		if field.ColumnName == primaryKey {
			continue // 跳过主键
		}
		if val, ok := data[field.ColumnName]; ok {
			setClauses = append(setClauses, fmt.Sprintf("`%s` = ?", field.ColumnName))
			args = append(args, val)
		}
	}

	if len(setClauses) == 0 {
		return "", nil, errors.New("no columns to update")
	}

	tableName := s.getMainTableName(md)
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE `%s` = ?",
		tableName,
		strings.Join(setClauses, ", "),
		primaryKey)
	args = append(args, id)

	return sql, args, nil
}

func (s *crudService) buildDeleteSQL(md *engine.ModelData, id string) (string, []any, error) {
	tableName := s.getMainTableName(md)
	primaryKey := s.getPrimaryKey(md)

	sql := fmt.Sprintf("DELETE FROM %s WHERE `%s` = ?", tableName, primaryKey)
	return sql, []any{id}, nil
}

func (s *crudService) buildListSQL(md *engine.ModelData, params map[string]any) (string, string, []any, error) {
	// 使用 SQLBuilder 构建查询
	sql, args, err := s.sqlBuilder.BuildSQL(md.Model.ID, params)
	if err != nil {
		return "", "", nil, err
	}

	// 构建计数SQL
	countSQL := fmt.Sprintf("SELECT COUNT(*) as count FROM (%s) AS t", sql)

	return sql, countSQL, args, nil
}

func (s *crudService) getMainTableName(md *engine.ModelData) string {
	for _, t := range md.Tables {
		if t.IsMain {
			if t.TableSchema != "" {
				return fmt.Sprintf("`%s`.`%s`", t.TableSchema, t.TableNameStr)
			}
			return "`" + t.TableNameStr + "`"
		}
	}
	if len(md.Tables) > 0 {
		t := md.Tables[0]
		if t.TableSchema != "" {
			return fmt.Sprintf("`%s`.`%s`", t.TableSchema, t.TableNameStr)
		}
		return "`" + t.TableNameStr + "`"
	}
	return ""
}
