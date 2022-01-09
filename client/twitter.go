package client

import (
	"fmt"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Client struct {
	TweetBot *twitter.Client
}

func NewClient() *Client {
	oauthConsumerKey := os.Getenv("OAUTH_CONSUMER_KEY")
	oauthConsumerSecret := os.Getenv("OAUTH_CONSUMER_SECRET")
	oauthAccessToken := os.Getenv("OAUTH_ACCESS_TOKEN")
	oauthAccessSecret := os.Getenv("OAUTH_ACCESS_SECRET")
	config := oauth1.NewConfig(oauthConsumerKey, oauthConsumerSecret)
	token := oauth1.NewToken(oauthAccessToken, oauthAccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	return &Client{
		TweetBot: client,
	}
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
