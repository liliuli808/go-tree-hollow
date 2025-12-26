package chat

import (
	"go-tree-hollow/internal/models"
)

// Service defines the interface for chat business logic
type Service interface {
	// SendMessage sends a message from one user to another
	SendMessage(senderID, receiverID uint, content string) (*models.Message, error)
	// GetMessages retrieves messages between two users
	GetMessages(currentUserID, otherUserID uint, page, pageSize int) ([]*models.Message, error)
	// GetConversations retrieves all conversations for a user
	GetConversations(userID uint) ([]*ConversationResponse, error)
	// MarkAsRead marks a message as read
	MarkAsRead(messageID, userID uint) error
	// MarkConversationAsRead marks all messages in a conversation as read
	MarkConversationAsRead(currentUserID, otherUserID uint) error
	// GetUnreadCount returns the total unread message count for a user
	GetUnreadCount(userID uint) (int64, error)
}

// ConversationResponse represents a conversation with additional metadata
type ConversationResponse struct {
	ID            uint        `json:"id"`
	OtherUser     UserSummary `json:"other_user"`
	LastMessage   string      `json:"last_message"`
	LastMessageAt int64       `json:"last_message_at"`
	UnreadCount   int         `json:"unread_count"`
}

// UserSummary represents a simplified user object
type UserSummary struct {
	ID        uint   `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

type service struct {
	repo Repository
}

// NewService creates a new chat service instance
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// SendMessage sends a message and updates the conversation
func (s *service) SendMessage(senderID, receiverID uint, content string) (*models.Message, error) {
	// Create the message
	message := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}
	if err := s.repo.CreateMessage(message); err != nil {
		return nil, err
	}

	// Get or create conversation and update last message
	conv, err := s.repo.GetOrCreateConversation(senderID, receiverID)
	if err != nil {
		return message, err // Return message even if conversation update fails
	}

	if err := s.repo.UpdateConversationLastMessage(conv.ID, message.ID); err != nil {
		return message, err
	}

	// Reload message with relations
	return s.repo.GetMessageByID(message.ID)
}

// GetMessages retrieves paginated messages between two users
func (s *service) GetMessages(currentUserID, otherUserID uint, page, pageSize int) ([]*models.Message, error) {
	offset := (page - 1) * pageSize
	return s.repo.GetMessagesBetweenUsers(currentUserID, otherUserID, pageSize, offset)
}

// GetConversations retrieves all conversations for a user with metadata
func (s *service) GetConversations(userID uint) ([]*ConversationResponse, error) {
	conversations, err := s.repo.GetConversationsByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*ConversationResponse, 0, len(conversations))
	for _, conv := range conversations {
		otherUser := conv.GetOtherUser(userID)

		response := &ConversationResponse{
			ID: conv.ID,
			OtherUser: UserSummary{
				ID:        otherUser.ID,
				Nickname:  otherUser.Nickname,
				AvatarURL: otherUser.AvatarURL,
			},
		}

		if conv.LastMessage != nil {
			response.LastMessage = conv.LastMessage.Content
			if conv.LastMessageAt != nil {
				response.LastMessageAt = conv.LastMessageAt.UnixMilli()
			}
		}

		// TODO: Calculate unread count per conversation
		// For now, we'll set it to 0 and improve later
		response.UnreadCount = 0

		responses = append(responses, response)
	}

	return responses, nil
}

// MarkAsRead marks a single message as read
func (s *service) MarkAsRead(messageID, userID uint) error {
	message, err := s.repo.GetMessageByID(messageID)
	if err != nil {
		return err
	}

	// Only the receiver can mark a message as read
	if message.ReceiverID != userID {
		return nil
	}

	return s.repo.MarkMessageAsRead(messageID)
}

// MarkConversationAsRead marks all messages from another user as read
func (s *service) MarkConversationAsRead(currentUserID, otherUserID uint) error {
	return s.repo.MarkAllMessagesAsRead(otherUserID, currentUserID)
}

// GetUnreadCount returns the total unread message count
func (s *service) GetUnreadCount(userID uint) (int64, error) {
	return s.repo.GetUnreadCount(userID)
}
