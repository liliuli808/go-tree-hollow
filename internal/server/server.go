package server

import (
	"context"
	"errors"
	"fmt"
	"go-tree-hollow/configs"
	"go-tree-hollow/internal/middleware"
	"go-tree-hollow/internal/modules/auth"
	"go-tree-hollow/internal/modules/user"
	"go-tree-hollow/internal/repository/redis"
	"go-tree-hollow/pkg/database"
	"go-tree-hollow/pkg/email"
	"go-tree-hollow/pkg/utils"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	config *configs.Config
	db     *gorm.DB
	router *gin.Engine
}

func NewServer(config *configs.Config) (*Server, error) {
	// 初始化数据库
	db, err := database.NewDB(config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	// 初始化Gin
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	server := &Server{
		config: config,
		db:     db,
		router: router,
	}

	// 注册路由
	server.setupRoutes()

	return server, nil
}

func (s *Server) setupRoutes() {
	// API v1路由组
	v1 := s.router.Group("/api/v1")

	// 认证模块
	authRepo := auth.NewRepository(s.db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)
	auth.RegisterRoutes(v1, authHandler)

	// 用户模块（需要认证）
	userRepo := user.NewRepository(s.db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	user.RegisterRoutes(v1, userHandler)

	// 健康检查
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})
}

func (s *Server) Start() error {
	addr := ":" + s.config.ServerPort
	log.Printf("Server starting on %s", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号
	quit := make(chan struct{})
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
	return nil
}

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
	sender   *email.Sender
	codeRepo *redis.CodeRepository
	cfg      *configs.Config
}

func NewEmailService(sender *email.Sender, codeRepo *redis.CodeRepository, cfg *configs.Config) EmailService {
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
