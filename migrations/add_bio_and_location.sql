-- Add bio and location fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS location VARCHAR(100);
