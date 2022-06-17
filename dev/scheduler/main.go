package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SoyPete/tweet_automated_bot/client"
	database "github.com/SoyPete/tweet_automated_bot/db"
	"github.com/SoyPete/tweet_automated_bot/internal/botguts"
)

func main() {

	// creates a ticker that pushished events every 5 seconds
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	// setup with api and database
	ctx := context.Background()
	client, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// start the twitter bot
	bot := botguts.NewAutoBot(db, client)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				err := bot.TweetYoutubeVideo(ctx)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}
