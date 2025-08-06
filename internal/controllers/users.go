package controllers

import (
	"github.com/gin-gonic/gin"

	"hope_blog/config"
	"hope_blog/models"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type Users struct{}

func checkUser(username string) (bool, error) {
	var count int64
	if err := config.GetDB().Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

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

// 注册函数
func (u Users) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JsonStruct_bad{}.ReturnError(c, 400, "参数错误", err.Error())
		return
	}

	// 检查用户名是否已存在
	exists, err := checkUser(req.Username)
	if err != nil {
		JsonStruct_bad{}.ReturnError(c, 500, "数据库错误", err.Error())
		return
	}
	if exists {
		JsonStruct_bad{}.ReturnError(c, 400, "用户名已存在", nil)
		return
	}

	// 创建新用户
	user := models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	// 保存到数据库
	if err := config.GetDB().Create(&user).Error; err != nil {
		JsonStruct_bad{}.ReturnError(c, 500, "注册失败_写入数据库出错", err.Error())
		return
	}

	JsonStruct{}.ReturnSuccess(c, 200, "注册成功", nil, 0)
}

// Login 登录函数
func (u Users) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JsonStruct_bad{}.ReturnError(c, 400, "参数错误", err.Error())
		return
	}
	exits, err := checkUser(req.Username)
	if err != nil {
		JsonStruct_bad{}.ReturnError(c, 500, "数据库错误", err.Error())
		return
	}
	if !exits {
		JsonStruct_bad{}.ReturnError(c, 400, "用户名不存在,请注册", nil)
		return
	}
	err = models.User.ValidatePassword(req.Password)
}
