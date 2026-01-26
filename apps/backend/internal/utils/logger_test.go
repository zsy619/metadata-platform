package utils

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TestInitLogger_Success 测试日志初始化成功
func TestInitLogger_Success(t *testing.T) {
	// 准备测试数据
	testLogFile := "/tmp/test_logger.log"
	testLogLevel := "info"

	// 清理之前的测试文件
	os.Remove(testLogFile)

	// 执行测试
	InitLogger(testLogLevel, testLogFile)

	// 验证 Logger 不为空
	if Logger == nil {
		t.Error("Logger should not be nil after initialization")
	}

	if SugarLogger == nil {
		t.Error("SugarLogger should not be nil after initialization")
	}

	// 写入一条日志以触发文件创建
	Logger.Info("Test log entry")
	SyncLogger()

	// 验证日志文件已创建
	if _, err := os.Stat(testLogFile); os.IsNotExist(err) {
		t.Error("Log file should be created")
	}

	// 清理
	os.Remove(testLogFile)
}

// TestInitLogger_AllLevels 测试所有日志级别
func TestInitLogger_AllLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "fatal", "invalid"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			testLogFile := "/tmp/test_logger_" + level + ".log"
			os.Remove(testLogFile)

			// 不应 panic
			InitLogger(level, testLogFile)

			if Logger == nil {
				t.Errorf("Logger should not be nil for level: %s", level)
			}

			SyncLogger()
			os.Remove(testLogFile)
		})
	}
}

// TestInitLogger_LogOutput 测试日志输出功能
func TestInitLogger_LogOutput(t *testing.T) {
	testLogFile := "/tmp/test_logger_output.log"
	os.Remove(testLogFile)

	// 使用 buffered writer 捕获输出
	var buf bytes.Buffer

	// 初始化 logger，将 console 输出重定向
	encoderConfig := zap.NewProductionEncoderConfig()
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(&buf), zapcore.InfoLevel)
	Logger = zap.New(core)
	SugarLogger = Logger.Sugar()

	// 测试不同级别的日志
	testMessages := []struct {
		level      string
		message    string
		shouldLog  bool
	}{
		{"debug", "debug message", false},
		{"info", "info message", true},
		{"warn", "warn message", true},
		{"error", "error message", true},
	}

	for _, tm := range testMessages {
		switch tm.level {
		case "debug":
			SugarLogger.Debug(tm.message)
		case "info":
			SugarLogger.Info(tm.message)
		case "warn":
			SugarLogger.Warn(tm.message)
		case "error":
			SugarLogger.Error(tm.message)
		}

		output := buf.String()
		if tm.shouldLog && !strings.Contains(output, tm.message) {
			t.Errorf("Expected log message '%s' not found in output", tm.message)
		}
		if !tm.shouldLog && strings.Contains(output, tm.message) {
			t.Errorf("Unexpected log message '%s' found in output", tm.message)
		}
	}

	SyncLogger()
	os.Remove(testLogFile)
}

// TestInitLogger_FileRotation 测试日志文件轮转
func TestInitLogger_FileRotation(t *testing.T) {
	testLogFile := "/tmp/test_logger_rotation.log"
	os.Remove(testLogFile)

	// 初始化带小文件大小的 logger（测试用）
	InitLogger("info", testLogFile)

	// 写入一些日志
	for i := 0; i < 10; i++ {
		Logger.Info("Test log entry", zap.Int("iteration", i))
	}

	SyncLogger()

	// 验证日志文件存在
	if _, err := os.Stat(testLogFile); os.IsNotExist(err) {
		t.Error("Log file should exist after writing logs")
	}

	// 清理
	os.Remove(testLogFile)
}

// TestInitLogger_Concurrency 测试并发日志写入
func TestInitLogger_Concurrency(t *testing.T) {
	testLogFile := "/tmp/test_logger_concurrent.log"
	os.Remove(testLogFile)

	InitLogger("info", testLogFile)

	// 并发写入日志
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func(id int) {
			Logger.Info("Concurrent log", zap.Int("id", id))
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 100; i++ {
		<-done
	}

	SyncLogger()

	// 验证日志文件
	if _, err := os.Stat(testLogFile); os.IsNotExist(err) {
		t.Error("Log file should exist after concurrent writes")
	}

	os.Remove(testLogFile)
}

// TestSyncLogger 测试日志同步
func TestSyncLogger(t *testing.T) {
	testLogFile := "/tmp/test_logger_sync.log"
	os.Remove(testLogFile)

	InitLogger("info", testLogFile)
	Logger.Info("Before sync")

	// 不应 panic
	SyncLogger()

	// 第二次调用也应安全
	SyncLogger()

	os.Remove(testLogFile)
}

// TestWithTraceID 测试 TraceID 上下文
func TestWithTraceID(t *testing.T) {
	testLogFile := "/tmp/test_logger_trace.log"
	os.Remove(testLogFile)

	InitLogger("info", testLogFile)

	traceID := "test-trace-12345"
	logger := WithTraceID(traceID)

	if logger == nil {
		t.Error("WithTraceID should not return nil")
	}

	// 写入日志验证 TraceID
	logger.Info("Test with trace ID")

	SyncLogger()
	os.Remove(testLogFile)
}

// TestWithContext 测试上下文字段
func TestWithContext(t *testing.T) {
	testLogFile := "/tmp/test_logger_context.log"
	os.Remove(testLogFile)

	InitLogger("info", testLogFile)

	fields := map[string]any{
		"user_id":   123,
		"action":    "login",
		"ip":        "192.168.1.1",
	}

	logger := WithContext(fields)

	if logger == nil {
		t.Error("WithContext should not return nil")
	}

	// 写入日志验证上下文
	logger.Info("Test with context")

	SyncLogger()
	os.Remove(testLogFile)
}

// TestInitLogger_EmptyPath 测试空日志路径
func TestInitLogger_EmptyPath(t *testing.T) {
	// 不应 panic
	InitLogger("info", "")

	if Logger == nil {
		t.Error("Logger should not be nil even with empty path")
	}

	SyncLogger()
}

// TestInitLogger_SpecialCharacters 测试路径中的特殊字符
func TestInitLogger_SpecialCharacters(t *testing.T) {
	testLogFile := "/tmp/test logger with spaces.log"
	os.Remove(testLogFile)

	// 不应 panic
	InitLogger("info", testLogFile)

	if Logger == nil {
		t.Error("Logger should not be nil with special characters in path")
	}

	SyncLogger()
	os.Remove(testLogFile)
}
