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
	UserID      uint     `json:"user_id" binding:"required"`
	TextContent string   `json:"text_content"`
	Images      []string `json:"images"`
	Video       string   `json:"video"`
	Audio       string   `json:"audio"`
	Status      string   `json:"status"`
	TagID       *uint    `json:"tag_id"` // Single tag ID for one-to-one relationship
}

// UpdatePostDto 定义了更新现有帖子的数据结构。
// It uses pointers to represent optional fields, allowing for partial updates without overwriting existing data with zero values.
type UpdatePostDto struct {
	TextContent *string  `json:"text_content"`
	Images      []string `json:"images"`
	Video       *string  `json:"video"`
	Audio       *string  `json:"audio"`
	Status      *string  `json:"status"`
	TagID       *uint    `json:"tag_id"` // Single tag ID
}

// Service defines the interface for post business logic operations.
// It abstracts the implementation details of data manipulation and sensitive content filtering.
type Service interface {
	// CreatePost handles the creation of a new post, including sensitive word filtering and validation.
	CreatePost(dto *CreatePostDto) (*models.Post, error)
	// GetPost retrieves a single post by its ID from the database.
	GetPost(id uint) (*models.Post, error)
	// UpdatePost handles updates to an existing post, applying sensitive word filtering and partial updates.
	UpdatePost(id uint, dto *UpdatePostDto) (*models.Post, error)
	// DeletePost handles the soft deletion of a post by its ID.
	DeletePost(id uint) error
	// ListPosts retrieves a paginated list of posts associated with a specific user ID.
	ListPosts(userID uint, page, pageSize int) ([]*models.Post, int64, error)
}

// service implements the Service interface, encapsulating business rules and interacting with the repository layer.
type service struct {
	db     *gorm.DB
	repo   Repository
	filter *sensitive.Filter
}

// NewService creates a new post service instance.
func NewService(db *gorm.DB, repo Repository) Service {
	filter := sensitive.New()
	err := filter.LoadWordDict("dict.txt")
	if err != nil {
		log.Printf("Warning: Failed to load sensitive word dictionary from 'dict.txt': %v", err)
	}
	return &service{
		db:     db,
		repo:   repo,
		filter: filter,
	}
}

// handleTagsInTx manages the finding and creation of tags within a database transaction.
func (s *service) handleTagsInTx(tx *gorm.DB, tagNames []string) ([]*models.Tag, error) {
	var tags []*models.Tag
	if len(tagNames) == 0 {
		return tags, nil
	}

	for _, name := range tagNames {
		var tag models.Tag
		err := tx.Clauses(clause.OnConflict{DoNothing: true}).FirstOrCreate(&tag, models.Tag{Name: name}).Error
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}

// CreatePost handles the logic for creating a new post.
func (s *service) CreatePost(dto *CreatePostDto) (*models.Post, error) {
	// Determine post type and media URLs
	var postType string
	var mediaUrls []string
	if len(dto.Images) > 0 {
		postType = "text_image"
		mediaUrls = dto.Images
		if len(mediaUrls) > 9 {
			return nil, errors.New("最多只能上传9张图片")
		}
	} else if dto.Video != "" {
		postType = "video"
		mediaUrls = []string{dto.Video}
	} else if dto.Audio != "" {
		postType = "audio"
		mediaUrls = []string{dto.Audio}
	} else {
		postType = "text"
	}

	// Sensitive word validation
	filteredText := s.filter.Replace(dto.TextContent, '*')

	mediaUrlsJSON, err := json.Marshal(mediaUrls)
	if err != nil {
		return nil, err
	}

	post := &models.Post{
		UserID:      dto.UserID,
		Type:        postType,
		TextContent: filteredText,
		MediaURLs:   datatypes.JSON(mediaUrlsJSON),
		Status:      "draft",
	}
	if dto.Status != "" {
		post.Status = dto.Status
	}

	// Create post with tag reference
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Assign tag ID directly if provided
		if dto.TagID != nil {
			post.TagID = dto.TagID
		}

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

// GetPost retrieves a single post by ID.
func (s *service) GetPost(id uint) (*models.Post, error) {
	return s.repo.FindByID(id)
}

// UpdatePost handles updating an existing post.
func (s *service) UpdatePost(id uint, dto *UpdatePostDto) (*models.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Apply updates from DTO
	if dto.TextContent != nil {
		found, _ := s.filter.Validate(*dto.TextContent)
		if found {
			return nil, errors.New("更新的帖子内容包含不允许的敏感内容")
		}
		post.TextContent = s.filter.Replace(*dto.TextContent, '*')
	}

	var mediaUrls []string
	if dto.Images != nil {
		post.Type = "text_image"
		mediaUrls = dto.Images
		if len(mediaUrls) > 9 {
			return nil, errors.New("最多只能上传9张图片")
		}
	} else if dto.Video != nil {
		post.Type = "video"
		mediaUrls = []string{*dto.Video}
	} else if dto.Audio != nil {
		post.Type = "audio"
		mediaUrls = []string{*dto.Audio}
	}
	if mediaUrls != nil {
		mediaUrlsJSON, err := json.Marshal(mediaUrls)
		if err != nil {
			return nil, err
		}
		post.MediaURLs = datatypes.JSON(mediaUrlsJSON)
	}

	if dto.Status != nil {
		post.Status = *dto.Status
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Update tag ID if provided
		if dto.TagID != nil {
			post.TagID = dto.TagID
		}

		if err := tx.Save(post).Error; err != nil {
			return err
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
