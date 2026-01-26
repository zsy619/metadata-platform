package engine

import (
	"metadata-platform/internal/module/metadata/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLBuilder_ValidateSQL(t *testing.T) {
	builder := &SQLBuilder{}

	tests := []struct {
		name    string
		sql     string
		wantErr bool
		error   string
	}{
		{
			name:    "Safe SELECT",
			sql:     "SELECT * FROM users WHERE id = ?",
			wantErr: false,
		},
		{
			name:    "Dangerous DROP",
			sql:     "SELECT * FROM users; DROP TABLE accounts",
			wantErr: true,
			error:   "multiple SQL statements are not allowed",
		},
		{
			name:    "Dangerous TRUNCATE",
			sql:     "TRUNCATE TABLE users",
			wantErr: true,
			error:   "dangerous SQL keyword detected: TRUNCATE",
		},
		{
			name:    "Unbalanced Parentheses",
			sql:     "SELECT * FROM users WHERE (id = 1",
			wantErr: true,
			error:   "unbalanced parentheses in SQL",
		},
		{
			name:    "Safe with multiple valid spaces",
			sql:     "SELECT name FROM products  WHERE  price > 100",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := builder.validateSQL(tt.sql)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.error)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSQLBuilder_BuildFieldExpression(t *testing.T) {
	builder := &SQLBuilder{}

	tests := []struct {
		name   string
		field  *model.MdModelField
		expect string
	}{
		{
			name: "Simple column",
			field: &model.MdModelField{
				ColumnName: "user_name",
			},
			expect: "`user_name`",
		},
		{
			name: "Table prefixed column",
			field: &model.MdModelField{
				TableNameStr: "users",
				ColumnName:   "id",
			},
			expect: "`users`.`id`",
		},
		{
			name: "Column with function",
			field: &model.MdModelField{
				ColumnName: "created_at",
				Func:       "DATE_FORMAT(%s, '%%Y-%%m-%%d')",
			},
			expect: "DATE_FORMAT(`created_at`, '%Y-%m-%d')",
		},
		{
			name: "Column with AggFunc",
			field: &model.MdModelField{
				ColumnName: "price",
				AggFunc:    "SUM",
			},
			expect: "SUM(`price`)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := builder.buildFieldExpression(tt.field)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestSQLBuilder_BuildLimitClause(t *testing.T) {
	builder := &SQLBuilder{}

	tests := []struct {
		name   string
		limit  *model.MdModelLimit
		expect string
	}{
		{
			name:   "No limit",
			limit:  nil,
			expect: "",
		},
		{
			name: "Basic limit",
			limit: &model.MdModelLimit{
				Limit: 10,
			},
			expect: "LIMIT 10",
		},
		{
			name: "Limit and offset",
			limit: &model.MdModelLimit{
				Limit: 20,
				Page:  2,
			},
			expect: "LIMIT 20 OFFSET 20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &ModelData{Limit: tt.limit}
			got, err := builder.buildLimitClause(data, nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func BenchmarkSQLBuilder_BuildFromMetadata(b *testing.B) {
	builder := &SQLBuilder{}

	// Prepare complex data
	data := &ModelData{
		Model: &model.MdModel{ID: "benchmark", ConnID: "c1"},
		Tables: []*model.MdModelTable{
			{TableNameStr: "orders", IsMain: true},
		},
		Fields: []*model.MdModelField{
			{ColumnName: "id"},
			{ColumnName: "user_id"},
			{ColumnName: "total_amount", AggFunc: "SUM"},
		},
		Joins: []*model.MdModelJoin{
			{JoinType: "left", JoinTableNameStr: "users", TableNameStr: "orders", ColumnName: "user_id", JoinColumnName: "id"},
		},
		Wheres: []*model.MdModelWhere{
			{TableNameStr: "orders", ColumnName: "status", Operator2: "=", Value1: "paid"},
			{TableNameStr: "orders", ColumnName: "created_at", Operator2: ">", Value1: "2023-01-01"},
		},
		Groups: []*model.MdModelGroup{
			{TableNameStr: "orders", ColumnName: "user_id"},
		},
		Orders: []*model.MdModelOrder{
			{TableNameStr: "orders", ColumnName: "total_amount", OrderType: "DESC"},
		},
		Limit: &model.MdModelLimit{
			Limit: 20,
			Page:  1,
		},
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _ = builder.BuildFromMetadata(data, nil)
		}
	})
}
