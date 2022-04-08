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
	go client.RunDiscordBot()
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	signal.Notify(client.ShutDown, syscall.SIGINT, syscall.SIGTERM)
	bot := botguts.NewAutoBot(db, client)
	go bot.ScheduleVideoTweet(ctx)

	go func() {
		<-client.ShutDown
		fmt.Println("Bot is now stopped.")
		client.DiscordBot.Close()
		db.Close(ctx)
		os.Exit(0)
	}()

	http.HandleFunc("/health", healthCheck)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("we are live"))
}
