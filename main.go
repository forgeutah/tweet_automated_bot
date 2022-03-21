package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/SoyPete/tweet_automated_bot/client"
	database "github.com/SoyPete/tweet_automated_bot/db"
	"github.com/SoyPete/tweet_automated_bot/internal/botguts"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("we are live"))
}

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

	defer db.Close(ctx)

	// TODO: ここでbotを作成する
	bot := botguts.NewAutoBot(db, client)
	err = bot.TweetYoutubeVideo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	client.RunDiscordBot()
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	signal.Notify(client.ShutDown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-client.ShutDown
		fmt.Println("Bot is now stopped.")
		os.Exit(0)
	}()

}
