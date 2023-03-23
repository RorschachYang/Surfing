package util

import (
	"fmt"

	"github.com/spf13/viper" // 导入viper包，用于读取配置文件
)

var configDir string

// GetConfigString函数从config.yaml配置文件中获取指定键的字符串值
func GetConfigString(key string) string {
	// 设置配置文件路径
	viper.SetConfigFile(configDir)
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		// 如果读取配置文件失败，抛出致命错误并终止程序
		panic(fmt.Errorf("fatal error reading config file: %s", err))
	}

	// 返回指定键的字符串值
	return viper.GetString(key)
}

// GetConfigInt函数从config.yaml配置文件中获取指定键的整数值
func GetConfigInt(key string) int64 {
	// 设置配置文件路径
	viper.SetConfigFile(configDir)
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		// 如果读取配置文件失败，抛出致命错误并终止程序
		panic(fmt.Errorf("fatal error reading config file: %s", err))
	}

	// 返回指定键的整数值
	return viper.GetInt64(key)
}

func SetConfigDir(dir string) {
	configDir = dir
}
