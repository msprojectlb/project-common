package config

import (
	"time"
)

type ViperConf struct {
	AutomaticEnv bool
	EnvPrefix    string
	ConfigType   string
	ConfName     string
	ConfigPath   []string
}

func NewViperConf(automaticEnv bool, envPrefix string, configType string, confName string, configPath []string) *ViperConf {
	return &ViperConf{
		AutomaticEnv: automaticEnv,
		EnvPrefix:    envPrefix,
		ConfigType:   configType,
		ConfName:     confName,
		ConfigPath:   configPath,
	}
}

type JWTConfig struct {
	AccessExp     time.Duration
	RefreshExp    time.Duration
	AccessSecret  string
	RefreshSecret string
}
