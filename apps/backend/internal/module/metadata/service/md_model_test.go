package service

import (
	"errors"
	"metadata-platform/internal/module/metadata/adapter"
	"metadata-platform/internal/module/metadata/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockModelRepo
type MockModelRepo struct {
	mock.Mock
}

func (m *MockModelRepo) CreateModel(model *model.MdModel) error { return m.Called(model).Error(0) }
func (m *MockModelRepo) GetModelByID(id string) (*model.MdModel, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MdModel), args.Error(1)
}

func (m *MockModelRepo) GetModelByCode(code string) (*model.MdModel, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MdModel), args.Error(1)
}
func (m *MockModelRepo) UpdateModel(model *model.MdModel) error { return m.Called(model).Error(0) }
func (m *MockModelRepo) DeleteModel(id string) error            { return m.Called(id).Error(0) }
func (m *MockModelRepo) GetModelsByConnID(connID string) ([]model.MdModel, error) {
	return m.Called(connID).Get(0).([]model.MdModel), m.Called(connID).Error(1)
}

func (m *MockModelRepo) GetModels(tenantID string, offset, limit int, search string, modelKind int) ([]model.MdModel, int64, error) {
	args := m.Called(tenantID, offset, limit, search, modelKind)
	return args.Get(0).([]model.MdModel), args.Get(1).(int64), args.Error(2)
}

func (m *MockModelRepo) GetAllModels(tenantID string) ([]model.MdModel, error) {
	return m.Called(tenantID).Get(0).([]model.MdModel), m.Called(tenantID).Error(1)
}

// MockModelSqlRepo
type MockModelSqlRepo struct {
	mock.Mock
}

func (m *MockModelSqlRepo) Create(sql *model.MdModelSql) error { return m.Called(sql).Error(0) }
func (m *MockModelSqlRepo) GetByModelID(modelID string) (*model.MdModelSql, error) {
	args := m.Called(modelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MdModelSql), args.Error(1)
}
func (m *MockModelSqlRepo) Update(sql *model.MdModelSql) error   { return m.Called(sql).Error(0) }
func (m *MockModelSqlRepo) DeleteByModelID(modelID string) error { return m.Called(modelID).Error(0) }

// MockFieldRepo
type MockFieldRepo struct {
	mock.Mock
}

func (m *MockFieldRepo) CreateField(field *model.MdModelField) error { return m.Called(field).Error(0) }
func (m *MockFieldRepo) GetFieldByID(id string) (*model.MdModelField, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MdModelField), args.Error(1)
}

func (m *MockFieldRepo) GetFieldsByModelID(modelID string) ([]model.MdModelField, error) {
	return m.Called(modelID).Get(0).([]model.MdModelField), m.Called(modelID).Error(1)
}
func (m *MockFieldRepo) UpdateField(field *model.MdModelField) error { return m.Called(field).Error(0) }
func (m *MockFieldRepo) DeleteField(id string) error                 { return m.Called(id).Error(0) }
func (m *MockFieldRepo) DeleteFieldsByModelID(modelID string) error {
	return m.Called(modelID).Error(0)
}

func (m *MockFieldRepo) GetAllFields(tenantID string) ([]model.MdModelField, error) {
	return m.Called(tenantID).Get(0).([]model.MdModelField), m.Called(tenantID).Error(1)
}

// MockConnService
type MockConnService struct {
	mock.Mock
}

func (m *MockConnService) CreateConn(conn *model.MdConn) error { return m.Called(conn).Error(0) }
func (m *MockConnService) GetConnByID(id string) (*model.MdConn, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MdConn), args.Error(1)
}

func (m *MockConnService) GetConnByName(name string) (*model.MdConn, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MdConn), args.Error(1)
}
func (m *MockConnService) UpdateConn(conn *model.MdConn) error { return m.Called(conn).Error(0) }
func (m *MockConnService) DeleteConn(id string) error          { return m.Called(id).Error(0) }
func (m *MockConnService) GetAllConns(tenantID string) ([]model.MdConn, error) {
	return m.Called(tenantID).Get(0).([]model.MdConn), m.Called(tenantID).Error(1)
}

