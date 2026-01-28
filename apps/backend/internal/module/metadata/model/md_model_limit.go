package model

import "time"

// MdModelLimit 模型-获取条数模型
type MdModelLimit struct {
	ID        string    `json:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID  string    `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';comment:租户ID"`
	ModelID   string    `json:"model_id" gorm:"index;type:varchar(64);not null;default:'0';comment:模型ID"`
	Page      int       `json:"page" gorm:"not null;default:0"`
	Limit           int       `json:"limit" gorm:"not null;default:0"`
	Remark          string    `json:"remark" gorm:"size:1024;default:'';comment:备注"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID  string    `json:"create_id" gorm:"size:64;default:'0';comment:创建人ID"`
	CreateBy  string    `json:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt  time.Time `json:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID  string    `json:"update_id" gorm:"size:64;default:'0';comment:更新人ID"`
	UpdateBy  string    `json:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt  time.Time `json:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdModelLimit) TableName() string {
	return "md_model_limit"
}
