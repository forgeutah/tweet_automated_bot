package main

import (
	"context"
	"log"

	"github.com/SoyPete/tweet_automated_bot/client"
	database "github.com/SoyPete/tweet_automated_bot/db"
)

func main() {
	ctx := context.Background()
	client := client.NewClient()
	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTweet("are you excited for #gowest2022? more information coming soon!")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Close(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
