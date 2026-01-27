package model

import "time"

// SsoUser 用户模型
type SsoUser struct {
	ID              string     `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id"`
	AccountID       string     `json:"account_id" form:"account_id" gorm:"type:varchar(64);column:account_id"`
	ApplicationCode string     `json:"application_code" form:"application_code" gorm:"size:64;column:application_code"`
	Account         string     `json:"account" form:"account" gorm:"size:128;uniqueIndex;column:account"`
	Password        string     `json:"password" form:"password" gorm:"size:64;column:password"`
	Salt            string     `json:"salt" form:"salt" gorm:"size:64;column:salt"` // Keep 64 for security upgrade, SQL said 32
	Name            string     `json:"name" form:"name" gorm:"size:128;column:name"`
	Code            string     `json:"code" form:"code" gorm:"size:64;column:code"`
	Sex             string     `json:"sex" form:"sex" gorm:"size:6;default:'男';column:sex"`
	Idcard          string     `json:"idcard" form:"idcard" gorm:"size:32;column:idcard"`
	Mobile          string     `json:"mobile" form:"mobile" gorm:"size:32;column:mobile"`
	Email           string     `json:"email" form:"email" gorm:"size:128;column:email"`
	Avatar          string     `json:"avatar" form:"avatar" gorm:"size:128;column:avatar"`
	OrganizationID  string     `json:"organization_id" form:"organization_id" gorm:"type:varchar(64);column:organization_id"`
	State           int        `json:"state" form:"state" gorm:"not null;default:1;column:state"`
	EndTime         *time.Time `json:"end_time" form:"end_time" gorm:"type:datetime;column:end_time"`
	Kind            int        `json:"kind" form:"kind" gorm:"default:2;column:kind"`
	Remark          string     `json:"remark" form:"remark" gorm:"size:512;column:remark"`
	Sort            int        `json:"sort" form:"sort" gorm:"default:0;column:sort"`
	FirstLogin      int        `json:"first_login" form:"first_login" gorm:"default:0;column:first_login"`
	LastLoginTime   *time.Time `json:"last_login_time" form:"last_login_time" gorm:"type:datetime;column:last_login_time"`
	LastIP          string     `json:"last_ip" form:"last_ip" gorm:"size:32;column:last_ip"`
	LoginErrorCount int        `json:"login_error_count" form:"login_error_count" gorm:"default:0;column:login_error_count"`
	IsDeleted       bool       `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;column:is_deleted"`
	TenantID        string     `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';column:tenant_id"`
	CreateID        string     `json:"create_id" form:"create_id" gorm:"size:64;default:'0';column:create_id"`
	CreateBy        string     `json:"create_by" form:"create_by" gorm:"size:64;column:create_by"`
	CreateAt        time.Time  `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime"`
	UpdateID        string     `json:"update_id" form:"update_id" gorm:"size:64;default:'0';column:update_id"`
	UpdateBy        string     `json:"update_by" form:"update_by" gorm:"size:64;column:update_by"`
	UpdateAt        time.Time  `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime"`

	// 关联数据 (不存入 sso_user 表)
	Roles         []SsoRole         `json:"roles" gorm:"many2many:sso_user_role;joinForeignKey:UserID;joinReferences:RoleID"`
	Positions     []SsoPosition     `json:"positions" gorm:"many2many:sso_user_position;joinForeignKey:UserID;joinReferences:PositionID"`
	Organizations []SsoOrganization `json:"organizations" gorm:"many2many:sso_organization_user;joinForeignKey:UserID;joinReferences:OrganizationID"`
}

// TableName 指定表名
func (SsoUser) TableName() string {
	return "sso_user"
}
