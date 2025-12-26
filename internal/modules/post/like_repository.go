package post

import (
	"go-tree-hollow/internal/models"

	"gorm.io/gorm"
)

// Repository defines like data access operations
type LikeRepository interface {
	Create(like *models.Like) error
	Delete(userID, postID uint) error
	FindByUserAndPost(userID, postID uint) (*models.Like, error)
	CountByPost(postID uint) (int64, error)
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(like *models.Like) error {
	return r.db.Create(like).Error
}

func (r *likeRepository) Delete(userID, postID uint) error {
	return r.db.Unscoped().Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.Like{}).Error
}

func (r *likeRepository) FindByUserAndPost(userID, postID uint) (*models.Like, error) {
	var like models.Like
	err := r.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error
	return &like, err
}

func (r *likeRepository) CountByPost(postID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Like{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}
