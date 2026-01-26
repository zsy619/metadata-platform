package model

import "time"

// MdQueryTemplate 查询模板模型
type MdQueryTemplate struct {
	ID           string    `json:"id" gorm:"primary_key;type:varchar(64)"`
	TenantID     string    `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	ModelID      string    `json:"model_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	TemplateName string    `json:"template_name" gorm:"size:256;not null;default:''"`
	TemplateCode string    `json:"template_code" gorm:"size:128;not null;default:''"`
	IsDefault    bool      `json:"is_default" gorm:"not null;default:false"`
	Remark       string    `json:"remark" gorm:"type:text"`
	IsDeleted    bool      `json:"is_deleted" gorm:"default:false"`
	CreateID     string    `json:"create_id" gorm:"size:64;default:'0'"`
	CreateBy     string    `json:"create_by" gorm:"size:64;default:''"`
	CreateAt     time.Time `json:"create_at"`
	UpdateID     string    `json:"update_id" gorm:"size:64;default:'0'"`
	UpdateBy     string    `json:"update_by" gorm:"size:64;default:''"`
	UpdateAt     time.Time `json:"update_at"`

	// Associations
	Conditions []MdQueryCondition `json:"conditions" gorm:"foreignKey:TemplateID"`
}

// TableName 指定表名
func (MdQueryTemplate) TableName() string {
	return "md_query_template"
}
