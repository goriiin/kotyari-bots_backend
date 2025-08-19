CREATE TABLE IF NOT EXISTS integrations (
    "id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    "provider" text NOT NULL,
    "url" text NOT NULL
)