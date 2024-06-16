package logs

import (
	"github.com/gin-gonic/gin"
	"github.com/msprojectlb/project-common/config"
	"github.com/msprojectlb/project-common/logs/writer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// ZapLogger 通用日志
type ZapLogger struct {
	*zap.Logger
}

// DBLogger 业务数据库日志
type DBLogger ZapLogger

// GormLogger Gorm日志
type GormLogger ZapLogger

// HttpLogger gin日志
type HttpLogger ZapLogger

func NewZapLogger(conf config.ZapLogConf, w *writer.ZapWriter) *ZapLogger {
	var level zapcore.Level
	//debug<info<warn<error<fatal<panic
	switch conf.Level {
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
	if conf.Encoding == "console" {
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
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	}
	var core zapcore.Core
	core = zapcore.NewCore(
		encoder,                                  // 编码器配置
		zapcore.AddSync((*lumberjack.Logger)(w)), // 仅打印到文件
		level,                                    // 日志级别
	)
	if conf.Environment == "debug" {
		core = zapcore.NewCore(
			encoder, // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync((*lumberjack.Logger)(w))), // 打印到控制台和文件
			level, // 日志级别
		)
		//开发环境
		return &ZapLogger{zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
	}
	//正式环境
	return &ZapLogger{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
}
func NewDBLogger(conf config.DBLogConf, w *writer.DBWriter) *DBLogger {
	return (*DBLogger)(NewZapLogger(config.ZapLogConf(conf), (*writer.ZapWriter)(w)))
}
func NewGormLogger(conf config.GormLogConf, w *writer.GormWriter) *GormLogger {
	return (*GormLogger)(NewZapLogger(config.ZapLogConf(conf), (*writer.ZapWriter)(w)))
}
func NewHttpLogger(conf config.HttpLogConf, w *writer.HttpWriter) *HttpLogger {
	return (*HttpLogger)(NewZapLogger(config.ZapLogConf(conf), (*writer.ZapWriter)(w)))
}

// 自定义时间编码器
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//enc.AppendString(t.Format("2006-01-02 15:04:05"))
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// GinLogger 接收gin框架默认的日志
func GinLogger(log *HttpLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		log.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool, log *HttpLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					log.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					log.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					log.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
