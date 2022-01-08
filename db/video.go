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
			created_at TIMESTAMP DEFAULT now(),
			last_sent_at TIMESTAMP 
		);

		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'All Types of Golang Types', 'GoWest Conference', 'https://youtu.be/1RYYsLy9bg8');

		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Field Report: Building a game engine for 300 DEFCON hackers to smash', 'GoWest Conference', 'https://youtu.be/ZtP-IggAlB8');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'The Standard Library Bootcamp', 'GoWest Conference', 'https://youtu.be/CyhmhY7aI-s');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Micro Machine Learning in Go', 'GoWest Conference', 'https://youtu.be/Fq5-KmNr_D8');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Go Without Generics, a Retrospective', 'GoWest Conference', 'https://youtu.be/ZtP-IggAlB8');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Functional Programming with Go', 'GoWest Conference', 'https://youtu.be/SSO78QkmMLs');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'So you think you know Go?', 'GoWest Conference', 'https://youtu.be/U_qVSHYgVSE');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Anatomy of a Gopher - Binary Analysis of Go Binaries', 'GoWest Conference', 'https://youtu.be/ou_B5YZzEeU');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Going Serverless', 'GoWest Conference', 'https://youtu.be/M1AvdQVjhx8');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Writing REST Services for the gRPC-curious', 'GoWest Conference', 'https://youtu.be/rbwvY0YpRPI');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Go Templates: Great Library, Bad Rap', 'GoWest Conference', 'https://youtu.be/m8mKqdyD_C8');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Practical Tips to Creating a Great Engineering Culture', 'GoWest Conference', 'https://youtu.be/TxBz7L9ev18');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Callstacks at Giverny: A Go Graphics CLI', 'GoWest Conference', 'https://youtu.be/4LmTe0P7qpg');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'The beauty of Go for building cross-platform graphical applications
', 'GoWest Conference', 'https://youtu.be/VESXrgCW50g');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Building a scalable API platform for live streaming using GoLang', 'GoWest Conference', 'https://youtu.be/jWUhpOnhOQ0');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Go Channels Demystified', 'GoWest Conference', 'https://youtu.be/KdHqx9Bx4jI');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Ent: Making Data Easy in Go', 'GoWest Conference', 'https://youtu.be/NvjvzYacgQg');
		INSERT into yt_videos (twitter_username, video_title, video_playlist, video_url) VALUES ('gowestconf', 'Build a home monitoring system with Go', 'GoWest Conference', 'https://youtu.be/PolzSYuwaQg');`

// YoutubeVideo has all of the information about forge youtube videos
type YoutubeVideo struct {
	Username      string    `db:"twitter_username"`
	VideoName     string    `db:"video_title"`
	VideoURL      string    `db:"video_url"`
	VideoPlaylist string    `db:"video_playlist"`
	LastSentAt    time.Time `db:"last_sent_time"`
}

func (c *Connection) SelectOneRandomVideo(ctx context.Context, videoPlaylist string) (YoutubeVideo, error) {
	var video YoutubeVideo
	// not sent in the last 3 months
	rows, err := c.DB.QueryxContext(ctx, `SELECT
	twitter_username,
	video_title,
	video_url,
	video_playlist
	FROM yt_videos WHERE video_playlist = $1 ORDER BY random() LIMIT 1`, videoPlaylist)
	if err != nil {
		return video, fmt.Errorf("err: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.StructScan(&video)
		if err != nil {
			return video, fmt.Errorf("rows.ScanStruct: %w", err)
		}
	}
	if rows.Err() != nil {
		return video, fmt.Errorf("rows.Err: %w", rows.Err())
	}

	return video, nil
}
