package botguts

import (
	"log"

	"github.com/SoyPete/tweet_automated_bot/client"
)

type AutoBot struct {
	twitterClient *client.Client
}

// NewAutoBot is a function that will create a new AutoBot. This auto bot is responsible for tweeting
// youtube videos on a weekly basis
func NewAutoBot(twitterClient *client.Client) *AutoBot {
	log.Println("creating new AutoBot")
	return &AutoBot{
		twitterClient: twitterClient,
	}
}
