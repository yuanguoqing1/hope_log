package controllers

import (
	"context"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"hope_blog/config"
)

type Metrics struct{}

// IncrementVisit 使用 Redis 计数器为网站访问量自增并返回当前计数
func (m Metrics) IncrementVisit(c *gin.Context) {
	ctx := context.Background()
	if config.RedisClient == nil {
		JsonStruct_bad{}.ReturnError(c, 500, "Redis 未初始化", nil)
		return
	}
	// 仅登录用户参与统计；每个用户仅统计一次
	session := sessions.Default(c)
	userID := session.Get("user_id")

	// 获取当前计数（用于未登录或已统计过时返回）
	currentCount, _ := config.RedisClient.Get(ctx, "site:visit_count").Int64()

	if userID == nil {
		JsonStruct{}.ReturnSuccess(c, 200, "未登录，不计入统计", gin.H{"count": currentCount}, currentCount)
		return
	}

	// 使用集合记录已统计过的用户
	uniqueSetKey := "site:visit_users"
	// SAdd 返回 1 表示新加入，0 表示已存在
	added, err := config.RedisClient.SAdd(ctx, uniqueSetKey, userID).Result()
	if err != nil {
		JsonStruct_bad{}.ReturnError(c, 500, "写入唯一访客集合失败", err.Error())
		return
	}

	if added > 0 {
		// 首次统计该用户，计数 +1
		newCount, err := config.RedisClient.Incr(ctx, "site:visit_count").Result()
		if err != nil {
			JsonStruct_bad{}.ReturnError(c, 500, "访问计数失败", err.Error())
			return
		}
		JsonStruct{}.ReturnSuccess(c, 200, "计数+1", gin.H{"count": newCount}, newCount)
		return
	}

	// 已统计过该用户，不再累加
	JsonStruct{}.ReturnSuccess(c, 200, "已统计过，不重复累加", gin.H{"count": currentCount}, currentCount)
}
