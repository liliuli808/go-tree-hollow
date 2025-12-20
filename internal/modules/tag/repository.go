package tag

import (
	"go-tree-hollow/internal/models"
	"gorm.io/gorm"
)

// Repository defines the interface for tag data access operations.
type Repository interface {
	FindAll() ([]*models.Tag, error)
	FindByID(id uint) (*models.Tag, error)
	FindByName(name string) (*models.Tag, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new tag repository instance.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// FindAll retrieves all tags from the database.
func (r *repository) FindAll() ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.Order("id ASC").Find(&tags).Error
	return tags, err
}

// FindByID retrieves a tag by its ID.
func (r *repository) FindByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, id).Error
	return &tag, err
}

// FindByName retrieves a tag by its name.
func (r *repository) FindByName(name string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	return &tag, err
}
