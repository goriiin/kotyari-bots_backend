CREATE TABLE IF NOT EXISTS bots (
                                    id UUID PRIMARY KEY,
                                    bot_name TEXT NOT NULL,
                                    system_prompt TEXT,
                                    moderation_required BOOLEAN NOT NULL DEFAULT FALSE,
                                    auto_publish BOOLEAN NOT NULL DEFAULT FALSE,
                                    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
                                    profile_ids UUID[] NOT NULL DEFAULT '{}',
                                    profiles_count INTEGER NOT NULL DEFAULT 0,
                                    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bots_name ON bots(bot_name);
CREATE INDEX IF NOT EXISTS idx_bots_profile_ids ON bots USING GIN (profile_ids);