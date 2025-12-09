package post

import "github.com/gin-gonic/gin"

// Routes 为帖子模块在给定的 Gin 路由组中设置 API 路由。
// 这里定义的所有路由都受提供的认证中间件保护。
func Routes(r *gin.RouterGroup, handler *Handler, authMiddleware gin.HandlerFunc) {
	// 用于帖子特定操作（单个帖子的 CRUD）的路由组。
	// 这些路由允许认证用户管理他们自己的或可访问的帖子。
	posts := r.Group("/posts")
	posts.Use(authMiddleware) // 对此组中的所有路由应用认证中间件。
	{
		posts.POST("", handler.CreatePost)       // POST /api/v1/posts - 创建新帖子。
		posts.GET("/:id", handler.GetPost)       // GET /api/v1/posts/:id - 根据ID检索单个帖子。
		posts.PUT("/:id", handler.UpdatePost)     // PUT /api/v1/posts/:id - 根据ID更新现有帖子。
		posts.DELETE("/:id", handler.DeletePost) // DELETE /api/v1/posts/:id - 根据ID软删除帖子。
	}

	// 用于用户特定帖子列表的路由组。
	// 此路由允许获取属于特定用户的所有帖子。
	users := r.Group("/users")
	users.Use(authMiddleware) // 对此组中的所有路由应用认证中间件。
	{
		users.GET("/:userID/posts", handler.ListPosts) // GET /api/v1/users/:userID/posts - 列出特定用户的所有帖子，支持分页。
	}
}
