package post

import "go-tree-hollow/internal/models"

type CreateCommentDto struct {
	UserID  uint   `json:"user_id" binding:"required"`
	PostID  uint   `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type CommentService interface {
	CreateComment(dto *CreateCommentDto) (*models.Comment, error)
	GetCommentsByPost(postID uint, page, pageSize int) ([]*models.Comment, int64, error)
	DeleteComment(id, userID uint) error
}

type commentService struct {
	repo CommentRepository
}

func NewCommentService(repo CommentRepository) CommentService {
	return &commentService{repo: repo}
}

func (s *commentService) CreateComment(dto *CreateCommentDto) (*models.Comment, error) {
	comment := &models.Comment{
		UserID:  dto.UserID,
		PostID:  dto.PostID,
		Content: dto.Content,
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}

	return s.repo.FindByID(comment.ID)
}

func (s *commentService) GetCommentsByPost(postID uint, page, pageSize int) ([]*models.Comment, int64, error) {
	return s.repo.FindByPost(postID, page, pageSize)
}

func (s *commentService) DeleteComment(id, userID uint) error {
	comment, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Only allow user to delete their own comments
	if comment.UserID != userID {
		return err
	}

	return s.repo.Delete(id)
}
