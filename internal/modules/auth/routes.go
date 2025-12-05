package auth

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册认证模块路由
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	// 创建 /auth 子路由组
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", handler.Register)
		authGroup.POST("/login", handler.Login)
	}
}
