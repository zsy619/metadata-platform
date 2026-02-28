package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/module/sso/repository"
	"metadata-platform/internal/utils"
	"time"

	"github.com/google/uuid"
)

type ssoKeyService struct {
	repo repository.SsoKeyRepository
}

func NewSsoKeyService(repo repository.SsoKeyRepository) SsoKeyService {
	return &ssoKeyService{repo: repo}
}

func (s *ssoKeyService) CreateKey(key *model.SsoKey) error {
	return s.repo.CreateKey(key)
}

func (s *ssoKeyService) GetKeyByID(id string) (*model.SsoKey, error) {
	return s.repo.GetKeyByID(id)
}

func (s *ssoKeyService) GetKeyByKeyID(keyID string) (*model.SsoKey, error) {
	return s.repo.GetKeyByKeyID(keyID)
}

func (s *ssoKeyService) GetKeysByProtocolConfigID(protocolConfigID string) ([]model.SsoKey, error) {
	return s.repo.GetKeysByProtocolConfigID(protocolConfigID)
}

func (s *ssoKeyService) GetPrimaryKey(protocolConfigID string) (*model.SsoKey, error) {
	return s.repo.GetPrimaryKey(protocolConfigID)
}

func (s *ssoKeyService) GetValidKeys(protocolConfigID string) ([]model.SsoKey, error) {
	return s.repo.GetValidKeys(protocolConfigID)
}

func (s *ssoKeyService) UpdateKey(key *model.SsoKey) error {
	return s.repo.UpdateKey(key)
}

func (s *ssoKeyService) UpdateKeyFields(id string, fields map[string]any) error {
	return s.repo.UpdateKeyFields(id, fields)
}

func (s *ssoKeyService) DeleteKey(id string) error {
	return s.repo.DeleteKey(id)
}

func (s *ssoKeyService) GetAllKeys() ([]model.SsoKey, error) {
	return s.repo.GetAllKeys()
}

func (s *ssoKeyService) GenerateKeyPair(keyType model.KeyType, algorithm string) (*model.SsoKey, error) {
	keyID := uuid.New().String()
	now := time.Now()

	var key *model.SsoKey

	switch keyType {
	case model.KeyTypeRSA:
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, fmt.Errorf("failed to generate RSA key: %w", err)
		}

		privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		})

		publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal public key: %w", err)
		}
		publicKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		})

		key = &model.SsoKey{
			ID:         uuid.New().String(),
			KeyID:      keyID,
			KeyName:    fmt.Sprintf("RSA Key %s", now.Format("2006-01-02")),
			KeyType:    model.KeyTypeRSA,
			KeyUsage:   model.KeyUsageBoth,
			Algorithm:  "RS256",
			PublicKey:  string(publicKeyPEM),
			PrivateKey: string(privateKeyPEM),
			IsPrimary:  true,
			IsEnabled:  true,
			ValidFrom:  now,
			ValidTo:    now.AddDate(2, 0, 0),
		}

	case model.KeyTypeOctet:
		secret := make([]byte, 32)
		_, err := rand.Read(secret)
		if err != nil {
			return nil, fmt.Errorf("failed to generate secret key: %w", err)
		}

		key = &model.SsoKey{
			ID:        uuid.New().String(),
			KeyID:     keyID,
			KeyName:   fmt.Sprintf("Octet Key %s", now.Format("2006-01-02")),
			KeyType:   model.KeyTypeOctet,
			KeyUsage:  model.KeyUsageSigning,
			Algorithm: "HS256",
			SecretKey: utils.EncodeBase64(secret),
			IsPrimary: true,
			IsEnabled: true,
			ValidFrom: now,
			ValidTo:   now.AddDate(1, 0, 0),
		}

	default:
		return nil, fmt.Errorf("unsupported key type: %s", keyType)
	}

	return key, nil
}
