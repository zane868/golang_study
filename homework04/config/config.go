package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServiceConfig
}

type ServiceConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}

func newViper() *viper.Viper {

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
		log.Fatalf("读取配置失败: %v", err)
	}

	return v

}

func Load() *Config {

	v := newViper()
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}
	return &cfg

}
