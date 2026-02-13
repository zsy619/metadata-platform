package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"metadata-platform/internal/module/metadata/engine"
	"metadata-platform/internal/module/metadata/model"
)

// MockQueryTemplateService 用于测试的模拟查询模板服务
type MockQueryTemplateService struct {
	CreateTemplateFunc      func(template *model.MdQueryTemplate) error
	GetTemplateByIDFunc    func(id string) (*model.MdQueryTemplate, error)
	GetTemplatesByModelIDFunc func(modelID string) ([]model.MdQueryTemplate, error)
	UpdateTemplateFunc     func(template *model.MdQueryTemplate) error
	DeleteTemplateFunc     func(id string) error
	SetDefaultFunc         func(modelID, templateID string) error
	GetDefaultTemplateFunc func(modelID string) (*model.MdQueryTemplate, error)
	ApplyTemplateFunc      func(templateID string, sqlData *engine.ModelData) error
	DuplicateTemplateFunc  func(id string) (*model.MdQueryTemplate, error)
}

func (m *MockQueryTemplateService) CreateTemplate(template *model.MdQueryTemplate) error {
	if m.CreateTemplateFunc != nil {
		return m.CreateTemplateFunc(template)
	}
	return nil
}

func (m *MockQueryTemplateService) GetTemplateByID(id string) (*model.MdQueryTemplate, error) {
	if m.GetTemplateByIDFunc != nil {
		return m.GetTemplateByIDFunc(id)
	}
	return nil, nil
}

func (m *MockQueryTemplateService) GetTemplatesByModelID(modelID string) ([]model.MdQueryTemplate, error) {
	if m.GetTemplatesByModelIDFunc != nil {
		return m.GetTemplatesByModelIDFunc(modelID)
	}
	return nil, nil
}

func (m *MockQueryTemplateService) UpdateTemplate(template *model.MdQueryTemplate) error {
	if m.UpdateTemplateFunc != nil {
		return m.UpdateTemplateFunc(template)
	}
	return nil
}

func (m *MockQueryTemplateService) DeleteTemplate(id string) error {
	if m.DeleteTemplateFunc != nil {
		return m.DeleteTemplateFunc(id)
	}
	return nil
}

func (m *MockQueryTemplateService) SetDefault(modelID, templateID string) error {
	if m.SetDefaultFunc != nil {
		return m.SetDefaultFunc(modelID, templateID)
	}
	return nil
}

func (m *MockQueryTemplateService) GetDefaultTemplate(modelID string) (*model.MdQueryTemplate, error) {
	if m.GetDefaultTemplateFunc != nil {
		return m.GetDefaultTemplateFunc(modelID)
	}
	return nil, nil
}

func (m *MockQueryTemplateService) ApplyTemplate(templateID string, sqlData *engine.ModelData) error {
	if m.ApplyTemplateFunc != nil {
		return m.ApplyTemplateFunc(templateID, sqlData)
	}
	return nil
}

func (m *MockQueryTemplateService) DuplicateTemplate(id string) (*model.MdQueryTemplate, error) {
	if m.DuplicateTemplateFunc != nil {
		return m.DuplicateTemplateFunc(id)
	}
	return nil, nil
}

func TestQueryTemplateService_CreateTemplate(t *testing.T) {
	mock := &MockQueryTemplateService{
		CreateTemplateFunc: func(template *model.MdQueryTemplate) error {
			assert.NotNil(t, template)
			assert.NotEmpty(t, template.ModelID)
			return nil
		},
	}

	err := mock.CreateTemplate(&model.MdQueryTemplate{
		ModelID:      "model_123",
		TemplateName: "用户查询模板",
		TemplateCode: "user_query",
		Conditions: []model.MdQueryCondition{
			{
				ColumnName: "status",
				Operator2:  "=",
				Value1:      "1",
			},
		},
	})

	assert.NoError(t, err)
}

