package config

type ZapWriterConf struct {
	FileName string
	Maxsize  int
	MaxBack  int
	MaxAge   int
	Compress bool
}
type GormWriterConf ZapWriterConf
type DBWriterConf ZapWriterConf
type HttpWriterConf ZapWriterConf

type ZapLogConf struct {
	Level       string //日志文件等级
	Encoding    string //日志文件编码方式 console|json
	Environment string //开发/生产环境 debug|prod
}
type GormLogConf ZapLogConf
type DBLogConf ZapLogConf
type HttpLogConf ZapLogConf
