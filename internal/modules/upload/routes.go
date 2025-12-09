package upload

import "github.com/gin-gonic/gin"

// Routes sets up the routes for the upload module.
func Routes(r *gin.RouterGroup, handler *Handler, authMiddleware gin.HandlerFunc) {
	upload := r.Group("/upload")
	upload.Use(authMiddleware)
	{
		upload.POST("", handler.UploadFile)
	}
}
