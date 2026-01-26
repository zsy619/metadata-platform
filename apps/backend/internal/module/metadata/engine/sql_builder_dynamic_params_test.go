package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"metadata-platform/internal/module/metadata/model"
)

func TestSQLBuilder_DynamicParams(t *testing.T) {
	builder := &SQLBuilder{}

	// 创建测试数据
	data := &ModelData{
		Model: &model.MdModel{ID: "test", ConnID: "c1"},
		Tables: []*model.MdModelTable{
			{TableNameStr: "users", IsMain: true},
		},
		Fields: []*model.MdModelField{
			{ColumnName: "id"},
			{ColumnName: "name"},
		},
		Wheres: []*model.MdModelWhere{
			{
				TableNameStr: "users",
				ColumnName:   "status",
				Operator2:    "=",
				Value1:       "active",
				ParamKey:     "status", // 使用动态参数
			},
			{
				TableNameStr: "users",
				ColumnName:   "created_at",
				Operator2:    ">",
				Value1:       "2023-01-01",
				ParamKey:     "start_date", // 使用动态参数
			},
		},
		Orders: []*model.MdModelOrder{
			{
				ColumnName: "id",
				OrderType:  "DESC",
			},
		},
		Limit: &model.MdModelLimit{
			Limit: 20,
			Page:  1,
		},
	}

	// 测试 1: 使用动态参数
	t.Run("With dynamic params", func(t *testing.T) {
		params := map[string]any{
			"status":     "pending",
			"start_date": "2023-06-01",
		}

		sql, args, err := builder.BuildFromMetadata(data, params)
		assert.NoError(t, err)
		assert.Contains(t, sql, "`users`.`status` = ?")
		assert.Contains(t, sql, "`users`.`created_at` > ?")
		assert.Contains(t, sql, "LIMIT 20")
		assert.Equal(t, args, []any{"pending", "2023-06-01"})
	})

	// 测试 2: 部分参数缺失（回退到配置的固定值）
	t.Run("With partial params", func(t *testing.T) {
		params := map[string]any{
			"status": "inactive",
			// 缺少 start_date，应该使用配置的固定值 "2023-01-01"
		}

		sql, args, err := builder.BuildFromMetadata(data, params)
		assert.NoError(t, err)
		assert.Contains(t, sql, "`users`.`status` = ?")
		assert.Contains(t, sql, "`users`.`created_at` > ?")
		assert.Equal(t, args, []any{"inactive", "2023-01-01"})
	})

	// 测试 3: 没有参数（全部使用配置的固定值）
	t.Run("With no params", func(t *testing.T) {
		params := map[string]any{}

		sql, args, err := builder.BuildFromMetadata(data, params)
		assert.NoError(t, err)
		assert.Contains(t, sql, "`users`.`status` = ?")
		assert.Contains(t, sql, "`users`.`created_at` > ?")
		assert.Equal(t, args, []any{"active", "2023-01-01"})
	})

	// 测试 4: IN 条件
	t.Run("IN condition with params", func(t *testing.T) {
		data.Wheres = []*model.MdModelWhere{
			{
				TableNameStr: "users",
				ColumnName:   "id",
				Operator2:    "IN",
				Value1:       "1,2,3",
				ParamKey:     "ids",
			},
		}

		params := map[string]any{
			"ids": "4,5,6",
		}

		sql, args, err := builder.BuildFromMetadata(data, params)
		assert.NoError(t, err)
		assert.Contains(t, sql, "`users`.`id` IN (?, ?, ?)")
		assert.Equal(t, args, []any{"4", "5", "6"})
	})

	// 测试 5: LIKE 条件
	t.Run("LIKE condition with params", func(t *testing.T) {
		data.Wheres = []*model.MdModelWhere{
			{
				TableNameStr: "users",
				ColumnName:   "name",
				Operator2:    "LIKE",
				Value1:       "%test%",
				ParamKey:     "name",
			},
		}

		params := map[string]any{
			"name": "admin",
		}

		sql, args, err := builder.BuildFromMetadata(data, params)
		assert.NoError(t, err)
		assert.Contains(t, sql, "`users`.`name` LIKE ?")
		assert.Equal(t, args, []any{"%admin%"})
	})

	// 测试 6: BETWEEN 条件
	t.Run("BETWEEN condition with params", func(t *testing.T) {
		data.Wheres = []*model.MdModelWhere{
			{
				TableNameStr: "users",
				ColumnName:   "age",
				Operator2:    "BETWEEN",
				Value1:       "18",
				Value2:       "30",
				ParamKey:     "age_range",
			},
		}

		params := map[string]any{
			"age_range": map[string]any{
				"min": 20,
				"max": 35,
			},
		}

		sql, args, err := builder.BuildFromMetadata(data, params)
		assert.NoError(t, err)
		assert.Contains(t, sql, "`users`.`age` BETWEEN ? AND ?")
		assert.Equal(t, args, []any{"20", "35"})
	})
}

func TestSQLBuilder_BuildSQLWithRawSQL(t *testing.T) {
	builder := &SQLBuilder{}

	// 测试原始 SQL 的参数替换

	t.Run("Raw SQL with params", func(t *testing.T) {
		params := map[string]any{
			"status":     "active",
			"start_date": "2023-01-01",
		}

		sql, args, err := builder.BuildSQL("test", params)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM users WHERE status = ? AND created_at > ?", sql)
		assert.Equal(t, args, []any{"active", "2023-01-01"})
	})

	t.Run("Raw SQL with no params", func(t *testing.T) {
		params := map[string]any{}

		sql, args, err := builder.BuildSQL("test", params)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM users WHERE status = :status AND created_at > :start_date", sql)
		assert.Empty(t, args)
	})
}
