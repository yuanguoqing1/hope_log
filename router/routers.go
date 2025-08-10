package router

import (
	"hope_blog/internal/controllers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	//初始化
	router := gin.Default()

	// 配置session存储
	store := cookie.NewStore([]byte("hope-blog-secret-key-2025-very-secure-random-string-32chars"))
	store.Options(sessions.Options{
		MaxAge:   3600 * 24, // 24小时
		HttpOnly: false,     // 临时设为false便于调试
		Secure:   false,     // 开发环境设为false
		Path:     "/",       // 设置cookie路径
	})
	router.Use(sessions.Sessions("blog-session", store))

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

		//获取当前用户信息API
		user.GET("/api/current", controllers.Users{}.GetCurrentUser)
		
		//退出登录API
		user.POST("/api/logout", controllers.Users{}.Logout)
	}
	
	// 博客API路由组
	blogAPI := router.Group("/api/blogs")
	{
		// 添加调试中间件
		blogAPI.Use(func(c *gin.Context) {
			session := sessions.Default(c)
			userID := session.Get("user_id")
			username := session.Get("username")

			if userID != nil {
				c.Set("user_id", userID)
			}
			if username != nil {
				c.Set("username", username)
			}

			c.Next()
		})

		blogAPI.GET("", controllers.Blogs{}.GetBlogs)        // 获取博客列表
		blogAPI.GET("/:id", controllers.Blogs{}.GetBlog)     // 获取单个博客
		blogAPI.POST("", controllers.Blogs{}.CreateBlog)     // 创建博客
		blogAPI.PUT("/:id", controllers.Blogs{}.UpdateBlog)  // 更新博客
		blogAPI.DELETE("/:id", controllers.Blogs{}.DeleteBlog) // 删除博客
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
