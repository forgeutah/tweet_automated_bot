-- +goose Up
-- Multi-Platform Social Media Bot Schema
-- Phase 1: Initial Supabase Schema Creation

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =============================================================================
-- YOUTUBE CHANNELS TABLE
-- Stores RSS feed sources for video discovery
-- =============================================================================
CREATE TABLE IF NOT EXISTS youtube_channels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id VARCHAR(255) NOT NULL UNIQUE,
    channel_name VARCHAR(255),
    rss_feed_url TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_fetched_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX idx_youtube_channels_active ON youtube_channels(is_active) WHERE is_active = true;
CREATE INDEX idx_youtube_channels_last_fetched ON youtube_channels(last_fetched_at);

-- =============================================================================
-- VIDEOS TABLE
-- Stores all YouTube videos (migrated from yt_videos + RSS-discovered)
-- =============================================================================
CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    video_id VARCHAR(255) NOT NULL UNIQUE,  -- YouTube video ID (e.g., "1RYYsLy9bg8")
    title VARCHAR(500) NOT NULL,
    description TEXT,
    url TEXT NOT NULL,
    channel_id VARCHAR(255) NOT NULL,
    youtube_channel_id UUID REFERENCES youtube_channels(id) ON DELETE SET NULL,
    published_at TIMESTAMP WITH TIME ZONE,
    thumbnail_url TEXT,
    duration_seconds INTEGER,
    view_count INTEGER,
    -- Legacy fields from old schema
    playlist VARCHAR(255),
    conference_year VARCHAR(4),
    presenter_twitter VARCHAR(64),
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    -- Flexible metadata for platform-specific info
    metadata JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX idx_videos_video_id ON videos(video_id);
CREATE INDEX idx_videos_channel_id ON videos(channel_id);
CREATE INDEX idx_videos_youtube_channel_id ON videos(youtube_channel_id);
CREATE INDEX idx_videos_published_at ON videos(published_at DESC);

-- =============================================================================
-- PLATFORM ACCOUNTS TABLE
-- Configured accounts for each social media platform
-- =============================================================================
CREATE TYPE platform_type AS ENUM ('twitter', 'linkedin', 'discord', 'bluesky');

CREATE TABLE IF NOT EXISTS platform_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    platform platform_type NOT NULL,
    account_handle VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    -- Scheduling configuration
    posting_schedule_days INTEGER DEFAULT 7,  -- Days between posts
    cooldown_days INTEGER DEFAULT 90,  -- Days before reposting same video
    last_posted_at TIMESTAMP WITH TIME ZONE,
    next_post_due_at TIMESTAMP WITH TIME ZONE,
    -- Credentials reference (key to look up in credentials file)
    credential_key VARCHAR(255) NOT NULL,
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    -- Platform-specific configuration
    metadata JSONB DEFAULT '{}'::jsonb,
    -- Unique constraint on platform + handle
    CONSTRAINT unique_platform_account UNIQUE (platform, account_handle)
);

CREATE INDEX idx_platform_accounts_active ON platform_accounts(is_active) WHERE is_active = true;
CREATE INDEX idx_platform_accounts_next_due ON platform_accounts(next_post_due_at) WHERE is_active = true;
CREATE INDEX idx_platform_accounts_platform ON platform_accounts(platform);

-- =============================================================================
-- POST STATUS TYPE
-- Tracks the lifecycle of a post attempt
-- =============================================================================
CREATE TYPE post_status AS ENUM (
    'queued',           -- Waiting to be processed
    'generating',       -- LLM generating message
    'validating',       -- Platform validating message
    'posting',          -- Attempting to post
    'posted',           -- Successfully posted
    'failed_network',   -- Failed due to network/API error
    'failed_validation',-- Failed validation after max iterations
    'failed_permanent'  -- Permanent failure (dead letter queue)
);

-- =============================================================================
-- POSTS TABLE
-- Tracks individual posting attempts across all platforms
-- =============================================================================
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    platform_account_id UUID NOT NULL REFERENCES platform_accounts(id) ON DELETE CASCADE,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    status post_status NOT NULL DEFAULT 'queued',
    -- Message content
    message_text TEXT,
    message_metadata JSONB DEFAULT '{}'::jsonb,
    -- Platform response
    platform_post_id VARCHAR(255),  -- ID returned by platform (tweet ID, LinkedIn URN, etc.)
    platform_url TEXT,  -- Direct URL to the post
    -- Error tracking
    error_message TEXT,
    error_code VARCHAR(50),
    retry_count INTEGER DEFAULT 0,
    tier_2_iteration_count INTEGER DEFAULT 0,
    -- Timing
    queued_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    started_at TIMESTAMP WITH TIME ZONE,
    posted_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE,
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_platform_account ON posts(platform_account_id);
CREATE INDEX idx_posts_video_id ON posts(video_id);
CREATE INDEX idx_posts_posted_at ON posts(posted_at DESC);
CREATE INDEX idx_posts_queued_at ON posts(queued_at) WHERE status = 'queued';

-- =============================================================================
-- GENERATED MESSAGES TABLE
-- LLM generation history for debugging and analysis
-- =============================================================================
CREATE TABLE IF NOT EXISTS generated_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    platform platform_type NOT NULL,
    iteration_number INTEGER NOT NULL,
    -- Generation input
    prompt_template TEXT,
    prompt_variables JSONB DEFAULT '{}'::jsonb,
    -- Generation output
    message_text TEXT NOT NULL,
    token_count INTEGER,
    model_used VARCHAR(100),
    -- Validation result
    validation_status VARCHAR(50),  -- 'valid', 'invalid', 'not_validated'
    validation_errors JSONB DEFAULT '[]'::jsonb,
    character_count INTEGER,
    -- Metadata
    generation_time_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX idx_generated_messages_post_id ON generated_messages(post_id);
