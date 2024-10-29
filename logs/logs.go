package logs

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"sync"
	"time"
)

const MsgName = "msg"

var Helper *ZapLogger
var once sync.Once

// ZapLogger 通用日志
type ZapLogger struct {
	*zap.Logger
}

func InitHelper(log *ZapLogger) {
	once.Do(func() {
		Helper = log
	})
}

func NewZapLogger(viper *viper.Viper, w io.Writer) *ZapLogger {
	var level zapcore.Level
	switch viper.GetString("app.log.level") {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	var encoder zapcore.Encoder
	if viper.GetString("app.log.encode") == "console" {
		encoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "Logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		})
	} else {
		encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    "func",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	}
	var core zapcore.Core
	core = zapcore.NewCore(
		encoder,            // 编码器配置
		zapcore.AddSync(w), // 仅打印到文件
		level,              // 日志级别
	)
	var res ZapLogger
	if viper.GetString("app.env") == "dev" {
		core = zapcore.NewCore(
			encoder, // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(w)), // 打印到控制台和文件
			level, // 日志级别
		)
		//开发环境
		res.Logger = zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
		InitHelper(&res)
		return &res
	}
	res.Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	//正式环境
	InitHelper(&res)
	return &res
}

func (l *ZapLogger) LogInfo(msg any) {
	l.Info(MsgName, zap.Any("", msg))
}
func (l *ZapLogger) LogError(msg any) {
	l.Error(MsgName, zap.Any("", msg))
}
func (l *ZapLogger) LogWarn(msg any) {
	l.Warn(MsgName, zap.Any("", msg))
}
func (l *ZapLogger) LogDebug(msg any) {
	l.Debug(MsgName, zap.Any("", msg))
}
func (l *ZapLogger) LogFatal(msg any) {
	l.Fatal(MsgName, zap.Any("", msg))
}
func (l *ZapLogger) LogPanic(msg any) {
	l.Panic(MsgName, zap.Any("", msg))
}
func (l *ZapLogger) LogInfoF(fmtStr string, args ...any) {
	l.Info(MsgName, zap.String("", fmt.Sprintf(fmtStr, args...)))
}
func (l *ZapLogger) LogErrorF(fmtStr string, args ...any) {
	l.Error(MsgName, zap.String("", fmt.Sprintf(fmtStr, args...)))
}
func (l *ZapLogger) LogWarnF(fmtStr string, args ...any) {
	l.Warn(MsgName, zap.String("", fmt.Sprintf(fmtStr, args...)))
}
func (l *ZapLogger) LogDebugF(fmtStr string, args ...any) {
	l.Debug(MsgName, zap.String("", fmt.Sprintf(fmtStr, args...)))
}
func (l *ZapLogger) LogFatalF(fmtStr string, args ...any) {
	l.Fatal(MsgName, zap.String("", fmt.Sprintf(fmtStr, args...)))
}
func (l *ZapLogger) LogPanicF(fmtStr string, args ...any) {
	l.Panic(MsgName, zap.String("", fmt.Sprintf(fmtStr, args...)))
}

// 自定义时间编码器
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.DateTime))
}
