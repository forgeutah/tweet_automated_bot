-- +goose Up
-- Data Migration from CockroachDB to Supabase
-- This script migrates the 18 hardcoded GoWest Conference videos
-- and creates platform accounts for existing Twitter handles

-- =============================================================================
-- STEP 1: Insert YouTube Channel (GoWest Conference)
-- =============================================================================
INSERT INTO youtube_channels (
    channel_id,
    channel_name,
    rss_feed_url,
    is_active,
    metadata
) VALUES (
    'UC_placeholder_gowest',  -- Will be updated with actual channel ID
    'GoWest Conference',
    'https://www.youtube.com/feeds/videos.xml?channel_id=UC_placeholder_gowest',
    true,
    '{"original_source": "hardcoded_migration", "conference": "GoWest"}'::jsonb
) ON CONFLICT (channel_id) DO NOTHING;

-- =============================================================================
-- STEP 2: Insert Platform Accounts (Twitter handles from old system)
-- =============================================================================

-- GoWest Conference Twitter account
INSERT INTO platform_accounts (
    platform,
    account_handle,
    display_name,
    is_active,
    posting_schedule_days,
    cooldown_days,
    credential_key,
    metadata
) VALUES (
    'twitter',
    'gowestconf',
    'GoWest Conference',
    true,
    7,  -- Post every 7 days
    90,  -- 3-month cooldown (90 days)
    'twitter_gowestconf',
    '{"migrated_from": "cockroachdb", "original_twitter_username": "gowestconf"}'::jsonb
) ON CONFLICT (platform, account_handle) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    posting_schedule_days = EXCLUDED.posting_schedule_days,
    cooldown_days = EXCLUDED.cooldown_days;

-- ForgeUtah Bot Twitter account
INSERT INTO platform_accounts (
    platform,
    account_handle,
    display_name,
    is_active,
    posting_schedule_days,
    cooldown_days,
    credential_key,
    metadata
) VALUES (
    'twitter',
    'forgeutahbot',
    'Forge Utah Bot',
    true,
    7,  -- Post every 7 days
    90,  -- 3-month cooldown (90 days)
    'twitter_forgeutahbot',
    '{"migrated_from": "cockroachdb", "original_twitter_username": "forgeutahbot"}'::jsonb
) ON CONFLICT (platform, account_handle) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    posting_schedule_days = EXCLUDED.posting_schedule_days,
    cooldown_days = EXCLUDED.cooldown_days;

-- =============================================================================
-- STEP 3: Insert Videos (18 GoWest Conference videos)
-- Extract video_id from URL using regex: youtu.be/VIDEO_ID or youtube.com/watch?v=VIDEO_ID
-- =============================================================================

