package database

import (
	"context"
	"fmt"
	"time"
)

// CREATE TYPE valid_status AS ENUM ('queued', 'sent', 'failed');

// CREATE TABLE IF NOT EXISTS tweets (
// 	id serial,
// 	twitter_username VARCHAR(15),
// 	tweet_text text,
// 	links text,
// 	send_time timestamp,
// 	status valid_status,
// 	created_at TIMESTAMP DEFAULT now()
// );

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
// TODO: filter out videos that have been sent recently.
func (c *Connection) SelectOneRandomVideo(ctx context.Context, videoPlaylist string) (*YoutubeVideo, error) {
	var video YoutubeVideo
	// not sent in the last 3 months
	rows, err := c.DB.QueryxContext(ctx, `SELECT
	twitter_username,
	video_title,
	video_url,
	video_playlist,
	conference_year,
	presenter_twitter_username
	FROM yt_videos WHERE video_playlist = $1 ORDER BY random() LIMIT 1`, videoPlaylist)
	if err != nil {
		return &video, fmt.Errorf("err: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.StructScan(&video)
		if err != nil {
			return &video, fmt.Errorf("rows.ScanStruct: %w", err)
		}
	}
	if rows.Err() != nil {
		return &video, fmt.Errorf("rows.Err: %w", rows.Err())
	}

	return &video, nil
}

// UpdateSentAt updates the last sent time for a video.
func (c *Connection) UpdateSentAt(ctx context.Context, video *YoutubeVideo) error {
	_, err := c.DB.ExecContext(ctx, `UPDATE yt_videos SET last_sent_at = $1 WHERE video_url = $2`, time.Now(), video.URL)
	if err != nil {
		return fmt.Errorf("err: %w", err)
	}
	return nil
}
