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
	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: remove and setup permanent datastore
	defer db.Close(ctx)

	bot := botguts.NewAutoBot(db, client)
	s := server.NewTweeterServer(bot)

	// This calls ListenAndServer, which blocks
	s.ServeHTTP()

}
