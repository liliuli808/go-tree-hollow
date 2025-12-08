package email

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"go-tree-hollow/configs"
	"go-tree-hollow/pkg/utils"
	"math/rand"
	"time"

	"gopkg.in/gomail.v2"
)

var (
	ErrInvalidEmail    = errors.New("邮箱格式不正确")
	ErrCodeSendTooFast = errors.New("发送过于频繁，请稍后再试")
	ErrCodeExpired     = errors.New("验证码已过期")
	ErrCodeInvalid     = errors.New("验证码错误")
)

type EmailService interface {
	SendVerificationCode(ctx context.Context, email string) error
	VerifyCode(ctx context.Context, email, code string) error
}

type emailServiceImpl struct {
	sender   *Sender
	codeRepo *CodeRepository
	cfg      *configs.Config
}

func NewEmailService(sender *Sender, codeRepo *CodeRepository, cfg *configs.Config) EmailService {
	return &emailServiceImpl{
		sender:   sender,
		codeRepo: codeRepo,
		cfg:      cfg,
	}
}

// SendVerificationCode 发送验证码
func (s *emailServiceImpl) SendVerificationCode(ctx context.Context, email string) error {
	// 1. 验证邮箱格式
	if !utils.ValidateEmail(email) {
		return ErrInvalidEmail
	}

	// 2. 检查发送频率（1分钟内只能发送一次）
	lockKey := fmt.Sprintf("lock:%s", email)
	if exists, _ := s.codeRepo.SetNX(ctx, email, "1", s.cfg.Code.SendInterval); !exists {
		return ErrCodeSendTooFast
	}

	// 3. 生成验证码
	code := utils.GenerateCode(s.cfg.Code.Length)

	// 4. 发送邮件
	if err := s.sender.SendVerificationCode(email, code); err != nil {
		// 发送失败，清除频率限制
		s.codeRepo.Delete(ctx, lockKey)
		return fmt.Errorf("邮件发送失败: %w", err)
	}

	// 5. 存储验证码到Redis
	if err := s.codeRepo.Set(ctx, email, code, s.cfg.Code.ExpireTime); err != nil {
		return fmt.Errorf("存储验证码失败: %w", err)
	}

	return nil
}

// VerifyCode 验证验证码
func (s *emailServiceImpl) VerifyCode(ctx context.Context, email, code string) error {
	// 1. 获取存储的验证码
	storedCode, err := s.codeRepo.Get(ctx, email)
	if err != nil {
		return ErrCodeExpired
	}

	// 2. 比对验证码
	if storedCode != code {
		return ErrCodeInvalid
	}

	// 3. 验证成功，删除验证码（防止重复使用）
	s.codeRepo.Delete(ctx, email)

	return nil
}

// GenerateCode 生成随机验证码
func GenerateCode(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	const digits = "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}

type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// SendCodeResponse 响应结构
type SendCodeResponse struct {
	Message string `json:"message"`
}

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
	fmt.Print(config.Password)
	dialer := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.Username, config.Password)

	// SSL/TLS 配置
	if config.Secure {
		dialer.TLSConfig = &tls.Config{
			ServerName:         config.SMTPHost,
			InsecureSkipVerify: false}
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
