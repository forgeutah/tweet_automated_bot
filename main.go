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
	// check tha the database us upto dat with video files
	db.Migrate(ctx)

	go func() {
		err = client.RunDiscordBot()
		if err != nil {
			shutDown(ctx, client, db)
		}
	}()
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	signal.Notify(client.ShutDown, syscall.SIGINT, syscall.SIGTERM)
	bot := botguts.NewAutoBot(db, client)

	go func() {
		err = bot.ScheduleVideoTweet(ctx)
		if err != nil {
			shutDown(ctx, client, db)
		}
	}()

	http.HandleFunc("/health", healthCheck)

	//handle for ctrl+c
	go shutDown(ctx, client, db)

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

// todo: =maybe run this once with a shutdown chanel
func shutDown(ctx context.Context, client *client.Client, db *database.Connection) {
	<-ctx.Done()
	fmt.Println("Bot is now stopped.")
	client.DiscordBot.Close()
	db.Close(ctx)
	os.Exit(0)
}

// healthCheck is a http handler for health check to make sure the server is up.
func healthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("we are live"))
	if err != nil {
		log.Println(err)
	}
}
