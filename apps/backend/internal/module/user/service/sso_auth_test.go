package service

import (
	"context"
	"testing"

	auditModel "metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *MockSsoUserRepository) UpdateUser(user *model.SsoUser) error {
	args := m.Called(user)
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

func TestSsoAuthService_Login(t *testing.T) {
	mockRepo := new(MockSsoUserRepository)
	mockAudit := new(MockAuditService)
	authSvc := NewSsoAuthService(mockRepo, mockAudit)

	t.Run("Success", func(t *testing.T) {
		salt := utils.GenerateSalt()
		hashedPassword := utils.EncryptPasswordSM3("password123", salt)
		user := &model.SsoUser{
			ID:       "1",
			Account:  "admin",
			Password: hashedPassword,
			Salt:     salt,
			State:    1,
		}
		mockRepo.On("GetUserByAccount", "admin").Return(user, nil).Twice() // Once in Login body, once in defer
		mockRepo.On("UpdateLoginInfo", "1", "127.0.0.1").Return(nil).Once()
		mockAudit.On("RecordLogin", mock.Anything, mock.Anything).Return().Once()


		clientInfo := utils.ClientInfo{IP: "127.0.0.1", Browser: "Chrome", OS: "Mac", Platform: "Web"}
		access, refresh, err := authSvc.Login("admin", "password123", 1, clientInfo)
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
			State:    1,
		}
		mockRepo.On("GetUserByAccount", "admin").Return(user, nil).Twice()
		mockRepo.On("IncrementLoginError", "1").Return(nil).Once()
		mockAudit.On("RecordLogin", mock.Anything, mock.Anything).Return().Once()

		clientInfo := utils.ClientInfo{IP: "127.0.0.1", Browser: "Chrome", OS: "Mac", Platform: "Web"}
		_, _, err := authSvc.Login("admin", "wrongpassword", 1, clientInfo)
		assert.Error(t, err)
		assert.Equal(t, "用户名或密码错误", err.Error())
	})
}

func TestSsoAuthService_Refresh(t *testing.T) {
	mockRepo := new(MockSsoUserRepository)
	mockAudit := new(MockAuditService)
	authSvc := NewSsoAuthService(mockRepo, mockAudit)

	user := &model.SsoUser{ID: "1", Account: "admin"}
	
	t.Run("Success", func(t *testing.T) {
		token, _ := utils.GenerateRefreshToken("1")
		mockRepo.On("GetUserByID", "1").Return(user, nil).Once()

		newAccess, err := authSvc.Refresh(token)
		assert.NoError(t, err)
		assert.NotEmpty(t, newAccess)
	})
}
