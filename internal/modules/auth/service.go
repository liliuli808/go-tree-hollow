package auth

import (
	"errors"
	"go-tree-hollow/internal/models"
	"go-tree-hollow/pkg/utils"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register 注册用户
func (s *Service) Register(req *RegisterRequest) (*models.User, error) {
	// 检查用户是否已存在
	existingUser, err := s.repo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("用户已存在")
	}

	// 创建新用户
	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	return user, nil
}

// Login 用户登录
func (s *Service) Login(req *LoginRequest) (string, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return "", errors.New("邮箱或密码错误")
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return "", errors.New("邮箱或密码错误")
	}

	// 生成JWT
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", errors.New("生成令牌失败")
	}

	return token, nil
}

// GetProfile 获取用户信息（示例业务）
func (s *Service) GetProfile(userID uint) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}
