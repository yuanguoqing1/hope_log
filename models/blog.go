package models

import (
	"time"

	"gorm.io/gorm"
)

// Blog 博客文章模型
type Blog struct {
	gorm.Model
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Summary   string    `gorm:"type:varchar(500)" json:"summary"`
	AuthorID  uint      `gorm:"not null" json:"author_id"`
	Author    User      `gorm:"foreignKey:AuthorID" json:"author"`
	Status    int       `gorm:"type:tinyint(1);comment:'1:已发布 0:草稿'" json:"status"`
	ViewCount int       `gorm:"default:0" json:"view_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Blog) TableName() string {
	return "blogs"
}

// IncrementViewCount 增加阅读量
func (b *Blog) IncrementViewCount(db *gorm.DB) error {
	return db.Model(b).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// 
