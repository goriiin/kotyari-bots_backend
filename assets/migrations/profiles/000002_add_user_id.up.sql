ALTER TABLE profiles ADD COLUMN IF NOT EXISTS user_id UUID;
CREATE INDEX IF NOT EXISTS idx_profiles_user_id ON profiles(user_id);