package post

import (
	"encoding/json"
	"errors"
	"go-tree-hollow/internal/models"
	"log"

	"github.com/importcjj/sensitive"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreatePostDto 定义了创建新帖子的数据结构。
// 它用于输入验证以及从处理程序到服务层的数据传输。
type CreatePostDto struct {
	UserID      uint            `json:"user_id" binding:"required"` // 创建帖子的用户ID。必需。
	Type        string          `json:"type" binding:"required"`    // 内容类型（例如："text_image", "video", "audio", "live_photo"）。必需。
	TextContent string          `json:"text_content"`               // 帖子的主要文本内容。
	MediaURLs   json.RawMessage `json:"media_urls"`                 // 图像、视频、音频等的URL JSON数组。
	Status      string          `json:"status"`                     // 帖子的初始状态（例如："draft" 草稿, "published" 已发布）。
	Tags        []string        `json:"tags"`                       // 与帖子关联的标签名称列表。
}

// UpdatePostDto 定义了更新现有帖子的数据结构。
// 它使用指针表示可选字段，允许进行部分更新而不会用零值覆盖现有数据。
type UpdatePostDto struct {
	TextContent *string         `json:"text_content"` // 帖子新的文本内容。指针为nil表示不更新。
	MediaURLs   json.RawMessage `json:"media_urls"`   // 媒体URL的新的JSON数组。
	Status      *string         `json:"status"`       // 帖子新的状态（例如："draft" 草稿, "published" 已发布）。指针为nil表示不更新。
	Tags        []string        `json:"tags"`         // 要替换现有标签的新的标签名称列表。
}

// Service 定义了帖子业务逻辑操作的接口。
// 它抽象了数据操作和敏感内容过滤的实现细节。
type Service interface {
	// CreatePost 处理新帖子的创建，包括敏感词过滤和验证。
	CreatePost(dto *CreatePostDto) (*models.Post, error)
	// GetPost 从数据库中根据ID检索单个帖子。
	GetPost(id uint) (*models.Post, error)
	// UpdatePost 处理现有帖子的更新，应用敏感词过滤和部分更新。
	UpdatePost(id uint, dto *UpdatePostDto) (*models.Post, error)
	// DeletePost 处理根据ID对帖子进行软删除。
	DeletePost(id uint) error
	// ListPosts 检索与特定用户ID关联的分页帖子列表。
	ListPosts(userID uint, page, pageSize int) ([]*models.Post, int64, error)
}

// service 实现了 Service 接口，封装了业务规则并与仓库层交互。
type service struct {
	db     *gorm.DB          // db 是 GORM 数据库客户端实例，用于事务等操作。
	repo   Repository        // repo 是帖子的数据仓库，处理数据库交互。
	filter *sensitive.Filter // filter 用于文本内容中的敏感词检测和替换。
}

// NewService 创建一个新的帖子服务实例。
// 它使用给定的仓库和数据库连接初始化服务，并从 "dict.txt" 加载敏感词词典。
func NewService(db *gorm.DB, repo Repository) Service {
	filter := sensitive.New()
	err := filter.LoadWordDict("dict.txt") // 从应用程序根目录的文件中加载敏感词。
	if err != nil {
		log.Printf("Warning: 无法从 'dict.txt' 加载敏感词词典: %v", err)
	}
	return &service{
		db:     db,
		repo:   repo,
		filter: filter,
	}
}

// handleTagsInTx 在一个数据库事务中处理标签的查找和创建。
func (s *service) handleTagsInTx(tx *gorm.DB, tagNames []string) ([]*models.Tag, error) {
	var tags []*models.Tag
	if len(tagNames) == 0 {
		return tags, nil
	}

	for _, name := range tagNames {
		var tag models.Tag
		// 查找或创建标签，确保标签名的唯一性。
		err := tx.Clauses(clause.OnConflict{DoNothing: true}).FirstOrCreate(&tag, models.Tag{Name: name}).Error
		if err != nil {
			return nil, err // 如果发生错误，回滚事务。
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}

// CreatePost 处理创建新帖子的逻辑。
func (s *service) CreatePost(dto *CreatePostDto) (*models.Post, error) {
	// 验证图片数量（如果适用）。
	if dto.Type == "text_image" && dto.MediaURLs != nil {
		var urls []string
		if err := json.Unmarshal(dto.MediaURLs, &urls); err == nil {
			if len(urls) > 9 {
				return nil, errors.New("最多只能上传9张图片")
			}
		}
	}

	// 验证敏感词。
	found, _ := s.filter.Validate(dto.TextContent)
	if found {
		return nil, errors.New("帖子包含不允许的敏感内容")
	}
	filteredText := s.filter.Replace(dto.TextContent, '*')

	post := &models.Post{
		UserID:      dto.UserID,
		Type:        dto.Type,
		TextContent: filteredText,
		MediaURLs:   datatypes.JSON(dto.MediaURLs),
		Status:      "draft",
	}
	if dto.Status != "" {
		post.Status = dto.Status
	}

	// 在事务中处理帖子和标签的创建。
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 处理标签
		tags, err := s.handleTagsInTx(tx, dto.Tags)
		if err != nil {
			return err
		}
		post.Tags = tags

		// 2. 创建帖子
		if err := tx.Create(post).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(post.ID)
}

// GetPost 根据ID检索单个帖子。
func (s *service) GetPost(id uint) (*models.Post, error) {
	return s.repo.FindByID(id)
}

// UpdatePost 处理更新现有帖子的逻辑。
func (s *service) UpdatePost(id uint, dto *UpdatePostDto) (*models.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Apply updates from DTO. For TextContent, perform sensitive word checks.
	if dto.TextContent != nil {
		found, _ := s.filter.Validate(*dto.TextContent)
		if found {
			return nil, errors.New("更新的帖子内容包含不允许的敏感内容")
		}
		post.TextContent = s.filter.Replace(*dto.TextContent, '*')
	}
	if dto.MediaURLs != nil {
		if post.Type == "text_image" {
			var urls []string
			if err := json.Unmarshal(dto.MediaURLs, &urls); err == nil {
				if len(urls) > 9 {
					return nil, errors.New("最多只能上传9张图片")
				}
			}
		}
		post.MediaURLs = datatypes.JSON(dto.MediaURLs)
	}
	if dto.Status != nil {
		post.Status = *dto.Status
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 更新帖子基本信息
		if err := tx.Save(post).Error; err != nil {
			return err
		}

		// 如果 DTO 中提供了标签，则替换所有现有标签
		if dto.Tags != nil {
			tags, err := s.handleTagsInTx(tx, dto.Tags)
			if err != nil {
				return err
			}
			// Replace 会替换掉所有关联，并设置新的关联
			if err := tx.Model(post).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(post.ID)
}

// DeletePost 处理根据ID对帖子进行软删除。
func (s *service) DeletePost(id uint) error {
	return s.repo.Delete(id)
}

// ListPosts 检索特定用户的分页帖子列表。
func (s *service) ListPosts(userID uint, page, pageSize int) ([]*models.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.FindAllByUserID(userID, page, pageSize)
}
