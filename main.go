package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/SoyPete/tweet_automated_bot/auth"
)

// Tweet is twitter payload for tweets endpoint.
type Tweet struct {
	Text string `json:"text"`
}

func main() {

	// define http client
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	newTweet := Tweet{
		Text: "We love #golang",
	}

	tweetPayload, err := json.Marshal(&newTweet)
	if err != nil {
		panic(err)
	}

	tweetUrl := "https://api.twitter.com/2/tweets"
	req, err := http.NewRequest("POST", tweetUrl, bytes.NewReader(tweetPayload))
	if err != nil {
		panic(err)
	}
	oauthString, err := auth.GetTwitterOauthHeader(tweetUrl, newTweet.Text)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", oauthString)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode, string(body))
}
