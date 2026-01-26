package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var SugarLogger *zap.SugaredLogger

// InitLogger 初始化日志系统
func InitLogger(logLevel, logFilePath string) {
	// 设置日志级别
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	// 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 控制台输出配置
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 文件输出配置
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100, // 每个日志文件100MB
		MaxBackups: 10,  // 保留10个备份
		MaxAge:     30,  // 保留30天
		Compress:   true, // 压缩旧日志
	})

	// 创建核心
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriter, level),
		zapcore.NewCore(fileEncoder, fileWriter, level),
	)

	// 创建日志记录器
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	SugarLogger = Logger.Sugar()
}

// SyncLogger 同步日志，确保所有日志都写入文件
func SyncLogger() {
	_ = Logger.Sync()
	_ = SugarLogger.Sync()
}

// WithTraceID 添加链路追踪ID到日志上下文
func WithTraceID(traceID string) *zap.SugaredLogger {
	return SugarLogger.With("trace_id", traceID)
}

// WithContext 添加上下文信息到日志
func WithContext(fields map[string]any) *zap.SugaredLogger {
	return SugarLogger.With(fields)
}
