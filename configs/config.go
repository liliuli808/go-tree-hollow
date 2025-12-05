package configs

import (
    "os"
    "strconv"
)

type Config struct {
    ServerPort   string
    DatabaseDSN  string
    JWTSecret    string
    JWTExpireDays int
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
    return &Config{
        ServerPort:    getEnv("SERVER_PORT", "8080"),
        DatabaseDSN:   getEnv("DATABASE_DSN", "test.db"),
        JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
        JWTExpireDays: getEnvAsInt("JWT_EXPIRE_DAYS", 7),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := os.Getenv(key)
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}