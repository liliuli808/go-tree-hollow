package post

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	service LikeService
}

func NewLikeHandler(service LikeService) *LikeHandler {
	return &LikeHandler{service: service}
}

// ToggleLike handles POST /api/v1/posts/:id/like
func (h *LikeHandler) ToggleLike(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	liked, err := h.service.ToggleLike(userID.(uint), uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get updated like count
	count, _ := h.service.GetLikeCount(uint(postID))

	c.JSON(http.StatusOK, gin.H{
		"liked": liked,
		"count": count,
	})
}

// GetLikeStatus handles GET /api/v1/posts/:id/like/status
func (h *LikeHandler) GetLikeStatus(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	liked, err := h.service.IsLikedByUser(userID.(uint), uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	count, err := h.service.GetLikeCount(uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"liked": liked,
		"count": count,
	})
}
