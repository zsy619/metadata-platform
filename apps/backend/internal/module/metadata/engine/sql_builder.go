package engine

import (
	"fmt"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/repository"
	"strings"

	"gorm.io/gorm"
)

// ModelData 聚合模型所有配置的结构体
type ModelData struct {
	Model   *model.MdModel
	Tables  []*model.MdModelTable
	Fields  []*model.MdModelField
	Joins   []*model.MdModelJoin
	Wheres  []*model.MdModelWhere
	Groups  []*model.MdModelGroup
	Havings []*model.MdModelHaving
	Orders  []*model.MdModelOrder
	Limit   *model.MdModelLimit
	SQL     *model.MdModelSql
}

// SQLBuilder SQL生成引擎主类
type SQLBuilder struct {
	db        *gorm.DB
	modelRepo repository.MdModelRepository
}

// NewSQLBuilder 创建一个新的 SQLBuilder 实例
func NewSQLBuilder(db *gorm.DB, modelRepo repository.MdModelRepository) *SQLBuilder {
	return &SQLBuilder{
		db:        db,
		modelRepo: modelRepo,
	}
}

// BuildSQL 主入口：分发元数据构建或原始 SQL 构建
func (b *SQLBuilder) BuildSQL(modelID string, params map[string]any) (sql string, args []any, err error) {
	data, err := b.LoadModelData(modelID)
	if err != nil {
		return "", nil, err
	}

	if data.Model.ModelKind == 1 {
		// 原始 SQL
		sql, args, err = b.buildFromSQL(data, params)
	} else {
		// 元数据构建
		sql, args, err = b.BuildFromMetadata(data, params)
	}

	if err != nil {
		return "", nil, err
	}

	// SQL 注入与安全性验证
	if err := b.validateSQL(sql); err != nil {
		return "", nil, err
	}

	return sql, args, nil
}

