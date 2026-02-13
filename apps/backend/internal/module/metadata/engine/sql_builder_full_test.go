package engine

import (
	"metadata-platform/internal/module/metadata/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLBuilder_BuildFromMetadata(t *testing.T) {
	builder := &SQLBuilder{}

	data := &ModelData{
		Model: &model.MdModel{
			ID:         "123",
			ModelName:  "User",
			ModelCode:  "user",
			ModelKind:  2,
		},
		Tables: []*model.MdModelTable{
			{
				TableNameStr: "users",
				IsMain:       true,
			},
		},
		Fields: []*model.MdModelField{
			{
				ColumnName:   "id",
				ShowTitle:    "id",
				TableNameStr: "users",
			},
			{
				ColumnName:   "name",
				ShowTitle:    "name",
				TableNameStr: "users",
			},
		},
	}

	params := map[string]any{}

	sql, args, err := builder.BuildFromMetadata(data, params)
	assert.NoError(t, err)
	assert.NotEmpty(t, sql)
	assert.Contains(t, sql, "SELECT")
	assert.Contains(t, sql, "FROM")
	assert.NotNil(t, args)
}

func TestSQLBuilder_BuildSelectClause(t *testing.T) {
	builder := &SQLBuilder{}

	data := &ModelData{
		Fields: []*model.MdModelField{
			{ColumnName: "id", ShowTitle: "id"},
			{ColumnName: "name", ShowTitle: "name"},
		},
	}

	result, err := builder.buildSelectClause(data)
	assert.NoError(t, err)
	assert.Contains(t, result, "`id`")
	assert.Contains(t, result, "`name`")
}

func TestSQLBuilder_BuildSelectClause_WithTablePrefix(t *testing.T) {
	builder := &SQLBuilder{}

	data := &ModelData{
		Fields: []*model.MdModelField{
			{ColumnName: "id", TableNameStr: "users"},
		},
	}

	result, err := builder.buildSelectClause(data)
	assert.NoError(t, err)
	assert.Contains(t, result, "`users`.`id`")
}

func TestSQLBuilder_BuildSelectClause_WithAggregate(t *testing.T) {
	builder := &SQLBuilder{}

	data := &ModelData{
		Fields: []*model.MdModelField{
			{ColumnName: "id", AggFunc: "COUNT", ShowTitle: "total"},
		},
	}

	result, err := builder.buildSelectClause(data)
	assert.NoError(t, err)
	assert.Contains(t, result, "COUNT")
	assert.Contains(t, result, "`total`")
}

func TestSQLBuilder_BuildFromClause(t *testing.T) {
	builder := &SQLBuilder{}

	tests := []struct {
		name   string
		tables []*model.MdModelTable
		expect string
	}{
		{
			name: "Simple table",
			tables: []*model.MdModelTable{
				{TableNameStr: "users", IsMain: true},
			},
			expect: "FROM `users`",
		},
		{
			name: "With schema",
			tables: []*model.MdModelTable{
				{TableNameStr: "users", TableSchema: "public", IsMain: true},
			},
			expect: "FROM `public`.`users`",
		},
		{
			name: "First table as main",
			tables: []*model.MdModelTable{
				{TableNameStr: "users"},
			},
			expect: "FROM `users`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &ModelData{Tables: tt.tables}
			result, err := builder.buildFromClause(data)
			assert.NoError(t, err)
			assert.Contains(t, result, tt.expect)
		})
	}
}

func TestSQLBuilder_BuildWhereClause(t *testing.T) {
	builder := &SQLBuilder{}

	tests := []struct {
		name   string
		wheres []*model.MdModelWhere
		expect string
	}{
		{
			name: "Simple equals",
			wheres: []*model.MdModelWhere{
				{
					ColumnName:   "id",
					TableNameStr: "users",
					Operator2:    "=",
					Value1:       "1",
				},
			},
			expect: "`users`.`id` = ?",
		},
		{
			name: "Like operator",
			wheres: []*model.MdModelWhere{
				{
					ColumnName:   "name",
					TableNameStr: "users",
					Operator2:    "LIKE",
					Value1:       "%test%",
				},
			},
			expect: "`users`.`name` LIKE ?",
		},
		{
			name: "IN operator",
			wheres: []*model.MdModelWhere{
				{
					ColumnName:   "status",
					TableNameStr: "users",
					Operator2:    "IN",
					Value1:       "1,2,3",
				},
			},
			expect: "`users`.`status` IN (?, ?, ?)",
		},
		{
			name: "Between operator",
			wheres: []*model.MdModelWhere{
				{
					ColumnName:   "age",
					TableNameStr: "users",
					Operator2:    "BETWEEN",
					Value1:       "18",
					Value2:       "30",
				},
			},
			expect: "`users`.`age` BETWEEN ? AND ?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &ModelData{Wheres: tt.wheres}
			result, args, err := builder.buildWhereClause(data, nil)
			assert.NoError(t, err)
			assert.Contains(t, result, tt.expect)
			assert.NotEmpty(t, args)
		})
	}
}

