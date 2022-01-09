package botguts

import (
	"context"
	"fmt"
	"log"

	"github.com/SoyPete/tweet_automated_bot/client"
	database "github.com/SoyPete/tweet_automated_bot/db"
)

type AutoBot struct {
	dbclient      *database.Connection
	twitterClient *client.Client
}

// NewAutoBot is a function that will create a new AutoBot. This auto bot is responsible for tweeting
// youtube videos on a weekly basis
func NewAutoBot(dbclient *database.Connection, twitterClient *client.Client) *AutoBot {
	log.Println("creating new AutoBot")
	return &AutoBot{
		dbclient:      dbclient,
		twitterClient: twitterClient,
	}
}

// TweetYoutubeVideo is a function that will tweet a random youtube video
func (ab *AutoBot) TweetYoutubeVideo(ctx context.Context) error {
	log.Println("tweeting youtube video")
	// get all the videos that have not been in last 3 months
	fmt.Println(ab.dbclient.Ping())
	video, err := ab.dbclient.SelectOneRandomVideo(ctx, "GoWest Conference")
	fmt.Println(video, err)
	// select random video from the list maybe via sql

	// create a tweet with the video link

	// send the tweet

	return nil
}
