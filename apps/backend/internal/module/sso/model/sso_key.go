package model

import "time"

// KeyType 密钥类型
type KeyType string

const (
	KeyTypeRSA     KeyType = "rsa"     // RSA密钥对
	KeyTypeEC      KeyType = "ec"      // ECDSA密钥对
	KeyTypeOctet   KeyType = "octet"   // 对称密钥
	KeyTypeCertificate KeyType = "cert" // 证书
)

// KeyUsage 密钥用途
type KeyUsage string

const (
	KeyUsageSigning    KeyUsage = "signing"    // 签名
	KeyUsageEncryption KeyUsage = "encryption" // 加密
	KeyUsageBoth       KeyUsage = "both"       // 两者都用
)

// SsoKey SSO密钥管理模型
type SsoKey struct {
	ID          string    `json:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	KeyName     string    `json:"key_name" gorm:"size:128;not null;column:key_name;comment:密钥名称"`
	KeyID       string    `json:"key_id" gorm:"size:128;uniqueIndex;not null;column:key_id;comment:密钥ID(KID)"`
	KeyType     KeyType   `json:"key_type" gorm:"size:32;not null;column:key_type;comment:密钥类型"`
	KeyUsage    KeyUsage  `json:"key_usage" gorm:"size:32;not null;column:key_usage;comment:密钥用途"`
	Algorithm   string    `json:"algorithm" gorm:"size:64;not null;column:algorithm;comment:算法(RS256, ES256等)"`
	PublicKey   string    `json:"public_key" gorm:"type:text;column:public_key;comment:公钥(PEM格式)"`
	PrivateKey  string    `json:"private_key" gorm:"type:text;column:private_key;comment:私钥(加密存储,PEM格式)"`
	Certificate string    `json:"certificate" gorm:"type:text;column:certificate;comment:证书(PEM格式)"`
	SecretKey   string    `json:"secret_key" gorm:"size:512;column:secret_key;comment:对称密钥(加密存储)"`
	IsPrimary   bool      `json:"is_primary" gorm:"default:false;column:is_primary;comment:是否为主密钥"`
	IsEnabled   bool      `json:"is_enabled" gorm:"default:true;column:is_enabled;comment:是否启用"`
	ValidFrom   time.Time `json:"valid_from" gorm:"column:valid_from;comment:有效期开始"`
	ValidTo     time.Time `json:"valid_to" gorm:"column:valid_to;comment:有效期结束"`
	ProtocolConfigID string `json:"protocol_config_id" gorm:"type:varchar(64);index;column:protocol_config_id;comment:关联的协议配置ID"`
	Remark      string    `json:"remark" gorm:"size:512;default:'';column:remark;comment:备注"`
	IsDeleted   bool      `json:"is_deleted" gorm:"default:false;column:is_deleted;comment:是否删除"`
	TenantID    string    `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'1';column:tenant_id;comment:租户ID"`
	CreateID    string    `json:"create_id" gorm:"size:64;default:'';column:create_id;comment:创建人ID"`
	CreateBy    string    `json:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt    time.Time `json:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID    string    `json:"update_id" gorm:"size:64;default:'';column:update_id;comment:更新人ID"`
	UpdateBy    string    `json:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt    time.Time `json:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoKey) TableName() string {
	return "sso_key"
}

// IsValid 检查密钥是否在有效期内
func (k *SsoKey) IsValid() bool {
	if !k.IsEnabled {
		return false
	}
	now := time.Now()
	if !k.ValidFrom.IsZero() && now.Before(k.ValidFrom) {
		return false
	}
	if !k.ValidTo.IsZero() && now.After(k.ValidTo) {
		return false
	}
	return true
}
