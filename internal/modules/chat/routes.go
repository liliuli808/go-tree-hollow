package chat

import (
	"go-tree-hollow/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers chat routes on the given router group
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	chatGroup := rg.Group("/chat")
	chatGroup.Use(middleware.AuthRequired())
	{
		chatGroup.GET("/conversations", handler.GetConversations)
		chatGroup.GET("/messages/:userId", handler.GetMessages)
		chatGroup.POST("/messages", handler.SendMessage)
		chatGroup.PUT("/messages/:messageId/read", handler.MarkAsRead)
		chatGroup.GET("/unread-count", handler.GetUnreadCount)
	}

	// WebSocket endpoint (also requires auth)
	rg.GET("/ws/chat", middleware.AuthRequired(), handler.HandleWebSocket)
}
