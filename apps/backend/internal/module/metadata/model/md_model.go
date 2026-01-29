package model

import "time"

// MdModel 模型定义模型
type MdModel struct {
	ID          string     `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID    string     `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';comment:租户ID"`
	ParentID    string     `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);not null;default:'0';comment:父ID"`
	ConnID      string     `json:"conn_id" form:"conn_id" gorm:"type:varchar(64);not null;default:'0';comment:连接ID"`
	ConnName    string     `json:"conn_name" form:"conn_name" gorm:"size:256;default:'';comment:连接名称"`
	ModelName   string     `json:"model_name" form:"model_name" gorm:"size:128;not null;default:'';comment:模型名称"`
	ModelCode   string     `json:"model_code" form:"model_code" gorm:"size:128;not null;default:'';uniqueIndex:uix_md_model_title_creator;comment:模型编码"`
	ModelVersion string    `json:"model_version" form:"model_version" gorm:"size:64;not null;default:'1.0.0';comment:模型版本"`
	ModelLogo   string     `json:"model_logo" form:"model_logo" gorm:"size:512;not null;default:'';comment:模型Logo"`
	ModelKind   int        `json:"model_kind" form:"model_kind" gorm:"not null;default:0;comment:模型类型：1sql语句、2视图/表、3存储过程、4关联"`
	IsPublic    bool       `json:"is_public" form:"is_public" gorm:"not null;default:false;comment:是否公开"`
	IsLocked    bool       `json:"is_locked" form:"is_locked" gorm:"default:false;comment:是否锁定"`
	IsTree          bool       `json:"is_tree" form:"is_tree" gorm:"default:false;comment:是否树形结构"`                   // 是否树形结构
	TreeParentField string     `json:"tree_parent_field" form:"tree_parent_field" gorm:"size:64;default:'';comment:父节点字段名"`    // 父节点字段名
	TreePathField   string     `json:"tree_path_field" form:"tree_path_field" gorm:"size:64;default:'';comment:路径字段名"`      // 路径字段名
	TreeLevelField  string     `json:"tree_level_field" form:"tree_level_field" gorm:"size:64;default:'';comment:层级字段名"`     // 层级字段名
	Parameters      string     `json:"parameters" form:"parameters" gorm:"type:text;comment:模型参数(JSON)"`
	Remark          string     `json:"remark" form:"remark" gorm:"size:1024;default:'';comment:备注"`
	IsDeleted   bool       `json:"is_deleted" form:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID    string     `json:"create_id" form:"create_id" gorm:"size:64;default:'0';comment:创建人ID"`
	CreateBy    string     `json:"create_by" form:"create_by" gorm:"size:64;default:'';uniqueIndex:uix_md_model_title_creator;comment:创建人"`
	CreateAt    time.Time  `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID    string     `json:"update_id" form:"update_id" gorm:"size:64;default:'0';comment:更新人ID"`
	UpdateBy    string     `json:"update_by" form:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt    time.Time  `json:"update_at" form:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdModel) TableName() string {
	return "md_model"
}
