package main

import (
	"context"
	"log"

	"github.com/SoyPete/tweet_automated_bot/client"
	database "github.com/SoyPete/tweet_automated_bot/db"
)

func main() {
	client := client.NewClient()
	database.Connect(context.Background())
	err := client.SendTweet("like if #golang is the best language")
	if err != nil {
		log.Fatal(err)
	}
}
