package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ViperConf struct {
	AutomaticEnv bool
	EnvPrefix    string
	ConfigType   string
	ConfName     string
	ConfigPath   []string
}

func NewViperConf(name, confType, prefix string, autoEnv bool, path ...string) *ViperConf {
	return &ViperConf{
		AutomaticEnv: autoEnv,
		EnvPrefix:    prefix,
		ConfigType:   confType,
		ConfName:     name,
		ConfigPath:   path,
	}
}
func NewViper(c *ViperConf) *viper.Viper {
	//设置viper
	conf := viper.New()
	conf.SetEnvPrefix(c.EnvPrefix)
	conf.AutomaticEnv()
	conf.SetConfigType(c.ConfigType)
	conf.SetConfigName(c.ConfName)
	for _, path := range c.ConfigPath {
		conf.AddConfigPath(path)
	}
	if err := conf.ReadInConfig(); err != nil {
		panic(fmt.Errorf("加载日志文件失败: %w", err))
	}
	return conf
}
