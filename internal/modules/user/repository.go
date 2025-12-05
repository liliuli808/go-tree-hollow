package user

import (
	"go-tree-hollow/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetByID 根据ID获取用户
func (r *Repository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// Update 更新用户信息
func (r *Repository) Update(user *models.User) error {
	return r.db.Save(user).Error
}