func (m *MockConnService) GetConnsByParentID(parentID string) ([]model.MdConn, error) {
	return m.Called(parentID).Get(0).([]model.MdConn), m.Called(parentID).Error(1)
}
func (m *MockConnService) TestConnection(conn *model.MdConn) error { return m.Called(conn).Error(0) }
func (m *MockConnService) GetTables(conn *model.MdConn, schema string) ([]adapter.TableInfo, error) {
	return m.Called(conn, schema).Get(0).([]adapter.TableInfo), m.Called(conn, schema).Error(1)
}

func (m *MockConnService) GetViews(conn *model.MdConn, schema string) ([]adapter.ViewInfo, error) {
	return m.Called(conn, schema).Get(0).([]adapter.ViewInfo), m.Called(conn, schema).Error(1)
}

func (m *MockConnService) GetTableStructure(conn *model.MdConn, schema, table string) ([]adapter.ColumnInfo, error) {
	return m.Called(conn, schema, table).Get(0).([]adapter.ColumnInfo), m.Called(conn, schema, table).Error(1)
}

func (m *MockConnService) PreviewTableData(conn *model.MdConn, schema, table string, limit int) ([]map[string]interface{}, error) {
	return m.Called(conn, schema, table, limit).Get(0).([]map[string]interface{}), m.Called(conn, schema, table, limit).Error(1)
}

func (m *MockConnService) GetSchemas(conn *model.MdConn) ([]string, error) {
	return m.Called(conn).Get(0).([]string), m.Called(conn).Error(1)
}

func (m *MockConnService) ExecuteSQLForColumns(conn *model.MdConn, query string, params map[string]interface{}) ([]adapter.ColumnInfo, error) {
	return m.Called(conn, query, params).Get(0).([]adapter.ColumnInfo), m.Called(conn, query, params).Error(1)
}

func TestMdModelService_BuildFromTable(t *testing.T) {
	mockModelRepo := new(MockModelRepo)
	mockFieldRepo := new(MockFieldRepo)
	mockModelSqlRepo := new(MockModelSqlRepo)
	mockConnSvc := new(MockConnService)
	svc := NewMdModelService(mockModelRepo, mockFieldRepo, mockModelSqlRepo, mockConnSvc)

	t.Run("Success", func(t *testing.T) {
		req := &BuildFromTableRequest{
			ConnID:    "conn123",
			Schema:    "test_db",
			Table:     "users",
			ModelName: "User Model",
			ModelCode: "user_model",
			TenantID:  "tenant1",
			UserID:    "user1",
			Username:  "admin",
		}

		conn := &model.MdConn{ID: "conn123", ConnName: "MySQL Test"}
		cols := []adapter.ColumnInfo{
			{Name: "id", Type: "int", Comment: "Primary Key"},
			{Name: "name", Type: "varchar(50)", Comment: "User Name"},
		}

		mockConnSvc.On("GetConnByID", "conn123").Return(conn, nil)
		mockConnSvc.On("GetTableStructure", conn, "test_db", "users").Return(cols, nil)
		mockModelRepo.On("GetModelByCode", "user_model").Return(nil, errors.New("not found"))
		mockModelRepo.On("CreateModel", mock.AnythingOfType("*model.MdModel")).Return(nil)
		mockFieldRepo.On("CreateField", mock.AnythingOfType("*model.MdModelField")).Return(nil).Twice()

		err := svc.BuildFromTable(req)
		assert.NoError(t, err)

		mockConnSvc.AssertExpectations(t)
		mockModelRepo.AssertExpectations(t)
		mockFieldRepo.AssertExpectations(t)
	})
}
