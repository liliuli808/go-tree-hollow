package post

import (
	"go-tree-hollow/internal/models"
	"gorm.io/gorm"
)

// Repository 定义了帖子数据操作的接口，抽象了数据库交互。
type Repository interface {
	// Create 在数据库中持久化一条新的帖子记录。
	Create(post *models.Post) error
	// FindByID 从数据库中根据ID检索帖子，并预加载关联的用户信息。
	FindByID(id uint) (*models.Post, error)
	// Update 保存数据库中现有帖子记录的更改。
	Update(post *models.Post) error
	// Delete 通过设置 'deleted_at' 时间戳将帖子标记为删除（软删除）。
	Delete(id uint) error
	// FindAllByUserID 检索特定用户的分页帖子列表。
	FindAllByUserID(userID uint, page, pageSize int) ([]*models.Post, int64, error)
}

// repository 使用 GORM 实现了 Repository 接口。
type repository struct {
	db *gorm.DB // db 表示 GORM 数据库客户端实例。
}

// NewRepository 创建一个新的帖子仓库实例。
// 它接收一个 *gorm.DB 实例用于数据库交互。
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create 在数据库中创建一条新的帖子记录。
// 它接收一个指向 models.Post 结构体的指针并将其持久化。
func (r *repository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

// FindByID 从数据库中根据ID检索帖子。
// 它预加载关联的 User 模型以避免 N+1 查询问题。
// 如果找到则返回 models.Post，否则返回错误。
func (r *repository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("User").Preload("Tags").First(&post, id).Error
	return &post, err
}

// Update 保存数据库中现有帖子记录的更改。
// 它接收一个带有更新字段的 models.Post 结构体指针。
// GORM 将更新所有非零值或明确标记的字段。
func (r *repository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

// Delete 通过设置 'deleted_at' 时间戳将帖子标记为删除（软删除）。
// 它接收要删除的帖子的ID。
func (r *repository) Delete(id uint) error {
	// GORM 的 Delete 方法，如果模型包含 gorm.DeletedAt，则对结构体和ID执行软删除。
	return r.db.Delete(&models.Post{}, id).Error
}

// FindAllByUserID 检索特定用户的分页帖子列表。
// 它接收用户ID、页码和页面大小，返回帖子切片以及该用户帖子的总数和遇到的任何错误。
func (r *repository) FindAllByUserID(userID uint, page, pageSize int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	// 构建按用户ID过滤帖子的基本查询。
	query := r.db.Model(&models.Post{}).Where("user_id = ?", userID)

	// 获取与查询匹配的帖子总数，用于分页元数据。
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 计算分页的偏移量。
	offset := (page - 1) * pageSize
	// 执行分页查询，预加载 User 并按创建日期排序。
	err := query.Preload("User").Preload("Tags").Offset(offset).Limit(pageSize).Order("created_at desc").Find(&posts).Error

	return posts, total, err
}
