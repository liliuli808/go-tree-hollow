package models

import "gorm.io/gorm"

// Tag represents a tag that can be associated with a post.
type Tag struct {
	gorm.Model
	Name  string  `json:"name" gorm:"unique;not null;index"` // Name of the tag, must be unique.
	Posts []*Post `json:"posts,omitempty" gorm:"many2many:post_tags;"` // Posts associated with this tag.
}