func TestSQLBuilder_BuildJoinClause(t *testing.T) {
	builder := &SQLBuilder{}

	data := &ModelData{
		Joins: []*model.MdModelJoin{
			{
				ID:              "join1",
				JoinType:        "LEFT",
				JoinTableNameStr: "orders",
				TableNameStr:     "users",
				ParentID:         "0",
			},
		},
		JoinFields: []*model.MdModelJoinField{
			{
				JoinID:        "join1",
				ColumnName:    "id",
				JoinColumnName: "user_id",
				Operator2:      "=",
			},
		},
	}

	result, err := builder.buildJoinClause(data)
	assert.NoError(t, err)
	assert.Contains(t, result, "LEFT JOIN")
	assert.Contains(t, result, "`orders`")
}

func TestSQLBuilder_BuildJoinClause_Empty(t *testing.T) {
	builder := &SQLBuilder{}

	data := &ModelData{}

	result, err := builder.buildJoinClause(data)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestSQLBuilder_BuildGroupByClause(t *testing.T) {
	builder := &SQLBuilder{}

	data := &ModelData{
		Groups: []*model.MdModelGroup{
			{
				ColumnName:   "status",
				TableNameStr: "users",
			},
		},
	}

	result, err := builder.buildGroupByClause(data)
	assert.NoError(t, err)
	assert.Contains(t, result, "GROUP BY")
	assert.Contains(t, result, "`status`")
}

func TestSQLBuilder_BuildGroupByClause_Empty(t *testing.T) {
	builder := &SQLBuilder{}

	data := &ModelData{}

	result, err := builder.buildGroupByClause(data)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestSQLBuilder_BuildOrderByClause(t *testing.T) {
	builder := &SQLBuilder{}

	tests := []struct {
		name   string
		orders []*model.MdModelOrder
		expect string
	}{
		{
			name: "Ascending order",
			orders: []*model.MdModelOrder{
				{
					ColumnName:   "id",
					TableNameStr: "users",
					OrderType:    "ASC",
				},
			},
			expect: "ORDER BY `users`.`id` ASC",
		},
		{
			name: "Descending order",
			orders: []*model.MdModelOrder{
				{
					ColumnName:   "created_at",
					TableNameStr: "users",
					OrderType:    "DESC",
				},
			},
			expect: "ORDER BY `users`.`created_at` DESC",
		},
		{
			name: "Default order",
			orders: []*model.MdModelOrder{
				{
					ColumnName: "id",
				},
			},
			expect: "ORDER BY `id` ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &ModelData{Orders: tt.orders}
			result, err := builder.buildOrderByClause(data)
			assert.NoError(t, err)
			assert.Contains(t, result, tt.expect)
		})
	}
}

func TestSQLBuilder_BuildOrderByClause_Empty(t *testing.T) {
	builder := &SQLBuilder{}

	data := &ModelData{}

	result, err := builder.buildOrderByClause(data)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestSQLBuilder_BuildFullSQL(t *testing.T) {
	builder := &SQLBuilder{}

	data := &ModelData{
		Model: &model.MdModel{
			ID:         "123",
			ModelName:  "User",
			ModelCode:  "user",
			ModelKind:  2,
		},
		Tables: []*model.MdModelTable{
			{
				TableNameStr: "users",
				IsMain:       true,
			},
		},
		Fields: []*model.MdModelField{
			{
				ColumnName: "id",
				ShowTitle:  "ID",
			},
			{
				ColumnName: "name",
				ShowTitle:  "用户名",
			},
		},
		Wheres: []*model.MdModelWhere{
			{
				ColumnName:   "status",
				Operator2:    "=",
				Value1:       "1",
				Operator1:    "AND",
			},
		},
		Orders: []*model.MdModelOrder{
			{
				ColumnName: "id",
				OrderType:  "DESC",
			},
		},
		Limit: &model.MdModelLimit{
			Limit: 10,
			Page:  1,
		},
	}

	sql, args, err := builder.BuildFromMetadata(data, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, sql)
	assert.Contains(t, sql, "SELECT")
	assert.Contains(t, sql, "FROM")
	assert.Contains(t, sql, "WHERE")
	assert.Contains(t, sql, "ORDER BY")
	assert.Contains(t, sql, "LIMIT")
	assert.NotNil(t, args)
}
