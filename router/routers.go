package router

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	//初始化
	router := gin.Default()

	//加载HTML模板
	router.LoadHTMLGlob("static/html/*")
	//用户路由组
	user := router.Group("/user")
	{
		//登录页面
		user.GET("/login", func(c *gin.Context) {
			c.HTML(200, "login.html", gin.H{})
		})
		// 忘记密码界面
		user.GET("/forgot-password", func(c *gin.Context) {
			c.HTML(200, "forgot-password.html", gin.H{})
		})
		//注册页面
		user.GET("/register", func(c *gin.Context) {
			c.HTML(200, "register.html", gin.H{})
		})

		//登录API
		user.POST("/user/login", func(c *gin.Context) {
			// TODO: 实现登录逻辑
			c.JSON(200, gin.H{
				"success": true,
				"message": "登录成功",
			})
		})

		//注册API
		user.POST("/api/register", func(c *gin.Context) {
			// TODO: 实现注册逻辑
			c.JSON(200, gin.H{
				"success": true,
				"message": "注册成功",
			})
		})
	}
	// blog页面
	blog := router.Group("/")
	{
		blog.GET("", func(c *gin.Context) {
			c.HTML(200, "index.html", gin.H{})
		})
	}
	return router
}
