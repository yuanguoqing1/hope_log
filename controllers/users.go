package controllers

import "github.com/gin-gonic/gin"

type Users struct{}

func (u Users) GetUserDone(c *gin.Context) {
	uid := c.Param("uid")
	name := c.Param("name")

	JsonStruct{}.ReturnSuccess(c, 200, uid, name, 12)
}

func (u Users) GetUserNone(c *gin.Context) {
	uid := c.PostForm("uid")
	name := c.PostForm("name")
	JsonStruct_bad{}.ReturnError(c, 200, name, uid)
}
