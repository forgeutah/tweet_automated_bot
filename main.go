package main

import (
	"context"
	"log"

	"github.com/SoyPete/tweet_automated_bot/client"
	database "github.com/SoyPete/tweet_automated_bot/db"
	"github.com/SoyPete/tweet_automated_bot/internal/botguts"
	"github.com/SoyPete/tweet_automated_bot/server"
)

func main() {
	ctx := context.Background()
	client, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}
<<<<<<< HEAD
	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: remove and setup permanent datastore
	defer db.Close(ctx)

	// TODO: ここでbotを作成する
	bot := botguts.NewAutoBot(db, client)
	err = bot.TweetYoutubeVideo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	client.RunDiscordBot()
=======
	bot := botguts.NewAutoBot(db, client)
	s := server.NewTweeterServer(bot)
	s.ServeHTTP()
>>>>>>> e8cfe33 (ad server for gcp cloud run)
}