// LoadModelData 加载指定模型的所有配置数据
func (b *SQLBuilder) LoadModelData(modelID string) (*ModelData, error) {
	data := &ModelData{}

	// 1. 加载模型基本信息
	md, err := b.modelRepo.GetModelByID(modelID)
	if err != nil {
		return nil, err
	}
	data.Model = md

	// 2. 加载其他关联配置
	// 加载字段
	var fields []*model.MdModelField
	if err := b.db.Where("model_id = ?", modelID).Order("id asc").Find(&fields).Error; err != nil {
		return nil, err
	}
	data.Fields = fields

	// 加载表
	var tables []*model.MdModelTable
	if err := b.db.Where("model_id = ?", modelID).Find(&tables).Error; err != nil {
		return nil, err
	}
	data.Tables = tables

	// 加载关联
	var joins []*model.MdModelJoin
	if err := b.db.Where("model_id = ?", modelID).Find(&joins).Error; err != nil {
		return nil, err
	}
	data.Joins = joins

	// 加载条件
	var wheres []*model.MdModelWhere
	if err := b.db.Where("model_id = ?", modelID).Order("id asc").Find(&wheres).Error; err != nil {
		return nil, err
	}
	data.Wheres = wheres

	// 加载分组
	var groups []*model.MdModelGroup
	if err := b.db.Where("model_id = ?", modelID).Order("id asc").Find(&groups).Error; err != nil {
		return nil, err
	}
	data.Groups = groups

	// 加载 Having
	var havings []*model.MdModelHaving
	if err := b.db.Where("model_id = ?", modelID).Order("id asc").Find(&havings).Error; err != nil {
		return nil, err
	}
	data.Havings = havings

	// 加载排序
	var orders []*model.MdModelOrder
	if err := b.db.Where("model_id = ?", modelID).Order("id asc").Find(&orders).Error; err != nil {
		return nil, err
	}
	data.Orders = orders

	// 加载分页
	var limit model.MdModelLimit
	if err := b.db.Where("model_id = ?", modelID).First(&limit).Error; err == nil {
		data.Limit = &limit
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 加载原始 SQL
	var modelSQL model.MdModelSql
	if err := b.db.Where("model_id = ?", modelID).First(&modelSQL).Error; err == nil {
		data.SQL = &modelSQL
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return data, nil
}

// BuildFromMetadata 从元数据配置构建完整 SQL
func (b *SQLBuilder) BuildFromMetadata(data *ModelData, params map[string]any) (sql string, args []any, err error) {
	// 按顺序构建各子句
	selectClause, err := b.buildSelectClause(data)
	if err != nil {
		return "", nil, err
	}

	fromClause, err := b.buildFromClause(data)
	if err != nil {
		return "", nil, err
	}

	joinClause, err := b.buildJoinClause(data)
	if err != nil {
		return "", nil, err
	}

	whereClause, whereArgs, err := b.buildWhereClause(data, params)
	if err != nil {
		return "", nil, err
	}
	args = append(args, whereArgs...)

	groupByClause, err := b.buildGroupByClause(data)
	if err != nil {
		return "", nil, err
	}

	havingClause, havingArgs, err := b.buildHavingClause(data, params)
	if err != nil {
		return "", nil, err
	}
	args = append(args, havingArgs...)

	orderByClause, err := b.buildOrderByClause(data)
	if err != nil {
		return "", nil, err
	}

	limitClause, err := b.buildLimitClause(data, params)
	if err != nil {
		return "", nil, err
	}

	// 组装最终 SQL
	var sb strings.Builder
	sb.WriteString(selectClause)
	sb.WriteString(" ")
	sb.WriteString(fromClause)
	if joinClause != "" {
		sb.WriteString(" ")
		sb.WriteString(joinClause)
	}
	if whereClause != "" {
		sb.WriteString(" ")
		sb.WriteString(whereClause)
	}
	if groupByClause != "" {
		sb.WriteString(" ")
		sb.WriteString(groupByClause)
	}
	if havingClause != "" {
		sb.WriteString(" ")
		sb.WriteString(havingClause)
	}
	if orderByClause != "" {
		sb.WriteString(" ")
		sb.WriteString(orderByClause)
	}
	if limitClause != "" {
		sb.WriteString(" ")
		sb.WriteString(limitClause)
	}

	return sb.String(), args, nil
}

// buildSelectClause 构建 SELECT 子句
func (b *SQLBuilder) buildSelectClause(data *ModelData) (string, error) {
	if len(data.Fields) == 0 {
		return "SELECT *", nil
	}

	var expressions []string
	for _, field := range data.Fields {
		expr := b.buildFieldExpression(field)

		// 添加别名
		if field.AggFunc != "" || field.Func != "" || (field.ShowTitle != "" && field.ShowTitle != field.ColumnName) {
			alias := field.ShowTitle
			if alias == "" {
				alias = field.ColumnName
			}
			expr += " AS `" + alias + "`"
		}

		expressions = append(expressions, expr)
	}

	return "SELECT " + strings.Join(expressions, ", "), nil
}

// buildFieldExpression 构建单个字段的表达式
func (b *SQLBuilder) buildFieldExpression(field *model.MdModelField) string {
	columnExpr := ""
	if field.TableNameStr != "" {
		columnExpr = "`" + field.TableNameStr + "`.`" + field.ColumnName + "`"
	} else {
		columnExpr = "`" + field.ColumnName + "`"
	}

	if field.Func != "" {
		if strings.Contains(field.Func, "%s") {
			columnExpr = fmt.Sprintf(field.Func, columnExpr)
		} else {
			columnExpr = field.Func + "(" + columnExpr + ")"
		}
	}

	if field.AggFunc != "" {
		columnExpr = field.AggFunc + "(" + columnExpr + ")"
	}

	return columnExpr
}

// buildFromClause 构建 FROM 子句
func (b *SQLBuilder) buildFromClause(data *ModelData) (string, error) {
	var mainTable *model.MdModelTable
	for _, t := range data.Tables {
		if t.IsMain {
			mainTable = t
			break
		}
	}

	if mainTable == nil && len(data.Tables) > 0 {
		mainTable = data.Tables[0]
	}

	if mainTable == nil {
		return "", fmt.Errorf("no table defined for model %s", data.Model.ID)
	}

	tableName := ""
	if mainTable.TableSchema != "" {
		tableName = "`" + mainTable.TableSchema + "`.`" + mainTable.TableNameStr + "`"
	} else {
		tableName = "`" + mainTable.TableNameStr + "`"
	}

	return "FROM " + tableName, nil
}

// buildJoinClause 构建 JOIN 子句
func (b *SQLBuilder) buildJoinClause(data *ModelData) (string, error) {
	if len(data.Joins) == 0 {
		return "", nil
	}

	joinMap := make(map[string][]*model.MdModelJoin)
	for _, j := range data.Joins {
		joinMap[j.ParentID] = append(joinMap[j.ParentID], j)
	}

	var sb strings.Builder
	if err := b.generateJoinSQL(&sb, "0", joinMap); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func (b *SQLBuilder) generateJoinSQL(sb *strings.Builder, parentID string, joinMap map[string][]*model.MdModelJoin) error {
	joins, ok := joinMap[parentID]
	if !ok {
		return nil
	}

	for _, j := range joins {
		if sb.Len() > 0 {
			sb.WriteString(" ")
		}

		joinType := strings.ToUpper(j.JoinType)
		if !strings.HasSuffix(joinType, "JOIN") {
			joinType += " JOIN"
		}
		sb.WriteString(joinType)
		sb.WriteString(" ")

		if j.JoinTableSchema != "" {
			sb.WriteString("`" + j.JoinTableSchema + "`.`" + j.JoinTableNameStr + "`")
		} else {
			sb.WriteString("`" + j.JoinTableNameStr + "`")
		}
		sb.WriteString(" ON ")

		b.buildJoinConditions(sb, j)

		if err := b.generateJoinSQL(sb, j.ID, joinMap); err != nil {
			return err
		}
	}

	return nil
}

func (b *SQLBuilder) buildJoinConditions(sb *strings.Builder, j *model.MdModelJoin) {
	leftExpr := ""
	if j.TableNameStr != "" {
		leftExpr = "`" + j.TableNameStr + "`.`" + j.ColumnName + "`"
	} else {
		leftExpr = "`" + j.ColumnName + "`"
	}

	rightExpr := ""
	if j.JoinTableNameStr != "" {
		rightExpr = "`" + j.JoinTableNameStr + "`.`" + j.JoinColumnName + "`"
	} else {
		rightExpr = "`" + j.JoinColumnName + "`"
	}

	op := j.Operator2
	if op == "" {
		op = "="
	}

	sb.WriteString(leftExpr)
	sb.WriteString(" ")
	sb.WriteString(op)
	sb.WriteString(" ")
	sb.WriteString(rightExpr)
}

// buildWhereClause 构建 WHERE 子句
func (b *SQLBuilder) buildWhereClause(data *ModelData, params map[string]any) (string, []any, error) {
	if len(data.Wheres) == 0 {
		return "", nil, nil
	}

	var sb strings.Builder
	var args []any

	sb.WriteString("WHERE ")
	for i, w := range data.Wheres {
		if i > 0 {
			op := w.Operator1
			if op == "" {
				op = "AND"
			}
			sb.WriteString(" ")
			sb.WriteString(strings.ToUpper(op))
			sb.WriteString(" ")
		}

		if w.Brackets1 != "" {
			sb.WriteString(w.Brackets1)
		}

		condSQL, condArgs := b.buildSingleCondition(w, params)
		sb.WriteString(condSQL)
		args = append(args, condArgs...)

		if w.Brackets2 != "" {
			sb.WriteString(w.Brackets2)
		}
	}

	return sb.String(), args, nil
}

func (b *SQLBuilder) buildSingleCondition(w *model.MdModelWhere, params map[string]any) (string, []any) {
	leftExpr := ""
	if w.TableNameStr != "" {
		leftExpr = "`" + w.TableNameStr + "`.`" + w.ColumnName + "`"
	} else {
		leftExpr = "`" + w.ColumnName + "`"
	}
	if w.Func != "" {
		if strings.Contains(w.Func, "%s") {
			leftExpr = fmt.Sprintf(w.Func, leftExpr)
		} else {
			leftExpr = w.Func + "(" + leftExpr + ")"
		}
	}

	op := strings.ToUpper(w.Operator2)
	if op == "" {
		op = "="
	}

	var args []any
	rightExpr := "?"

	switch op {
	case "IS NULL", "IS NOT NULL":
		return leftExpr + " " + op, nil
	case "IN", "NOT IN":
		value1 := w.Value1
		if w.ParamKey != "" && params != nil {
			if val, ok := params[w.ParamKey]; ok {
				value1 = fmt.Sprintf("%v", val)
			}
		}
		vals := strings.Split(value1, ",")
		placeholders := make([]string, len(vals))
		for i, v := range vals {
			placeholders[i] = "?"
			args = append(args, strings.TrimSpace(v))
		}
		rightExpr = "(" + strings.Join(placeholders, ", ") + ")"
	case "BETWEEN", "NOT BETWEEN":
		value1 := w.Value1
		value2 := w.Value2
		if w.ParamKey != "" && params != nil {
			if val, ok := params[w.ParamKey]; ok {
				switch v := val.(type) {
				case []any:
					if len(v) >= 2 {
						value1 = fmt.Sprintf("%v", v[0])
						value2 = fmt.Sprintf("%v", v[1])
					}
				case map[string]any:
					if minVal, ok := v["min"]; ok {
						value1 = fmt.Sprintf("%v", minVal)
					}
					if maxVal, ok := v["max"]; ok {
						value2 = fmt.Sprintf("%v", maxVal)
					}
				default:
					value1 = fmt.Sprintf("%v", val)
				}
			}
		}
		rightExpr = "? AND ?"
		args = append(args, value1, value2)
	case "LIKE", "NOT LIKE":
		value1 := w.Value1
		if w.ParamKey != "" && params != nil {
			if val, ok := params[w.ParamKey]; ok {
				value1 = fmt.Sprintf("%v", val)
			}
		}
		args = append(args, "%"+value1+"%")
	default:
		value1 := w.Value1
		if w.ParamKey != "" && params != nil {
			if val, ok := params[w.ParamKey]; ok {
				value1 = fmt.Sprintf("%v", val)
			}
		}
		args = append(args, value1)
	}

	return leftExpr + " " + op + " " + rightExpr, args
}

// buildGroupByClause 构建 GROUP BY 子句
func (b *SQLBuilder) buildGroupByClause(data *ModelData) (string, error) {
	if len(data.Groups) == 0 {
		return "", nil
	}

	var groups []string
	for _, g := range data.Groups {
		expr := ""
		if g.TableNameStr != "" {
			expr = "`" + g.TableNameStr + "`.`" + g.ColumnName + "`"
		} else {
			expr = "`" + g.ColumnName + "`"
		}

		if g.Func != "" {
			if strings.Contains(g.Func, "%s") {
				expr = fmt.Sprintf(g.Func, expr)
			} else {
				expr = g.Func + "(" + expr + ")"
			}
		}
		groups = append(groups, expr)
	}

	return "GROUP BY " + strings.Join(groups, ", "), nil
}

// buildHavingClause 构建 HAVING 子句
func (b *SQLBuilder) buildHavingClause(data *ModelData, params map[string]any) (string, []any, error) {
	if len(data.Havings) == 0 {
		return "", nil, nil
	}

	var sb strings.Builder
	var args []any

	sb.WriteString("HAVING ")
	for i, h := range data.Havings {
		if i > 0 {
			op := h.Operator1
			if op == "" {
				op = "AND"
			}
			sb.WriteString(" ")
			sb.WriteString(strings.ToUpper(op))
			sb.WriteString(" ")
		}

		if h.Brackets1 != "" {
			sb.WriteString(h.Brackets1)
		}

		condSQL, condArgs := b.buildHavingCondition(h, params)
		sb.WriteString(condSQL)
		args = append(args, condArgs...)

		if h.Brackets2 != "" {
			sb.WriteString(h.Brackets2)
		}
	}

	return sb.String(), args, nil
}

func (b *SQLBuilder) buildHavingCondition(h *model.MdModelHaving, params map[string]any) (string, []any) {
	leftExpr := ""
	if h.TableNameStr != "" {
		leftExpr = "`" + h.TableNameStr + "`.`" + h.ColumnName + "`"
	} else {
		leftExpr = "`" + h.ColumnName + "`"
	}

	if h.Func != "" {
		if strings.Contains(h.Func, "%s") {
			leftExpr = fmt.Sprintf(h.Func, leftExpr)
		} else {
			leftExpr = h.Func + "(" + leftExpr + ")"
		}
	}

	op := strings.ToUpper(h.Operator2)
	if op == "" {
		op = "="
	}

	var args []any
	rightExpr := "?"

	switch op {
	case "IS NULL", "IS NOT NULL":
		return leftExpr + " " + op, nil
	case "IN", "NOT IN":
		value1 := h.Value1
		if h.ParamKey != "" && params != nil {
			if val, ok := params[h.ParamKey]; ok {
				value1 = fmt.Sprintf("%v", val)
			}
		}
		vals := strings.Split(value1, ",")
		placeholders := make([]string, len(vals))
		for i, v := range vals {
			placeholders[i] = "?"
			args = append(args, strings.TrimSpace(v))
		}
		rightExpr = "(" + strings.Join(placeholders, ", ") + ")"
	default:
		value1 := h.Value1
		if h.ParamKey != "" && params != nil {
			if val, ok := params[h.ParamKey]; ok {
				value1 = fmt.Sprintf("%v", val)
			}
		}
		args = append(args, value1)
	}

	return leftExpr + " " + op + " " + rightExpr, args
}

// buildOrderByClause 构建 ORDER BY 子句
func (b *SQLBuilder) buildOrderByClause(data *ModelData) (string, error) {
	if len(data.Orders) == 0 {
		return "", nil
	}

	var orders []string
	for _, o := range data.Orders {
		expr := ""
		if o.TableNameStr != "" {
			expr = "`" + o.TableNameStr + "`.`" + o.ColumnName + "`"
		} else {
			expr = "`" + o.ColumnName + "`"
		}

		if o.Func != "" {
			if strings.Contains(o.Func, "%s") {
				expr = fmt.Sprintf(o.Func, expr)
			} else {
				expr = o.Func + "(" + expr + ")"
			}
		}

		orderType := strings.ToUpper(o.OrderType)
		if orderType == "" {
			orderType = "ASC"
		}

		orders = append(orders, expr+" "+orderType)
	}

	return "ORDER BY " + strings.Join(orders, ", "), nil
}

// buildLimitClause 构建 LIMIT 子句
func (b *SQLBuilder) buildLimitClause(data *ModelData, params map[string]any) (string, error) {
	if data.Limit == nil || (data.Limit.Limit == 0 && data.Limit.Page == 0) {
		return "", nil
	}

	limit := data.Limit.Limit
	page := data.Limit.Page

	if params != nil {
		if val, ok := params["limit"]; ok {
			if v, ok := val.(float64); ok && v > 0 {
				limit = int(v)
			}
		}
		if val, ok := params["page"]; ok {
			if v, ok := val.(float64); ok && v > 0 {
				page = int(v)
			}
		}
	}

	if limit <= 0 {
		return "", nil
	}

	offset := 0
	if page > 1 {
		offset = (page - 1) * limit
	}

	if offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset), nil
	}
	return fmt.Sprintf("LIMIT %d", limit), nil
}

// buildFromSQL 处理原始 SQL 模型
func (b *SQLBuilder) buildFromSQL(data *ModelData, params map[string]any) (string, []any, error) {
	if data.SQL == nil || data.SQL.Content == "" {
		return "", nil, fmt.Errorf("raw SQL content is empty for model %s", data.Model.ID)
	}

	sql := data.SQL.Content
	var args []any

	if params != nil && len(params) > 0 {
		for key, val := range params {
			placeholder := ":" + key
			if strings.Contains(sql, placeholder) {
				args = append(args, val)
				sql = strings.Replace(sql, placeholder, "?", 1)
			}
		}
	}

	return sql, args, nil
}

// validateSQL 执行基本的 SQL 安全检查
func (b *SQLBuilder) validateSQL(sql string) error {
	// 检查语句分隔符，防止多重语句注入
	if strings.Contains(sql, ";") {
		trimmed := strings.TrimSpace(sql)
		if strings.Count(trimmed, ";") > 1 || (!strings.HasSuffix(trimmed, ";") && strings.Contains(trimmed, ";")) {
			return fmt.Errorf("multiple SQL statements are not allowed")
		}
	}

	dangerKeywords := []string{"DROP", "TRUNCATE", "ALTER", "GRANT", "REVOKE", "SHUTDOWN", "EXEC"}
	upperSQL := strings.ToUpper(sql)

	for _, kw := range dangerKeywords {
		if strings.Contains(upperSQL, " "+kw+" ") || strings.HasPrefix(upperSQL, kw+" ") {
			return fmt.Errorf("dangerous SQL keyword detected: %s", kw)
		}
	}

	if strings.Count(sql, "(") != strings.Count(sql, ")") {
		return fmt.Errorf("unbalanced parentheses in SQL")
	}

	return nil
}
