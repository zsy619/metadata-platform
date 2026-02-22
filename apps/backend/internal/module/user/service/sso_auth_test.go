package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	auditModel "metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/utils"
)

// MockSsoUserRepository is a mock implementation of repository.SsoUserRepository
type MockSsoUserRepository struct {
	mock.Mock
}

func (m *MockSsoUserRepository) CreateUser(user *model.SsoUser) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockSsoUserRepository) GetUserByID(id string) (*model.SsoUser, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SsoUser), args.Error(1)
}

func (m *MockSsoUserRepository) GetUserByAccount(account string) (*model.SsoUser, error) {
	args := m.Called(account)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SsoUser), args.Error(1)
}

func (m *MockSsoUserRepository) UpdateUser(id string, updates map[string]any) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockSsoUserRepository) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSsoUserRepository) GetAllUsers() ([]model.SsoUser, error) {
	args := m.Called()
	return args.Get(0).([]model.SsoUser), args.Error(1)
}

func (m *MockSsoUserRepository) GetUserWithDetails(id string) (*model.SsoUser, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SsoUser), args.Error(1)
}

func (m *MockSsoUserRepository) UpdateLoginInfo(id string, ip string) error {
	args := m.Called(id, ip)
	return args.Error(0)
}

func (m *MockSsoUserRepository) IncrementLoginError(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSsoUserRepository) GetUserByMobile(mobile string) (*model.SsoUser, error) {
	args := m.Called(mobile)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SsoUser), args.Error(1)
}

func (m *MockSsoUserRepository) GetUserByEmail(email string) (*model.SsoUser, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SsoUser), args.Error(1)
}

// MockAuditService
type MockAuditService struct {
	mock.Mock
}

func (m *MockAuditService) RecordOperation(ctx context.Context, log *auditModel.SysOperationLog) {
	m.Called(ctx, log)
}

func (m *MockAuditService) RecordDataChange(ctx context.Context, log *auditModel.SysDataChangeLog) {
	m.Called(ctx, log)
}

func (m *MockAuditService) RecordLogin(ctx context.Context, log *auditModel.SysLoginLog) {
	m.Called(ctx, log)
}

// Stub implementations for new methods
func (m *MockAuditService) GetLoginLogs(page, pageSize int, filters map[string]interface{}) ([]auditModel.SysLoginLog, int64, error) {
	return nil, 0, nil
}

func (m *MockAuditService) GetOperationLogs(page, pageSize int, filters map[string]interface{}) ([]auditModel.SysOperationLog, int64, error) {
	return nil, 0, nil
}

func (m *MockAuditService) GetDataChangeLogs(page, pageSize int, filters map[string]interface{}) ([]auditModel.SysDataChangeLog, int64, error) {
	return nil, 0, nil
}

func (m *MockAuditService) GetAllLoginLogs(filters map[string]interface{}) ([]auditModel.SysLoginLog, error) {
	return nil, nil
}

func (m *MockAuditService) GetAllOperationLogs(filters map[string]interface{}) ([]auditModel.SysOperationLog, error) {
	return nil, nil
}

func (m *MockAuditService) GetAllDataChangeLogs(filters map[string]interface{}) ([]auditModel.SysDataChangeLog, error) {
	return nil, nil
}

// MockSsoRoleRepository is a mock implementation of repository.SsoRoleRepository
type MockSsoRoleRepository struct {
	mock.Mock
}

func (m *MockSsoRoleRepository) CreateRole(role *model.SsoRole) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockSsoRoleRepository) GetRoleByID(id string) (*model.SsoRole, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SsoRole), args.Error(1)
}

func (m *MockSsoRoleRepository) GetRoleByCode(code string) (*model.SsoRole, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SsoRole), args.Error(1)
}

func (m *MockSsoRoleRepository) UpdateRole(role *model.SsoRole) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockSsoRoleRepository) DeleteRole(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSsoRoleRepository) GetAllRoles() ([]model.SsoRole, error) {
	args := m.Called()
	return args.Get(0).([]model.SsoRole), args.Error(1)
}

func (m *MockSsoRoleRepository) GetMaxSort() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockSsoRoleRepository) GetRolesByUserID(userID string) ([]model.SsoRole, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.SsoRole), args.Error(1)
}

func (m *MockSsoRoleRepository) HasChildren(parentID string) (bool, error) {
	args := m.Called(parentID)
	return args.Bool(0), args.Error(1)
}

// MockSsoUserRoleRepository is a mock implementation of repository.SsoUserRoleRepository
type MockSsoUserRoleRepository struct {
	mock.Mock
}

