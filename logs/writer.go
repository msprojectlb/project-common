package logs

import (
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

type ZapWriter struct {
	lumberjack.Logger
}

func NewZapWriter(viper *viper.Viper) io.Writer {
	return &ZapWriter{
		lumberjack.Logger{
			Filename:  viper.GetString("app.log.file"),   // 日志文件路径
			MaxSize:   viper.GetInt("app.log.maxSize"),   // 每个日志文件保存的最大尺寸 单位：M
			MaxAge:    viper.GetInt("app.log.maxAge"),    // 文件最多保存多少天
			Compress:  viper.GetBool("app.log.compress"), // 是否压缩
			LocalTime: true,
		},
	}
}
