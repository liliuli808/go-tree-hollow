package post

import (
	"go-tree-hollow/internal/models"
	"gorm.io/gorm"
)

type LikeService interface {
	ToggleLike(userID, postID uint) (bool, error) // Returns true if liked, false if unliked
	GetLikeCount(postID uint) (int64, error)
	IsLikedByUser(userID, postID uint) (bool, error)
}

type likeService struct {
	repo LikeRepository
}

func NewLikeService(repo LikeRepository) LikeService {
	return &likeService{repo: repo}
}

func (s *likeService) ToggleLike(userID, postID uint) (bool, error) {
	// Check if already liked
	_, err := s.repo.FindByUserAndPost(userID, postID)
	
	if err == gorm.ErrRecordNotFound {
		// Not liked yet, create like
		like := &models.Like{
			UserID: userID,
			PostID: postID,
		}
		if err := s.repo.Create(like); err != nil {
			return false, err
		}
		return true, nil
	} else if err != nil {
		return false, err
	}
	
	// Already liked, remove like
	if err := s.repo.Delete(userID, postID); err != nil {
		return false, err
	}
	return false, nil
}

func (s *likeService) GetLikeCount(postID uint) (int64, error) {
	return s.repo.CountByPost(postID)
}

func (s *likeService) IsLikedByUser(userID, postID uint) (bool, error) {
	_, err := s.repo.FindByUserAndPost(userID, postID)
	if err == gorm.ErrRecordNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
