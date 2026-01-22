package database

import (
	"context"
	"fmt"
)

// VideoRepository handles all video-related database operations
type VideoRepository struct {
	conn *SupabaseConnection
}

// NewVideoRepository creates a new VideoRepository
func NewVideoRepository(conn *SupabaseConnection) *VideoRepository {
	return &VideoRepository{conn: conn}
}

// GetByID retrieves a video by its UUID
func (r *VideoRepository) GetByID(ctx context.Context, id string) (*Video, error) {
	var video Video
	query := `
		SELECT id, video_id, title, description, url, channel_id, youtube_channel_id,
		       published_at, thumbnail_url, duration_seconds, view_count,
		       playlist, conference_year, presenter_twitter,
		       created_at, updated_at, metadata
		FROM videos
		WHERE id = $1`

	err := r.conn.DB.GetContext(ctx, &video, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get video by ID: %w", err)
	}
	return &video, nil
}

// GetByVideoID retrieves a video by its YouTube video ID
func (r *VideoRepository) GetByVideoID(ctx context.Context, videoID string) (*Video, error) {
	var video Video
	query := `
		SELECT id, video_id, title, description, url, channel_id, youtube_channel_id,
		       published_at, thumbnail_url, duration_seconds, view_count,
		       playlist, conference_year, presenter_twitter,
		       created_at, updated_at, metadata
		FROM videos
		WHERE video_id = $1`

	err := r.conn.DB.GetContext(ctx, &video, query, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get video by video_id: %w", err)
	}
	return &video, nil
}

// ListAll retrieves all videos
func (r *VideoRepository) ListAll(ctx context.Context) ([]Video, error) {
	var videos []Video
	query := `
		SELECT id, video_id, title, description, url, channel_id, youtube_channel_id,
		       published_at, thumbnail_url, duration_seconds, view_count,
		       playlist, conference_year, presenter_twitter,
		       created_at, updated_at, metadata
		FROM videos
		ORDER BY published_at DESC, title ASC`

	err := r.conn.DB.SelectContext(ctx, &videos, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list all videos: %w", err)
	}
	return videos, nil
}

// ListByChannel retrieves all videos for a specific YouTube channel
func (r *VideoRepository) ListByChannel(ctx context.Context, channelID string) ([]Video, error) {
	var videos []Video
	query := `
		SELECT id, video_id, title, description, url, channel_id, youtube_channel_id,
		       published_at, thumbnail_url, duration_seconds, view_count,
		       playlist, conference_year, presenter_twitter,
		       created_at, updated_at, metadata
		FROM videos
		WHERE channel_id = $1
		ORDER BY published_at DESC`

	err := r.conn.DB.SelectContext(ctx, &videos, query, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to list videos by channel: %w", err)
	}
	return videos, nil
}

// GetEligibleVideos retrieves videos eligible for posting to a specific platform account
// A video is eligible if it hasn't been posted to this account within the cooldown period
func (r *VideoRepository) GetEligibleVideos(ctx context.Context, platformAccountID string) ([]EligibleVideo, error) {
	var videos []EligibleVideo
	query := `
		SELECT *
		FROM v_eligible_videos
		WHERE platform_account_id = $1
		ORDER BY RANDOM()`

	err := r.conn.DB.SelectContext(ctx, &videos, query, platformAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get eligible videos: %w", err)
	}
	return videos, nil
}

// GetRandomEligibleVideo retrieves a random video eligible for posting
func (r *VideoRepository) GetRandomEligibleVideo(ctx context.Context, platformAccountID string) (*EligibleVideo, error) {
	videos, err := r.GetEligibleVideos(ctx, platformAccountID)
	if err != nil {
		return nil, err
	}

	if len(videos) == 0 {
		return nil, nil // No eligible videos
	}

	return &videos[0], nil // Already randomized by query
}

// Create inserts a new video
func (r *VideoRepository) Create(ctx context.Context, video *Video) error {
	query := `
		INSERT INTO videos (
			video_id, title, description, url, channel_id, youtube_channel_id,
			published_at, thumbnail_url, duration_seconds, view_count,
			playlist, conference_year, presenter_twitter, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		) RETURNING id, created_at, updated_at`

	err := r.conn.DB.QueryRowContext(ctx, query,
		video.VideoID,
		video.Title,
		video.Description,
		video.URL,
		video.ChannelID,
		video.YouTubeChannelID,
		video.PublishedAt,
		video.ThumbnailURL,
		video.DurationSeconds,
		video.ViewCount,
		video.Playlist,
		video.ConferenceYear,
		video.PresenterTwitter,
		video.Metadata,
	).Scan(&video.ID, &video.CreatedAt, &video.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create video: %w", err)
	}
	return nil
}

// Update updates an existing video
func (r *VideoRepository) Update(ctx context.Context, video *Video) error {
	query := `
		UPDATE videos SET
			title = $1,
			description = $2,
			url = $3,
			channel_id = $4,
			youtube_channel_id = $5,
			published_at = $6,
			thumbnail_url = $7,
			duration_seconds = $8,
			view_count = $9,
			playlist = $10,
			conference_year = $11,
			presenter_twitter = $12,
			metadata = $13,
			updated_at = NOW()
		WHERE id = $14
		RETURNING updated_at`

	err := r.conn.DB.QueryRowContext(ctx, query,
		video.Title,
		video.Description,
		video.URL,
		video.ChannelID,
		video.YouTubeChannelID,
		video.PublishedAt,
		video.ThumbnailURL,
		video.DurationSeconds,
		video.ViewCount,
		video.Playlist,
		video.ConferenceYear,
		video.PresenterTwitter,
		video.Metadata,
		video.ID,
	).Scan(&video.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}
	return nil
}

// VideoExists checks if a video with the given video_id already exists
func (r *VideoRepository) VideoExists(ctx context.Context, videoID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM videos WHERE video_id = $1)`

	err := r.conn.DB.GetContext(ctx, &exists, query, videoID)
	if err != nil {
		return false, fmt.Errorf("failed to check if video exists: %w", err)
	}
	return exists, nil
}

// Count returns the total number of videos
func (r *VideoRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM videos`

	err := r.conn.DB.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count videos: %w", err)
	}
	return count, nil
}
