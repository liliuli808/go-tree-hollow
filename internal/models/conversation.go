package models

import (
	"time"

	"gorm.io/gorm"
)

// Conversation represents a chat conversation between two users
type Conversation struct {
	gorm.Model
	User1ID       uint       `gorm:"not null;index:idx_users,unique" json:"user1_id"`
	User2ID       uint       `gorm:"not null;index:idx_users,unique" json:"user2_id"`
	LastMessageID *uint      `json:"last_message_id"`
	LastMessageAt *time.Time `json:"last_message_at"`

	// Relations
	User1       User     `gorm:"foreignKey:User1ID" json:"user1,omitempty"`
	User2       User     `gorm:"foreignKey:User2ID" json:"user2,omitempty"`
	LastMessage *Message `gorm:"foreignKey:LastMessageID" json:"last_message,omitempty"`
}

// GetOtherUserID returns the ID of the other user in the conversation
func (c *Conversation) GetOtherUserID(currentUserID uint) uint {
	if c.User1ID == currentUserID {
		return c.User2ID
	}
	return c.User1ID
}

// GetOtherUser returns the other user in the conversation
func (c *Conversation) GetOtherUser(currentUserID uint) User {
	if c.User1ID == currentUserID {
		return c.User2
	}
	return c.User1
}
