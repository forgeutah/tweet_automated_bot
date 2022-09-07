package database

import (
	"context"
	"fmt"
	"time"
)

var create_query = `
		CREATE TABLE IF NOT EXISTS yt_videos (
			id serial,
			twitter_username VARCHAR(15),
			video_title VARCHAR(255),
			video_playlist VARCHAR(255),
			video_url text,
			conference_year VARCHAR(4),
			presenter_twitter_username VARCHAR(64),
			created_at TIMESTAMP DEFAULT now(),
			last_sent_at TIMESTAMP 
		);
	`

// YoutubeVideo has all of the information about forge youtube videos
type YoutubeVideo struct {
	Username         string    `db:"twitter_username"`
	Title            string    `db:"video_title"`
	URL              string    `db:"video_url"`
	Playlist         string    `db:"video_playlist"`
	ConferenceYear   string    `db:"conference_year"`
	PresenterTwitter string    `db:"presenter_twitter_username"`
	LastSentAt       time.Time `db:"last_sent_at"`
}

// SelectOneRandomVideo selects a random video from the database.
func (c *Connection) SelectOneRandomVideo(ctx context.Context, accountName string) (*YoutubeVideo, error) {
	var video YoutubeVideo
	// not sent in the last 3 months
	rows, err := c.DB.QueryxContext(ctx, `SELECT
	twitter_username,
	video_title,
	video_url,
	video_playlist,
	conference_year,
	presenter_twitter_username,
	last_sent_at
	FROM yt_videos WHERE twitter_username = $1 ORDER BY random()`, accountName)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.StructScan(&video)
		if err != nil {
			return nil, fmt.Errorf("rows.ScanStruct: %w", err)
		}
		// filter out videos posted within the last 3 months
		if canSendVideo(&video) {
			return &video, nil
		}
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err: %w", rows.Err())
	}

	// all videos have been sent recently
	return nil, nil
}

// UpdateSentAt updates the last sent time for a video.
func (c *Connection) UpdateSentAt(ctx context.Context, video *YoutubeVideo) error {
	_, err := c.DB.ExecContext(ctx, `UPDATE yt_videos SET last_sent_at = $1 WHERE video_url = $2`, time.Now(), video.URL)
	if err != nil {
		return fmt.Errorf("err: %w", err)
	}
	return nil
}

// canSendVideo checks the last time a video was sent. If it was sent less than 3 months ago, it returns false.
func canSendVideo(video *YoutubeVideo) bool {
	return time.Since(video.LastSentAt) > (time.Hour * 24 * 30 * 3)
}
