package service

import (
	"metadata-platform/internal/module/metadata/engine"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/repository"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock repositories
type MockMdModelRepo struct{ mock.Mock }

func (m *MockMdModelRepo) CreateModel(md *model.MdModel) error { return m.Called(md).Error(0) }
func (m *MockMdModelRepo) GetModelByID(id string) (*model.MdModel, error) {
	args := m.Called(id)
	return args.Get(0).(*model.MdModel), args.Error(1)
}

func (m *MockMdModelRepo) GetModelByCode(code string) (*model.MdModel, error) {
	args := m.Called(code)
	return args.Get(0).(*model.MdModel), args.Error(1)
}
func (m *MockMdModelRepo) UpdateModel(md *model.MdModel) error                      { return m.Called(md).Error(0) }
func (m *MockMdModelRepo) DeleteModel(id string) error                              { return m.Called(id).Error(0) }
func (m *MockMdModelRepo) GetModelsByConnID(connID string) ([]model.MdModel, error) { return nil, nil }
func (m *MockMdModelRepo) GetAllModels(tenantID string) ([]model.MdModel, error)    { return nil, nil }

type MockMdConnRepo struct{ mock.Mock }

func (m *MockMdConnRepo) CreateConn(conn *model.MdConn) error { return nil }
func (m *MockMdConnRepo) GetConnByID(id string) (*model.MdConn, error) {
	args := m.Called(id)
	return args.Get(0).(*model.MdConn), args.Error(1)
}
func (m *MockMdConnRepo) GetConnByName(name string) (*model.MdConn, error)           { return nil, nil }
func (m *MockMdConnRepo) UpdateConn(conn *model.MdConn) error                        { return nil }
func (m *MockMdConnRepo) DeleteConn(id string) error                                 { return nil }
func (m *MockMdConnRepo) GetAllConns(tenantID string) ([]model.MdConn, error)        { return nil, nil }
func (m *MockMdConnRepo) GetConnsByParentID(parentID string) ([]model.MdConn, error) { return nil, nil }
func (m *MockMdConnRepo) GetMDConnByID(id string) (*model.MdConn, error) {
	args := m.Called(id)
	return args.Get(0).(*model.MdConn), args.Error(1)
}

func TestCRUDService_Lifecycle(t *testing.T) {
	// 1. Setup in-memory DB for both Metadata and Target
	metaDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	targetDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// Create target table
	targetDB.Exec("CREATE TABLE test_users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER)")

	// 2. Setup mocks & Engine
	modelRepo := new(MockMdModelRepo)
	connRepo := new(MockMdConnRepo)

	modelID := "m1"
	connID := "c1"

	builder := engine.NewSQLBuilder(metaDB, modelRepo)
	executor := engine.NewSQLExecutor(metaDB, connRepo)
	executor.SetCustomConnection(connID, targetDB) // Inject target DB
	validator := NewDataValidator()
	queryTemplateRepo := repository.NewMdQueryTemplateRepository(metaDB)
	queryConditionRepo := repository.NewMdQueryConditionRepository(metaDB)
	queryTemplateService := NewQueryTemplateService(queryTemplateRepo, queryConditionRepo)

	svc := NewCRUDService(builder, executor, validator, queryTemplateService)

	// 3. Prepare Metadata in metaDB (since SQLBuilder.LoadModelData queries metaDB)
	metaDB.AutoMigrate(&model.MdModelTable{}, &model.MdModelField{}, &model.MdModelSql{})
	metaDB.Create(&model.MdModelTable{ModelID: modelID, TableNameStr: "test_users", IsMain: true, ConnID: connID})
	metaDB.Create(&model.MdModelField{ModelID: modelID, ColumnName: "id", IsPrimaryKey: true})
	metaDB.Create(&model.MdModelField{ModelID: modelID, ColumnName: "name"})
	metaDB.Create(&model.MdModelField{ModelID: modelID, ColumnName: "age"})

	modelRepo.On("GetModelByID", modelID).Return(&model.MdModel{ID: modelID, ConnID: connID, ModelCode: "test_user"}, nil)

	// 4. Test Create
	t.Run("Create", func(t *testing.T) {
		data := map[string]any{"id": 1, "name": "Alice", "age": 25}
		res, err := svc.Create(modelID, data)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	// 5. Test Get
	t.Run("Get", func(t *testing.T) {
		res, err := svc.Get(modelID, "1")
		assert.NoError(t, err)
		assert.Equal(t, "Alice", res["name"])
		// SQLite might return int64 for INTEGER
		assert.Equal(t, int64(25), res["age"])
	})

	// 6. Test Update
	t.Run("Update", func(t *testing.T) {
		err := svc.Update(modelID, "1", map[string]any{"name": "Alice Smith"})
		assert.NoError(t, err)

		res, _ := svc.Get(modelID, "1")
		assert.Equal(t, "Alice Smith", res["name"])
	})

	// 7. Test Delete
	t.Run("Delete", func(t *testing.T) {
		err := svc.Delete(modelID, "1")
		assert.NoError(t, err)

		res, _ := svc.Get(modelID, "1")
		assert.Nil(t, res)
	})
}
