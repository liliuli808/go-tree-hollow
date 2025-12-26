package chat

import (
	"go-tree-hollow/internal/models"
	"time"

	"gorm.io/gorm"
)

// Repository defines the interface for chat data operations
type Repository interface {
	// Message operations
	CreateMessage(message *models.Message) error
	GetMessageByID(id uint) (*models.Message, error)
	GetMessagesBetweenUsers(user1ID, user2ID uint, limit, offset int) ([]*models.Message, error)
	MarkMessageAsRead(messageID uint) error
	MarkAllMessagesAsRead(senderID, receiverID uint) error
	GetUnreadCount(userID uint) (int64, error)

	// Conversation operations
	GetOrCreateConversation(user1ID, user2ID uint) (*models.Conversation, error)
	GetConversationsByUserID(userID uint) ([]*models.Conversation, error)
	UpdateConversationLastMessage(conversationID, messageID uint) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new chat repository instance
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateMessage creates a new message in the database
func (r *repository) CreateMessage(message *models.Message) error {
	return r.db.Create(message).Error
}

// GetMessageByID retrieves a message by its ID
func (r *repository) GetMessageByID(id uint) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("Sender").Preload("Receiver").First(&message, id).Error
	return &message, err
}

// GetMessagesBetweenUsers retrieves messages between two users with pagination
func (r *repository) GetMessagesBetweenUsers(user1ID, user2ID uint, limit, offset int) ([]*models.Message, error) {
	var messages []*models.Message
	err := r.db.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			user1ID, user2ID, user2ID, user1ID).
		Preload("Sender").
		Preload("Receiver").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

// MarkMessageAsRead marks a single message as read
func (r *repository) MarkMessageAsRead(messageID uint) error {
	now := time.Now()
	return r.db.Model(&models.Message{}).
		Where("id = ? AND read_at IS NULL", messageID).
		Update("read_at", now).Error
}

// MarkAllMessagesAsRead marks all messages from sender to receiver as read
func (r *repository) MarkAllMessagesAsRead(senderID, receiverID uint) error {
	now := time.Now()
	return r.db.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND read_at IS NULL", senderID, receiverID).
		Update("read_at", now).Error
}

// GetUnreadCount returns the count of unread messages for a user
func (r *repository) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Message{}).
		Where("receiver_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

// GetOrCreateConversation gets an existing conversation or creates a new one
func (r *repository) GetOrCreateConversation(user1ID, user2ID uint) (*models.Conversation, error) {
	// Ensure consistent ordering (lower ID first)
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	var conversation models.Conversation
	err := r.db.
		Where("user1_id = ? AND user2_id = ?", user1ID, user2ID).
		Preload("User1").
		Preload("User2").
		Preload("LastMessage").
		First(&conversation).Error

	if err == gorm.ErrRecordNotFound {
		conversation = models.Conversation{
			User1ID: user1ID,
			User2ID: user2ID,
		}
		if err := r.db.Create(&conversation).Error; err != nil {
			return nil, err
		}
		// Reload with preloads
		r.db.Preload("User1").Preload("User2").First(&conversation, conversation.ID)
	} else if err != nil {
		return nil, err
	}

	return &conversation, nil
}

// GetConversationsByUserID retrieves all conversations for a user
func (r *repository) GetConversationsByUserID(userID uint) ([]*models.Conversation, error) {
	var conversations []*models.Conversation
	err := r.db.
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Preload("User1").
		Preload("User2").
		Preload("LastMessage").
		Order("last_message_at DESC").
		Find(&conversations).Error
	return conversations, err
}

// UpdateConversationLastMessage updates the last message of a conversation
func (r *repository) UpdateConversationLastMessage(conversationID, messageID uint) error {
	now := time.Now()
	return r.db.Model(&models.Conversation{}).
		Where("id = ?", conversationID).
		Updates(map[string]interface{}{
			"last_message_id": messageID,
			"last_message_at": now,
		}).Error
}
