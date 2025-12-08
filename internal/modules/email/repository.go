package email

import (
	"context"
	"fmt"
	"go-tree-hollow/configs"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewClient(cfg *configs.RedisConfig) (*redis.Client, error) {
	fmt.Println(cfg.Password, 1)
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis 连接失败: %w", err)
	}

	return client, nil
}

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
