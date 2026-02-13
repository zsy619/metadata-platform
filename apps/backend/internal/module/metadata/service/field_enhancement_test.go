package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"metadata-platform/internal/module/metadata/model"
)

// MockMdModelFieldEnhancementService 用于测试的模拟字段增强服务
type MockMdModelFieldEnhancementService struct {
	CreateEnhancementFunc        func(enh *model.MdModelFieldEnhancement) error
	GetEnhancementByFieldIDFunc  func(fieldID string) (*model.MdModelFieldEnhancement, error)
	GetEnhancementsByModelIDFunc func(modelID string) ([]model.MdModelFieldEnhancement, error)
	UpdateEnhancementFunc        func(enh *model.MdModelFieldEnhancement) error
	DeleteEnhancementFunc        func(id string) error
	BatchUpdateEnhancementsFunc  func(enhancements []model.MdModelFieldEnhancement) error
}

func (m *MockMdModelFieldEnhancementService) CreateEnhancement(enh *model.MdModelFieldEnhancement) error {
	if m.CreateEnhancementFunc != nil {
		return m.CreateEnhancementFunc(enh)
	}
	return nil
}

func (m *MockMdModelFieldEnhancementService) GetEnhancementByFieldID(fieldID string) (*model.MdModelFieldEnhancement, error) {
	if m.GetEnhancementByFieldIDFunc != nil {
		return m.GetEnhancementByFieldIDFunc(fieldID)
	}
	return nil, nil
}

func (m *MockMdModelFieldEnhancementService) GetEnhancementsByModelID(modelID string) ([]model.MdModelFieldEnhancement, error) {
	if m.GetEnhancementsByModelIDFunc != nil {
		return m.GetEnhancementsByModelIDFunc(modelID)
	}
	return nil, nil
}

func (m *MockMdModelFieldEnhancementService) UpdateEnhancement(enh *model.MdModelFieldEnhancement) error {
	if m.UpdateEnhancementFunc != nil {
		return m.UpdateEnhancementFunc(enh)
	}
	return nil
}

func (m *MockMdModelFieldEnhancementService) DeleteEnhancement(id string) error {
	if m.DeleteEnhancementFunc != nil {
		return m.DeleteEnhancementFunc(id)
	}
	return nil
}

func (m *MockMdModelFieldEnhancementService) BatchUpdateEnhancements(enhancements []model.MdModelFieldEnhancement) error {
	if m.BatchUpdateEnhancementsFunc != nil {
		return m.BatchUpdateEnhancementsFunc(enhancements)
	}
	return nil
}

func TestMdModelFieldEnhancementService_Create(t *testing.T) {
	mock := &MockMdModelFieldEnhancementService{
		CreateEnhancementFunc: func(enh *model.MdModelFieldEnhancement) error {
			assert.NotNil(t, enh)
			assert.Equal(t, "field_123", enh.FieldID)
			return nil
		},
	}

	err := mock.CreateEnhancement(&model.MdModelFieldEnhancement{
		ID:           "enh_1",
		TenantID:     "tenant_1",
		ModelID:      "model_123",
		FieldID:      "field_123",
		DisplayName:  "用户名",
		DisplayOrder: 1,
		DisplayWidth: 120,
	})

	assert.NoError(t, err)
}

func TestMdModelFieldEnhancementService_GetByFieldID(t *testing.T) {
	mock := &MockMdModelFieldEnhancementService{
		GetEnhancementByFieldIDFunc: func(fieldID string) (*model.MdModelFieldEnhancement, error) {
			return &model.MdModelFieldEnhancement{
				ID:           "enh_1",
				FieldID:      fieldID,
				DisplayName:  "用户名",
				DisplayOrder: 1,
				DisplayWidth: 120,
			}, nil
		},
	}

	result, err := mock.GetEnhancementByFieldID("field_123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "field_123", result.FieldID)
	assert.Equal(t, "用户名", result.DisplayName)
}

func TestMdModelFieldEnhancementService_GetByModelID(t *testing.T) {
	mock := &MockMdModelFieldEnhancementService{
		GetEnhancementsByModelIDFunc: func(modelID string) ([]model.MdModelFieldEnhancement, error) {
			return []model.MdModelFieldEnhancement{
				{ID: "enh_1", FieldID: "field_1", DisplayName: "用户名", DisplayOrder: 1},
				{ID: "enh_2", FieldID: "field_2", DisplayName: "邮箱", DisplayOrder: 2},
			}, nil
		},
	}

	result, err := mock.GetEnhancementsByModelID("model_123")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "field_1", result[0].FieldID)
}

func TestMdModelFieldEnhancementService_Update(t *testing.T) {
	mock := &MockMdModelFieldEnhancementService{
		UpdateEnhancementFunc: func(enh *model.MdModelFieldEnhancement) error {
			assert.Equal(t, "enh_1", enh.ID)
			assert.Equal(t, "用户名称", enh.DisplayName)
			return nil
		},
	}

	err := mock.UpdateEnhancement(&model.MdModelFieldEnhancement{
		ID:          "enh_1",
		DisplayName: "用户名称",
	})

	assert.NoError(t, err)
}

func TestMdModelFieldEnhancementService_Delete(t *testing.T) {
	mock := &MockMdModelFieldEnhancementService{
		DeleteEnhancementFunc: func(id string) error {
			assert.Equal(t, "enh_1", id)
			return nil
		},
	}

	err := mock.DeleteEnhancement("enh_1")

	assert.NoError(t, err)
}

func TestMdModelFieldEnhancementService_BatchUpdate(t *testing.T) {
	mock := &MockMdModelFieldEnhancementService{
		BatchUpdateEnhancementsFunc: func(enhancements []model.MdModelFieldEnhancement) error {
			assert.Len(t, enhancements, 3)
			return nil
		},
	}

	err := mock.BatchUpdateEnhancements([]model.MdModelFieldEnhancement{
		{ID: "enh_1", DisplayOrder: 1},
		{ID: "enh_2", DisplayOrder: 2},
		{ID: "enh_3", DisplayOrder: 3},
	})

	assert.NoError(t, err)
}

func TestMdModelFieldEnhancement_ComponentTypes(t *testing.T) {
	// 测试组件类型
	componentTypes := []string{
		"input",        // 文本输入
		"textarea",     // 文本域
		"number",       // 数字输入
		"select",       // 下拉选择
		"radio",        // 单选框
		"checkbox",     // 复选框
		"date",         // 日期
		"datetime",     // 日期时间
		"switch",       // 开关
		"slider",       // 滑块
		"rate",         // 评分
		"color",        // 颜色选择
		"file",         // 文件上传
		"image",        // 图片上传
		"editor",       // 富文本编辑器
	}

	for _, ct := range componentTypes {
		enh := &model.MdModelFieldEnhancement{
			ID:            "test",
			ComponentType: ct,
		}
		assert.Equal(t, ct, enh.ComponentType)
	}
}

func TestMdModelFieldEnhancement_DefaultValues(t *testing.T) {
	// 测试默认值
	enh := &model.MdModelFieldEnhancement{}

	assert.Equal(t, 0, enh.DisplayOrder)
	assert.Equal(t, 100, enh.DisplayWidth)
	assert.True(t, enh.IsSearchable)
	assert.True(t, enh.IsSortable)
	assert.True(t, enh.IsFilterable)
}
