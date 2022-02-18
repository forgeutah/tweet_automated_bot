package client

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Client struct {
	TweetBot   *twitter.Client
	DiscordBot *discordgo.Session
}

func NewClient() (*Client, error) {
	oauthConsumerKey := os.Getenv("OAUTH_CONSUMER_KEY")
	oauthConsumerSecret := os.Getenv("OAUTH_CONSUMER_SECRET")
	oauthAccessToken := os.Getenv("OAUTH_ACCESS_TOKEN")
	oauthAccessSecret := os.Getenv("OAUTH_ACCESS_SECRET")
	config := oauth1.NewConfig(oauthConsumerKey, oauthConsumerSecret)
	token := oauth1.NewToken(oauthAccessToken, oauthAccessSecret)
	discordToken := os.Getenv("DISCORD_TOKEN")
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	//Discord client
	dgclient, err := setupDiscord(discordToken)
	if err != nil {
		return nil, fmt.Errorf("failed to setup discord: %w", err)
	}

	c := &Client{
		TweetBot:   client,
		DiscordBot: dgclient,
	}

	return c, nil

}

func (c *Client) SendTweet(message string) error {
	_, resp, err := c.TweetBot.Statuses.Update(message, nil)
	if err != nil {
		return fmt.Errorf("failed to send tweet: %w", err)
	}
	if resp.StatusCode > 300 {
		return fmt.Errorf("status code: %d\n %v", resp.StatusCode, resp.Body)
	}
	return nil
}
