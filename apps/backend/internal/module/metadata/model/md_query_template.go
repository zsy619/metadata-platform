package model

import "time"

// MdQueryTemplate 查询模板模型
type MdQueryTemplate struct {
	ID           string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID status: 1 有效 0 无效"`
	TenantID     string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';comment:租户ID"`
	ModelID      string    `json:"model_id" form:"model_id" gorm:"index;type:varchar(64);not null;default:'0';comment:模型ID"`
	TemplateName string    `json:"template_name" form:"template_name" gorm:"size:256;not null;default:''"`
	TemplateCode string    `json:"template_code" form:"template_code" gorm:"size:128;not null;default:''"`
	IsDefault    bool      `json:"is_default" form:"is_default" gorm:"default:false"`
	Remark       string    `json:"remark" form:"remark" gorm:"size:1024;default:'';comment:备注"`
	IsDeleted    bool      `json:"is_deleted" form:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID     string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';comment:创建人ID"`
	CreateBy     string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt     time.Time `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID     string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';comment:更新人ID"`
	UpdateBy     string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt     time.Time `json:"update_at" form:"update_at" gorm:"autoUpdateTime;comment:更新时间"`

	// Associations
	Conditions []MdQueryCondition `json:"conditions" gorm:"foreignKey:TemplateID"`
}

// TableName 指定表名
func (MdQueryTemplate) TableName() string {
	return "md_query_template"
}
