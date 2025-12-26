-- PostgreSQL DDL for the 'posts' table

-- This DDL is based on the GORM model defined in internal/models/post.go
-- It includes the table creation, foreign key constraints, and indexes.

-- Assumes the 'users' table already exists.

CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    user_id BIGINT NOT NULL,
    type VARCHAR(255) NOT NULL,
    tag VARCHAR(255) NOT NULL,
    text_content TEXT,
    media_urls JSONB,
    status VARCHAR(255) NOT NULL DEFAULT 'draft',
    CONSTRAINT fk_posts_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);
CREATE INDEX IF NOT EXISTS idx_posts_type ON posts (type);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts (status);

-- GORM uses this index for soft deletes
CREATE INDEX IF NOT EXISTS idx_posts_deleted_at ON posts (deleted_at);


ALTER TABLE posts ADD COLUMN cover_url VARCHAR(1024) DEFAULT NULL;