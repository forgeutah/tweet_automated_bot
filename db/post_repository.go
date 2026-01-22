package database

import (
	"context"
	"fmt"
	"time"
)

// PostRepository handles all post-related database operations
type PostRepository struct {
	conn *SupabaseConnection
}

// NewPostRepository creates a new PostRepository
func NewPostRepository(conn *SupabaseConnection) *PostRepository {
	return &PostRepository{conn: conn}
}

// GetByID retrieves a post by its UUID
func (r *PostRepository) GetByID(ctx context.Context, id string) (*Post, error) {
	var post Post
	query := `
		SELECT id, platform_account_id, video_id, status, message_text, message_metadata,
		       platform_post_id, platform_url, error_message, error_code,
		       retry_count, tier_2_iteration_count,
		       queued_at, started_at, posted_at, failed_at, created_at, updated_at
		FROM posts
		WHERE id = $1`

	err := r.conn.DB.GetContext(ctx, &post, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}
	return &post, nil
}

// ListByStatus retrieves all posts with a specific status
func (r *PostRepository) ListByStatus(ctx context.Context, status PostStatus) ([]Post, error) {
	var posts []Post
	query := `
		SELECT id, platform_account_id, video_id, status, message_text, message_metadata,
		       platform_post_id, platform_url, error_message, error_code,
		       retry_count, tier_2_iteration_count,
		       queued_at, started_at, posted_at, failed_at, created_at, updated_at
		FROM posts
		WHERE status = $1
		ORDER BY queued_at ASC`

	err := r.conn.DB.SelectContext(ctx, &posts, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts by status: %w", err)
	}
	return posts, nil
}

// GetQueuedPosts retrieves all posts waiting to be processed
func (r *PostRepository) GetQueuedPosts(ctx context.Context, limit int) ([]Post, error) {
	var posts []Post
	query := `
		SELECT id, platform_account_id, video_id, status, message_text, message_metadata,
		       platform_post_id, platform_url, error_message, error_code,
		       retry_count, tier_2_iteration_count,
		       queued_at, started_at, posted_at, failed_at, created_at, updated_at
		FROM posts
		WHERE status = 'queued'
		ORDER BY queued_at ASC
		LIMIT $1`

	err := r.conn.DB.SelectContext(ctx, &posts, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get queued posts: %w", err)
	}
	return posts, nil
}

// Create inserts a new post
func (r *PostRepository) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (
			platform_account_id, video_id, status, message_text, message_metadata,
			platform_post_id, platform_url, error_message, error_code,
			retry_count, tier_2_iteration_count
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING id, queued_at, created_at, updated_at`

	err := r.conn.DB.QueryRowContext(ctx, query,
		post.PlatformAccountID,
		post.VideoID,
		post.Status,
		post.MessageText,
		post.MessageMetadata,
		post.PlatformPostID,
		post.PlatformURL,
		post.ErrorMessage,
		post.ErrorCode,
		post.RetryCount,
		post.Tier2IterationCount,
	).Scan(&post.ID, &post.QueuedAt, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

// UpdateStatus updates the status of a post
func (r *PostRepository) UpdateStatus(ctx context.Context, postID string, status PostStatus) error {
	query := `
		UPDATE posts SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2`

	_, err := r.conn.DB.ExecContext(ctx, query, status, postID)
	if err != nil {
		return fmt.Errorf("failed to update post status: %w", err)
	}
	return nil
}

// UpdateWithMessage updates a post with generated message text
func (r *PostRepository) UpdateWithMessage(ctx context.Context, postID string, messageText string, metadata JSONB) error {
	query := `
		UPDATE posts SET
			message_text = $1,
			message_metadata = $2,
			updated_at = NOW()
		WHERE id = $3`

	_, err := r.conn.DB.ExecContext(ctx, query, messageText, metadata, postID)
	if err != nil {
		return fmt.Errorf("failed to update post with message: %w", err)
	}
	return nil
}

// MarkAsPosted marks a post as successfully posted
func (r *PostRepository) MarkAsPosted(ctx context.Context, postID string, platformPostID string, platformURL string) error {
	now := time.Now()
	query := `
		UPDATE posts SET
			status = 'posted',
			platform_post_id = $1,
			platform_url = $2,
			posted_at = $3,
			updated_at = NOW()
		WHERE id = $4`

	_, err := r.conn.DB.ExecContext(ctx, query, platformPostID, platformURL, now, postID)
	if err != nil {
		return fmt.Errorf("failed to mark post as posted: %w", err)
	}
	return nil
}

// MarkAsFailed marks a post as failed with error details
func (r *PostRepository) MarkAsFailed(ctx context.Context, postID string, status PostStatus, errorMessage string, errorCode string) error {
	now := time.Now()
	query := `
		UPDATE posts SET
			status = $1,
			error_message = $2,
			error_code = $3,
			failed_at = $4,
			updated_at = NOW()
		WHERE id = $5`

	_, err := r.conn.DB.ExecContext(ctx, query, status, errorMessage, errorCode, now, postID)
	if err != nil {
		return fmt.Errorf("failed to mark post as failed: %w", err)
	}
	return nil
}

// IncrementRetryCount increments the retry counter (Tier 1 network retries)
func (r *PostRepository) IncrementRetryCount(ctx context.Context, postID string) error {
	query := `
		UPDATE posts SET
			retry_count = retry_count + 1,
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.conn.DB.ExecContext(ctx, query, postID)
	if err != nil {
		return fmt.Errorf("failed to increment retry count: %w", err)
	}
	return nil
}

