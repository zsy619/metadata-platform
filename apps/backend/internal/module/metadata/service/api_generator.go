package service

import (
	"fmt"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/repository"
	"metadata-platform/internal/utils"
	"strings"
)

// APIGenerator API 生成器接口
type APIGenerator interface {
	BatchGenerate(modelID string, userID string, tenantID string) ([]*model.API, error)
}

type apiGenerator struct {
	modelRepo repository.MdModelRepository
	apiRepo   repository.APIRepository
	snowflake *utils.Snowflake
}

// NewAPIGenerator 创建 API 生成器实例
func NewAPIGenerator(modelRepo repository.MdModelRepository, apiRepo repository.APIRepository) APIGenerator {
	return &apiGenerator{
		modelRepo: modelRepo,
		apiRepo:   apiRepo,
		snowflake: utils.NewSnowflake(1, 1),
	}
}

// BatchGenerate 批量生成 CRUD 接口配置
func (g *apiGenerator) BatchGenerate(modelID string, userID string, tenantID string) ([]*model.API, error) {
	md, err := g.modelRepo.GetModelByID(modelID)
	if err != nil {
		return nil, err
	}

	// 基础路径，例如 /api/data/user
	basePath := "/api/data/" + strings.ToLower(md.ModelCode)
	
	// 定义标准 CRUD 模板及扩展接口
	templates := []struct {
		Name       string
		Suffix     string
		Method     string
		CodeSuffix string // 可选，默认使用Method
		Remark     string
	}{
		{"创建" + md.ModelName, "", "POST", "", "自动生成的创建接口"},
		{"查询" + md.ModelName + "列表", "", "GET", "", "自动生成的列表查询接口"},
		{"获取" + md.ModelName + "详情", "/:id", "GET", "", "自动生成的单条查询接口"},
		{"更新" + md.ModelName, "/:id", "PUT", "", "自动生成的更新接口"},
		{"删除" + md.ModelName, "/:id", "DELETE", "", "自动生成的删除接口"},
		
		// 统一查询与批量操作接口
		{"通用查询" + md.ModelName, "/query", "POST", "QUERY", "自动生成的通用查询接口"},
		{"批量创建" + md.ModelName, "/batch-create", "POST", "BATCH_CREATE", "自动生成的批量创建接口"},
		{"批量删除" + md.ModelName, "/batch-delete", "POST", "BATCH_DELETE", "自动生成的批量删除接口"},
		{"数据统计" + md.ModelName, "/statistics", "POST", "STATISTICS", "自动生成的数据统计接口"},
		{"聚合查询" + md.ModelName, "/aggregate", "POST", "AGGREGATE", "自动生成的聚合查询接口"},
	}

	apis := make([]*model.API, 0)
	for _, t := range templates {
		codeSuffix := t.CodeSuffix
		if codeSuffix == "" {
			codeSuffix = t.Method
		}

		api := &model.API{
			ID:        g.snowflake.GenerateIDString(),
			TenantID:  tenantID,
			Name:      t.Name,
			Code:      fmt.Sprintf("%s_%s", md.ModelCode, codeSuffix),
			Path:      basePath + t.Suffix,
			Method:    t.Method,
			IsPublic:  false,
			Status:     1,
			Remark:    t.Remark,
		}

		if err := g.apiRepo.CreateAPI(api); err != nil {
			return nil, err
		}
		apis = append(apis, api)
	}

	return apis, nil
}
