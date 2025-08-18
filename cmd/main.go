package main

import (
	"hope_blog/config"
	"hope_blog/pkg/logger"
	"hope_blog/router"
)

func main() {
	//初始化logger
	if err := logger.InitLogger(); err != nil {
		logger.Error("初始化Logger失败", "err", err)
		return
	}
	defer logger.Close()
	//初始化配置文件
	if err := config.LoadConfig(); err != nil {
		logger.Error("加载配置文件失败", "err", err)
		return
	}
	//初始化数据库
	if err := config.InitDB(); err != nil {
		logger.Error("初始化数据库失败", "err", err)
		return
	}
	//初始化redis
	if err := config.InitRedis(); err != nil {
		logger.Error("初始化Redis失败", "err", err)
		return
	}
	//初始化路由
	route := router.InitRouter()

	route.Run(":" + config.Global.Server.Port)

}
