package log

import (
	"flag"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	Logger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	Logger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	Logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	Logger.Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	Logger.Fatalf(template, args...)
}

func DPanic(args ...interface{}) {
	Logger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	Logger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	Logger.Panicf(template, args...)
}

// 全局变量
var Logger *zap.SugaredLogger

// 日志级别
var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"fatal":  zapcore.FatalLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
}

func getLoggerLevel(logLevel string) zapcore.Level {
	if level, ok := levelMap[logLevel]; ok {
		return level
	}
	return zapcore.InfoLevel
}

// main 相关启动配置
var Addr = flag.String("addr", ":8080", "http service address")

// log 相关启动配置
var logPath = flag.String("log_path", "./logs/ligthAPI.log", "set app log path.")
var logLevel = flag.String("log_level", "debug", "set app log level. 默认debug")

// Redis 相关启动配置
var RedisIP = flag.String("redis_ip", "192.168.0.1:6379", "redis address")
var RedisPsw = flag.String("redis_passwd", "", "redis passwd")
var RedisDB = flag.Int("redis_db", 0, "redis db")

// 接口加密密钥
var KeyAPI = flag.String("key_api", "key_test", "接口加/解密密钥")
var KeySplit = flag.String("key_split", "_|_", "接口加/解密数据分割符")

// 放在初始化方法中执行日志的初始化
func init() {
	flag.Parse()
	level := getLoggerLevel(*logLevel)
	hook := lumberjack.Logger{
		Filename:   *logPath,
		MaxSize:    128,
		MaxBackups: 14,
		MaxAge:     1,
		LocalTime:  true,
		Compress:   true,
	}
	syncWriter := zapcore.AddSync(&hook)
	encoderConfig := zap.NewProductionEncoderConfig()
	// 时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		syncWriter,
		zap.NewAtomicLevelAt(level),
	)
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}
