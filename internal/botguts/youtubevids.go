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

// TODO: move to twitter client
func makeVideoTweet(video *database.YoutubeVideo) string {
	hashtagYear := "#gowest2021"
	if video.ConferenceYear == "2020" {
		hashtagYear = "#gowest2020"
	}
	return fmt.Sprintf("Video Drop from %s\n\n%s\nby @%s\n\n%s\n\n#golang #gowestconf", hashtagYear, video.Title, video.PresenterTwitter, video.URL)
}

// TweetYoutubeVideo is a function that will tweet a random youtube video
func (ab *AutoBot) TweetYoutubeVideo(ctx context.Context) error {
	log.Println("tweeting youtube video")
	// get all the videos that have not been in last 3 months
	fmt.Println(ab.dbclient.Ping())
	video, err := ab.dbclient.SelectOneRandomVideo(ctx, "GoWest Conference")
	if err != nil {
		return fmt.Errorf("error getting random video: %w", err)
	}
	// create and send tweet
	err = ab.twitterClient.SendTweet(makeVideoTweet(video))
	if err != nil {
		return fmt.Errorf("error sending tweet: %w", err)
	}

	// last send time of video
	err = ab.dbclient.UpdateSentAt(ctx, video)
	if err != nil {
		return fmt.Errorf("error updating sent at: %w", err)
	}

	return nil
}
