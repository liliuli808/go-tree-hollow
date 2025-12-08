// pkg/email/sender.go
package email

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

type Sender struct {
	config *EmailConfig
	dialer *gomail.Dialer
}

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
	Secure   bool
}

func NewSender(config *EmailConfig) *Sender {
	dialer := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.Username, config.Password)

	// SSL/TLS 配置
	if config.Secure {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false}
	}

	return &Sender{
		config: config,
		dialer: dialer,
	}
}

// Send 发送邮件
func (s *Sender) Send(to, subject, body string, contentType ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// 默认使用 HTML 格式
	ct := "text/html"
	if len(contentType) > 0 {
		ct = contentType[0]
	}
	m.SetBody(ct, body)

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}
	return nil
}

// SendVerificationCode 发送验证码邮件（使用模板）
func (s *Sender) SendVerificationCode(to, code string) error {
	subject := fmt.Sprintf("验证码 - %s", code)

	// HTML 模板（可提取到单独文件）
	body := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <title>验证码</title>
    </head>
    <body>
        <div style="padding: 20px;">
            <h2>您的验证码</h2>
            <p>验证码：<strong style="color: #1890ff; font-size: 24px;">%s</strong></p>
            <p>有效期：5分钟，请勿泄露给他人</p>
            <p>如非本人操作，请忽略此邮件</p>
        </div>
    </body>
    </html>
    `, code)

	return s.Send(to, subject, body)
}
