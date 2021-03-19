package log

import (
	"github.com/IanVzs/lightAPI/flag_parse"
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

func init() {
	level := getLoggerLevel(*flag_parse.LogLevel)
	hook := lumberjack.Logger{
		Filename:   *flag_parse.LogPath,
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
