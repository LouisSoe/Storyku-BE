CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE story_category AS ENUM ('Financial', 'Technology', 'Health');
CREATE TYPE story_status    AS ENUM ('publish', 'draft');

CREATE TABLE IF NOT EXISTS stories (
    story_id   UUID           PRIMARY KEY DEFAULT uuid_generate_v4(),
    title      VARCHAR(500)   NOT NULL,
    author     VARCHAR(255)   NOT NULL,
    synopsis   TEXT           NOT NULL DEFAULT '',
    category   story_category NOT NULL,
    cover_url  TEXT           NOT NULL DEFAULT '',
    tags       TEXT[]         NOT NULL DEFAULT '{}',
    status     story_status   NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chapters (
    chapter_id  UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    story_id    UUID         NOT NULL REFERENCES stories(story_id) ON DELETE CASCADE,
    title       VARCHAR(500) NOT NULL,
    content     TEXT         NOT NULL DEFAULT '',
    order_index INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_stories_author   ON stories(author);
CREATE INDEX IF NOT EXISTS idx_stories_category ON stories(category);
CREATE INDEX IF NOT EXISTS idx_stories_status   ON stories(status);
CREATE INDEX IF NOT EXISTS idx_chapters_story   ON chapters(story_id, order_index);