// IncrementIterationCount increments the iteration counter (Tier 2 validation iterations)
func (r *PostRepository) IncrementIterationCount(ctx context.Context, postID string) error {
	query := `
		UPDATE posts SET
			tier_2_iteration_count = tier_2_iteration_count + 1,
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.conn.DB.ExecContext(ctx, query, postID)
	if err != nil {
		return fmt.Errorf("failed to increment iteration count: %w", err)
	}
	return nil
}

// =============================================================================
// GENERATED MESSAGES REPOSITORY
// =============================================================================

// SaveGeneratedMessage saves a generated message for debugging
func (r *PostRepository) SaveGeneratedMessage(ctx context.Context, msg *GeneratedMessage) error {
	query := `
		INSERT INTO generated_messages (
			post_id, platform, iteration_number, prompt_template, prompt_variables,
			message_text, token_count, model_used, validation_status, validation_errors,
			character_count, generation_time_ms, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id, created_at`

	err := r.conn.DB.QueryRowContext(ctx, query,
		msg.PostID,
		msg.Platform,
		msg.IterationNumber,
		msg.PromptTemplate,
		msg.PromptVariables,
		msg.MessageText,
		msg.TokenCount,
		msg.ModelUsed,
		msg.ValidationStatus,
		msg.ValidationErrors,
		msg.CharacterCount,
		msg.GenerationTimeMs,
		msg.Metadata,
	).Scan(&msg.ID, &msg.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to save generated message: %w", err)
	}
	return nil
}

// GetGeneratedMessages retrieves all generated messages for a post
func (r *PostRepository) GetGeneratedMessages(ctx context.Context, postID string) ([]GeneratedMessage, error) {
	var messages []GeneratedMessage
	query := `
		SELECT id, post_id, platform, iteration_number, prompt_template, prompt_variables,
		       message_text, token_count, model_used, validation_status, validation_errors,
		       character_count, generation_time_ms, created_at, metadata
		FROM generated_messages
		WHERE post_id = $1
		ORDER BY iteration_number ASC`

	err := r.conn.DB.SelectContext(ctx, &messages, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get generated messages: %w", err)
	}
	return messages, nil
}

// =============================================================================
// POST HISTORY REPOSITORY
// =============================================================================

// CreatePostHistory creates a post history record for a successful post
func (r *PostRepository) CreatePostHistory(ctx context.Context, history *PostHistory) error {
	query := `
		INSERT INTO post_history (
			platform_account_id, video_id, post_id, posted_at,
			platform_post_id, platform_url, message_text, engagement_metrics
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING id, created_at`

	err := r.conn.DB.QueryRowContext(ctx, query,
		history.PlatformAccountID,
		history.VideoID,
		history.PostID,
		history.PostedAt,
		history.PlatformPostID,
		history.PlatformURL,
		history.MessageText,
		history.EngagementMetrics,
	).Scan(&history.ID, &history.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create post history: %w", err)
	}
	return nil
}

// GetPostHistoryForAccount retrieves post history for a specific account
func (r *PostRepository) GetPostHistoryForAccount(ctx context.Context, accountID string, limit int) ([]PostHistory, error) {
	var history []PostHistory
	query := `
		SELECT id, platform_account_id, video_id, post_id, posted_at,
		       platform_post_id, platform_url, message_text, engagement_metrics, created_at
		FROM post_history
		WHERE platform_account_id = $1
		ORDER BY posted_at DESC
		LIMIT $2`

	err := r.conn.DB.SelectContext(ctx, &history, query, accountID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get post history: %w", err)
	}
	return history, nil
}

// GetLastPostForVideo retrieves the most recent post of a video to a specific account
func (r *PostRepository) GetLastPostForVideo(ctx context.Context, accountID string, videoID string) (*PostHistory, error) {
	var history PostHistory
	query := `
		SELECT id, platform_account_id, video_id, post_id, posted_at,
		       platform_post_id, platform_url, message_text, engagement_metrics, created_at
		FROM post_history
		WHERE platform_account_id = $1 AND video_id = $2
		ORDER BY posted_at DESC
		LIMIT 1`

	err := r.conn.DB.GetContext(ctx, &history, query, accountID, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last post for video: %w", err)
	}
	return &history, nil
}

// UpdateEngagementMetrics updates the engagement metrics for a post
func (r *PostRepository) UpdateEngagementMetrics(ctx context.Context, historyID string, metrics JSONB) error {
	query := `
		UPDATE post_history SET
			engagement_metrics = $1
		WHERE id = $2`

	_, err := r.conn.DB.ExecContext(ctx, query, metrics, historyID)
	if err != nil {
		return fmt.Errorf("failed to update engagement metrics: %w", err)
	}
	return nil
}

// =============================================================================
// STATISTICS AND REPORTING
// =============================================================================

// GetPostStats retrieves posting statistics
func (r *PostRepository) GetPostStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total posts by status
	var statusCounts []struct {
		Status PostStatus `db:"status"`
		Count  int        `db:"count"`
	}
	statusQuery := `
		SELECT status, COUNT(*) as count
		FROM posts
		GROUP BY status`

	err := r.conn.DB.SelectContext(ctx, &statusCounts, statusQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get status counts: %w", err)
	}
	stats["by_status"] = statusCounts

	// Posts in last 24 hours
	var last24h int
	err = r.conn.DB.GetContext(ctx, &last24h, `
		SELECT COUNT(*)
		FROM posts
		WHERE created_at > NOW() - interval '24 hours'`)
	if err != nil {
		return nil, fmt.Errorf("failed to get 24h count: %w", err)
	}
	stats["last_24_hours"] = last24h

	// Average retry count
	var avgRetries float64
	err = r.conn.DB.GetContext(ctx, &avgRetries, `
		SELECT AVG(retry_count)
		FROM posts
		WHERE status = 'posted'`)
	if err != nil {
		return nil, fmt.Errorf("failed to get avg retries: %w", err)
	}
	stats["avg_retries"] = avgRetries

	// Average iteration count
	var avgIterations float64
	err = r.conn.DB.GetContext(ctx, &avgIterations, `
		SELECT AVG(tier_2_iteration_count)
		FROM posts
		WHERE status = 'posted'`)
	if err != nil {
		return nil, fmt.Errorf("failed to get avg iterations: %w", err)
	}
	stats["avg_iterations"] = avgIterations

	return stats, nil
}
