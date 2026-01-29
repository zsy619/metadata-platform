package model

import "time"

// MdModelRelation 模型关联定义
type MdModelRelation struct {
	ID            string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	MasterModelID string    `json:"master_model_id" form:"master_model_id" gorm:"type:varchar(64);not null;index"`
	DetailModelID string    `json:"detail_model_id" form:"detail_model_id" gorm:"type:varchar(64);not null;index"`
	ForeignKey    string    `json:"foreign_key" form:"foreign_key" gorm:"size:64;not null"`      // 子表中指向主表ID的字段名
	RelationType  string    `json:"relation_type" form:"relation_type" gorm:"size:32;default:'OneToMany'"` // OneToMany, OneToOne
	Remark        string    `json:"remark" form:"remark" gorm:"size:1024;default:'';comment:备注"`
	CreateBy      string    `json:"create_by" form:"create_by" gorm:"size:64;default:''"`
	CreateAt      time.Time `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
}

// TableName 指定表名
func (MdModelRelation) TableName() string {
	return "md_model_relation"
}
