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
	Status      string         `json:"status" gorm:"not null;default:'draft';index"` // "draft", "published"
	Tags        []*Tag         `json:"tags,omitempty" gorm:"many2many:post_tags;"`
}
