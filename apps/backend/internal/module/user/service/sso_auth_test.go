package service

import (
	"testing"

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

func TestSsoAuthService_Login(t *testing.T) {
	mockRepo := new(MockSsoUserRepository)
	authSvc := NewSsoAuthService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		salt := utils.GenerateSalt()
		hashedPassword := utils.EncryptPasswordSM3("password123", salt)
		user := &model.SsoUser{
			ID:       "1",
			Account:  "admin",
			Password: hashedPassword,
			Salt:     salt,
		}
		mockRepo.On("GetUserByAccount", "admin").Return(user, nil).Once()

		access, refresh, err := authSvc.Login("admin", "password123", 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, access)
		assert.NotEmpty(t, refresh)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		salt := utils.GenerateSalt()
		hashedPassword := utils.EncryptPasswordSM3("password123", salt)
		user := &model.SsoUser{
			ID:       "1",
			Account:  "admin",
			Password: hashedPassword,
			Salt:     salt,
		}
		mockRepo.On("GetUserByAccount", "admin").Return(user, nil).Once()

		_, _, err := authSvc.Login("admin", "wrongpassword", 1)
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})
}

func TestSsoAuthService_Refresh(t *testing.T) {
	mockRepo := new(MockSsoUserRepository)
	authSvc := NewSsoAuthService(mockRepo)

	user := &model.SsoUser{ID: "1", Account: "admin"}
	
	t.Run("Success", func(t *testing.T) {
		token, _ := utils.GenerateRefreshToken("1")
		mockRepo.On("GetUserByID", "1").Return(user, nil).Once()

		newAccess, err := authSvc.Refresh(token)
		assert.NoError(t, err)
		assert.NotEmpty(t, newAccess)
	})
}
