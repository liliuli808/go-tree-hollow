package chat

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for chat operations
type Handler struct {
	service Service
	hub     *Hub
}

// NewHandler creates a new chat handler instance
func NewHandler(service Service, hub *Hub) *Handler {
	return &Handler{service: service, hub: hub}
}

// SendMessageRequest represents the request body for sending a message
type SendMessageRequest struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// GetConversations returns a list of conversations for the authenticated user
// @Summary Get user's conversations
// @Tags Chat
// @Security BearerAuth
// @Success 200 {array} ConversationResponse
// @Router /chat/conversations [get]
func (h *Handler) GetConversations(c *gin.Context) {
	userID := c.GetUint("userID")

	conversations, err := h.service.GetConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": conversations})
}

// GetMessages returns messages between the authenticated user and another user
// @Summary Get messages with a user
// @Tags Chat
// @Security BearerAuth
// @Param userId path int true "Other user's ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(50)
// @Success 200 {array} models.Message
// @Router /chat/messages/{userId} [get]
func (h *Handler) GetMessages(c *gin.Context) {
	currentUserID := c.GetUint("userID")
	otherUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	messages, err := h.service.GetMessages(currentUserID, uint(otherUserID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	// Mark messages as read
	_ = h.service.MarkConversationAsRead(currentUserID, uint(otherUserID))

	c.JSON(http.StatusOK, gin.H{"data": messages})
}

// SendMessage sends a message to another user (REST fallback)
// @Summary Send a message
// @Tags Chat
// @Security BearerAuth
// @Param request body SendMessageRequest true "Message details"
// @Success 201 {object} models.Message
// @Router /chat/messages [post]
func (h *Handler) SendMessage(c *gin.Context) {
	userID := c.GetUint("userID")

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.service.SendMessage(userID, req.ReceiverID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Broadcast via WebSocket if available
	if h.hub != nil {
		h.hub.SendToUser(req.ReceiverID, &WebSocketMessage{
			Type:    "message",
			From:    userID,
			Content: req.Content,
			Message: message,
		})
	}

	c.JSON(http.StatusCreated, gin.H{"data": message})
}

// MarkAsRead marks a message as read
// @Summary Mark message as read
// @Tags Chat
// @Security BearerAuth
// @Param messageId path int true "Message ID"
// @Success 200
// @Router /chat/messages/{messageId}/read [put]
func (h *Handler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("userID")
	messageID, err := strconv.ParseUint(c.Param("messageId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	if err := h.service.MarkAsRead(uint(messageID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetUnreadCount returns the total unread message count
// @Summary Get unread message count
// @Tags Chat
// @Security BearerAuth
// @Success 200 {object} map[string]int64
// @Router /chat/unread-count [get]
func (h *Handler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.service.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get unread count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// HandleWebSocket handles WebSocket connection upgrade
func (h *Handler) HandleWebSocket(c *gin.Context) {
	userID := c.GetUint("userID")
	h.hub.HandleConnection(c.Writer, c.Request, userID)
}
