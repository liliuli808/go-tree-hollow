-- Chat System Migration (PostgreSQL)
-- Creates tables for real-time messaging functionality

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    sender_id BIGINT NOT NULL,
    receiver_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    read_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for messages
CREATE INDEX IF NOT EXISTS idx_messages_sender_receiver ON messages(sender_id, receiver_id);
CREATE INDEX IF NOT EXISTS idx_messages_deleted_at ON messages(deleted_at);

-- Add foreign keys for messages
ALTER TABLE messages 
    ADD CONSTRAINT fk_messages_sender FOREIGN KEY (sender_id) REFERENCES users(id),
    ADD CONSTRAINT fk_messages_receiver FOREIGN KEY (receiver_id) REFERENCES users(id);

-- Conversations table
CREATE TABLE IF NOT EXISTS conversations (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    user1_id BIGINT NOT NULL,
    user2_id BIGINT NOT NULL,
    last_message_id BIGINT,
    last_message_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for conversations
CREATE INDEX IF NOT EXISTS idx_conversations_deleted_at ON conversations(deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_conversations_users ON conversations(user1_id, user2_id);

-- Add foreign keys for conversations
ALTER TABLE conversations 
    ADD CONSTRAINT fk_conversations_user1 FOREIGN KEY (user1_id) REFERENCES users(id),
    ADD CONSTRAINT fk_conversations_user2 FOREIGN KEY (user2_id) REFERENCES users(id),
    ADD CONSTRAINT fk_conversations_last_message FOREIGN KEY (last_message_id) REFERENCES messages(id);
