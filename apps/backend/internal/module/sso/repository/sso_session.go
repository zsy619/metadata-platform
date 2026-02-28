package repository

import (
	"metadata-platform/internal/module/sso/model"
	"time"

	"gorm.io/gorm"
)

type ssoSessionRepository struct {
	db *gorm.DB
}

func NewSsoSessionRepository(db *gorm.DB) SsoSessionRepository {
	return &ssoSessionRepository{db: db}
}

func (r *ssoSessionRepository) CreateSession(session *model.SsoSession) error {
	return r.db.Create(session).Error
}

func (r *ssoSessionRepository) GetSessionByID(id string) (*model.SsoSession, error) {
	var session model.SsoSession
	result := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func (r *ssoSessionRepository) GetSessionBySessionID(sessionID string) (*model.SsoSession, error) {
	var session model.SsoSession
	result := r.db.Where("session_id = ? AND is_deleted = ?", sessionID, false).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func (r *ssoSessionRepository) GetSessionsByUserID(userID string) ([]model.SsoSession, error) {
	var sessions []model.SsoSession
	result := r.db.Where("user_id = ? AND is_deleted = ?", userID, false).
		Order("create_at DESC").Find(&sessions)
	if result.Error != nil {
		return nil, result.Error
	}
	return sessions, nil
}

func (r *ssoSessionRepository) GetSessionsByClientID(clientID string) ([]model.SsoSession, error) {
	var sessions []model.SsoSession
	result := r.db.Where("client_id = ? AND is_deleted = ?", clientID, false).
		Order("create_at DESC").Find(&sessions)
	if result.Error != nil {
		return nil, result.Error
	}
	return sessions, nil
}

func (r *ssoSessionRepository) GetActiveSessionsByUserID(userID string) ([]model.SsoSession, error) {
	var sessions []model.SsoSession
	result := r.db.Where("user_id = ? AND status = ? AND is_deleted = ?",
		userID, model.SessionStatusActive, false).
		Order("create_at DESC").Find(&sessions)
	if result.Error != nil {
		return nil, result.Error
	}
	var activeSessions []model.SsoSession
	for _, session := range sessions {
		if session.IsActive() {
			activeSessions = append(activeSessions, session)
		}
	}
	return activeSessions, nil
}

func (r *ssoSessionRepository) UpdateSession(session *model.SsoSession) error {
	return r.db.Save(session).Error
}

func (r *ssoSessionRepository) UpdateSessionFields(id string, fields map[string]any) error {
	return r.db.Model(&model.SsoSession{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ssoSessionRepository) UpdateSessionStatus(sessionID string, status model.SessionStatus) error {
	return r.db.Model(&model.SsoSession{}).Where("session_id = ?", sessionID).
		Updates(map[string]any{
			"status":    status,
			"update_at": time.Now(),
		}).Error
}

func (r *ssoSessionRepository) DeleteSession(id string) error {
	return r.db.Model(&model.SsoSession{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoSessionRepository) DeleteSessionsByUserID(userID string) error {
	return r.db.Model(&model.SsoSession{}).Where("user_id = ?", userID).Update("is_deleted", true).Error
}

func (r *ssoSessionRepository) RevokeSession(sessionID string) error {
	return r.db.Model(&model.SsoSession{}).Where("session_id = ?", sessionID).
		Updates(map[string]any{
			"status":    model.SessionStatusRevoked,
			"update_at": time.Now(),
		}).Error
}

func (r *ssoSessionRepository) GetAllSessions() ([]model.SsoSession, error) {
	var sessions []model.SsoSession
	result := r.db.Where("is_deleted = ?", false).
		Order("create_at DESC").Find(&sessions)
	if result.Error != nil {
		return nil, result.Error
	}
	return sessions, nil
}
