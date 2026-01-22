package database

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// =============================================================================
// ENUMS
// =============================================================================

// PlatformType represents the social media platform
type PlatformType string

const (
	PlatformTwitter  PlatformType = "twitter"
	PlatformLinkedIn PlatformType = "linkedin"
	PlatformDiscord  PlatformType = "discord"
	PlatformBluesky  PlatformType = "bluesky"
)

// PostStatus represents the lifecycle status of a post
type PostStatus string

const (
	PostStatusQueued           PostStatus = "queued"
	PostStatusGenerating       PostStatus = "generating"
	PostStatusValidating       PostStatus = "validating"
	PostStatusPosting          PostStatus = "posting"
	PostStatusPosted           PostStatus = "posted"
	PostStatusFailedNetwork    PostStatus = "failed_network"
	PostStatusFailedValidation PostStatus = "failed_validation"
	PostStatusFailedPermanent  PostStatus = "failed_permanent"
)

// =============================================================================
// MODELS
// =============================================================================

// YouTubeChannel represents a YouTube channel for video discovery
type YouTubeChannel struct {
	ID            string     `db:"id" json:"id"`
	ChannelID     string     `db:"channel_id" json:"channel_id"`
	RSSFeedURL    string     `db:"rss_feed_url" json:"rss_feed_url"`
	ChannelName   *string    `db:"channel_name" json:"channel_name,omitempty"`
	LastFetchedAt *time.Time `db:"last_fetched_at" json:"last_fetched_at,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	Metadata      JSONB      `db:"metadata" json:"metadata"`
	IsActive      bool       `db:"is_active" json:"is_active"`
}

// Video represents a YouTube video
type Video struct {
	ID               string     `db:"id" json:"id"`
	VideoID          string     `db:"video_id" json:"video_id"`
	Title            string     `db:"title" json:"title"`
	URL              string     `db:"url" json:"url"`
	ChannelID        string     `db:"channel_id" json:"channel_id"`
	Description      *string    `db:"description" json:"description,omitempty"`
	YouTubeChannelID *string    `db:"youtube_channel_id" json:"youtube_channel_id,omitempty"`
	ThumbnailURL     *string    `db:"thumbnail_url" json:"thumbnail_url,omitempty"`
	Playlist         *string    `db:"playlist" json:"playlist,omitempty"`
	ConferenceYear   *string    `db:"conference_year" json:"conference_year,omitempty"`
	PresenterTwitter *string    `db:"presenter_twitter" json:"presenter_twitter,omitempty"`
	PublishedAt      *time.Time `db:"published_at" json:"published_at,omitempty"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
	Metadata         JSONB      `db:"metadata" json:"metadata"`
	DurationSeconds  *int       `db:"duration_seconds" json:"duration_seconds,omitempty"`
	ViewCount        *int       `db:"view_count" json:"view_count,omitempty"`
}

// PlatformAccount represents a configured social media account
type PlatformAccount struct {
	ID                  string       `db:"id" json:"id"`
	Platform            PlatformType `db:"platform" json:"platform"`
	AccountHandle       string       `db:"account_handle" json:"account_handle"`
	CredentialKey       string       `db:"credential_key" json:"credential_key"`
	DisplayName         *string      `db:"display_name" json:"display_name,omitempty"`
	LastPostedAt        *time.Time   `db:"last_posted_at" json:"last_posted_at,omitempty"`
	NextPostDueAt       *time.Time   `db:"next_post_due_at" json:"next_post_due_at,omitempty"`
	CreatedAt           time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time    `db:"updated_at" json:"updated_at"`
	Metadata            JSONB        `db:"metadata" json:"metadata"`
	PostingScheduleDays int          `db:"posting_schedule_days" json:"posting_schedule_days"`
	CooldownDays        int          `db:"cooldown_days" json:"cooldown_days"`
	IsActive            bool         `db:"is_active" json:"is_active"`
}

// Post represents a posting attempt
type Post struct {
	ID                  string     `db:"id" json:"id"`
	PlatformAccountID   string     `db:"platform_account_id" json:"platform_account_id"`
	VideoID             string     `db:"video_id" json:"video_id"`
	Status              PostStatus `db:"status" json:"status"`
	MessageText         *string    `db:"message_text" json:"message_text,omitempty"`
	PlatformPostID      *string    `db:"platform_post_id" json:"platform_post_id,omitempty"`
	PlatformURL         *string    `db:"platform_url" json:"platform_url,omitempty"`
	ErrorMessage        *string    `db:"error_message" json:"error_message,omitempty"`
	ErrorCode           *string    `db:"error_code" json:"error_code,omitempty"`
	StartedAt           *time.Time `db:"started_at" json:"started_at,omitempty"`
	PostedAt            *time.Time `db:"posted_at" json:"posted_at,omitempty"`
	FailedAt            *time.Time `db:"failed_at" json:"failed_at,omitempty"`
	QueuedAt            time.Time  `db:"queued_at" json:"queued_at"`
	CreatedAt           time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at" json:"updated_at"`
	MessageMetadata     JSONB      `db:"message_metadata" json:"message_metadata"`
	RetryCount          int        `db:"retry_count" json:"retry_count"`
	Tier2IterationCount int        `db:"tier_2_iteration_count" json:"tier_2_iteration_count"`
}

