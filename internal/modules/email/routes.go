package email

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册认证模块路由
func RegisterRoutes(router *gin.RouterGroup, handler *EmailHandler) {
	// 创建 /auth 子路由组
	authGroup := router.Group("/email")
	{
		authGroup.POST("/send", handler.SendVerificationCode)
	}
}
