package service

import (
	"metadata-platform/internal/module/metadata/engine"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/repository"
	"metadata-platform/internal/utils"
)

// QueryTemplateService 查询模板服务接口
type QueryTemplateService interface {
	CreateTemplate(template *model.MdQueryTemplate) error
	GetTemplateByID(id string) (*model.MdQueryTemplate, error)
	GetTemplatesByModelID(modelID string) ([]model.MdQueryTemplate, error)
	UpdateTemplate(template *model.MdQueryTemplate) error
	DeleteTemplate(id string) error
	SetDefault(modelID, templateID string) error
	GetDefaultTemplate(modelID string) (*model.MdQueryTemplate, error)
	ApplyTemplate(templateID string, sqlData *engine.ModelData) error
}

type queryTemplateService struct {
	templateRepo repository.MdQueryTemplateRepository
	conditionRepo repository.MdQueryConditionRepository
	snowflake *utils.Snowflake
}

func NewQueryTemplateService(
	templateRepo repository.MdQueryTemplateRepository,
	conditionRepo repository.MdQueryConditionRepository,
) QueryTemplateService {
	snowflake := utils.NewSnowflake(1, 1)
	return &queryTemplateService{
		templateRepo:  templateRepo,
		conditionRepo: conditionRepo,
		snowflake:     snowflake,
	}
}

func (s *queryTemplateService) CreateTemplate(template *model.MdQueryTemplate) error {
	template.ID = s.snowflake.GenerateIDString()
	for i := range template.Conditions {
		template.Conditions[i].ID = s.snowflake.GenerateIDString()
		template.Conditions[i].TemplateID = template.ID
	}
	return s.templateRepo.CreateTemplate(template)
}

func (s *queryTemplateService) GetTemplateByID(id string) (*model.MdQueryTemplate, error) {
	return s.templateRepo.GetTemplateByID(id)
}

func (s *queryTemplateService) GetTemplatesByModelID(modelID string) ([]model.MdQueryTemplate, error) {
	return s.templateRepo.GetTemplatesByModelID(modelID)
}

func (s *queryTemplateService) UpdateTemplate(template *model.MdQueryTemplate) error {
	// Simple update logic: delete old conditions and add new ones (standard procedure for small sets)
	err := s.conditionRepo.DeleteConditionsByTemplateID(template.ID)
	if err != nil {
		return err
	}
	
	for i := range template.Conditions {
		if template.Conditions[i].ID == "" {
			template.Conditions[i].ID = s.snowflake.GenerateIDString()
		}
		template.Conditions[i].TemplateID = template.ID
	}
	
	return s.templateRepo.UpdateTemplate(template)
}

func (s *queryTemplateService) DeleteTemplate(id string) error {
	return s.templateRepo.DeleteTemplate(id)
}

func (s *queryTemplateService) SetDefault(modelID, templateID string) error {
	return s.templateRepo.SetDefault(modelID, templateID)
}

func (s *queryTemplateService) GetDefaultTemplate(modelID string) (*model.MdQueryTemplate, error) {
	// 查找标识为默认的模板
	templates, err := s.templateRepo.GetTemplatesByModelID(modelID)
	if err != nil {
		return nil, err
	}
	for _, t := range templates {
		if t.IsDefault {
			// 加载详情（含条件）
			return s.templateRepo.GetTemplateByID(t.ID)
		}
	}
	return nil, nil
}

func (s *queryTemplateService) ApplyTemplate(templateID string, sqlData *engine.ModelData) error {
	template, err := s.templateRepo.GetTemplateByID(templateID)
	if err != nil {
		return err
	}
	if template == nil {
		return nil
	}

	// 将模板条件转换为引擎 WHERE 条件
	for _, c := range template.Conditions {
		where := &model.MdModelWhere{
			Operator1:    c.Operator1,
			Brackets1:    c.Brackets1,
			TableSchema:  c.TableSchema,
			TableNameStr: c.TableNameStr,
			ColumnName:   c.ColumnName,
			Func:         c.Func,
			Operator2:    c.Operator2,
			Value1:       c.Value1,
			Value2:       c.Value2,
			Brackets2:    c.Brackets2,
		}
		// 转换 logic (这里假设 MdModelWhere 的字段名一致)
		// 注意：MdModelWhere 的 TableNameStr 是 md_model_where 表里的 column:table_name
		sqlData.Wheres = append(sqlData.Wheres, where)
	}

	return nil
}
