package config

import (
	"errors"
	"os"
	"registry-proxy/pkg/console"
	"strings"

	"github.com/spf13/viper"
)

type server struct {
	Listen   string `toml:"listen"`
	LogLevel string `toml:"logLevel"`
	TLS      *tls   `toml:"tls"`
}

type tls struct {
	Enable         bool   `toml:"enable"`
	Listen         string `toml:"listen"`
	UseLetsEncrypt bool   `toml:"useLetsEncrypt"`
	CertPath       string `toml:"certPath"`
	KeyPath        string `toml:"keyPath"`
}

type proxy struct {
	Binding map[string]string `toml:"binding"`
}

func Load() {
	console.Log().Info("[config] 正在初始化配置")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	defaultConfig := map[string]any{
		"server": Server,
		"proxy":  Proxy,
	}

	for key, val := range defaultConfig {
		viper.SetDefault(key, val)
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			err := viper.SafeWriteConfig()
			if err != nil {
				console.Log().Error("[config] 默认配置文件写出失败: %s", err)
				os.Exit(1)
			}

			console.Log().Info("[config] 已生成默认配置文件，请配置完成后再次运行。")
			os.Exit(0)
		}

		console.Log().Error("[config] 配置文件读取失败: %s", err)
		os.Exit(1)
	}

	for key, val := range defaultConfig {
		err := viper.UnmarshalKey(key, val)
		if err != nil {
			console.Log().Error("[config] 配置文件解析失败, key: %s, error: %s", key, err)
		}
	}

	// 重设log等级
	switch strings.ToUpper(Server.LogLevel) {
	case "DEBUG":
		console.Level = console.LevelDebug
	case "INFO":
		console.Level = console.LevelInformational
	case "WARN":
		console.Level = console.LevelWarning
	case "ERROR":
		console.Level = console.LevelError
	default:
		console.Level = console.LevelInformational
	}

	console.GlobalLogger = nil
	console.Log()
}
