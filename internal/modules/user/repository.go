package user

import (
	"go-tree-hollow/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetByID 根据ID获取用户
func (r *Repository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// Update 更新用户信息
func (r *Repository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// GetFollowCount 获取关注数
func (r *Repository) GetFollowCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Follow{}).Where("follower_id = ?", userID).Count(&count).Error
	return count, err
}

// GetFanCount 获取粉丝数
func (r *Repository) GetFanCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Follow{}).Where("followed_id = ?", userID).Count(&count).Error
	return count, err
}

// GetReceivedLikeCount 获取收到的赞数 (用户发布的所有帖子的获赞总和)
func (r *Repository) GetReceivedLikeCount(userID uint) (int64, error) {
	var count int64
	// 关联 posts 和 likes 表
	// likes 表有 post_id, posts 表有 id 和 user_id
	// 统计所有 posts.user_id = userID 的 likes 数量
	err := r.db.Model(&models.Like{}).
		Joins("JOIN posts ON likes.post_id = posts.id").
		Where("posts.user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// GetReceivedCollectionCount 获取收到的收藏数 (用户发布的所有帖子的被收藏总和)
func (r *Repository) GetReceivedCollectionCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Collection{}).
		Joins("JOIN posts ON collections.post_id = posts.id").
		Where("posts.user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// GetUserPosts 获取用户发布的帖子列表
func (r *Repository) GetUserPosts(userID uint) ([]models.Post, error) {
	var posts []models.Post
	// 预加载 User 和 Tag, 以及计算 LikesCount
	// 注意：Post 结构体重 LikesCount 是 gorm:"-"，需要手动计算或者 SQL 映射
	// 这里简化，只返回帖子基本信息，如果需要 LikesCount 可能需要额外处理
	// 或者在 Post model 中如果 LikesCount 是数据库字段则直接取，如果是 gorm:"-" 则需要 service 层处理
	// Post model definition shows LikesCount is -, let's populate it via subquery or separate call if needed.
	// For now, standard find. Service can enhance it if needed.
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Preload("User").Preload("Tag").Find(&posts).Error
	return posts, err
}
