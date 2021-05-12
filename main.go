package main

import (
	"github.com/LainNetWork/go-boot/application"
	log "github.com/sirupsen/logrus"
	"time"
)

/**
基础设施搭建。
使用yaml作为配置，拥有全局context，支持环境配置
使用gin作为web框架
*/

func main() {

	//读取运行时传参，可覆盖配置文件中的配置
	var config = struct {
		Server struct {
			Port int `yaml:"port"`
		}
	}{}

	var config2 = struct {
		Proxy struct {
			Users []int  `yaml:"users"`
			Host  string `yaml:"host"`
		}
	}{}
	application.Context.RegisterConfigs(&config, &config2)
	application.Context.Init()
	ticker := time.NewTicker(time.Millisecond * 10)
	for range ticker.C {
		log.Warning("测试日志 WARNING")
		log.Info("测试日志 INFO")
		log.Error("测试等级ERROR")
	}
}
