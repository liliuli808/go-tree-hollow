package tag

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// Handler handles tag-related HTTP requests.
type Handler struct {
	service Service
}

// NewHandler creates a new tag handler instance.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetAllTags handles the HTTP GET request to retrieve all tags.
// @Summary Get all tags
// @Description Retrieve a list of all available tags
// @Tags tags
// @Produce json
// @Success 200 {array} models.Tag "Successfully retrieved tags"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /tags [get]
func (h *Handler) GetAllTags(c *gin.Context) {
	tags, err := h.service.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}
