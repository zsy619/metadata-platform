package model

import "time"

// SsoUser 用户模型
type SsoUser struct {
	ID              string     `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	AccountID       string     `json:"account_id" form:"account_id" gorm:"type:varchar(64);column:account_id;comment:账号ID"`
	AppCode         string     `json:"app_code" form:"app_code" gorm:"size:64;column:app_code;comment:应用编码"`
	Account         string     `json:"account" form:"account" gorm:"size:128;uniqueIndex;column:account;comment:账号"`
	Password        string     `json:"password" form:"password" gorm:"size:64;column:password;comment:密码"`
	Salt            string     `json:"salt" form:"salt" gorm:"size:64;column:salt;comment:盐值"` // Keep 64 for security upgrade, SQL said 32
	Name            string     `json:"name" form:"name" gorm:"size:128;column:name;comment:姓名"`
	Code            string     `json:"code" form:"code" gorm:"size:64;column:code;comment:编码"`
	Sex             string     `json:"sex" form:"sex" gorm:"size:8;default:'男';column:sex;comment:性别"`
	Idcard          string     `json:"idcard" form:"idcard" gorm:"size:32;column:idcard;comment:身份证号"`
	Mobile          string     `json:"mobile" form:"mobile" gorm:"size:32;column:mobile;comment:手机号"`
	Email           string     `json:"email" form:"email" gorm:"size:128;column:email;comment:电子邮箱"`
	Avatar          string     `json:"avatar" form:"avatar" gorm:"size:128;column:avatar;comment:头像"`
	OrgID           string     `json:"org_id" form:"org_id" gorm:"type:varchar(64);column:org_id;comment:所属组织ID"`
	Status          int        `json:"status" form:"status" gorm:"not null;default:1;column:status;comment:状态"`
	EndTime         *time.Time `json:"end_time" form:"end_time" gorm:";column:end_time;comment:账号截止日期"`
	Kind            int        `json:"kind" form:"kind" gorm:"default:2;column:kind;comment:账号类型(1:管理员, 2:普通用户)"`
	Remark          string     `json:"remark" form:"remark" gorm:"size:512;column:remark;comment:备注"`
	Sort            int        `json:"sort" form:"sort" gorm:"default:0;column:sort;comment:序号"`
	FirstLogin      int        `json:"first_login" form:"first_login" gorm:"default:0;column:first_login;comment:是否首次登录"`
	LastLoginTime   *time.Time `json:"last_login_time" form:"last_login_time" gorm:";column:last_login_time;comment:最后登录时间"`
	LastIP          string     `json:"last_ip" form:"last_ip" gorm:"size:32;column:last_ip;comment:最后登录IP"`
	LoginErrorCount int        `json:"login_error_count" form:"login_error_count" gorm:"default:0;column:login_error_count;comment:登录错误次数"`
	IsDeleted       bool       `json:"is_deleted" form:"is_deleted" gorm:";default:0;column:is_deleted;comment:是否删除"`
	IsSystem        bool       `json:"is_system" gorm:";default:0;comment:是否系统内置"`
	TenantID        string     `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'';column:tenant_id;comment:租户ID"`
	CreateID        string     `json:"create_id" form:"create_id" gorm:"size:64;default:'';column:create_id;comment:创建人ID"`
	CreateBy        string     `json:"create_by" form:"create_by" gorm:"size:64;column:create_by;comment:创建人"`
	CreateAt        time.Time  `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID        string     `json:"update_id" form:"update_id" gorm:"size:64;default:'';column:update_id;comment:更新人ID"`
	UpdateBy        string     `json:"update_by" form:"update_by" gorm:"size:64;column:update_by;comment:更新人"`
	UpdateAt        time.Time  `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`

	// 关联数据 (不存入 sso_user 表)
	Roles         []SsoRole    `json:"roles" gorm:"many2many:sso_user_role;joinForeignKey:UserID;joinReferences:RoleID"`
	Positions     []SsoPos     `json:"positions" gorm:"many2many:sso_user_pos;joinForeignKey:UserID;joinReferences:PosID"`
	Organizations []SsoOrgUser `json:"organizations" gorm:"many2many:sso_org_user;joinForeignKey:UserID;joinReferences:OrgID"`
}

// TableName 指定表名
func (SsoUser) TableName() string {
	return "sso_user"
}
