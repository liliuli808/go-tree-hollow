package server

import (
	"context"
	"go-tree-hollow/configs"
	"go-tree-hollow/internal/middleware"
	"go-tree-hollow/internal/modules/auth"
	"go-tree-hollow/internal/modules/email"
	"go-tree-hollow/internal/modules/post"
	"go-tree-hollow/internal/modules/tag"
	"go-tree-hollow/internal/modules/upload"
	"go-tree-hollow/internal/modules/user"
	"go-tree-hollow/pkg/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Server struct {
	config      *configs.Config
	db          *gorm.DB
	router      *gin.Engine
	redisClient *redis.Client
}

func NewServer(config *configs.Config) (*Server, error) {
	// 初始化数据库
	db, err := database.NewDB(config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	redisClient, err := email.NewClient(&config.Redis)
	if err != nil {
		return nil, err
	}

	// 初始化Gin
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	server := &Server{
		config:      config,
		db:          db,
		router:      router,
		redisClient: redisClient,
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

	// 邮箱模块
	emailSender := email.NewSender((*email.EmailConfig)(&s.config.Email))
	emailRepo := email.NewCodeRepository(s.redisClient, "app:email")
	emailService := email.NewEmailService(
		emailSender,
		emailRepo,
		s.config,
	)
	emailHandler := email.NewEmailHandler(emailService)
	email.RegisterRoutes(v1, emailHandler)

	// 内容模块 (需要认证)
	postRepo := post.NewRepository(s.db)
	postService := post.NewService(s.db, postRepo)
	postHandler := post.NewHandler(postService)
	
	// 点赞功能
	likeRepo := post.NewLikeRepository(s.db)
	likeService := post.NewLikeService(likeRepo)
	likeHandler := post.NewLikeHandler(likeService)
	
	// 评论功能
	commentRepo := post.NewCommentRepository(s.db)
	commentService := post.NewCommentService(commentRepo)
	commentHandler := post.NewCommentHandler(commentService)
	
	post.Routes(v1, postHandler, likeHandler, commentHandler, middleware.AuthRequired())

	// 标签模块
	tagRepo := tag.NewRepository(s.db)
	tagService := tag.NewService(s.db, tagRepo)
	tagHandler := tag.NewHandler(tagService)
	tag.RegisterRoutes(v1, tagHandler)

	// 文件上传模块 (需要认证)
	uploadHandler := upload.NewHandler()
	upload.Routes(v1, uploadHandler, middleware.AuthRequired())

	// 提供静态文件访问
	s.router.Static("/uploads", "./uploads")

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
