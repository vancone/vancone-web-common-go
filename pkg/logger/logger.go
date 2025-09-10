package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var SugaredLogger *zap.SugaredLogger

func InitLogger() {
	// 配置日志写入：控制台 + 文件（带轮转）
	fileWriter := &lumberjack.Logger{
		Filename:   "logs/app.log", // 日志文件路径
		MaxSize:    10,             // 每个日志文件最大 10MB
		MaxBackups: 5,              // 最多保留 5 个备份
		MaxAge:     30,             // 文件最多保存 30 天
		LocalTime:  true,
		Compress:   false, // 是否压缩（gzip）
	}

	// 设置日志级别（生产环境用 InfoLevel，开发可用 DebugLevel）
	level := zapcore.InfoLevel

	// 编码配置：JSON 格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 可读时间格式
	encoderConfig.StacktraceKey = ""                      // 在访问日志中不记录 stack trace

	// 3. 创建多个写入目标：文件 + 控制台
	//    zapcore.AddSync 是必要的，因为 lumberjack 不是并发安全的
	ws := []zapcore.WriteSyncer{
		zapcore.AddSync(fileWriter), // 写入文件
		zapcore.AddSync(os.Stdout),  // 写入控制台
	}

	// 使用 zapcore.NewTee 将日志“分发”到多个目标
	multiWriteSyncer := zapcore.NewMultiWriteSyncer(ws...)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		multiWriteSyncer,
		level,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	// 全局设置
	SugaredLogger = logger.Sugar()
	defer logger.Sync() // 确保程序退出时刷新日志
}

// --- 简化函数 ---
func Debug(args ...interface{}) {
	SugaredLogger.Debug(args...)
}

func Info(args ...interface{}) {
	SugaredLogger.Info(args...)
}

func Warn(args ...interface{}) {
	SugaredLogger.Warn(args...)
}

func Error(args ...interface{}) {
	SugaredLogger.Error(args...)
}

func DPanic(args ...interface{}) {
	SugaredLogger.DPanic(args...)
}

func Panic(args ...interface{}) {
	SugaredLogger.Panic(args...)
}

func Fatal(args ...interface{}) {
	SugaredLogger.Fatal(args...)
}

// 支持格式化输出：Infof, Errorf 等
func Debugf(template string, args ...interface{}) {
	SugaredLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	SugaredLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	SugaredLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	SugaredLogger.Errorf(template, args...)
}

func DPanicf(template string, args ...interface{}) {
	SugaredLogger.DPanicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	SugaredLogger.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	SugaredLogger.Fatalf(template, args...)
}
