package main

import (
	"hope_blog/pkg/logger"
	"hope_blog/router"
)

func main() {
	//初始化logger
	logger.InitLogger()
	defer logger.Close()
	//初始化数据库

	//初始化路由
	route := router.Router()
	if route != nil {
		logger.Write("初始化路由成功")
	}
	route.Run(":9999")

}
