package tag

import (
	"go-tree-hollow/internal/models"
	"gorm.io/gorm"
)

// Service defines the interface for tag business logic operations.
type Service interface {
	GetAllTags() ([]*models.Tag, error)
	GetTagByID(id uint) (*models.Tag, error)
}

type service struct {
	repo Repository
}

// NewService creates a new tag service instance.
func NewService(db *gorm.DB, repo Repository) Service {
	return &service{repo: repo}
}

// GetAllTags retrieves all available tags.
func (s *service) GetAllTags() ([]*models.Tag, error) {
	return s.repo.FindAll()
}

// GetTagByID retrieves a tag by its ID.
func (s *service) GetTagByID(id uint) (*models.Tag, error) {
	return s.repo.FindByID(id)
}
