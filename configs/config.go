package configs

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort    string
	DatabaseDSN   string
	JWTSecret     string
	JWTExpireDays int
	Email         EmailConfig
	Code          CodeConfig
}

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string // 使用授权码而非密码
	From     string // 发件人格式： "昵称 <邮箱>"
	Secure   bool   // 是否使用SSL/TLS
}

type CodeConfig struct {
	Length       int           // 验证码长度
	ExpireTime   time.Duration // 验证码有效期
	SendInterval time.Duration // 发送间隔限制
	MaxAttempts  int           // 最大尝试次数

}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {

	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		DatabaseDSN:   getEnv("DATABASE_DSN", "test.db"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpireDays: getEnvAsInt("JWT_EXPIRE_DAYS", 7),
		Email: EmailConfig{
			SMTPHost: getEnv("EMAIL_SMTP_HOST", "smtp.qq.com"),
			SMTPPort: getEnvAsInt("EMAIL_SMTP_PORT", 587),
			Username: getEnv("EMAIL_USERNAME", "your-email@qq.com"),
			Password: getEnv("EMAIL_PASSWORD", "your-email-password"),
			From:     getEnv("EMAIL_FROM", "Your Name <your-email@qq.com>"),
			Secure:   true,
		},
		Code: CodeConfig{
			Length:       6,
			ExpireTime:   5 * time.Minute,
			SendInterval: 1 * time.Minute,
			MaxAttempts:  5,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
