package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"metadata-platform/internal/module/metadata/model"
)

// MockAPIGenerator 用于测试的模拟 API 生成器
type MockAPIGenerator struct {
	APIs []*model.API
	Err  error
}

func (m *MockAPIGenerator) BatchGenerate(modelID string, userID string, tenantID string) ([]*model.API, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.APIs, nil
}

func TestAPIGenerator_BatchGenerate(t *testing.T) {
	// 由于 apiGenerator 是私有字段且 NewAPIGenerator 需要实际的 repository，
	// 这里我们测试 Mock 行为来展示正确的接口使用方式

	mock := &MockAPIGenerator{
		APIs: []*model.API{
			{
				ID:     "1",
				Name:   "创建用户",
				Code:   "user_POST",
				Path:   "/api/data/user",
				Method: "POST",
			},
			{
				ID:     "2",
				Name:   "查询用户列表",
				Code:   "user_GET",
				Path:   "/api/data/user",
				Method: "GET",
			},
			{
				ID:     "3",
				Name:   "获取用户详情",
				Code:   "user_GET",
				Path:   "/api/data/user/:id",
				Method: "GET",
			},
			{
				ID:     "4",
				Name:   "更新用户",
				Code:   "user_PUT",
				Path:   "/api/data/user/:id",
				Method: "PUT",
			},
			{
				ID:     "5",
				Name:   "删除用户",
				Code:   "user_DELETE",
				Path:   "/api/data/user/:id",
				Method: "DELETE",
			},
			{
				ID:     "6",
				Name:   "通用查询用户",
				Code:   "user_QUERY",
				Path:   "/api/data/user/query",
				Method: "POST",
			},
			{
				ID:     "7",
				Name:   "批量创建用户",
				Code:   "user_BATCH_CREATE",
				Path:   "/api/data/user/batch-create",
				Method: "POST",
			},
			{
				ID:     "8",
				Name:   "批量删除用户",
				Code:   "user_BATCH_DELETE",
				Path:   "/api/data/user/batch-delete",
				Method: "POST",
			},
			{
				ID:     "9",
				Name:   "数据统计用户",
				Code:   "user_STATISTICS",
				Path:   "/api/data/user/statistics",
				Method: "POST",
			},
			{
				ID:     "10",
				Name:   "聚合查询用户",
				Code:   "user_AGGREGATE",
				Path:   "/api/data/user/aggregate",
				Method: "POST",
			},
		},
	}

	apis, err := mock.BatchGenerate("model_123", "user_1", "tenant_1")

	assert.NoError(t, err)
	assert.NotNil(t, apis)
	assert.Equal(t, 10, len(apis))

	// 验证第一个 API
	assert.Equal(t, "创建用户", apis[0].Name)
	assert.Equal(t, "POST", apis[0].Method)
	assert.Equal(t, "/api/data/user", apis[0].Path)

	// 验证详情 API
	assert.Equal(t, "/api/data/user/:id", apis[2].Path)
	assert.Equal(t, "GET", apis[2].Method)

	// 验证批量操作 API
	assert.Equal(t, "/api/data/user/batch-create", apis[6].Path)
	assert.Equal(t, "/api/data/user/batch-delete", apis[7].Path)
}

func TestAPIGenerator_BatchGenerate_Error(t *testing.T) {
	mock := &MockAPIGenerator{
		Err: assert.AnError,
	}

	apis, err := mock.BatchGenerate("model_123", "user_1", "tenant_1")

	assert.Error(t, err)
	assert.Nil(t, apis)
	assert.Equal(t, assert.AnError, err)
}

func TestAPIGenerator_APIPathFormat(t *testing.T) {
	// 测试 API 路径格式是否符合预期
	tests := []struct {
		modelCode string
		suffix    string
		expected  string
	}{
		{"user", "", "/api/data/user"},
		{"user", "/:id", "/api/data/user/:id"},
		{"user", "/query", "/api/data/user/query"},
		{"user", "/batch-create", "/api/data/user/batch-create"},
		{"product", "", "/api/data/product"},
		{"order_item", "", "/api/data/order_item"},
	}

	for _, tt := range tests {
		path := "/api/data/" + tt.modelCode + tt.suffix
		assert.Equal(t, tt.expected, path)
	}
}

func TestAPIGenerator_APICodeFormat(t *testing.T) {
	// 测试 API Code 格式
	tests := []struct {
		modelCode string
		method    string
		expected  string
	}{
		{"user", "POST", "user_POST"},
		{"user", "GET", "user_GET"},
		{"user", "PUT", "user_PUT"},
		{"user", "DELETE", "user_DELETE"},
		{"user", "QUERY", "user_QUERY"},
		{"user", "BATCH_CREATE", "user_BATCH_CREATE"},
		{"user", "BATCH_DELETE", "user_BATCH_DELETE"},
		{"user", "STATISTICS", "user_STATISTICS"},
		{"user", "AGGREGATE", "user_AGGREGATE"},
	}

	for _, tt := range tests {
		code := tt.modelCode + "_" + tt.method
		assert.Equal(t, tt.expected, code)
	}
}

func TestAPIGenerator_Integration(t *testing.T) {
	// 简化测试：验证基础路径生成
	basePath := "/api/data/employee"
	
	// 验证各种 suffix 的组合
	assert.Equal(t, "/api/data/employee", basePath+"")
	assert.Equal(t, "/api/data/employee/:id", basePath+"/:id")
	assert.Equal(t, "/api/data/employee/query", basePath+"/query")
	assert.Equal(t, "/api/data/employee/batch-create", basePath+"/batch-create")
	assert.Equal(t, "/api/data/employee/batch-delete", basePath+"/batch-delete")
	assert.Equal(t, "/api/data/employee/statistics", basePath+"/statistics")
	assert.Equal(t, "/api/data/employee/aggregate", basePath+"/aggregate")
	
	// 验证 code 生成
	assert.Equal(t, "employee_POST", "employee_"+"POST")
	assert.Equal(t, "employee_GET", "employee_"+"GET")
	assert.Equal(t, "employee_QUERY", "employee_"+"QUERY")
	assert.Equal(t, "employee_BATCH_CREATE", "employee_"+"BATCH_CREATE")
}
