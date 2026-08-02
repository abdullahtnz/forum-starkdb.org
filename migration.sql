-- ============================================================
-- STARKDB FORUM — PostgreSQL Migration
-- Run: psql -U your_user -d your_db -f migration.sql
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- ============================================================
-- USERS
-- ============================================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(30) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    bio TEXT DEFAULT '',
    avatar_url TEXT DEFAULT '',
    is_admin BOOLEAN DEFAULT FALSE,
    accepted_terms BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login TIMESTAMPTZ
);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- ============================================================
-- REFRESH TOKENS
-- ============================================================
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at);

-- ============================================================
-- TAGS
-- ============================================================
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL UNIQUE,
    slug VARCHAR(50) NOT NULL UNIQUE,
    description TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- POSTS
-- ============================================================
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(300) NOT NULL,
    content TEXT NOT NULL,
    view_count INT DEFAULT 0,
    is_pinned BOOLEAN DEFAULT FALSE,
    is_closed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_posts_user ON posts(user_id);
CREATE INDEX idx_posts_created ON posts(created_at DESC);
CREATE INDEX idx_posts_title_trgm ON posts USING GIN (title gin_trgm_ops);

-- ============================================================
-- POST TAGS (junction)
-- ============================================================
CREATE TABLE post_tags (
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);
CREATE INDEX idx_post_tags_tag ON post_tags(tag_id);

-- ============================================================
-- POST ATTACHMENTS (files/images)
-- ============================================================
CREATE TABLE post_attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_attachments_post ON post_attachments(post_id);

-- ============================================================
-- COMMENTS (unlimited nesting via parent_id)
-- ============================================================
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_comments_post ON comments(post_id, created_at ASC);
CREATE INDEX idx_comments_parent ON comments(parent_id);
CREATE INDEX idx_comments_user ON comments(user_id);

-- ============================================================
-- POST LIKES
-- ============================================================
CREATE TABLE post_likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, post_id)
);
CREATE INDEX idx_post_likes_post ON post_likes(post_id);

-- ============================================================
-- COMMENT LIKES
-- ============================================================
CREATE TABLE comment_likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    comment_id UUID NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, comment_id)
);
CREATE INDEX idx_comment_likes_comment ON comment_likes(comment_id);

-- ============================================================
-- BAD WORDS (content filtering)
-- ============================================================
CREATE TABLE bad_words (
    id SERIAL PRIMARY KEY,
    word VARCHAR(100) NOT NULL UNIQUE
);

-- ============================================================
-- FULL-TEXT SEARCH ON POSTS
-- ============================================================
ALTER TABLE posts ADD COLUMN IF NOT EXISTS search_vector tsvector
    GENERATED ALWAYS AS (
        setweight(to_tsvector('english', COALESCE(title, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(content, '')), 'B')
    ) STORED;
CREATE INDEX IF NOT EXISTS idx_posts_search ON posts USING GIN (search_vector);

-- ============================================================
-- SEED TAGS
-- ============================================================
INSERT INTO tags (name, slug, description) VALUES
    ('b-tree', 'b-tree', 'B-Tree structure, indexing, and traversal'),
    ('cli', 'cli', 'Command-line interface usage'),
    ('cpp', 'cpp', 'C++ API and integration'),
    ('c-api', 'c-api', 'C API and integration'),
    ('installation', 'installation', 'Installation and setup'),
    ('performance', 'performance', 'Performance tuning and benchmarking'),
    ('embedded', 'embedded', 'Embedded database usage'),
    ('key-value', 'key-value', 'Key-value operations'),
    ('typed-records', 'typed-records', 'Typed record operations'),
    ('transactions', 'transactions', 'ACID transactions'),
    ('mobile', 'mobile', 'Mobile app integration'),
    ('game-dev', 'game-dev', 'Game development with STARKDB')
ON CONFLICT (slug) DO NOTHING;
