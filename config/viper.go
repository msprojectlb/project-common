package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func NewViper(c *ViperConf) *viper.Viper {
	//设置viper
	conf := viper.New()
	conf.SetEnvPrefix(c.EnvPrefix)
	if c.AutomaticEnv {
		conf.AutomaticEnv()
	}
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
