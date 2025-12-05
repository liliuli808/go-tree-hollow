package user

import (
	"go-tree-hollow/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册用户模块路由
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	// 需要认证的用户相关路由
	userGroup := router.Group("/users")
	userGroup.Use(middleware.AuthRequired())
	{
		userGroup.GET("/profile", handler.GetProfile)
	}
}
