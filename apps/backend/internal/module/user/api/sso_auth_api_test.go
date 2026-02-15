package api

import (
	"context"
	"encoding/json"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/utils"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/stretchr/testify/mock"
)

// MockSsoAuthService
type MockSsoAuthService struct {
	mock.Mock
}

func (m *MockSsoAuthService) Login(account string, password string, tenantID uint, clientInfo utils.ClientInfo) (string, string, error) {
	args := m.Called(account, password, tenantID, clientInfo)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockSsoAuthService) Logout(ctx context.Context, userID string, clientInfo utils.ClientInfo) error {
	args := m.Called(ctx, userID, clientInfo)
	return args.Error(0)
}

func (m *MockSsoAuthService) Refresh(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}

func (m *MockSsoAuthService) GetUserInfo(userID string) (*model.SsoUser, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SsoUser), args.Error(1)
}

func (m *MockSsoAuthService) ChangePassword(userID string, oldPassword string, newPassword string) error {
	args := m.Called(userID, oldPassword, newPassword)
	return args.Error(0)
}

func BenchmarkSsoAuthLogin(b *testing.B) {
	mockSvc := new(MockSsoAuthService)
	mockSvc.On("Login", "admin", "123456", mock.Anything, mock.Anything).Return("access-token", "refresh-token", nil)

	handler := NewSsoAuthHandler(mockSvc)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c := context.Background()
			ctx := app.NewContext(0)

			body, _ := json.Marshal(SsoLoginRequest{
				Account:     "admin",
				Password:    "123456",
				TenantID:    1,
				CaptchaID:   "test-captcha-id",
				CaptchaCode: "1234",
			})
			ctx.Request.SetBody(body)
			ctx.Request.Header.SetContentTypeBytes([]byte("application/json"))

			handler.Login(c, ctx)

			if ctx.Response.StatusCode() != 200 {
				b.Errorf("expected 200, got %d", ctx.Response.StatusCode())
			}
		}
	})
}
