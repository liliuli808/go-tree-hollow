package models

import "gorm.io/gorm"

// Follow 关注关系
type Follow struct {
	gorm.Model
	FollowerID uint `json:"follower_id" gorm:"not null;index"`
	Follower   User `json:"follower" gorm:"foreignKey:FollowerID"`
	FollowedID uint `json:"followed_id" gorm:"not null;index"`
	Followed   User `json:"followed" gorm:"foreignKey:FollowedID"`
}

// Collection 收藏
type Collection struct {
	gorm.Model
	UserID uint `json:"user_id" gorm:"not null;index"`
	User   User `json:"user" gorm:"foreignKey:UserID"`
	PostID uint `json:"post_id" gorm:"not null;index"`
	Post   Post `json:"post" gorm:"foreignKey:PostID"`
}

func (Follow) TableName() string {
	return "follows"
}

func (Collection) TableName() string {
	return "collections"
}
