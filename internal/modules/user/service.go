package user

import (
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ProfileResponse 用户信息响应
type ProfileResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// GetProfile 获取用户资料
func (s *Service) GetProfile(userID uint) (*ProfileResponse, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, errors.New("获取用户信息失败")
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	return &ProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
