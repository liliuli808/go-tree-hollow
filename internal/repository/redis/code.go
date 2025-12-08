// internal/repository/redis/code.go
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type CodeRepository struct {
	client *redis.Client
	prefix string
}

func NewCodeRepository(client *redis.Client, prefix string) *CodeRepository {
	return &CodeRepository{
		client: client,
		prefix: prefix,
	}
}

// Set 存储验证码
func (r *CodeRepository) Set(ctx context.Context, email, code string, expire time.Duration) error {
	key := r.buildKey(email)
	return r.client.Set(ctx, key, code, expire).Err()
}

// Get 获取验证码
func (r *CodeRepository) Get(ctx context.Context, email string) (string, error) {
	key := r.buildKey(email)
	code, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("验证码已过期或不存在")
	}
	if err != nil {
		return "", fmt.Errorf("获取验证码失败: %w", err)
	}
	return code, nil
}

// Delete 删除验证码
func (r *CodeRepository) Delete(ctx context.Context, email string) error {
	key := r.buildKey(email)
	return r.client.Del(ctx, key).Err()
}

// Exists 检查是否存在
func (r *CodeRepository) Exists(ctx context.Context, email string) bool {
	key := r.buildKey(email)
	return r.client.Exists(ctx, key).Val() > 0
}

// SetNX 设置防重发标记
func (r *CodeRepository) SetNX(ctx context.Context, email string, value string, expire time.Duration) (bool, error) {
	key := fmt.Sprintf("%s:lock:%s", r.prefix, email)
	return r.client.SetNX(ctx, key, value, expire).Result()
}

func (r *CodeRepository) buildKey(email string) string {
	return fmt.Sprintf("%s:code:%s", r.prefix, email)
}
