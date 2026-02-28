package service

import (
	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/module/sso/repository"
)

// ssoSessionService 会话服务实现
type ssoSessionService struct {
	repo repository.SsoSessionRepository
}

// NewSsoSessionService 创建会话服务实例
func NewSsoSessionService(repo repository.SsoSessionRepository) SsoSessionService {
	return &ssoSessionService{repo: repo}
}

func (s *ssoSessionService) CreateSession(session *model.SsoSession) error {
	return s.repo.CreateSession(session)
}

func (s *ssoSessionService) GetSessionByID(id string) (*model.SsoSession, error) {
	return s.repo.GetSessionByID(id)
}

func (s *ssoSessionService) GetSessionBySessionID(sessionID string) (*model.SsoSession, error) {
	return s.repo.GetSessionBySessionID(sessionID)
}

func (s *ssoSessionService) GetSessionsByUserID(userID string) ([]model.SsoSession, error) {
	return s.repo.GetSessionsByUserID(userID)
}

func (s *ssoSessionService) GetActiveSessionsByUserID(userID string) ([]model.SsoSession, error) {
	return s.repo.GetActiveSessionsByUserID(userID)
}

func (s *ssoSessionService) UpdateSession(session *model.SsoSession) error {
	return s.repo.UpdateSession(session)
}

func (s *ssoSessionService) UpdateSessionFields(id string, fields map[string]any) error {
	return s.repo.UpdateSessionFields(id, fields)
}

func (s *ssoSessionService) UpdateSessionStatus(sessionID string, status model.SessionStatus) error {
	return s.repo.UpdateSessionStatus(sessionID, status)
}

func (s *ssoSessionService) DeleteSession(id string) error {
	return s.repo.DeleteSession(id)
}

func (s *ssoSessionService) DeleteSessionsByUserID(userID string) error {
	return s.repo.DeleteSessionsByUserID(userID)
}

func (s *ssoSessionService) RevokeSession(sessionID string) error {
	return s.repo.RevokeSession(sessionID)
}

func (s *ssoSessionService) GetAllSessions() ([]model.SsoSession, error) {
	return s.repo.GetAllSessions()
}
