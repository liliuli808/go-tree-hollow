package post

import (
	"go-tree-hollow/internal/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	FindByPost(postID uint, page, pageSize int) ([]*models.Comment, int64, error)
	FindByID(id uint) (*models.Comment, error)
	Delete(id uint) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) FindByPost(postID uint, page, pageSize int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	offset := (page - 1) * pageSize

	// Get total count
	if err := r.db.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get comments with pagination
	err := r.db.Where("post_id = ?", postID).
		Order("created_at desc").
		Limit(pageSize).
		Offset(offset).
		Preload("User").
		Find(&comments).Error

	return comments, total, err
}

func (r *commentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("User").First(&comment, id).Error
	return &comment, err
}

func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}
