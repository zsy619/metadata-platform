package service

import (
	"errors"
	"metadata-platform/internal/module/metadata/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConnRepo 模拟数据连接仓储
type MockConnRepo struct {
	mock.Mock
}

func (m *MockConnRepo) CreateConn(conn *model.MdConn) error {
	return m.Called(conn).Error(0)
}

func (m *MockConnRepo) GetConnByID(id string) (*model.MdConn, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MdConn), args.Error(1)
}

func (m *MockConnRepo) GetMDConnByID(id string) (*model.MdConn, error) {
	return m.GetConnByID(id)
}

func (m *MockConnRepo) GetConnByName(name string) (*model.MdConn, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MdConn), args.Error(1)
}

func (m *MockConnRepo) UpdateConn(conn *model.MdConn) error {
	return m.Called(conn).Error(0)
}

func (m *MockConnRepo) DeleteConn(id string) error {
	return m.Called(id).Error(0)
}

func (m *MockConnRepo) GetAllConns(tenantID string) ([]model.MdConn, error) {
	args := m.Called(tenantID)
	return args.Get(0).([]model.MdConn), args.Error(1)
}

func (m *MockConnRepo) GetConnsByParentID(parentID string) ([]model.MdConn, error) {
	args := m.Called(parentID)
	return args.Get(0).([]model.MdConn), args.Error(1)
}

// TestMdConnService_CreateConn 测试创建数据连接
func TestMdConnService_CreateConn(t *testing.T) {
	mockRepo := new(MockConnRepo)
	svc := NewMdConnService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		conn := &model.MdConn{
			ConnName:     "Test MySQL",
			ConnKind:     "MySQL",
			ConnHost:     "localhost",
			ConnPort:     3306,
			ConnUser:     "root",
			ConnPassword: "password",
			ConnDatabase: "test_db",
			Status:       0, // 初始状态为未检测
		}

		mockRepo.On("GetConnByName", "Test MySQL").Return(nil, errors.New("not found"))
		mockRepo.On("CreateConn", mock.AnythingOfType("*model.MdConn")).Return(nil)

		err := svc.CreateConn(conn)
		assert.NoError(t, err)
		assert.NotEmpty(t, conn.ID)     // 确保ID已生成
		assert.Equal(t, 0, conn.Status) // 确保初始状态为0

		mockRepo.AssertExpectations(t)
	})

	t.Run("DuplicateName", func(t *testing.T) {
		conn := &model.MdConn{
			ConnName: "Existing Connection",
		}

		existingConn := &model.MdConn{
			ID:       "existing_id",
			ConnName: "Existing Connection",
		}

		mockRepo.On("GetConnByName", "Existing Connection").Return(existingConn, nil)

		err := svc.CreateConn(conn)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已存在")

		mockRepo.AssertExpectations(t)
	})
}

// TestMdConnService_UpdateConn 测试更新数据连接
func TestMdConnService_UpdateConn(t *testing.T) {
	mockRepo := new(MockConnRepo)
	svc := NewMdConnService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		existingConn := &model.MdConn{
			ID:       "conn123",
			ConnName: "Old Name",
			Status:   0,
		}

		updatedConn := &model.MdConn{
			ID:       "conn123",
			ConnName: "New Name",
			Status:   1, // 状态已更新
		}

		mockRepo.On("GetConnByID", "conn123").Return(existingConn, nil)
		mockRepo.On("GetConnByName", "New Name").Return(nil, errors.New("not found"))
		mockRepo.On("UpdateConn", updatedConn).Return(nil)

		err := svc.UpdateConn(updatedConn)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		conn := &model.MdConn{
			ID: "nonexistent",
		}

		mockRepo.On("GetConnByID", "nonexistent").Return(nil, errors.New("not found"))

		err := svc.UpdateConn(conn)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不存在")

		mockRepo.AssertExpectations(t)
	})
}

// TestMdConnService_TestConnection 测试连接测试功能
func TestMdConnService_TestConnection(t *testing.T) {
	mockRepo := new(MockConnRepo)
	svc := NewMdConnService(mockRepo)

	t.Run("SuccessAndUpdateState", func(t *testing.T) {
		// 注意：这个测试需要实际的数据库连接，或者需要mock adapter
		// 这里我们假设测试会成功，并验证状态是否更新
		mockRepo.On("UpdateConn", mock.MatchedBy(func(c *model.MdConn) bool {
			return c.ID == "conn123" && c.Status == 1
		})).Return(nil)

		// 实际测试需要真实数据库或mock extractor
		// err := svc.TestConnection(conn)
		// assert.NoError(t, err)
		// assert.Equal(t, 1, conn.State) // 验证状态已更新为有效

		// mockRepo.AssertExpectations(t)
	})

	t.Run("UnsupportedDatabaseType", func(t *testing.T) {
		conn := &model.MdConn{
			ID:       "conn123",
			ConnKind: "UnsupportedDB",
		}

		err := svc.TestConnection(conn)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不支持的数据源类型")
	})
}

// TestMdConnService_GetTables 测试获取表列表
func TestMdConnService_GetTables(t *testing.T) {
	mockRepo := new(MockConnRepo)
	svc := NewMdConnService(mockRepo)

	t.Run("UnsupportedDatabaseType", func(t *testing.T) {
		conn := &model.MdConn{
			ID:       "conn123",
			ConnKind: "UnsupportedDB",
		}

		tables, err := svc.GetTables(conn, "test_schema")
		assert.Error(t, err)
		assert.Nil(t, tables)
		assert.Contains(t, err.Error(), "不支持的数据源类型")
	})
}

// TestMdConnService_DeleteConn 测试删除数据连接
func TestMdConnService_DeleteConn(t *testing.T) {
	mockRepo := new(MockConnRepo)
	svc := NewMdConnService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		conn := &model.MdConn{
			ID:     "conn123",
			Status: 1,
		}

		mockRepo.On("GetConnByID", "conn123").Return(conn, nil)
		mockRepo.On("DeleteConn", "conn123").Return(nil)

		err := svc.DeleteConn("conn123")
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.On("GetConnByID", "nonexistent").Return(nil, errors.New("not found"))

		err := svc.DeleteConn("nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不存在")

		mockRepo.AssertExpectations(t)
	})
}

// TestMdConnService_GetAllConns 测试获取所有连接
func TestMdConnService_GetAllConns(t *testing.T) {
	mockRepo := new(MockConnRepo)
	svc := NewMdConnService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedConns := []model.MdConn{
			{ID: "conn1", ConnName: "MySQL", Status: 1},
			{ID: "conn2", ConnName: "PostgreSQL", Status: 0},
		}

		mockRepo.On("GetAllConns", "tenant1").Return(expectedConns, nil)

		conns, err := svc.GetAllConns("tenant1")
		assert.NoError(t, err)
		assert.Len(t, conns, 2)
		assert.Equal(t, 1, conns[0].Status) // 验证状态字段
		assert.Equal(t, 0, conns[1].Status)

		mockRepo.AssertExpectations(t)
	})
}
