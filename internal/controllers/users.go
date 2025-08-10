package controllers

import (
	"github.com/gin-contrib/sessions"
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
	
	// 获取用户信息
	var user models.User
	if err := config.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		JsonStruct_bad{}.ReturnError(c, 400, "用户名不存在,请注册", nil)
		return
	}
	
	// 验证密码
	if !user.ValidatePassword(req.Password) {
		JsonStruct_bad{}.ReturnError(c, 400, "密码错误", nil)
		return
	}
	
	// 登录成功，更新最后登录时间
	if err := user.UpdateLastLogin(config.GetDB()); err != nil {
		// 记录错误但不影响登录流程
		// 可以在这里添加日志记录
	}
	
	// 保存用户信息到session
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("username", user.Username)
	session.Set("email", user.Email)
	
	
	if err := session.Save(); err != nil {

		JsonStruct_bad{}.ReturnError(c, 500, "Session保存失败", err.Error())
		return
	}
	



	
	// 返回登录成功信息
	JsonStruct{}.ReturnSuccess(c, 200, "登录成功", gin.H{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	}, int64(user.ID))
}

// GetCurrentUser 获取当前登录用户信息
func (u Users) GetCurrentUser(c *gin.Context) {
	
	// 从session中获取用户信息
	session := sessions.Default(c)
	userID := session.Get("user_id")
	username := session.Get("username")
	email := session.Get("email")

	
	
	if userID == nil {
		JsonStruct_bad{}.ReturnError(c, 401, "未登录", nil)
		return
	}


	JsonStruct{}.ReturnSuccess(c, 200, "获取用户信息成功", gin.H{
		"user_id":  userID,
		"username": username,
		"email":    email,
	}, int64(userID.(uint)))
}

// Logout 退出登录
func (u Users) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	
	JsonStruct{}.ReturnSuccess(c, 200, "退出登录成功", nil, 0)
}
