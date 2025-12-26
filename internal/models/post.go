package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Post represents the canned content created by a user.
type Post struct {
	gorm.Model
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Type        string         `json:"type" gorm:"not null;index"` // e.g., "text_image", "video", "audio", "live_photo"
	TextContent string         `json:"text_content,omitempty" gorm:"type:text"`
	MediaURLs   datatypes.JSON `json:"media_urls,omitempty" gorm:"type:json"`
	CoverURL    string         `json:"cover_url,omitempty" gorm:"type:varchar(1024)"`
	Status      string         `json:"status" gorm:"not null;default:'draft';index"` // "draft", "published"
	TagID       *uint          `json:"tag_id,omitempty" gorm:"index"`
	Tag         *Tag           `json:"tag,omitempty" gorm:"foreignKey:TagID"`
	LikesCount  int64          `json:"likes_count" gorm:"-"`
	IsLiked     bool           `json:"is_liked" gorm:"-"`
}
