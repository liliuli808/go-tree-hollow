package user

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

// ProfileResponse 用户信息响应
type ProfileResponse struct {
	ID            uint   `json:"id"`
	Email         string `json:"email"`
	Nickname      string `json:"nickname"`
	AvatarURL     string `json:"avatar_url"`
	BackgroundURL string `json:"background_url"`
	CreatedAt     string `json:"created_at"`
	Bio           string `json:"bio"`
	Location      string `json:"location"`
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Nickname      *string `json:"nickname"`
	AvatarURL     *string `json:"avatar_url"`
	BackgroundURL *string `json:"background_url"`
	Birthday      *string `json:"birthday"`
	Bio           *string `json:"bio"`
	Location      *string `json:"location"`
}

type MyProfileResponse struct {
	ProfileResponse
	Birthday       string        `json:"birthday"`
	Age            int           `json:"age"`
	Constellation  string        `json:"constellation"`
	FollowCount    int64         `json:"follow_count"`
	FanCount       int64         `json:"fan_count"`
	LikedCount     int64         `json:"liked_count"`
	CollectedCount int64         `json:"collected_count"`
	Posts          []models.Post `json:"posts"`
}

// GetMyProfile 获取当前登录用户的完整资料
func (s *Service) GetMyProfile(userID uint) (*MyProfileResponse, error) {
	// 1. 基本信息
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, errors.New("获取用户信息失败")
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 2. 统计数据
	followCount, _ := s.repo.GetFollowCount(userID)
	fanCount, _ := s.repo.GetFanCount(userID)
	likedCount, _ := s.repo.GetReceivedLikeCount(userID)
	colCount, _ := s.repo.GetReceivedCollectionCount(userID)

	// 3. 帖子列表
	posts, _ := s.repo.GetUserPosts(userID)

	// 4. 计算年龄星座
	age := utils.CalculateAge(user.Birthday)
	constellation := utils.GetConstellation(user.Birthday)

	return &MyProfileResponse{
		ProfileResponse: ProfileResponse{
			ID:            user.ID,
			Email:         user.Email,
			Nickname:      user.Nickname,
			AvatarURL:     user.AvatarURL,
			BackgroundURL: user.BackgroundURL,
			CreatedAt:     user.CreatedAt.Format("2006-01-02 15:04:05"),
			Bio:           user.Bio,
			Location:      user.Location,
		},
		Birthday:       user.Birthday,
		Age:            age,
		Constellation:  constellation,
		FollowCount:    followCount,
		FanCount:       fanCount,
		LikedCount:     likedCount,
		CollectedCount: colCount,
		Posts:          posts,
	}, nil
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
		ID:            user.ID,
		Email:         user.Email,
		Nickname:      user.Nickname,
		AvatarURL:     user.AvatarURL,
		BackgroundURL: user.BackgroundURL,
		CreatedAt:     user.CreatedAt.Format("2006-01-02 15:04:05"),
		Bio:           user.Bio,
		Location:      user.Location,
	}, nil
}

// UpdateProfile 更新用户资料
func (s *Service) UpdateProfile(userID uint, req *UpdateProfileRequest) (*ProfileResponse, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, errors.New("获取用户信息失败")
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 更新字段
	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.AvatarURL != nil {
		user.AvatarURL = *req.AvatarURL
	}
	if req.BackgroundURL != nil {
		user.BackgroundURL = *req.BackgroundURL
	}
	if req.Birthday != nil {
		user.Birthday = *req.Birthday
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}
	if req.Location != nil {
		user.Location = *req.Location
	}

	// 保存更新
	if err := s.repo.Update(user); err != nil {
		return nil, errors.New("更新用户信息失败")
	}

	return s.GetProfile(userID)
}
