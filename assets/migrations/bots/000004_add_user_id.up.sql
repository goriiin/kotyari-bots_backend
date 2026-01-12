ALTER TABLE bots ADD COLUMN IF NOT EXISTS user_id UUID;
CREATE INDEX IF NOT EXISTS idx_bots_user_id ON bots(user_id);