INSERT INTO videos (
    video_id,
    title,
    url,
    channel_id,
    youtube_channel_id,
    playlist,
    conference_year,
    presenter_twitter,
    published_at,
    metadata
) VALUES
    -- 2020 Videos
    ('1RYYsLy9bg8', 'All Types of Golang Types', 'https://youtu.be/1RYYsLy9bg8', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2020', '@carson_ops',
     '2020-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('ZtP-IggAlB8', 'Field Report: Building a game engine for 300 DEFCON hackers to smash', 'https://youtu.be/ZtP-IggAlB8', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2020', '@astockwell and @WarOnShrugs',
     '2020-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('CyhmhY7aI-s', 'The Standard Library Bootcamp', 'https://youtu.be/CyhmhY7aI-s', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2020', '@JeremyCMorgan',
     '2020-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('Fq5-KmNr_D8', 'Micro Machine Learning in Go', 'https://youtu.be/Fq5-KmNr_D8', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2020', 'Joshua Bowles',
     '2020-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('UqEtKx_9Wc8', 'Go Without Generics, a Retrospective', 'https://youtu.be/UqEtKx_9Wc8', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2020', '@lostluck',
     '2020-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('SSO78QkmMLs', 'Functional Programming with Go', 'https://youtu.be/SSO78QkmMLs', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2020', 'Dylan Meeus',
     '2020-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('U_qVSHYgVSE', 'So you think you know Go?', 'https://youtu.be/U_qVSHYgVSE', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2020', '@corylanou',
     '2020-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('ou_B5YZzEeU', 'Anatomy of a Gopher - Binary Analysis of Go Binaries', 'https://youtu.be/ou_B5YZzEeU', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2020', 'Alex Useche',
     '2020-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('M1AvdQVjhx8', 'Going Serverless', 'https://youtu.be/M1AvdQVjhx8', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2020', '@bogaczio',
     '2020-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('rbwvY0YpRPI', 'Writing REST Services for the gRPC-curious', 'https://youtu.be/rbwvY0YpRPI', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2020', '@JohanBrandhorst',
     '2020-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    -- 2021 Videos
    ('m8mKqdyD_C8', 'Go Templates: Great Library, Bad Rap', 'https://youtu.be/m8mKqdyD_C8', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2021', '@carson_ops',
     '2021-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('TxBz7L9ev18', 'Practical Tips to Creating a Great Engineering Culture', 'https://youtu.be/TxBz7L9ev18', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2021', '@clintberry',
     '2021-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('4LmTe0P7qpg', 'Callstacks at Giverny: A Go Graphics CLI', 'https://youtu.be/4LmTe0P7qpg', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2021', '@juliecoding',
     '2021-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('VESXrgCW50g', 'The beauty of Go for building cross-platform graphical applications', 'https://youtu.be/VESXrgCW50g', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2021', '@andydotxyz',
     '2021-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('jWUhpOnhOQ0', 'Building a scalable API platform for live streaming using GoLang', 'https://youtu.be/jWUhpOnhOQ0', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2021', 'Amit Mishra',
     '2021-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('KdHqx9Bx4jI', 'Go Channels Demystified', 'https://youtu.be/KdHqx9Bx4jI', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2021', '@moficodes',
     '2021-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('NvjvzYacgQg', 'Ent: Making Data Easy in Go', 'https://youtu.be/NvjvzYacgQg', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2021', '@DmitryVinnik',
     '2021-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb),

    ('PolzSYuwaQg', 'Build a home monitoring system with Go', 'https://youtu.be/PolzSYuwaQg', 'UC_placeholder_gowest',
     (SELECT id FROM youtube_channels WHERE channel_id = 'UC_placeholder_gowest'), 'GoWest Conference', '2021', '@dlsniper',
     '2021-01-01'::timestamp, '{"migrated_from": "cockroachdb"}'::jsonb)

ON CONFLICT (video_id) DO UPDATE SET
    title = EXCLUDED.title,
    url = EXCLUDED.url,
    conference_year = EXCLUDED.conference_year,
    presenter_twitter = EXCLUDED.presenter_twitter;

-- =============================================================================
-- STEP 4: Migrate last_sent_at data to post_history
-- NOTE: This requires data from CockroachDB yt_videos.last_sent_at field
-- This is a placeholder - actual migration requires querying CockroachDB
-- =============================================================================

-- Example of how to populate post_history if last_sent_at data exists
-- This would be executed by the Go migration tool after querying CockroachDB

-- INSERT INTO post_history (
--     platform_account_id,
--     video_id,
--     post_id,  -- Will need to create a placeholder post record
--     posted_at,
--     message_text
-- )
-- SELECT
--     pa.id,
--     v.id,
--     uuid_generate_v4(),  -- Placeholder post_id
--     -- [last_sent_at from CockroachDB],
--     'Migrated from legacy system'
-- FROM videos v
-- CROSS JOIN platform_accounts pa
-- WHERE v.metadata->>'migrated_from' = 'cockroachdb'
--   AND pa.metadata->>'original_twitter_username' IN ('gowestconf', 'forgeutahbot');

-- =============================================================================
-- VERIFICATION QUERIES
-- =============================================================================

-- Count migrated videos (should be 18)
-- SELECT COUNT(*) as migrated_video_count FROM videos WHERE metadata->>'migrated_from' = 'cockroachdb';

-- Count platform accounts (should be 2)
-- SELECT COUNT(*) as platform_account_count FROM platform_accounts WHERE metadata->>'migrated_from' = 'cockroachdb';

-- List all videos with their conference year
-- SELECT video_id, title, conference_year, presenter_twitter FROM videos ORDER BY conference_year, title;

-- List platform accounts
-- SELECT platform, account_handle, posting_schedule_days, cooldown_days FROM platform_accounts;

-- +goose Down
-- Rollback data migration: Remove migrated data

DELETE FROM post_history WHERE metadata->>'migrated_from' = 'cockroachdb';
DELETE FROM videos WHERE metadata->>'migrated_from' = 'cockroachdb';
DELETE FROM platform_accounts WHERE metadata->>'migrated_from' = 'cockroachdb';
DELETE FROM youtube_channels WHERE metadata->>'original_source' = 'hardcoded_migration';
