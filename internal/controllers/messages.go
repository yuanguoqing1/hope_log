package controllers

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"hope_blog/config"
	"hope_blog/models"
	"hope_blog/pkg/email"
	"hope_blog/pkg/logger"
)

type MessageRequest struct {
	Message string `json:"message" binding:"required"`
	Email   string `json:"email"`    // 可选，用于接收回复通知
	ReplyTo uint   `json:"reply_to"` // 回复的留言ID
}

type Messages struct{}

// CheckLoginStatus 检查用户登录状态
func (m Messages) CheckLoginStatus(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	username := session.Get("username")

	if userID == nil {
		JsonStruct_bad{}.ReturnError(c, 401, "未登录", nil)
		return
	}

	JsonStruct{}.ReturnSuccess(c, 200, "已登录", gin.H{
		"user_id":   userID,
		"username":  username,
		"logged_in": true,
	}, int64(userID.(uint)))
}

// 协程池，限制并发协程数量
var emailWorkerPool = make(chan struct{}, 10) // 最多10个并发邮件发送

// CreateMessage 接收前端留言并记录到数据库
func (m Messages) CreateMessage(c *gin.Context) {
	// 检查用户是否登录
	session := sessions.Default(c)
	userID := session.Get("user_id")
	username := session.Get("username")
	userEmail := session.Get("email")

	if userID == nil {
		JsonStruct_bad{}.ReturnError(c, 401, "请先登录后再留言", nil)
		return
	}

	var req MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JsonStruct_bad{}.ReturnError(c, 400, "参数错误", err.Error())
		return
	}

	// 构建留言数据（现在确保用户已登录）
	message := models.Message{
		Content:  req.Message,
		IP:       c.ClientIP(),
		Email:    req.Email,
		Status:   1, // 默认直接通过，可以根据需要改为0需要审核
		UserID:   userID.(uint),
		Username: username.(string),
	}

	// 如果是回复留言，设置ReplyTo
	if req.ReplyTo > 0 {
		message.ReplyTo = &req.ReplyTo
	}

	// 如果前端没有提供邮箱，使用用户登录的邮箱
	if message.Email == "" && userEmail != nil {
		message.Email = userEmail.(string)
	}

	// 保存到数据库
	if err := config.GetDB().Create(&message).Error; err != nil {
		logger.Error("保存留言失败", "error", err)
		JsonStruct_bad{}.ReturnError(c, 500, "保存留言失败", err.Error())
		return
	}

	// 写入日志
	logger.Info("收到用户留言",
		"message_id", message.ID,
		"content", message.Content,
		"user_id", message.UserID,
		"username", message.Username,
		"ip", message.IP,
	)

	// 使用协程异步发送邮件通知
	go m.sendNotificationAsync(message)

	// 如果是回复留言，异步通知原留言者
	if message.ReplyTo != nil && *message.ReplyTo > 0 {
		go m.notifyOriginalAuthor(message)
	}

	// 立即返回成功响应，不等待邮件发送
	JsonStruct{}.ReturnSuccess(c, 200, "留言已提交", gin.H{
		"message_id": message.ID,
		"status":     message.Status,
	}, int64(message.ID))
}

// sendNotificationAsync 异步发送邮件通知给管理员
func (m Messages) sendNotificationAsync(message models.Message) {
	// 获取协程池令牌
	emailWorkerPool <- struct{}{}
	defer func() {
		<-emailWorkerPool // 释放令牌
		// 捕获可能的panic
		if r := recover(); r != nil {
			logger.Error("发送邮件通知时发生panic", "error", r)
		}
	}()

	// 记录开始时间
	startTime := time.Now()

	// 创建邮件服务
	emailService := email.NewEmailService()

	// 发送邮件
	err := emailService.SendMessageNotification(
		message.Content,
		message.Username,
		message.Email,
		message.IP,
	)

	// 记录发送结果
	if err != nil {
		logger.Error("发送邮件通知失败",
			"message_id", message.ID,
			"error", err,
			"duration", time.Since(startTime),
		)
	} else {
		// 更新通知时间
		now := time.Now()
		config.GetDB().Model(&message).Update("notify_time", &now)

		logger.Info("邮件通知发送成功",
			"message_id", message.ID,
			"duration", time.Since(startTime),
		)
	}
}

// notifyOriginalAuthor 通知原留言的作者
func (m Messages) notifyOriginalAuthor(replyMessage models.Message) {
	emailWorkerPool <- struct{}{}
	defer func() {
		<-emailWorkerPool
		if r := recover(); r != nil {
			logger.Error("发送回复通知时发生panic", "error", r)
		}
	}()

	// 查找原留言
	var originalMessage models.Message
	if err := config.GetDB().First(&originalMessage, *replyMessage.ReplyTo).Error; err != nil {
		logger.Error("查找原留言失败", "reply_to", *replyMessage.ReplyTo, "error", err)
		return
	}

	// 如果原留言者没有邮箱，则无法通知
	if originalMessage.Email == "" {
		return
	}

	// 发送回复通知
	emailService := email.NewEmailService()
	err := emailService.SendReplyNotification(
		originalMessage.Email,
		originalMessage.Content,
		replyMessage.Content,
	)

	if err != nil {
		logger.Error("发送回复通知失败", "error", err)
	} else {
		logger.Info("回复通知发送成功", "to", originalMessage.Email)
	}
}

// GetMessages 获取留言列表
func (m Messages) GetMessages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var messages []models.Message
	var total int64

	// 只查询已通过的留言
	db := config.GetDB().Model(&models.Message{}).Where("status = 1")

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		JsonStruct_bad{}.ReturnError(c, 500, "获取留言总数失败", err.Error())
		return
	}

	// 获取列表，包含用户信息
	if err := db.Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&messages).Error; err != nil {
		JsonStruct_bad{}.ReturnError(c, 500, "获取留言列表失败", err.Error())
		return
	}

	JsonStruct{}.ReturnSuccess(c, 200, "获取成功", gin.H{
		"messages": messages,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}, total)
}

// BatchProcessMessages 批量处理留言（演示使用协程池处理批量任务）
func (m Messages) BatchProcessMessages(c *gin.Context) {
	var messageIDs []uint
	if err := c.ShouldBindJSON(&messageIDs); err != nil {
		JsonStruct_bad{}.ReturnError(c, 400, "参数错误", err.Error())
		return
	}

	// 使用 WaitGroup 等待所有协程完成
	var wg sync.WaitGroup

	// 创建结果通道
	results := make(chan string, len(messageIDs))

	// 为每个留言启动一个协程处理
	for _, id := range messageIDs {
		wg.Add(1)
		go func(msgID uint) {
			defer wg.Done()

			// 模拟处理留言的操作
			time.Sleep(100 * time.Millisecond)

			// 这里可以执行实际的处理逻辑，比如：
			// - 批量审核
			// - 批量发送通知
			// - 批量标记已读等

			results <- fmt.Sprintf("留言 %d 处理完成", msgID)
		}(id)
	}

	// 启动一个协程等待所有处理完成并关闭通道
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集所有结果
	var processedResults []string
	for result := range results {
		processedResults = append(processedResults, result)
	}

	JsonStruct{}.ReturnSuccess(c, 200, "批量处理完成", gin.H{
		"processed": len(processedResults),
		"results":   processedResults,
	}, int64(len(processedResults)))
}
