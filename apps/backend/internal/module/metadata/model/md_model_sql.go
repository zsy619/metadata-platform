package model

import "time"

// MdModelSql 模型-sql模型
type MdModelSql struct {
	ID        string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID  string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'';comment:租户ID"`
	ModelID   string    `json:"model_id" form:"model_id" gorm:"index;type:varchar(64);not null;default:'';comment:模型ID"`
	Content   string    `json:"content" form:"content" gorm:"type:text;comment:查询内容：sql语句"`
	Remark    string    `json:"remark" form:"remark" gorm:"size:1024;default:'';comment:备注"`
	IsDeleted bool      `json:"is_deleted" form:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID  string    `json:"create_id" form:"create_id" gorm:"size:64;default:'';comment:创建人ID"`
	CreateBy  string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt  time.Time `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID  string    `json:"update_id" form:"update_id" gorm:"size:64;default:'';comment:更新人ID"`
	UpdateBy  string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt  time.Time `json:"update_at" form:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdModelSql) TableName() string {
	return "md_model_sql"
}
