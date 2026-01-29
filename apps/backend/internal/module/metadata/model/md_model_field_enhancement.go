package model

import "time"

// MdModelFieldEnhancement 字段增强配置模型
type MdModelFieldEnhancement struct {
	ID              string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID        string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';comment:租户ID"`
	ModelID         string    `json:"model_id" form:"model_id" gorm:"index;type:varchar(64);not null;default:'0';comment:模型ID"`
	FieldID         string    `json:"field_id" form:"field_id" gorm:"index;type:varchar(64);not null;default:'0';comment:字段ID"`
	DisplayName     string    `json:"display_name" form:"display_name" gorm:"size:128;default:''"`    // 显示名称
	DisplayOrder    int       `json:"display_order" form:"display_order" gorm:"not null;default:0"`    // 显示顺序
	DisplayWidth    int       `json:"display_width" form:"display_width" gorm:"not null;default:100"`  // 显示宽度
	IsSearchable    bool      `json:"is_searchable" form:"is_searchable" gorm:"not null;default:true"` // 可搜索
	IsSortable      bool      `json:"is_sortable" form:"is_sortable" gorm:"not null;default:true"`   // 可排序
	IsFilterable    bool      `json:"is_filterable" form:"is_filterable" gorm:"not null;default:true"` // 可筛选
	Placeholder     string    `json:"placeholder" form:"placeholder" gorm:"size:256;default:''"`     // 占位符
	HelpText        string    `json:"help_text" form:"help_text" gorm:"size:512;default:''"`       // 帮助文本
	ComponentType   string    `json:"component_type" form:"component_type" gorm:"size:64;default:''"`   // 组件类型
	ComponentConfig string    `json:"component_config" form:"component_config" gorm:"type:text"`          // 组件配置 JSON
	IsDeleted       bool      `json:"is_deleted" form:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID        string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';comment:创建人ID"`
	CreateBy        string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt        time.Time `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID        string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';comment:更新人ID"`
	UpdateBy        string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt        time.Time `json:"update_at" form:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdModelFieldEnhancement) TableName() string {
	return "md_model_field_enhancements"
}
