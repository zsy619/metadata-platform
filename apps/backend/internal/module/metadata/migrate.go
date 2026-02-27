package metadata

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/utils"
)

// Migrate 自动迁移元数据模块表结构
func Migrate(db *gorm.DB) error {
	utils.SugarLogger.Info("Starting metadata migration...")

	helper := utils.NewMigrationHelper(db)
	var err error

	// 使用 MigrationHelper 自动迁移
	models := []any{
		&model.API{},
		&model.MdConn{},
		&model.MdTable{},
		&model.MdTableField{},
		&model.MdModel{},
		&model.MdModelField{},
		&model.MdModelGroup{},
		&model.MdModelHaving{},
		&model.MdModelJoin{},
		&model.MdModelJoinField{},
		&model.MdModelLimit{},
		&model.MdModelOrder{},
		&model.MdModelSql{},
		&model.MdModelTable{},
		&model.MdModelWhere{},
		&model.MdModelFieldEnhancement{},
		&model.MdModelRelation{},
		&model.MdQueryTemplate{},
		&model.MdQueryCondition{},
		&model.MdModelParam{},
		&model.MdModelProcedure{},
		&model.MdModelProcedureParam{},
	}

	if err = helper.AutoMigrate(models...); err != nil {
		return err
	}

	// 添加表注释
	comments := map[string]string{
		"md_api":                   "API接口",
		"md_conn":                  "数据连接",
		"md_table":                 "物理表",
		"md_table_field":           "物理表字段",
		"md_model":                 "模型定义",
		"md_model_field":           "模型字段",
		"md_model_group":           "模型字段分组",
		"md_model_having":          "模型聚合过滤",
		"md_model_join":            "模型关联JOIN",
		"md_model_join_field":      "模型关联JOIN字段",
		"md_model_limit":           "模型分页限制",
		"md_model_order":           "模型排序",
		"md_model_sql":             "模型自定义SQL",
		"md_model_table":           "模型关联物理表",
		"md_model_where":           "模型筛选条件",
		"md_model_field_enhance":   "模型字段增强",
		"md_model_relation":        "模型间关联",
		"md_query_template":        "查询模板",
		"md_query_condition":       "查询条件",
		"md_model_param":           "模型参数",
		"md_model_procedure":       "模型存储过程/函数",
		"md_model_procedure_param": "模型存储过程/函数参数",
	}
	helper.AddComments(comments)

	utils.SugarLogger.Info("Metadata migration completed successfully")
	return nil
}
