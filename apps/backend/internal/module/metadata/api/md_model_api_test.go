package api

import (
	"context"
	"encoding/json"
	"testing"

	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/stretchr/testify/mock"
)

// MockMdModelService
type MockMdModelService struct {
	mock.Mock
}

func (m *MockMdModelService) CreateModel(model *model.MdModel) error { return m.Called(model).Error(0) }
func (m *MockMdModelService) GetModelByID(id string) (*model.MdModel, error) { 
	args := m.Called(id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.MdModel), args.Error(1)
}
func (m *MockMdModelService) GetModelByCode(code string) (*model.MdModel, error) {
	args := m.Called(code)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.MdModel), args.Error(1)
}
func (m *MockMdModelService) UpdateModel(model *model.MdModel) error { return m.Called(model).Error(0) }
func (m *MockMdModelService) DeleteModel(id string) error { return m.Called(id).Error(0) }
func (m *MockMdModelService) GetModelsByConnID(connID string) ([]model.MdModel, error) {
	return m.Called(connID).Get(0).([]model.MdModel), m.Called(connID).Error(1)
}
func (m *MockMdModelService) GetAllModels(tenantID string) ([]model.MdModel, error) {
	return m.Called(tenantID).Get(0).([]model.MdModel), m.Called(tenantID).Error(1)
}
func (m *MockMdModelService) BuildFromTable(req *service.BuildFromTableRequest) error { return m.Called(req).Error(0) }
func (m *MockMdModelService) BuildFromView(req *service.BuildFromViewRequest) error { return m.Called(req).Error(0) }
func (m *MockMdModelService) GetFieldsByModelID(modelID string) ([]model.MdModelField, error) {
	return m.Called(modelID).Get(0).([]model.MdModelField), m.Called(modelID).Error(1)
}
func (m *MockMdModelService) CreateField(field *model.MdModelField) error { return m.Called(field).Error(0) }
func (m *MockMdModelService) UpdateField(field *model.MdModelField) error { return m.Called(field).Error(0) }
func (m *MockMdModelService) DeleteField(id string) error { return m.Called(id).Error(0) }

func BenchmarkBuildFromTable(b *testing.B) {
	mockSvc := new(MockMdModelService)
	handler := NewMdModelHandler(mockSvc)
	
	mockSvc.On("BuildFromTable", mock.Anything).Return(nil)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c := context.Background()
			ctx := app.NewContext(0)
			
			body, _ := json.Marshal(BuildFromTableRequest{
				ConnID:    "1",
				Schema:    "test",
				Table:     "users",
				ModelName: "User",
				ModelCode: "user",
			})
			ctx.Request.SetBody(body)
			ctx.Request.Header.SetContentTypeBytes([]byte("application/json"))
			// inject simulated tenant/user info
			ctx.Set("tenant_id", uint(1))
			ctx.Set("user_id", int64(1))
			ctx.Set("username", "admin")
			
			handler.BuildFromTable(c, ctx)
			
			if ctx.Response.StatusCode() != 201 {
				b.Errorf("expected 201, got %d. Body: %s", ctx.Response.StatusCode(), ctx.Response.Body())
			}
		}
	})
}
