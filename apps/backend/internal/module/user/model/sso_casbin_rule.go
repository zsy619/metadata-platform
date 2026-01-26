package model

// SsoCasbinRule 权限引擎规则模型
type SsoCasbinRule struct {
	ID    string `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id"`
	PType string `json:"ptype" form:"ptype" gorm:"size:100;column:ptype;uniqueIndex:idx_casbin_rule"`
	V0    string `json:"v0" form:"v0" gorm:"size:100;column:v0;uniqueIndex:idx_casbin_rule"`
	V1    string `json:"v1" form:"v1" gorm:"size:100;column:v1;uniqueIndex:idx_casbin_rule"`
	V2    string `json:"v2" form:"v2" gorm:"size:100;column:v2;uniqueIndex:idx_casbin_rule"`
	V3    string `json:"v3" form:"v3" gorm:"size:100;column:v3;uniqueIndex:idx_casbin_rule"`
	V4    string `json:"v4" form:"v4" gorm:"size:100;column:v4;uniqueIndex:idx_casbin_rule"`
	V5    string `json:"v5" form:"v5" gorm:"size:100;column:v5;uniqueIndex:idx_casbin_rule"`
}

// TableName 指定表名
func (SsoCasbinRule) TableName() string {
	return "casbin_rule"
}
