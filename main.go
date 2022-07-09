package main

import (
	"fmt"
	"go-redis/config"
	"go-redis/lib/logger"
	"go-redis/tcp"
	"os"
)

const configFile = "redis.conf"

// 默认配置
var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	// 设置日志初始化参数
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "go-redis",
		Ext:        "log",
		TimeFormat: "2022-01-02",
	})

	// 读取配置文件
	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}

	err := tcp.ListenAndServeWithSignal(
		&tcp.Config{
			Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
		}, tcp.MakeHandler())

	if err != nil {
		logger.Error(err)
	}

}
