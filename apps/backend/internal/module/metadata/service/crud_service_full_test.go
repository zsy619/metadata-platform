package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"metadata-platform/internal/module/metadata/engine"
	"metadata-platform/internal/module/metadata/model"
)

// MockCRUDService 用于测试的模拟 CRUD 服务
type MockCRUDService struct {
	CreateFunc        func(ctx context.Context, modelID string, data map[string]any) (map[string]any, error)
	GetFunc           func(modelID string, id string) (map[string]any, error)
	UpdateFunc        func(ctx context.Context, modelID string, id string, data map[string]any) error
	DeleteFunc        func(ctx context.Context, modelID string, id string) error
	ListFunc          func(modelID string, queryParams map[string]any) ([]map[string]any, int64, error)
	BatchCreateFunc   func(modelID string, dataList []map[string]any) ([]map[string]any, error)
	BatchDeleteFunc   func(modelID string, ids []string) error
	StatisticsFunc    func(modelID string, queryParams map[string]any) (map[string]int64, error)
	AggregateFunc     func(modelID string, queryParams map[string]any) ([]map[string]any, error)
	BuildSQLFunc      func(data *engine.ModelData, params map[string]any) (string, []any, error)
	ExecuteModelDataFunc func(data *engine.ModelData, params map[string]any) ([]map[string]any, int64, error)
}

func (m *MockCRUDService) Create(ctx context.Context, modelID string, data map[string]any) (map[string]any, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, modelID, data)
	}
	return nil, nil
}

func (m *MockCRUDService) Get(modelID string, id string) (map[string]any, error) {
	if m.GetFunc != nil {
		return m.GetFunc(modelID, id)
	}
	return nil, nil
}

func (m *MockCRUDService) Update(ctx context.Context, modelID string, id string, data map[string]any) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, modelID, id, data)
	}
	return nil
}

func (m *MockCRUDService) Delete(ctx context.Context, modelID string, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, modelID, id)
	}
	return nil
}

func (m *MockCRUDService) List(modelID string, queryParams map[string]any) ([]map[string]any, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(modelID, queryParams)
	}
	return nil, 0, nil
}

func (m *MockCRUDService) BatchCreate(modelID string, dataList []map[string]any) ([]map[string]any, error) {
	if m.BatchCreateFunc != nil {
		return m.BatchCreateFunc(modelID, dataList)
	}
	return nil, nil
}

func (m *MockCRUDService) BatchDelete(modelID string, ids []string) error {
	if m.BatchDeleteFunc != nil {
		return m.BatchDeleteFunc(modelID, ids)
	}
	return nil
}

func (m *MockCRUDService) Statistics(modelID string, queryParams map[string]any) (map[string]int64, error) {
	if m.StatisticsFunc != nil {
		return m.StatisticsFunc(modelID, queryParams)
	}
	return nil, nil
}

func (m *MockCRUDService) Aggregate(modelID string, queryParams map[string]any) ([]map[string]any, error) {
	if m.AggregateFunc != nil {
		return m.AggregateFunc(modelID, queryParams)
	}
	return nil, nil
}

func (m *MockCRUDService) BuildSQLFromData(data *engine.ModelData, params map[string]any) (string, []any, error) {
	if m.BuildSQLFunc != nil {
		return m.BuildSQLFunc(data, params)
	}
	return "", nil, nil
}

func (m *MockCRUDService) ExecuteModelData(data *engine.ModelData, params map[string]any) ([]map[string]any, int64, error) {
	if m.ExecuteModelDataFunc != nil {
		return m.ExecuteModelDataFunc(data, params)
	}
	return nil, 0, nil
}

