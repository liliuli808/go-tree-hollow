package upload

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler for the upload module.
type Handler struct{}

// NewHandler creates a new upload handler.
func NewHandler() *Handler {
	return &Handler{}
}

// UploadFile handles the file upload request.
func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Generate a unique filename
	ext := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + ext
	dst := filepath.Join("uploads", newFileName)

	// Save the file
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	// Return the file path
	c.JSON(http.StatusOK, gin.H{"path": fmt.Sprintf("/%s", dst)})
}