func (m *MockSsoUserRoleRepository) CreateUserRole(userRole *model.SsoUserRole) error {
	args := m.Called(userRole)
	return args.Error(0)
}

func (m *MockSsoUserRoleRepository) GetUserRoleByID(id string) (*model.SsoUserRole, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SsoUserRole), args.Error(1)
}

func (m *MockSsoUserRoleRepository) DeleteUserRole(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSsoUserRoleRepository) DeleteUserRolesByUserID(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockSsoUserRoleRepository) DeleteUserRolesByRoleID(roleID string) error {
	args := m.Called(roleID)
	return args.Error(0)
}

func (m *MockSsoUserRoleRepository) GetUserRolesByUserID(userID string) ([]model.SsoUserRole, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.SsoUserRole), args.Error(1)
}

func (m *MockSsoUserRoleRepository) GetUserRolesByRoleID(roleID string) ([]model.SsoUserRole, error) {
	args := m.Called(roleID)
	return args.Get(0).([]model.SsoUserRole), args.Error(1)
}

func (m *MockSsoUserRoleRepository) GetAllUserRoles() ([]model.SsoUserRole, error) {
	args := m.Called()
	return args.Get(0).([]model.SsoUserRole), args.Error(1)
}

func TestSsoAuthService_Login(t *testing.T) {
	mockRepo := new(MockSsoUserRepository)
	mockRoleRepo := new(MockSsoRoleRepository)
	mockUserRoleRepo := new(MockSsoUserRoleRepository)
	mockAudit := new(MockAuditService)
	authSvc := NewSsoAuthService(mockRepo, mockRoleRepo, mockUserRoleRepo, mockAudit)

	t.Run("Success", func(t *testing.T) {
		salt := utils.GenerateSalt()
		hashedPassword := utils.EncryptPasswordSM3("password123", salt)
		user := &model.SsoUser{
			ID:       "1",
			Account:  "admin",
			Password: hashedPassword,
			Salt:     salt,
			Status:   1,
		}
		mockRepo.On("GetUserByAccount", "admin").Return(user, nil).Twice()
		mockRepo.On("UpdateLoginInfo", "1", "127.0.0.1").Return(nil).Once()
		mockRoleRepo.On("GetRolesByUserID", "1").Return([]model.SsoRole{}, nil).Once()
		mockAudit.On("RecordLogin", mock.Anything, mock.Anything).Return().Once()

		clientInfo := utils.ClientInfo{IP: "127.0.0.1", Browser: "Chrome", OS: "Mac", Platform: "Web"}
		access, refresh, _, err := authSvc.Login("admin", "password123", 1, clientInfo)
		assert.NoError(t, err)
		assert.NotEmpty(t, access)
		assert.NotEmpty(t, refresh)
		mockRepo.AssertExpectations(t)
		mockAudit.AssertExpectations(t)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		salt := utils.GenerateSalt()
		hashedPassword := utils.EncryptPasswordSM3("password123", salt)
		user := &model.SsoUser{
			ID:       "1",
			Account:  "admin",
			Password: hashedPassword,
			Salt:     salt,
			Status:   1,
		}
		mockRepo.On("GetUserByAccount", "admin").Return(user, nil).Twice()
		mockRepo.On("IncrementLoginError", "1").Return(nil).Once()
		mockAudit.On("RecordLogin", mock.Anything, mock.Anything).Return().Once()

		clientInfo := utils.ClientInfo{IP: "127.0.0.1", Browser: "Chrome", OS: "Mac", Platform: "Web"}
		_, _, _, err := authSvc.Login("admin", "wrongpassword", 1, clientInfo)
		assert.Error(t, err)
		assert.Equal(t, "用户名或密码错误", err.Error())
	})
}

func TestSsoAuthService_Refresh(t *testing.T) {
	mockRepo := new(MockSsoUserRepository)
	mockRoleRepo := new(MockSsoRoleRepository)
	mockUserRoleRepo := new(MockSsoUserRoleRepository)
	mockAudit := new(MockAuditService)
	authSvc := NewSsoAuthService(mockRepo, mockRoleRepo, mockUserRoleRepo, mockAudit)

	user := &model.SsoUser{ID: "1", Account: "admin"}

	t.Run("Success", func(t *testing.T) {
		token, _ := utils.GenerateRefreshToken("1")
		mockRepo.On("GetUserByID", "1").Return(user, nil).Once()

		newAccess, err := authSvc.Refresh(token)
		assert.NoError(t, err)
		assert.NotEmpty(t, newAccess)
	})
}
