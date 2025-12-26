package post

import "github.com/gin-gonic/gin"

// Routes 为帖子模块在给定的 Gin 路由组中设置 API 路由。
// 这里定义的所有路由都受提供的认证中间件保护。
func Routes(r *gin.RouterGroup, handler *Handler, likeHandler *LikeHandler, commentHandler *CommentHandler, authMiddleware gin.HandlerFunc, optionalAuthMiddleware gin.HandlerFunc) {
	// 公开路由组（不需要认证）
	publicPosts := r.Group("/posts")
	{
		publicPosts.GET("", handler.GetAllPosts)                                               // GET /api/v1/posts - 获取所有帖子
		publicPosts.GET("/:id", handler.GetPost)                                               // GET /api/v1/posts/:id - 获取单个帖子
		publicPosts.GET("/:id/like/status", optionalAuthMiddleware, likeHandler.GetLikeStatus) // GET /api/v1/posts/:id/like/status - 获取点赞状态
		publicPosts.GET("/:id/comments", commentHandler.GetComments)                           // GET /api/v1/posts/:id/comments - 获取评论列表
	}

	// 需要认证的路由组
	authPosts := r.Group("/posts")
	authPosts.Use(authMiddleware)
	{
		authPosts.POST("", handler.CreatePost)       // POST /api/v1/posts - 创建新帖子
		authPosts.PUT("/:id", handler.UpdatePost)    // PUT /api/v1/posts/:id - 更新帖子
		authPosts.DELETE("/:id", handler.DeletePost) // DELETE /api/v1/posts/:id - 删除帖子

		// 点赞（需要登录）
		authPosts.POST("/:id/like", likeHandler.ToggleLike) // POST /api/v1/posts/:id/like - 切换点赞

		// 评论（需要登录）
		authPosts.POST("/:id/comments", commentHandler.CreateComment) // POST /api/v1/posts/:id/comments - 创建评论
	}

	// 删除评论（需要认证）
	r.DELETE("/comments/:id", authMiddleware, commentHandler.DeleteComment)

	// 用户帖子列表（需要认证）
	users := r.Group("/users")
	users.Use(authMiddleware)
	{
		users.GET("/:userID/posts", handler.ListPosts) // GET /api/v1/users/:userID/posts - 列出用户帖子
	}
}
