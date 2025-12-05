package utils

import (
	"fmt"
	"sync" // 必须引入 sync 包
	"time"

	"go-tree-hollow/configs"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT载荷结构
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// 定义全局变量，确保在包级别可见
var (
	jwtSecret     []byte
	jwtExpireDays int
	once          sync.Once
)

// setupConfig 内部辅助函数：确保配置只被安全加载一次
// GenerateToken 和 ParseToken 都要调用它，防止谁先谁后的问题
func setupConfig() {
	once.Do(func() {
		config := configs.LoadConfig()
		jwtSecret = []byte(config.JWTSecret)
		jwtExpireDays = config.JWTExpireDays
	})
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, email string) (string, error) {
	// 1. 初始化配置 (只会执行一次)
	setupConfig()

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			// 2. 这里必须使用全局变量 jwtExpireDays，而不是局部的 config
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * time.Duration(jwtExpireDays))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	// 3. 解析时也要调用 setupConfig，替代原来非线程安全的 if jwtSecret == nil
	setupConfig()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
