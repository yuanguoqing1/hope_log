package config

import (
	"context"
	"fmt"
	"os"

	"hope_blog/models"
	"hope_blog/pkg/logger"

	"github.com/redis/go-redis/v9"
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

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}

// AIConfig 大模型配置
type AIConfig struct {
	API_KEY  string `yaml:"api_key"`
	BASE_URL string `yaml:"base_url"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	Host        string `yaml:"host"`         // SMTP服务器地址
	Port        string `yaml:"port"`         // SMTP端口
	Username    string `yaml:"username"`     // 邮箱账号
	Password    string `yaml:"password"`     // 邮箱密码/授权码
	From        string `yaml:"from"`         // 发件人地址
	AdminEmails string `yaml:"admin_emails"` // 管理员邮箱，多个用逗号分隔
}

// Config 配置结构体
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DataBaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	AI       AIConfig       `yaml:"ai"`
	Email    EmailConfig    `yaml:"email"`
}

var (
	DB          *gorm.DB
	Global      Config
	RedisClient *redis.Client
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
	if err := DB.AutoMigrate(&models.User{}, &models.Blog{}, &models.Message{}); err != nil {
		logger.Error("数据库迁移失败", "error", err)
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}

func InitRedis() error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", Global.Redis.Host, Global.Redis.Port),
		Password: Global.Redis.Password,
		DB:       0,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("连接Redis失败", "error", err)
		return fmt.Errorf("连接Redis失败: %v", err)
	}

	RedisClient = redisClient
	// 初始化访问计数键，若不存在则设为0
	if _, err := RedisClient.SetNX(context.Background(), "site:visit_count", 0, 0).Result(); err != nil {
		logger.Error("初始化Redis键失败", "key", "site:visit_count", "error", err)
		return fmt.Errorf("初始化Redis键失败: %v", err)
	}
	return nil
}
