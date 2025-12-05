package database

import (
	"fmt"
	"go-tree-hollow/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewDB 创建数据库连接（支持SQLite和PostgreSQL）
func NewDB(dsn string) (*gorm.DB, error) {
    var db *gorm.DB
    var err error

    // 根据DSN前缀判断数据库类型
    if dsn[:6] == "postgres" {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    } else {
        db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
    }

    if err != nil {
        return nil, fmt.Errorf("failed to connect database: %w", err)
    }

    // 自动迁移
    err = db.AutoMigrate(&models.User{})
    if err != nil {
        return nil, fmt.Errorf("failed to migrate database: %w", err)
    }

    return db, nil
}