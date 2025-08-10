package controllers

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"hope_blog/config"
	"hope_blog/models"
	"hope_blog/pkg/logger"
)

type BlogRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Summary string `json:"summary"`
	Status  int    `json:"status"` // 1:发布 0:草稿
}

type Blogs struct{}

// GetBlogs 获取博客列表
func (b Blogs) GetBlogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var blogs []models.Blog
	var total int64

	// 读取筛选参数：status (0/1)，own=1 仅本人草稿
	statusParam := c.DefaultQuery("status", "1")
	ownParam := c.DefaultQuery("own", "0")

	db := config.GetDB().Model(&models.Blog{})
	listDB := config.GetDB().Preload("Author")

	if statusParam == "0" {
		// 草稿只允许本人查看
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if ownParam == "1" && userID != nil {
			db = db.Where("status = 0 AND author_id = ?", userID)
			listDB = listDB.Where("status = 0 AND author_id = ?", userID)
		} else {
			// 非本人/未登录不返回草稿
			db = db.Where("1=0")
			listDB = listDB.Where("1=0")
		}
	} else {
		// 已发布
		db = db.Where("status = 1")
		listDB = listDB.Where("status = 1")
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		JsonStruct_bad{}.ReturnError(c, 500, "获取博客总数失败", err.Error())
		return
	}

	// 列表
	if err := listDB.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&blogs).Error; err != nil {
		JsonStruct_bad{}.ReturnError(c, 500, "获取博客列表失败", err.Error())
		return
	}

	JsonStruct{}.ReturnSuccess(c, 200, "获取成功", gin.H{
		"blogs": blogs,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}, total)
}

// GetBlog 获取单个博客详情
func (b Blogs) GetBlog(c *gin.Context) {
	id := c.Param("id")

	var blog models.Blog
	// 允许作者查看自己的草稿；非作者只能查看已发布
	session := sessions.Default(c)
	currentUserID := session.Get("user_id")

	db := config.GetDB().Preload("Author")
	if currentUserID != nil {
		// 查找：若是作者，允许 status=0；否则仅 status=1
		if err := db.Where("id = ? AND (status = 1 OR (status = 0 AND author_id = ?))", id, currentUserID).First(&blog).Error; err != nil {
			JsonStruct_bad{}.ReturnError(c, 404, "博客不存在", nil)
			return
		}
	} else {
		if err := db.Where("id = ? AND status = 1", id).First(&blog).Error; err != nil {
			JsonStruct_bad{}.ReturnError(c, 404, "博客不存在", nil)
			return
		}
	}

	// 仅已发布文章增加阅读量
	if blog.Status == 1 {
		blog.IncrementViewCount(config.GetDB())
	}

	JsonStruct{}.ReturnSuccess(c, 200, "获取成功", blog, int64(blog.ID))
}

// CreateBlog 创建博客
func (b Blogs) CreateBlog(c *gin.Context) {
	// 检查用户是否登录
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {

		JsonStruct_bad{}.ReturnError(c, 401, "请先登录", nil)
		return
	}

	var req BlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JsonStruct_bad{}.ReturnError(c, 400, "参数错误", err.Error())
		return
	}
	logger.Info("CreateBlog: received payload", "status", req.Status, "title", req.Title)

	// 如果没有提供摘要，从内容中截取前200字符作为摘要
	if req.Summary == "" && len(req.Content) > 200 {
		req.Summary = req.Content[:200] + "..."
	} else if req.Summary == "" {
		req.Summary = req.Content
	}

	blog := models.Blog{
		Title:    req.Title,
		Content:  req.Content,
		Summary:  req.Summary,
		AuthorID: userID.(uint),
		Status:   req.Status,
	}
	logger.Info("CreateBlog: before create", "status", blog.Status)

	if err := config.GetDB().Select("Title", "Content", "Summary", "AuthorID", "Status").Create(&blog).Error; err != nil {
		JsonStruct_bad{}.ReturnError(c, 500, "创建博客失败", err.Error())
		return
	}
	logger.Info("CreateBlog: after create", "id", blog.ID, "status", blog.Status)

	// 重新查询获取完整信息（包括作者）
	config.GetDB().Preload("Author").First(&blog, blog.ID)
	logger.Info("CreateBlog: after reload", "id", blog.ID, "status", blog.Status)

	JsonStruct{}.ReturnSuccess(c, 200, "创建成功", blog, int64(blog.ID))
}

// UpdateBlog 更新博客
func (b Blogs) UpdateBlog(c *gin.Context) {
	id := c.Param("id")

	// 检查用户是否登录
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		JsonStruct_bad{}.ReturnError(c, 401, "请先登录", nil)
		return
	}

	var blog models.Blog
	if err := config.GetDB().Where("id = ?", id).First(&blog).Error; err != nil {
		JsonStruct_bad{}.ReturnError(c, 404, "博客不存在", nil)
		return
	}

	// 检查是否是作者本人
	if blog.AuthorID != userID.(uint) {
		JsonStruct_bad{}.ReturnError(c, 403, "无权限修改", nil)
		return
	}

	var req BlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JsonStruct_bad{}.ReturnError(c, 400, "参数错误", err.Error())
		return
	}

	// 更新博客信息
	blog.Title = req.Title
	blog.Content = req.Content
	blog.Summary = req.Summary
	blog.Status = req.Status

	if err := config.GetDB().Save(&blog).Error; err != nil {
		JsonStruct_bad{}.ReturnError(c, 500, "更新博客失败", err.Error())
		return
	}

	// 重新查询获取完整信息
	config.GetDB().Preload("Author").First(&blog, blog.ID)

	JsonStruct{}.ReturnSuccess(c, 200, "更新成功", blog, int64(blog.ID))
}

// DeleteBlog 删除博客
func (b Blogs) DeleteBlog(c *gin.Context) {
	id := c.Param("id")

	// 检查用户是否登录
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		JsonStruct_bad{}.ReturnError(c, 401, "请先登录", nil)
		return
	}

	var blog models.Blog
	if err := config.GetDB().Where("id = ?", id).First(&blog).Error; err != nil {
		JsonStruct_bad{}.ReturnError(c, 404, "博客不存在", nil)
		return
	}

	// 检查是否是作者本人
	if blog.AuthorID != userID.(uint) {
		JsonStruct_bad{}.ReturnError(c, 403, "无权限删除", nil)
		return
	}

	if err := config.GetDB().Delete(&blog).Error; err != nil {
		JsonStruct_bad{}.ReturnError(c, 500, "删除博客失败", err.Error())
		return
	}

	JsonStruct{}.ReturnSuccess(c, 200, "删除成功", nil, 0)
}
