CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE story_status AS ENUM ('publish', 'draft');

-- Master tables
CREATE TABLE IF NOT EXISTS categories (
    id         UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(100) NOT NULL UNIQUE,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tags (
    id         UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(100) NOT NULL UNIQUE,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Stories
CREATE TABLE IF NOT EXISTS stories (
    id          UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID         REFERENCES categories(id) ON DELETE SET NULL,
    title       VARCHAR(500) NOT NULL,
    author      VARCHAR(255) NOT NULL,
    synopsis    TEXT,
    cover_url   TEXT,
    status      story_status NOT NULL DEFAULT 'draft',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Chapters
CREATE TABLE IF NOT EXISTS chapters (
    id          UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    story_id    UUID         NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
    title       VARCHAR(500) NOT NULL,
    content     TEXT,
    order_index INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Junction
CREATE TABLE IF NOT EXISTS story_tags (
    story_id UUID NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
    tag_id   UUID NOT NULL REFERENCES tags(id)    ON DELETE CASCADE,
    PRIMARY KEY (story_id, tag_id)
);

-- Index
CREATE INDEX IF NOT EXISTS idx_stories_category_id ON stories(category_id);
CREATE INDEX IF NOT EXISTS idx_stories_status       ON stories(status);
CREATE INDEX IF NOT EXISTS idx_story_tags_story     ON story_tags(story_id);
CREATE INDEX IF NOT EXISTS idx_story_tags_tag       ON story_tags(tag_id);
CREATE INDEX IF NOT EXISTS idx_chapters_story       ON chapters(story_id);

-- Seed
INSERT INTO categories (name, slug) VALUES
    ('Financial',  'financial'),
    ('Technology', 'technology'),
    ('Health',     'health')
ON CONFLICT (slug) DO NOTHING;