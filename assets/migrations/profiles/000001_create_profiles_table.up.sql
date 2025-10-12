CREATE TABLE IF NOT EXISTS profiles (
                                        id UUID PRIMARY KEY,
                                        name TEXT NOT NULL,
                                        email TEXT NOT NULL,
                                        system_prompt TEXT,
                                        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_profiles_email ON profiles(email);