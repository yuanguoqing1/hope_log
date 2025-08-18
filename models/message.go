package models

import (
	"time"

	"gorm.io/gorm"
)

// Message 留言模型
type Message struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Content    string     `gorm:"type:text;not null" json:"content"` // 留言内容
	UserID     uint       `gorm:"not null" json:"user_id"`           // 留言者ID（必须登录）
	Username   string     `gorm:"size:50" json:"username"`           // 留言者用户名
	Email      string     `gorm:"size:100" json:"email"`             // 留言者邮箱
	IP         string     `gorm:"size:50" json:"ip"`                 // 留言者IP
	Status     int        `gorm:"default:0" json:"status"`           // 状态：0-待审核，1-已通过，2-已拒绝
	IsRead     bool       `gorm:"default:false" json:"is_read"`      // 管理员是否已读
	ReplyTo    *uint      `json:"reply_to"`                          // 回复的留言ID（null表示不是回复）
	NotifyTime *time.Time `json:"notify_time"`                       // 邮件通知发送时间

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}
