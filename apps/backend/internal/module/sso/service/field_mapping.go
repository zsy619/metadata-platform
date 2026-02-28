package service

import (
	"fmt"
	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/module/sso/repository"
)

// ssoFieldMappingService 字段映射服务实现
type ssoFieldMappingService struct {
	repo repository.SsoFieldMappingRepository
}

// NewSsoFieldMappingService 创建字段映射服务实例
func NewSsoFieldMappingService(repo repository.SsoFieldMappingRepository) SsoFieldMappingService {
	return &ssoFieldMappingService{repo: repo}
}

func (s *ssoFieldMappingService) CreateMapping(mapping *model.SsoFieldMapping) error {
	return s.repo.CreateMapping(mapping)
}

func (s *ssoFieldMappingService) GetMappingByID(id string) (*model.SsoFieldMapping, error) {
	return s.repo.GetMappingByID(id)
}

func (s *ssoFieldMappingService) GetMappingsByProtocolConfigID(protocolConfigID string) ([]model.SsoFieldMapping, error) {
	return s.repo.GetMappingsByProtocolConfigID(protocolConfigID)
}

func (s *ssoFieldMappingService) GetMappingsByClientID(clientID string) ([]model.SsoFieldMapping, error) {
	return s.repo.GetMappingsByClientID(clientID)
}

func (s *ssoFieldMappingService) UpdateMapping(mapping *model.SsoFieldMapping) error {
	return s.repo.UpdateMapping(mapping)
}

func (s *ssoFieldMappingService) UpdateMappingFields(id string, fields map[string]any) error {
	return s.repo.UpdateMappingFields(id, fields)
}

func (s *ssoFieldMappingService) DeleteMapping(id string) error {
	return s.repo.DeleteMapping(id)
}

func (s *ssoFieldMappingService) GetAllMappings() ([]model.SsoFieldMapping, error) {
	return s.repo.GetAllMappings()
}

// MapUserFields 映射用户字段
func (s *ssoFieldMappingService) MapUserFields(sourceData map[string]any, mappings []model.SsoFieldMapping) (map[string]any, error) {
	result := make(map[string]any)

	for _, mapping := range mappings {
		if !mapping.IsEnabled {
			continue
		}

		sourceValue, exists := sourceData[mapping.SourceField]
		if !exists {
			if mapping.IsRequired {
				return nil, fmt.Errorf("required field %s not found in source data", mapping.SourceField)
			}
			if mapping.DefaultValue != "" {
				result[mapping.TargetField] = mapping.DefaultValue
			}
			continue
		}

		transformedValue, err := s.transformField(sourceValue, mapping)
		if err != nil {
			return nil, fmt.Errorf("failed to transform field %s: %w", mapping.SourceField, err)
		}

		result[mapping.TargetField] = transformedValue
	}

	return result, nil
}

// transformField 转换字段值
func (s *ssoFieldMappingService) transformField(value any, mapping model.SsoFieldMapping) (any, error) {
	if mapping.TransformScript != "" {
		return value, nil
	}

	switch mapping.FieldType {
	case "string":
		return fmt.Sprintf("%v", value), nil
	case "int":
		return value, nil
	case "bool":
		return value, nil
	case "array":
		return value, nil
	default:
		return value, nil
	}
}
