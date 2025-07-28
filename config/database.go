package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB
func InitDB() {
	var err error
	dsn := "root:password@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("连接数据库失败: %v", err))
	}

	// 自动迁移数据库表结构
	// DB.AutoMigrate(&models.User{}) // 需要导入 models 包
}
