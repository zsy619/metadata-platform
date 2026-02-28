package model

import "time"

// SsoFieldMapping SSO字段映射配置模型
type SsoFieldMapping struct {
	ID               string    `json:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	MappingName      string    `json:"mapping_name" gorm:"size:128;not null;column:mapping_name;comment:映射配置名称"`
	ProtocolConfigID string    `json:"protocol_config_id" gorm:"type:varchar(64);index;column:protocol_config_id;comment:关联的协议配置ID"`
	ClientID         string    `json:"client_id" gorm:"type:varchar(128);index;column:client_id;comment:关联的客户端ID"`
	SourceField      string    `json:"source_field" gorm:"size:128;not null;column:source_field;comment:源字段(第三方)"`
	TargetField      string    `json:"target_field" gorm:"size:128;not null;column:target_field;comment:目标字段(本地)"`
	FieldType        string    `json:"field_type" gorm:"size:32;default:'string';column:field_type;comment:字段类型(string, int, bool, array)"`
	IsRequired       bool      `json:"is_required" gorm:"default:false;column:is_required;comment:是否必需"`
	DefaultValue     string    `json:"default_value" gorm:"size:512;default:'';column:default_value;comment:默认值"`
	TransformScript  string    `json:"transform_script" gorm:"type:text;column:transform_script;comment:转换脚本"`
	IsEnabled        bool      `json:"is_enabled" gorm:"default:true;column:is_enabled;comment:是否启用"`
	Sort             int       `json:"sort" gorm:"default:0;column:sort;comment:排序"`
	Remark           string    `json:"remark" gorm:"size:512;default:'';column:remark;comment:备注"`
	IsDeleted        bool      `json:"is_deleted" gorm:"default:false;column:is_deleted;comment:是否删除"`
	TenantID         string    `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'1';column:tenant_id;comment:租户ID"`
	CreateID         string    `json:"create_id" gorm:"size:64;default:'';column:create_id;comment:创建人ID"`
	CreateBy         string    `json:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt         time.Time `json:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID         string    `json:"update_id" gorm:"size:64;default:'';column:update_id;comment:更新人ID"`
	UpdateBy         string    `json:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt         time.Time `json:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoFieldMapping) TableName() string {
	return "sso_field_mapping"
}

// UserFieldMapping 用户字段映射常量
const (
	// 本地系统用户字段
	UserFieldID           = "id"
	UserFieldUsername     = "username"
	UserFieldEmail        = "email"
	UserFieldMobile       = "mobile"
	UserFieldNickname     = "nickname"
	UserFieldAvatar       = "avatar"
	UserFieldRealName     = "real_name"
	UserFieldFirstName    = "first_name"
	UserFieldLastName     = "last_name"
	UserFieldGender       = "gender"
	UserFieldBirthday     = "birthday"
	UserFieldDepartment   = "department"
	UserFieldPosition     = "position"
	UserFieldEmployeeNo   = "employee_no"
	UserFieldStatus       = "status"
)

// GetDefaultUserFieldMappings 获取默认的用户字段映射配置
func GetDefaultUserFieldMappings() []map[string]string {
	return []map[string]string{
		{"source": "sub", "target": UserFieldID, "type": "string", "required": "true"},
		{"source": "preferred_username", "target": UserFieldUsername, "type": "string", "required": "true"},
		{"source": "email", "target": UserFieldEmail, "type": "string", "required": "false"},
		{"source": "phone_number", "target": UserFieldMobile, "type": "string", "required": "false"},
		{"source": "name", "target": UserFieldNickname, "type": "string", "required": "false"},
		{"source": "picture", "target": UserFieldAvatar, "type": "string", "required": "false"},
		{"source": "given_name", "target": UserFieldFirstName, "type": "string", "required": "false"},
		{"source": "family_name", "target": UserFieldLastName, "type": "string", "required": "false"},
		{"source": "gender", "target": UserFieldGender, "type": "string", "required": "false"},
		{"source": "birthdate", "target": UserFieldBirthday, "type": "string", "required": "false"},
	}
}
