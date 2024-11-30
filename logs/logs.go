package logs

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"sync"
	"time"
)

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
	var res ZapLogger
	var level zapcore.Level
	var encoder zapcore.Encoder
	var encoderConfig zapcore.EncoderConfig

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

	encoderConfig = zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoder = zapcore.NewJSONEncoder(encoderConfig)
	if viper.GetString("app.log.encode") == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	if viper.GetString("app.env") == "dev" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
		core := zapcore.NewCore(
			encoder, // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(w)), // 打印到控制台和文件
			level, // 日志级别
		)
		res.Logger = zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
		InitHelper(&res)
		return &res
	}
	core := zapcore.NewCore(
		encoder,            // 编码器配置
		zapcore.AddSync(w), // 仅打印到文件
		level,              // 日志级别
	)
	res.Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	InitHelper(&res)
	return &res
}

// 自定义时间编码器
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.DateTime))
}
