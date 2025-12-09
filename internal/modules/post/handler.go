package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler 处理与帖子相关的 HTTP 请求。它是帖子模块 API 的入口点。
type Handler struct {
	service Service // service 提供对帖子模块业务逻辑的访问。
}

// NewHandler 创建并返回一个新的 Handler 实例。
// 它接收一个 Service 接口的实现作为依赖项，以解耦关注点。
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// CreatePost 处理创建新帖子的 HTTP POST 请求。
// 它期望一个符合 CreatePostDto 结构体的 JSON 请求体。
// @Summary 创建新帖子
// @Description 创建包含文本内容和媒体URL的新帖子。帖子初始状态为草稿。
// @Tags posts
// @Accept json
// @Produce json
// @Param post body CreatePostDto true "帖子创建数据"
// @Success 201 {object} models.Post "成功创建帖子"
// @Failure 400 {object} gin.H "无效的请求体或缺少必填字段"
// @Failure 500 {object} gin.H "内部服务器错误，例如数据库错误或检测到敏感内容"
// @Security BearerAuth
// @Router /posts [post]
func (h *Handler) CreatePost(c *gin.Context) {
	var dto CreatePostDto
	// 尝试将 JSON 请求体绑定到 CreatePostDto 结构体。
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层创建帖子。
	// 帖子的 UserID 通常应从已认证用户的上下文（例如 JWT claims）中获取。
	// 此处假设认证中间件已在 DTO 或上下文中设置了 UserID。
	post, err := h.service.CreatePost(&dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 响应新创建的帖子和 201 Created 状态。
	c.JSON(http.StatusCreated, post)
}

// GetPost 处理根据ID检索单个帖子的 HTTP GET 请求。
// 它期望帖子ID作为URL路径参数。
// @Summary 根据ID获取帖子
// @Description 检索由其唯一ID标识的单个帖子。
// @Tags posts
// @Produce json
// @Param id path int true "帖子ID"
// @Success 200 {object} models.Post "成功检索到帖子"
// @Failure 400 {object} gin.H "无效的帖子ID格式"
// @Failure 404 {object} gin.H "未找到帖子"
// @Security BearerAuth
// @Router /posts/{id} [get]
func (h *Handler) GetPost(c *gin.Context) {
	// 从URL路径参数中解析帖子ID。
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID格式"})
		return
	}

	// 调用服务层根据ID检索帖子。
	post, err := h.service.GetPost(uint(id))
	if err != nil {
		// 如果服务返回错误，通常意味着未找到帖子。
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到帖子"})
		return
	}

	// 响应检索到的帖子数据和 200 OK 状态。
	c.JSON(http.StatusOK, post)
}

// UpdatePost 处理更新现有帖子的 HTTP PUT 请求。
// 它期望帖子ID作为路径参数，并期望一个符合 UpdatePostDto 结构体的 JSON 请求体。
// @Summary 更新现有帖子
// @Description 更新由其ID标识的现有帖子，包括新的文本内容、媒体URL或状态。
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "帖子ID"
// @Param post body UpdatePostDto true "帖子更新数据"
// @Success 200 {object} models.Post "成功更新帖子"
// @Failure 400 {object} gin.H "无效的请求体、帖子ID格式或敏感内容"
// @Failure 404 {object} gin.H "未找到帖子"
// @Failure 500 {object} gin.H "内部服务器错误"
// @Security BearerAuth
// @Router /posts/{id} [put]
func (h *Handler) UpdatePost(c *gin.Context) {
	// 从URL路径参数中解析帖子ID。
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID格式"})
		return
	}

	var dto UpdatePostDto
	// 尝试将 JSON 请求体绑定到 UpdatePostDto 结构体。
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层更新帖子。
	post, err := h.service.UpdatePost(uint(id), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 响应更新后的帖子数据和 200 OK 状态。
	c.JSON(http.StatusOK, post)
}

// DeletePost 处理根据ID删除帖子的 HTTP DELETE 请求。
// 它期望帖子ID作为路径参数。此操作执行软删除。
// @Summary 删除帖子
// @Description 软删除由其ID标识的帖子。帖子不会被永久删除，而是被标记为已删除。
// @Tags posts
// @Param id path int true "帖子ID"
// @Success 204 "成功删除帖子 (无内容)"
// @Failure 400 {object} gin.H "无效的帖子ID格式"
// @Failure 500 {object} gin.H "内部服务器错误"
// @Security BearerAuth
// @Router /posts/{id} [delete]
func (h *Handler) DeletePost(c *gin.Context) {
	// 从URL路径参数中解析帖子ID。
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID格式"})
		return
	}

	// 调用服务层删除帖子。
	if err := h.service.DeletePost(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 成功软删除后，响应 204 No Content 状态。
	c.JSON(http.StatusNoContent, nil)
}

// ListPosts 处理列出特定用户帖子的 HTTP GET 请求。
// 它期望用户ID作为路径参数，以及可选的 'page' 和 'pageSize' 查询参数用于分页。
// @Summary 列出用户的帖子
// @Description 检索属于特定用户的分页帖子列表。
// @Tags posts
// @Produce json
// @Param userID path int true "用户ID"
// @Param page query int false "页码 (默认为1)"
// @Param pageSize query int false "每页项目数 (默认为10)"
// @Success 200 {object} gin.H{data=[]models.Post,total=int64,page=int} "成功检索到帖子列表"
// @Failure 400 {object} gin.H "无效的用户ID格式或查询参数"
// @Failure 500 {object} gin.H "内部服务器错误"
// @Security BearerAuth
// @Router /users/{userID}/posts [get]
func (h *Handler) ListPosts(c *gin.Context) {
	// 从URL路径参数中解析用户ID。
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	// 解析可选的 'page' 和 'pageSize' 查询参数，并设置默认值。
	// strconv.Atoi 将字符串转换为整数，'_' 丢弃默认查询参数的潜在错误。
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 调用服务层获取给定用户的分页帖子列表。
	posts, total, err := h.service.ListPosts(uint(userID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 响应帖子列表、总数和当前页，以及 200 OK 状态。
	c.JSON(http.StatusOK, gin.H{
		"data":  posts,
		"total": total,
		"page":  page,
	})
}
