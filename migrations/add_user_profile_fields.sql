-- 为 users 表添加 nickname、avatar_url 和 background_url 字段

-- 添加昵称字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS nickname VARCHAR(50);

-- 添加头像URL字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(1024);

-- 添加背景URL字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS background_url VARCHAR(1024);

-- 可选：为现有用户设置默认昵称（从邮箱提取）
-- UPDATE users SET nickname = split_part(email, '@', 1) WHERE nickname IS NULL OR nickname = '';

-- 可选：添加索引以提高查询性能
-- CREATE INDEX IF NOT EXISTS idx_users_nickname ON users(nickname);
