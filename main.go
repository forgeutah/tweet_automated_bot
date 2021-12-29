package main

import (
	"fmt"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	oauthConsumerKey := os.Getenv("OAUTH_CONSUMER_KEY")
	oauthConsumerSecret := os.Getenv("OAUTH_CONSUMER_SECRET")
	oauthAccessToken := os.Getenv("OAUTH_ACCESS_TOKEN")
	oauthAccessSecret := os.Getenv("OAUTH_ACCESS_SECRET")
	config := oauth1.NewConfig(oauthConsumerKey, oauthConsumerSecret)
	token := oauth1.NewToken(oauthAccessToken, oauthAccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Send a Tweet
	tweet, resp, err := client.Statuses.Update("this is a robot tweet", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(tweet, resp)

}
