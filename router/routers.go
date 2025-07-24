package router

import (
	"hope_blog/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	user := router.Group("/user")
	{
		user.GET("/info", controllers.GetUserInfo)
		user.GET("/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, "已经登入了")
		})

		user.GET("/registerr", func(c *gin.Context) {
			c.JSON(http.StatusOK, "注册")
		})
	}
	return router
}
