package email

import (
	"fmt"
	"hope_blog/config"
	"hope_blog/pkg/logger"
	"net/smtp"
	"strings"
	"time"
)

// EmailService 邮件服务
type EmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

// NewEmailService 创建邮件服务
func NewEmailService() *EmailService {
	return &EmailService{
		host:     config.Global.Email.Host,
		port:     config.Global.Email.Port,
		username: config.Global.Email.Username,
		password: config.Global.Email.Password,
		from:     config.Global.Email.From,
	}
}

// SendEmail 发送邮件
func (s *EmailService) SendEmail(to []string, subject, body string) error {
	// 如果邮件服务未配置，则跳过
	if s.host == "" || s.username == "" {
		logger.Warn("邮件服务未配置，跳过发送邮件")
		return nil
	}

	// 构建邮件内容
	message := fmt.Sprintf("From: %s\r\n", s.from)
	message += fmt.Sprintf("To: %s\r\n", strings.Join(to, ","))
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n"
	message += body

	// 连接SMTP服务器
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	err := smtp.SendMail(addr, auth, s.from, to, []byte(message))
	if err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	logger.Info("邮件发送成功", "to", strings.Join(to, ","), "subject", subject)
	return nil
}

// SendMessageNotification 发送留言通知邮件给管理员
func (s *EmailService) SendMessageNotification(messageContent, username, userEmail, userIP string) error {
	adminEmails := strings.Split(config.Global.Email.AdminEmails, ",")
	if len(adminEmails) == 0 {
		return nil
	}

	subject := "博客收到新留言"
	body := fmt.Sprintf(`
		<h3>您的博客收到新留言</h3>
		<p><strong>留言者：</strong>%s</p>
		<p><strong>邮箱：</strong>%s</p>
		<p><strong>IP地址：</strong>%s</p>
		<p><strong>时间：</strong>%s</p>
		<hr>
		<p><strong>留言内容：</strong></p>
		<p>%s</p>
		<hr>
		<p>请登录后台管理系统查看和处理留言。</p>
	`, username, userEmail, userIP, time.Now().Format("2006-01-02 15:04:05"), messageContent)

	return s.SendEmail(adminEmails, subject, body)
}

// SendReplyNotification 发送回复通知邮件给留言者
func (s *EmailService) SendReplyNotification(toEmail, originalMessage, replyContent string) error {
	if toEmail == "" {
		return nil
	}

	subject := "您的留言收到回复"
	body := fmt.Sprintf(`
		<h3>您在博客的留言收到了回复</h3>
		<p><strong>您的留言：</strong></p>
		<p>%s</p>
		<hr>
		<p><strong>回复内容：</strong></p>
		<p>%s</p>
		<hr>
		<p>感谢您的留言，欢迎继续访问我的博客！</p>
	`, originalMessage, replyContent)

	return s.SendEmail([]string{toEmail}, subject, body)
}
