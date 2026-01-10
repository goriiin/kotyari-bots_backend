DROP INDEX IF EXISTS idx_profiles_user_id;
ALTER TABLE profiles DROP COLUMN IF EXISTS user_id;