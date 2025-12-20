package models

import "gorm.io/gorm"

// Like represents a like on a post by a user
type Like struct {
	gorm.Model
	UserID uint  `json:"user_id" gorm:"not null;index"`
	User   User  `json:"user" gorm:"foreignKey:UserID"`
	PostID uint  `json:"post_id" gorm:"not null;index"`
	Post   Post  `json:"post" gorm:"foreignKey:PostID"`
}

// Ensure unique constraint: one user can only like a post once
// Add index for faster queries
func (Like) TableName() string {
	return "likes"
}
