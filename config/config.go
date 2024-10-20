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

type JWTConfig struct {
	AccessExp     time.Duration
	RefreshExp    time.Duration
	AccessSecret  string
	RefreshSecret string
}
