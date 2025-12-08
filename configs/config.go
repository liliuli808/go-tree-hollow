package configs

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	DatabaseDSN   string
	JWTSecret     string
	JWTExpireDays int
	Email         EmailConfig
	Code          CodeConfig
	Redis         RedisConfig
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

type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`           // 地址: "localhost:6379"
	Password     string        `mapstructure:"password"`       // 密码（默认为空）
	DB           int           `mapstructure:"db"`             // 数据库编号
	PoolSize     int           `mapstructure:"pool_size"`      // 连接池大小
	MinIdleConns int           `mapstructure:"min_idle_conns"` // 最小空闲连接
	MaxRetries   int           `mapstructure:"max_retries"`    // 最大重试次数
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`   // 连接超时
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`   // 读取超时
	WriteTimeout time.Duration `mapstructure:"write_timeout"`  // 写入超时
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，将直接使用系统环境变量")
		// 注意：这里不 return，继续执行
	}
	godotenv.Load(".env")
	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8081"),
		DatabaseDSN:   getEnv("DATABASE_DSN", "test.db"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpireDays: getEnvAsInt("JWT_EXPIRE_DAYS", 7),
		Email: EmailConfig{
			SMTPHost: getEnv("EMAIL_SMTP_HOST", "smtp.gmail.com"),
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
		Redis: RedisConfig{
			Addr:         getEnv("REDIS_ADDR", "localhost:6379"),
			Password:     getEnv("REDIS_PASSWORD", "secret"),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 0),
			MaxRetries:   getEnvAsInt("REDIS_MAX_RETRIES", 3),
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
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
