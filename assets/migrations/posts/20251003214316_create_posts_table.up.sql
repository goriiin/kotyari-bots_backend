CREATE TYPE post_type_enum as ENUM (
    'opinion',
    'knowledge',
    'history'
);

CREATE TYPE platform_type_enum as ENUM (
    'otveti'
);

CREATE TABLE IF NOT EXISTS categories(
    "id" UUID NOT NULL PRIMARY KEY,
    "category_name" TEXT NOT NULL UNIQUE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts(
    "id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "bot_id" UUID NOT NULL,
    "profile_id" UUID NOT NULL,
    "platform_type" platform_type_enum NOT NULL,
    "post_type" post_type_enum,
    "post_title" TEXT NOT NULL,
    "post_text" TEXT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS post_categories(
    "post_id" BIGINT NOT NULL,
    "category_id" UUID NOT NULL,
    PRIMARY KEY ("post_id", "category_id"),
    FOREIGN KEY("post_id") REFERENCES posts("id") ON DELETE CASCADE,
    FOREIGN KEY("category_id") REFERENCES categories("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

