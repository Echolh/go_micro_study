package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 配置
type Config struct {
	Env    string       `mapstructure:"env"`
	Server ServerConfig `mapstructure:"server"`
	Auth   AuthConfig   `mapstructure:"auth"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Name         string `mapstructure:"name"`
	Port         string `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_out"`
}

// 授权 配置
type AuthConfig struct {
	SecretToken string `mapstructure:"secret_token"`
}

func Load() (*Config, error) {
	var cfg Config

	// 配置viper
	rootDir, _ := os.Getwd()
	configPath := filepath.Join(rootDir, "etc", "http_svc_dev.yaml")
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// TODO:监控并重新读取配置文件 -WatchConfig()

	// 读取环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// 解析配置文件到结构体
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
