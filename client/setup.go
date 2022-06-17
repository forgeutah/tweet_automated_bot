// Package client is a used to wrapper the twitter client and discord client.
package client

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Client struct {
	TweetBot       *twitter.Client
	TwitterClients map[string]*twitter.Client
	DiscordBot     *discordgo.Session
	ShutDown       chan os.Signal
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
	sc := make(chan os.Signal, 1)

	c := &Client{
		TweetBot:   client,
		DiscordBot: dgclient,
		ShutDown:   sc,
	}

	return c, nil
}

// SetupTwitterClients will parse the local json file with twitter credentials and setup
// clients that connect to the twitter api. jsonFileName is an optional commandline flag.
func SetupTwitterClients(jsonFileName string) (map[string]*twitter.Client, error) {

	var credentials struct {
		ApiCredentials []struct {
			TwitterHandle       string `json:"twitter_handle"`
			OauthConsumerKey    string `json:"consumer_key"`
			OauthConsumerSecret string `json:"consumer_secret"`
			OauthAccessToken    string `json:"access_token"`
			OauthAccessSecret   string `json:"access_secret"`
		} `json:"twitter_clients"`
	}

	f, err := os.Open(jsonFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get json credential file: %w", err)
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&credentials)
	if err != nil {
		return nil, fmt.Errorf("failure to decode api credential json: %w", err)
	}
	clients := make(map[string]*twitter.Client)
	for _, cred := range credentials.ApiCredentials {
		config := oauth1.NewConfig(cred.OauthConsumerKey, cred.OauthConsumerSecret)
		token := oauth1.NewToken(cred.OauthAccessToken, cred.OauthAccessSecret)
		httpClient := config.Client(oauth1.NoContext, token)
		clients[cred.TwitterHandle] = twitter.NewClient(httpClient)
	}
	return clients, nil
}