CREATE INDEX idx_generated_messages_iteration ON generated_messages(post_id, iteration_number);
CREATE INDEX idx_generated_messages_validation_status ON generated_messages(validation_status);

-- =============================================================================
-- POST HISTORY TABLE
-- Audit trail for successful posts (used for cooldown tracking)
-- =============================================================================
CREATE TABLE IF NOT EXISTS post_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    platform_account_id UUID NOT NULL REFERENCES platform_accounts(id) ON DELETE CASCADE,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    posted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    platform_post_id VARCHAR(255),
    platform_url TEXT,
    message_text TEXT,
    engagement_metrics JSONB DEFAULT '{}'::jsonb,  -- likes, shares, comments, etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_post_history_video_platform ON post_history(video_id, platform_account_id);
CREATE INDEX idx_post_history_posted_at ON post_history(posted_at DESC);
CREATE INDEX idx_post_history_platform_account ON post_history(platform_account_id, posted_at DESC);

-- =============================================================================
-- FUNCTIONS AND TRIGGERS
-- =============================================================================

-- Updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply updated_at trigger to relevant tables
CREATE TRIGGER update_youtube_channels_updated_at
    BEFORE UPDATE ON youtube_channels
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_videos_updated_at
    BEFORE UPDATE ON videos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_platform_accounts_updated_at
    BEFORE UPDATE ON platform_accounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_posts_updated_at
    BEFORE UPDATE ON posts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =============================================================================
-- VIEWS FOR COMMON QUERIES
-- =============================================================================

-- View: Videos eligible for posting to a specific account
CREATE OR REPLACE VIEW v_eligible_videos AS
SELECT
    v.id,
    v.video_id,
    v.title,
    v.url,
    v.channel_id,
    v.published_at,
    pa.id as platform_account_id,
    pa.platform,
    pa.account_handle,
    pa.cooldown_days,
    COALESCE(MAX(ph.posted_at), '1970-01-01'::timestamp) as last_posted_to_account,
    EXTRACT(EPOCH FROM (NOW() - COALESCE(MAX(ph.posted_at), '1970-01-01'::timestamp)))/86400 as days_since_last_post
FROM videos v
CROSS JOIN platform_accounts pa
LEFT JOIN post_history ph ON ph.video_id = v.id AND ph.platform_account_id = pa.id
WHERE pa.is_active = true
GROUP BY v.id, v.video_id, v.title, v.url, v.channel_id, v.published_at,
         pa.id, pa.platform, pa.account_handle, pa.cooldown_days
HAVING COALESCE(MAX(ph.posted_at), '1970-01-01'::timestamp) < (NOW() - (pa.cooldown_days || ' days')::interval);

-- View: Platform account posting status
CREATE OR REPLACE VIEW v_platform_account_status AS
SELECT
    pa.id,
    pa.platform,
    pa.account_handle,
    pa.is_active,
    pa.posting_schedule_days,
    pa.cooldown_days,
    pa.last_posted_at,
    pa.next_post_due_at,
    EXTRACT(EPOCH FROM (NOW() - COALESCE(pa.last_posted_at, '1970-01-01'::timestamp)))/86400 as days_since_last_post,
    COUNT(DISTINCT ph.video_id) as total_videos_posted,
    COUNT(DISTINCT CASE WHEN ph.posted_at > NOW() - interval '30 days' THEN ph.video_id END) as videos_posted_last_30_days
FROM platform_accounts pa
LEFT JOIN post_history ph ON ph.platform_account_id = pa.id
GROUP BY pa.id, pa.platform, pa.account_handle, pa.is_active,
         pa.posting_schedule_days, pa.cooldown_days, pa.last_posted_at, pa.next_post_due_at;

-- =============================================================================
-- COMMENTS
-- =============================================================================

COMMENT ON TABLE youtube_channels IS 'RSS feed sources for automatic video discovery';
COMMENT ON TABLE videos IS 'All YouTube videos from hardcoded GoWest Conference videos and RSS feeds';
COMMENT ON TABLE platform_accounts IS 'Social media accounts configured for automated posting';
COMMENT ON TABLE posts IS 'Individual post attempts with full retry and error tracking';
COMMENT ON TABLE generated_messages IS 'LLM message generation history for debugging';
COMMENT ON TABLE post_history IS 'Audit trail of successful posts for cooldown enforcement';

COMMENT ON COLUMN videos.video_id IS 'YouTube video ID extracted from URL (e.g., dQw4w9WgXcQ)';
COMMENT ON COLUMN platform_accounts.cooldown_days IS 'Days before same video can be reposted to this account';
COMMENT ON COLUMN platform_accounts.credential_key IS 'Key to lookup credentials in credentials.json file';
COMMENT ON COLUMN posts.retry_count IS 'Tier 1: Network/API retry attempts (max 3)';
COMMENT ON COLUMN posts.tier_2_iteration_count IS 'Tier 2: LLM regeneration iterations (max 10)';

-- +goose Down
-- Rollback migration: Drop all tables and types in reverse order

DROP VIEW IF EXISTS v_platform_account_status;
DROP VIEW IF EXISTS v_eligible_videos;

DROP TABLE IF EXISTS post_history CASCADE;
DROP TABLE IF EXISTS generated_messages CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS platform_accounts CASCADE;
DROP TABLE IF EXISTS videos CASCADE;
DROP TABLE IF EXISTS youtube_channels CASCADE;

DROP TYPE IF EXISTS post_status;
DROP TYPE IF EXISTS platform_type;

DROP FUNCTION IF EXISTS update_updated_at_column();
