package logs

import (
	"github.com/msprojectlb/project-common/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapWriter struct {
	*lumberjack.Logger
}

func NewZapWriter(conf config.ZapWriterConf) *ZapWriter {
	hook := &lumberjack.Logger{
		Filename:   conf.FileName, // 日志文件路径
		MaxSize:    conf.Maxsize,  // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: conf.MaxBack,  // 日志文件最多保存多少个备份
		MaxAge:     conf.MaxAge,   // 文件最多保存多少天
		Compress:   conf.Compress, // 是否压缩
	}
	return &ZapWriter{hook}
}
