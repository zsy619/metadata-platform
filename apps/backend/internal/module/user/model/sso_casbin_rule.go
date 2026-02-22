package model

// SsoCasbinRule 权限引擎规则模型
type SsoCasbinRule struct {
	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Ptype string `json:"ptype" gorm:"size:100;column:ptype"`
	V0    string `json:"v0" gorm:"size:100;column:v0"`
	V1    string `json:"v1" gorm:"size:100;column:v1"`
	V2    string `json:"v2" gorm:"size:100;column:v2"`
	V3    string `json:"v3" gorm:"size:100;column:v3"`
	V4    string `json:"v4" gorm:"size:100;column:v4"`
	V5    string `json:"v5" gorm:"size:100;column:v5"`
}

// TableName 指定表名
func (SsoCasbinRule) TableName() string {
	return "casbin_rule"
}
