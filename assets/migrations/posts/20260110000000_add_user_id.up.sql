ALTER TABLE posts ADD COLUMN IF NOT EXISTS user_id UUID;
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);