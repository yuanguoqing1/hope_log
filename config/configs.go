package config

import (
	"fmt"
	"os"

	"hope_blog/models"
	"hope_blog/pkg/logger"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DataBaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

// Config 配置结构体
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DataBaseConfig `yaml:"database"`
}

var (
	DB     *gorm.DB
	Global Config
)

// LoadConfig 加载配置文件
func LoadConfig() error {
	// 读取配置文件 - 修正路径
	data, err := os.ReadFile("../config/config.yaml")
	if err != nil {
		logger.Error("读取配置文件失败", "error", err)
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析配置文件
	if err := yaml.Unmarshal(data, &Global); err != nil {
		logger.Error("解析配置文件失败", "error", err)
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	logger.Info("配置加载成功")
	return nil
}

// InitDB 初始化数据库连接
func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		Global.Database.User,
		Global.Database.Password,
		Global.Database.Host,
		Global.Database.Port,
		Global.Database.Database,
		Global.Database.Charset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("连接数据库失败", "error", err)
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 自动迁移模型
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		logger.Error("数据库迁移失败", "error", err)
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	logger.Info("数据库连接成功")
	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}
