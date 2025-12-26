package models

import (
	"time"

	"gorm.io/gorm"
)

// Message represents a chat message between two users
type Message struct {
	gorm.Model
	SenderID   uint       `gorm:"not null;index:idx_sender_receiver" json:"sender_id"`
	ReceiverID uint       `gorm:"not null;index:idx_sender_receiver" json:"receiver_id"`
	Content    string     `gorm:"type:text;not null" json:"content"`
	ReadAt     *time.Time `json:"read_at"`

	// Relations
	Sender   User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Receiver User `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
}

// IsRead returns whether the message has been read
func (m *Message) IsRead() bool {
	return m.ReadAt != nil
}

// MarkAsRead marks the message as read with the current timestamp
func (m *Message) MarkAsRead() {
	now := time.Now()
	m.ReadAt = &now
}
