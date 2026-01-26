package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

// DBConfig 数据库配置结构体
type DBConfig struct {
	Host          string `mapstructure:"HOST"`
	Port          int    `mapstructure:"PORT"`
	User          string `mapstructure:"USER"`
	Password      string `mapstructure:"PASSWORD"`
	Name          string `mapstructure:"NAME"`
	MaxIdleConns  int    `mapstructure:"MAX_IDLE_CONNS"`
	MaxOpenConns  int    `mapstructure:"MAX_OPEN_CONNS"`
}

// Config 应用配置结构体
type Config struct {
	// 应用基本配置
	AppName         string `mapstructure:"APP_NAME"`
	AppMode         string `mapstructure:"APP_MODE"`
	ServerHost      string `mapstructure:"SERVER_HOST"`
	ServerPort      int    `mapstructure:"SERVER_PORT"`

	// 元数据数据库配置
	MetadataDB      DBConfig `mapstructure:"METADATA_DB"`
	// 用户管理数据库配置
	UserDB          DBConfig `mapstructure:"USER_DB"`

	// JWT配置
	JWTSecret       string `mapstructure:"JWT_SECRET"`
	JWTExpireHours  int    `mapstructure:"JWT_EXPIRE_HOURS"`

	// 日志配置
	LogLevel        string `mapstructure:"LOG_LEVEL"`
	LogFilePath     string `mapstructure:"LOG_FILE_PATH"`
}

// LoadConfig 从环境变量或配置文件加载配置
func LoadConfig() (*Config, error) {
	// 初始化Viper
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// 配置文件不存在，使用默认值
	}

	// 设置默认值
	setDefaults()

	// 解析配置到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// 设置默认配置值
func setDefaults() {
	// 应用基本配置
	viper.SetDefault("APP_NAME", "metadata-platform")
	viper.SetDefault("APP_MODE", "debug")
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", 8080)

	// 元数据数据库配置
	viper.SetDefault("METADATA_DB.HOST", "localhost")
	viper.SetDefault("METADATA_DB.PORT", 3306)
	viper.SetDefault("METADATA_DB.USER", "root")
	viper.SetDefault("METADATA_DB.PASSWORD", "password")
	viper.SetDefault("METADATA_DB.NAME", "metadata_platform")
	viper.SetDefault("METADATA_DB.MAX_IDLE_CONNS", 10)
	viper.SetDefault("METADATA_DB.MAX_OPEN_CONNS", 100)

	// 用户管理数据库配置
	viper.SetDefault("USER_DB.HOST", "localhost")
	viper.SetDefault("USER_DB.PORT", 3306)
	viper.SetDefault("USER_DB.USER", "root")
	viper.SetDefault("USER_DB.PASSWORD", "password")
	viper.SetDefault("USER_DB.NAME", "metadata_sso")
	viper.SetDefault("USER_DB.MAX_IDLE_CONNS", 10)
	viper.SetDefault("USER_DB.MAX_OPEN_CONNS", 100)

	// JWT配置
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("JWT_EXPIRE_HOURS", 24)

	// 日志配置
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FILE_PATH", "logs/app.log")
}