func TestCRUDService_Create(t *testing.T) {
	mock := &MockCRUDService{
		CreateFunc: func(ctx context.Context, modelID string, data map[string]any) (map[string]any, error) {
			return data, nil
		},
	}

	ctx := context.Background()
	result, err := mock.Create(ctx, "model_123", map[string]any{
		"name": "张三",
		"age":  25,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "张三", result["name"])
}

func TestCRUDService_Get(t *testing.T) {
	mock := &MockCRUDService{
		GetFunc: func(modelID string, id string) (map[string]any, error) {
			return map[string]any{
				"id":   id,
				"name": "张三",
				"age":  25,
			}, nil
		},
	}

	result, err := mock.Get("model_123", "1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "1", result["id"])
	assert.Equal(t, "张三", result["name"])
}

func TestCRUDService_Update(t *testing.T) {
	mock := &MockCRUDService{
		UpdateFunc: func(ctx context.Context, modelID string, id string, data map[string]any) error {
			assert.Equal(t, "model_123", modelID)
			assert.Equal(t, "1", id)
			assert.Equal(t, "李四", data["name"])
			return nil
		},
	}

	ctx := context.Background()
	err := mock.Update(ctx, "model_123", "1", map[string]any{"name": "李四"})

	assert.NoError(t, err)
}

func TestCRUDService_Delete(t *testing.T) {
	mock := &MockCRUDService{
		DeleteFunc: func(ctx context.Context, modelID string, id string) error {
			assert.Equal(t, "model_123", modelID)
			assert.Equal(t, "1", id)
			return nil
		},
	}

	ctx := context.Background()
	err := mock.Delete(ctx, "model_123", "1")

	assert.NoError(t, err)
}

func TestCRUDService_List(t *testing.T) {
	mock := &MockCRUDService{
		ListFunc: func(modelID string, queryParams map[string]any) ([]map[string]any, int64, error) {
			return []map[string]any{
				{"id": 1, "name": "张三"},
				{"id": 2, "name": "李四"},
			}, 2, nil
		},
	}

	results, count, err := mock.List("model_123", map[string]any{
		"page": 1,
		"limit": 10,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
	assert.Len(t, results, 2)
}

func TestCRUDService_BatchCreate(t *testing.T) {
	mock := &MockCRUDService{
		BatchCreateFunc: func(modelID string, dataList []map[string]any) ([]map[string]any, error) {
			return dataList, nil
		},
	}

	dataList := []map[string]any{
		{"name": "张三", "age": 25},
		{"name": "李四", "age": 30},
	}

	result, err := mock.BatchCreate("model_123", dataList)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCRUDService_BatchDelete(t *testing.T) {
	mock := &MockCRUDService{
		BatchDeleteFunc: func(modelID string, ids []string) error {
			assert.ElementsMatch(t, []string{"1", "2", "3"}, ids)
			return nil
		},
	}

	err := mock.BatchDelete("model_123", []string{"1", "2", "3"})

	assert.NoError(t, err)
}

func TestCRUDService_Statistics(t *testing.T) {
	mock := &MockCRUDService{
		StatisticsFunc: func(modelID string, queryParams map[string]any) (map[string]int64, error) {
			return map[string]int64{"total": 100}, nil
		},
	}

	result, err := mock.Statistics("model_123", nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(100), result["total"])
}

func TestCRUDService_BuildSQLFromData(t *testing.T) {
	mock := &MockCRUDService{
		BuildSQLFunc: func(data *engine.ModelData, params map[string]any) (string, []any, error) {
			return "SELECT * FROM users WHERE id = ?", []any{1}, nil
		},
	}

	data := &engine.ModelData{
		Model: &model.MdModel{ID: "123"},
	}
	result, args, err := mock.BuildSQLFromData(data, map[string]any{"id": 1})

	assert.NoError(t, err)
	assert.Contains(t, result, "SELECT")
	assert.Equal(t, []any{1}, args)
}

func TestCRUDService_ExecuteModelData(t *testing.T) {
	mock := &MockCRUDService{
		ExecuteModelDataFunc: func(data *engine.ModelData, params map[string]any) ([]map[string]any, int64, error) {
			return []map[string]any{
				{"id": 1, "name": "张三"},
			}, 1, nil
		},
	}

	data := &engine.ModelData{
		Model: &model.MdModel{ID: "123"},
	}
	results, count, err := mock.ExecuteModelData(data, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Len(t, results, 1)
}

func TestCRUDService_List_WithPagination(t *testing.T) {
	mock := &MockCRUDService{
		ListFunc: func(modelID string, queryParams map[string]any) ([]map[string]any, int64, error) {
			page := 1
			limit := 10
			if p, ok := queryParams["page"].(float64); ok {
				page = int(p)
			}
			if l, ok := queryParams["limit"].(float64); ok {
				limit = int(l)
			}

			start := (page - 1) * limit
			allData := []map[string]any{
				{"id": 1}, {"id": 2}, {"id": 3}, {"id": 4}, {"id": 5},
			}

			if start >= len(allData) {
				return []map[string]any{}, int64(len(allData)), nil
			}

			end := start + limit
			if end > len(allData) {
				end = len(allData)
			}

			return allData[start:end], int64(len(allData)), nil
		},
	}

	results, count, err := mock.List("model_123", map[string]any{
		"page":  2.0,
		"limit": 2.0,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.Len(t, results, 2)
	assert.Equal(t, 3, results[0]["id"])
}

func TestCRUDService_List_WithFilters(t *testing.T) {
	mock := &MockCRUDService{
		ListFunc: func(modelID string, queryParams map[string]any) ([]map[string]any, int64, error) {
			filters, ok := queryParams["filters"].([]any)
			if !ok || len(filters) == 0 {
				return []map[string]any{}, 0, nil
			}

			return []map[string]any{
				{"id": 1, "name": "张三", "status": "active"},
			}, 1, nil
		},
	}

	results, count, err := mock.List("model_123", map[string]any{
		"filters": []any{
			map[string]any{
				"column_name": "status",
				"operator":   "=",
				"value":      "active",
			},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Len(t, results, 1)
}
