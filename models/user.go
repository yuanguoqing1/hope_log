package models

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model            // 内嵌 gorm.Model, 包含 ID、CreatedAt、UpdatedAt、DeletedAt
	Username   string     `gorm:"type:varchar(32);uniqueIndex;not null" json:"username"`
	Password   string     `gorm:"type:varchar(128);not null" json:"-"`
	Email      string     `gorm:"type:varchar(128);uniqueIndex;not null" json:"email"`
	Avatar     string     `gorm:"type:varchar(256)" json:"avatar"`
	Status     int        `gorm:"type:tinyint(1);default:1;comment:'1:正常 0:禁用'" json:"status"`
	LastLogin  *time.Time `gorm:"comment:'最后登录时间'" json:"last_login"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeSave GORM 钩子 - 保存前加密密码
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		// 检查密码是否已经被加密过（bcrypt hash 通常以 $2a$, $2b$, $2y$ 开头）
		if len(u.Password) < 60 || (!strings.HasPrefix(u.Password, "$2a$") && 
			!strings.HasPrefix(u.Password, "$2b$") && 
			!strings.HasPrefix(u.Password, "$2y$")) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			u.Password = string(hashedPassword)
		}
	}
	return nil
}

// ValidatePassword 验证密码
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ChangePassword 修改密码
func (u *User) ChangePassword(db *gorm.DB, newPassword string) error {
	u.Password = newPassword
	return db.Save(u).Error
}

// UpdateLastLogin 更新最后登录时间
func (u *User) UpdateLastLogin(db *gorm.DB) error {
	now := time.Now()
	u.LastLogin = &now
	return db.Model(u).UpdateColumn("last_login", now).Error
}
