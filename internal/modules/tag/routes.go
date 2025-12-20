package tag

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all tag-related routes.
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	tags := router.Group("/tags")
	{
		tags.GET("", handler.GetAllTags) // GET /api/v1/tags - Get all tags
	}
}
