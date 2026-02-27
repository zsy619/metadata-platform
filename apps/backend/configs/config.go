package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

// DBConfig 数据库配置结构体
type DBConfig struct {
	Type         string `mapstructure:"TYPE"`         // 数据库类型: mysql, postgres
	Host         string `mapstructure:"HOST"`
	Port         int    `mapstructure:"PORT"`
	User         string `mapstructure:"USER"`
	Password     string `mapstructure:"PASSWORD"`
	Name         string `mapstructure:"NAME"`
	MaxIdleConns int    `mapstructure:"MAX_IDLE_CONNS"`
	MaxOpenConns int    `mapstructure:"MAX_OPEN_CONNS"`
	SSLMode      string `mapstructure:"SSL_MODE"`     // SSL模式: disable, require, verify-ca, verify-full (PostgreSQL)
}

// Config 应用配置结构体
type Config struct {
	// 应用基本配置
	AppName    string `mapstructure:"APP_NAME"`
	AppMode    string `mapstructure:"APP_MODE"`
	ServerHost string `mapstructure:"SERVER_HOST"`
	ServerPort int    `mapstructure:"SERVER_PORT"`

	// 元数据数据库配置
	MetadataDB DBConfig `mapstructure:"METADATA_DB"`
	// 用户管理数据库配置
	UserDB DBConfig `mapstructure:"USER_DB"`
	// 审计日志数据库配置
	AuditDB DBConfig `mapstructure:"AUDIT_DB"`

	// JWT配置
	JWTSecret      string `mapstructure:"JWT_SECRET"`
	JWTExpireHours int    `mapstructure:"JWT_EXPIRE_HOURS"`

	// 日志配置
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	LogFilePath string `mapstructure:"LOG_FILE_PATH"`
}

// LoadConfig 从环境变量或配置文件加载配置
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	setDefaults()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	config.MetadataDB.Type = getString("METADATA_DB.TYPE", config.MetadataDB.Type)
	config.MetadataDB.Host = getString("METADATA_DB.HOST", config.MetadataDB.Host)
	config.MetadataDB.Port = getInt("METADATA_DB.PORT", config.MetadataDB.Port)
	config.MetadataDB.User = getString("METADATA_DB.USER", config.MetadataDB.User)
	config.MetadataDB.Password = getString("METADATA_DB.PASSWORD", config.MetadataDB.Password)
	config.MetadataDB.Name = getString("METADATA_DB.NAME", config.MetadataDB.Name)
	config.MetadataDB.MaxIdleConns = getInt("METADATA_DB.MAX_IDLE_CONNS", config.MetadataDB.MaxIdleConns)
	config.MetadataDB.MaxOpenConns = getInt("METADATA_DB.MAX_OPEN_CONNS", config.MetadataDB.MaxOpenConns)
	config.MetadataDB.SSLMode = getString("METADATA_DB.SSL_MODE", config.MetadataDB.SSLMode)

	config.UserDB.Type = getString("USER_DB.TYPE", config.UserDB.Type)
	config.UserDB.Host = getString("USER_DB.HOST", config.UserDB.Host)
	config.UserDB.Port = getInt("USER_DB.PORT", config.UserDB.Port)
	config.UserDB.User = getString("USER_DB.USER", config.UserDB.User)
	config.UserDB.Password = getString("USER_DB.PASSWORD", config.UserDB.Password)
	config.UserDB.Name = getString("USER_DB.NAME", config.UserDB.Name)
	config.UserDB.MaxIdleConns = getInt("USER_DB.MAX_IDLE_CONNS", config.UserDB.MaxIdleConns)
	config.UserDB.MaxOpenConns = getInt("USER_DB.MAX_OPEN_CONNS", config.UserDB.MaxOpenConns)
	config.UserDB.SSLMode = getString("USER_DB.SSL_MODE", config.UserDB.SSLMode)

	config.AuditDB.Type = getString("AUDIT_DB.TYPE", config.AuditDB.Type)
	config.AuditDB.Host = getString("AUDIT_DB.HOST", config.AuditDB.Host)
	config.AuditDB.Port = getInt("AUDIT_DB.PORT", config.AuditDB.Port)
	config.AuditDB.User = getString("AUDIT_DB.USER", config.AuditDB.User)
	config.AuditDB.Password = getString("AUDIT_DB.PASSWORD", config.AuditDB.Password)
	config.AuditDB.Name = getString("AUDIT_DB.NAME", config.AuditDB.Name)
	config.AuditDB.MaxIdleConns = getInt("AUDIT_DB.MAX_IDLE_CONNS", config.AuditDB.MaxIdleConns)
	config.AuditDB.MaxOpenConns = getInt("AUDIT_DB.MAX_OPEN_CONNS", config.AuditDB.MaxOpenConns)
	config.AuditDB.SSLMode = getString("AUDIT_DB.SSL_MODE", config.AuditDB.SSLMode)

	return &config, nil
}

func getString(key string, defaultValue string) string {
	if val := viper.GetString(key); val != "" {
		return val
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if val := viper.GetInt(key); val != 0 {
		return val
	}
	return defaultValue
}

// 设置默认配置值
func setDefaults() {
	// 应用基本配置
	viper.SetDefault("APP_NAME", "metadata-platform")
	viper.SetDefault("APP_MODE", "debug")
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", 8080)

	// 元数据数据库配置
	viper.SetDefault("METADATA_DB.TYPE", "mysql")
	viper.SetDefault("METADATA_DB.HOST", "localhost")
	viper.SetDefault("METADATA_DB.PORT", 3306)
	viper.SetDefault("METADATA_DB.USER", "root")
	viper.SetDefault("METADATA_DB.PASSWORD", "password")
	viper.SetDefault("METADATA_DB.NAME", "metadata_platform")
	viper.SetDefault("METADATA_DB.MAX_IDLE_CONNS", 10)
	viper.SetDefault("METADATA_DB.MAX_OPEN_CONNS", 100)
	viper.SetDefault("METADATA_DB.SSL_MODE", "disable")

	// 用户管理数据库配置
	viper.SetDefault("USER_DB.TYPE", "mysql")
	viper.SetDefault("USER_DB.HOST", "localhost")
	viper.SetDefault("USER_DB.PORT", 3306)
	viper.SetDefault("USER_DB.USER", "root")
	viper.SetDefault("USER_DB.PASSWORD", "password")
	viper.SetDefault("USER_DB.NAME", "metadata_sso")
	viper.SetDefault("USER_DB.MAX_IDLE_CONNS", 10)
	viper.SetDefault("USER_DB.MAX_OPEN_CONNS", 100)
	viper.SetDefault("USER_DB.SSL_MODE", "disable")

	// 审计日志数据库配置
	viper.SetDefault("AUDIT_DB.TYPE", "mysql")
	viper.SetDefault("AUDIT_DB.HOST", "localhost")
	viper.SetDefault("AUDIT_DB.PORT", 3306)
	viper.SetDefault("AUDIT_DB.USER", "root")
	viper.SetDefault("AUDIT_DB.PASSWORD", "password")
	viper.SetDefault("AUDIT_DB.NAME", "metadata_audit")
	viper.SetDefault("AUDIT_DB.MAX_IDLE_CONNS", 10)
	viper.SetDefault("AUDIT_DB.MAX_OPEN_CONNS", 100)
	viper.SetDefault("AUDIT_DB.SSL_MODE", "disable")

	// JWT配置
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("JWT_EXPIRE_HOURS", 24)

	// 日志配置
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FILE_PATH", "logs/app.log")
}
