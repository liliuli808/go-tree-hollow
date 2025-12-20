package models

import "gorm.io/gorm"

// Comment represents a comment on a post
type Comment struct {
	gorm.Model
	UserID  uint   `json:"user_id" gorm:"not null;index"`
	User    User   `json:"user" gorm:"foreignKey:UserID"`
	PostID  uint   `json:"post_id" gorm:"not null;index"`
	Post    Post   `json:"post" gorm:"foreignKey:PostID"`
	Content string `json:"content" gorm:"type:text;not null"`
}
