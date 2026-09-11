package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServiceConfig
	Jwt    JwtConfig
}

type ServiceConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}

type JwtConfig struct {
	Secret string `mapstructure:"secret"`
	Expire string `mapstructure:"expire"`
}

func newViper() (*viper.Viper, error) {

	v := viper.New()
	// 设置配置文件名称（不含扩展名）
	v.SetConfigName("config")
	// 设置配置文件类型
	v.SetConfigType("yaml")
	// 添加配置文件搜索路径
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/.app")
	// 读取环境变量
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")

	// 设置默认值
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.mode", "debug")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}

	return v, nil

}

func Load() (*Config, error) {

	v, err := newViper()
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}
	return &cfg, nil

}