func TestQueryTemplateService_GetTemplateByID(t *testing.T) {
	mock := &MockQueryTemplateService{
		GetTemplateByIDFunc: func(id string) (*model.MdQueryTemplate, error) {
			return &model.MdQueryTemplate{
				ID:            id,
				TemplateName:  "用户查询模板",
				TemplateCode:  "user_query",
				ModelID:       "model_123",
				IsDefault:     true,
				Conditions: []model.MdQueryCondition{
					{ColumnName: "status", Operator2: "=", Value1: "1"},
				},
			}, nil
		},
	}

	result, err := mock.GetTemplateByID("template_1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "用户查询模板", result.TemplateName)
	assert.True(t, result.IsDefault)
	assert.Len(t, result.Conditions, 1)
}

func TestQueryTemplateService_GetTemplatesByModelID(t *testing.T) {
	mock := &MockQueryTemplateService{
		GetTemplatesByModelIDFunc: func(modelID string) ([]model.MdQueryTemplate, error) {
			return []model.MdQueryTemplate{
				{ID: "t1", TemplateName: "模板1", ModelID: modelID},
				{ID: "t2", TemplateName: "模板2", ModelID: modelID},
			}, nil
		},
	}

	result, err := mock.GetTemplatesByModelID("model_123")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestQueryTemplateService_UpdateTemplate(t *testing.T) {
	mock := &MockQueryTemplateService{
		UpdateTemplateFunc: func(template *model.MdQueryTemplate) error {
			assert.Equal(t, "template_1", template.ID)
			assert.Equal(t, "新名称", template.TemplateName)
			return nil
		},
	}

	err := mock.UpdateTemplate(&model.MdQueryTemplate{
		ID:           "template_1",
		TemplateName: "新名称",
	})

	assert.NoError(t, err)
}

func TestQueryTemplateService_DeleteTemplate(t *testing.T) {
	mock := &MockQueryTemplateService{
		DeleteTemplateFunc: func(id string) error {
			assert.Equal(t, "template_1", id)
			return nil
		},
	}

	err := mock.DeleteTemplate("template_1")

	assert.NoError(t, err)
}

func TestQueryTemplateService_SetDefault(t *testing.T) {
	mock := &MockQueryTemplateService{
		SetDefaultFunc: func(modelID, templateID string) error {
			assert.Equal(t, "model_123", modelID)
			assert.Equal(t, "template_1", templateID)
			return nil
		},
	}

	err := mock.SetDefault("model_123", "template_1")

	assert.NoError(t, err)
}

func TestQueryTemplateService_GetDefaultTemplate(t *testing.T) {
	mock := &MockQueryTemplateService{
		GetDefaultTemplateFunc: func(modelID string) (*model.MdQueryTemplate, error) {
			return &model.MdQueryTemplate{
				ID:         "default_template",
				TemplateName: "默认模板",
				IsDefault:  true,
			}, nil
		},
	}

	result, err := mock.GetDefaultTemplate("model_123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsDefault)
}

func TestQueryTemplateService_ApplyTemplate(t *testing.T) {
	mock := &MockQueryTemplateService{
		GetTemplateByIDFunc: func(id string) (*model.MdQueryTemplate, error) {
			return &model.MdQueryTemplate{
				ID:            id,
				TemplateName:  "测试模板",
				ModelID:       "model_123",
				Conditions: []model.MdQueryCondition{
					{
						Operator1:    "AND",
						ColumnName:   "status",
						Operator2:    "=",
						Value1:       "1",
					},
				},
			}, nil
		},
	}

	sqlData := &engine.ModelData{
		Model: &model.MdModel{ID: "model_123"},
		Wheres: []*model.MdModelWhere{},
	}

	err := mock.ApplyTemplate("template_1", sqlData)

	assert.NoError(t, err)
	assert.Len(t, sqlData.Wheres, 1)
	assert.Equal(t, "status", sqlData.Wheres[0].ColumnName)
	assert.Equal(t, "=", sqlData.Wheres[0].Operator2)
}

func TestQueryTemplateService_DuplicateTemplate(t *testing.T) {
	mock := &MockQueryTemplateService{
		GetTemplateByIDFunc: func(id string) (*model.MdQueryTemplate, error) {
			return &model.MdQueryTemplate{
				ID:            id,
				TemplateName:  "原始模板",
				TemplateCode:  "original",
				ModelID:       "model_123",
				Conditions: []model.MdQueryCondition{
					{ColumnName: "status", Operator2: "=", Value1: "1"},
				},
			}, nil
		},
		CreateTemplateFunc: func(template *model.MdQueryTemplate) error {
			return nil
		},
	}

	result, err := mock.DuplicateTemplate("template_1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.TemplateName, "Copy")
}

func TestQueryTemplateService_ApplyTemplate_MultipleConditions(t *testing.T) {
	mock := &MockQueryTemplateService{
		GetTemplateByIDFunc: func(id string) (*model.MdQueryTemplate, error) {
			return &model.MdQueryTemplate{
				ID:         id,
				ModelID:    "model_123",
				Conditions: []model.MdQueryCondition{
					{Operator1: "AND", ColumnName: "status", Operator2: "=", Value1: "1"},
					{Operator1: "AND", ColumnName: "type", Operator2: "=", Value1: "vip"},
					{Operator1: "OR", ColumnName: "deleted", Operator2: "=", Value1: "0"},
				},
			}, nil
		},
	}

	sqlData := &engine.ModelData{
		Model:  &model.MdModel{ID: "model_123"},
		Wheres: []*model.MdModelWhere{},
	}

	err := mock.ApplyTemplate("template_1", sqlData)

	assert.NoError(t, err)
	assert.Len(t, sqlData.Wheres, 3)
}

func TestQueryTemplateService_QueryCondition_Operators(t *testing.T) {
	// 测试各种运算符
	operators := []string{
		"=", "<>", "<", ">", "<=", ">=",
		"LIKE", "NOT LIKE", "IN", "NOT IN",
		"BETWEEN", "NOT BETWEEN", "IS NULL", "IS NOT NULL",
	}

	for _, op := range operators {
		cond := &model.MdQueryCondition{
			ColumnName: "test_field",
			Operator2:  op,
			Value1:     "test_value",
		}
		assert.Equal(t, op, cond.Operator2)
	}
}
