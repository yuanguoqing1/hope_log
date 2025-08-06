package router

import (
	"github.com/gin-gonic/gin"
	"hope_blog/internal/controllers"
)

func InitRouter() *gin.Engine {
	//初始化
	router := gin.Default()

	//加载HTML模板 - 修正模板路径
	router.LoadHTMLGlob("../static/html/*")
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
		user.POST("/api/login", controllers.Users{}.Login)

		//注册API
		user.POST("/api/register", controllers.Users{}.Register)
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