// GeneratedMessage represents an LLM-generated message
type GeneratedMessage struct {
	ID               string       `db:"id" json:"id"`
	PostID           string       `db:"post_id" json:"post_id"`
	Platform         PlatformType `db:"platform" json:"platform"`
	MessageText      string       `db:"message_text" json:"message_text"`
	PromptTemplate   *string      `db:"prompt_template" json:"prompt_template,omitempty"`
	ModelUsed        *string      `db:"model_used" json:"model_used,omitempty"`
	ValidationStatus *string      `db:"validation_status" json:"validation_status,omitempty"`
	CreatedAt        time.Time    `db:"created_at" json:"created_at"`
	PromptVariables  JSONB        `db:"prompt_variables" json:"prompt_variables"`
	ValidationErrors JSONB        `db:"validation_errors" json:"validation_errors"`
	Metadata         JSONB        `db:"metadata" json:"metadata"`
	IterationNumber  int          `db:"iteration_number" json:"iteration_number"`
	TokenCount       *int         `db:"token_count" json:"token_count,omitempty"`
	CharacterCount   *int         `db:"character_count" json:"character_count,omitempty"`
	GenerationTimeMs *int         `db:"generation_time_ms" json:"generation_time_ms,omitempty"`
}

// PostHistory represents the audit trail of successful posts
type PostHistory struct {
	ID                string    `db:"id" json:"id"`
	PlatformAccountID string    `db:"platform_account_id" json:"platform_account_id"`
	VideoID           string    `db:"video_id" json:"video_id"`
	PostID            string    `db:"post_id" json:"post_id"`
	PlatformPostID    *string   `db:"platform_post_id" json:"platform_post_id,omitempty"`
	PlatformURL       *string   `db:"platform_url" json:"platform_url,omitempty"`
	MessageText       *string   `db:"message_text" json:"message_text,omitempty"`
	PostedAt          time.Time `db:"posted_at" json:"posted_at"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	EngagementMetrics JSONB     `db:"engagement_metrics" json:"engagement_metrics"`
}

// =============================================================================
// JSONB TYPE (for PostgreSQL JSONB columns)
// =============================================================================

// JSONB is a custom type that allows storing JSON data in PostgreSQL JSONB columns
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface for JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), j)
	}
	return json.Unmarshal(bytes, j)
}

// =============================================================================
// VIEW MODELS (for database views)
// =============================================================================

// EligibleVideo represents a video eligible for posting to an account
type EligibleVideo struct {
	ID                  string       `db:"id" json:"id"`
	VideoID             string       `db:"video_id" json:"video_id"`
	Title               string       `db:"title" json:"title"`
	URL                 string       `db:"url" json:"url"`
	ChannelID           string       `db:"channel_id" json:"channel_id"`
	PlatformAccountID   string       `db:"platform_account_id" json:"platform_account_id"`
	Platform            PlatformType `db:"platform" json:"platform"`
	AccountHandle       string       `db:"account_handle" json:"account_handle"`
	PublishedAt         *time.Time   `db:"published_at" json:"published_at,omitempty"`
	LastPostedToAccount time.Time    `db:"last_posted_to_account" json:"last_posted_to_account"`
	DaysSinceLastPost   float64      `db:"days_since_last_post" json:"days_since_last_post"`
	CooldownDays        int          `db:"cooldown_days" json:"cooldown_days"`
}

// PlatformAccountStatus represents the posting status of an account
type PlatformAccountStatus struct {
	ID                     string       `db:"id" json:"id"`
	Platform               PlatformType `db:"platform" json:"platform"`
	AccountHandle          string       `db:"account_handle" json:"account_handle"`
	LastPostedAt           *time.Time   `db:"last_posted_at" json:"last_posted_at,omitempty"`
	NextPostDueAt          *time.Time   `db:"next_post_due_at" json:"next_post_due_at,omitempty"`
	DaysSinceLastPost      float64      `db:"days_since_last_post" json:"days_since_last_post"`
	PostingScheduleDays    int          `db:"posting_schedule_days" json:"posting_schedule_days"`
	CooldownDays           int          `db:"cooldown_days" json:"cooldown_days"`
	TotalVideosPosted      int          `db:"total_videos_posted" json:"total_videos_posted"`
	VideosPostedLast30Days int          `db:"videos_posted_last_30_days" json:"videos_posted_last_30_days"`
	IsActive               bool         `db:"is_active" json:"is_active"`
}
