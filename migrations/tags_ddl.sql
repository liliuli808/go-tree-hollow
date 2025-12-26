-- Tags table DDL for PostgreSQL
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(50) NOT NULL UNIQUE
);

-- No junction table needed for one-to-one relationship
-- The tag_id foreign key will be added directly to the posts table
-- Add this column if the posts table already exists:
-- ALTER TABLE posts ADD COLUMN tag_id INTEGER REFERENCES tags(id) ON DELETE SET NULL;
-- CREATE INDEX IF NOT EXISTS idx_posts_tag_id ON posts(tag_id);
