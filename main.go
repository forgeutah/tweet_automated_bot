package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SoyPete/tweet_automated_bot/client"
	database "github.com/SoyPete/tweet_automated_bot/db"
	"github.com/SoyPete/tweet_automated_bot/internal/botguts"
)

func main() {
	// TODO: setup command line flag for json file
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

	// make bot for each twitter account?
	gowestbot := botguts.NewAutoBot(db, client, "gowestconf")
	forgeutahbot := botguts.NewAutoBot(db, client, "forgeutahbot")

	go func() {
		err = gowestbot.ScheduleVideoTweet(ctx)
		if err != nil {
			shutDown(ctx, client, db)
		}
	}()
	go func() {
		err = forgeutahbot.ScheduleVideoTweet(ctx)
		if err != nil {
			shutDown(ctx, client, db)
		}
	}()

	http.HandleFunc("/health", healthCheck)

	//handle for ctrl+c
	go func() {
		<-client.ShutDown
		shutDown(ctx, client, db)
	}()

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

func shutDown(ctx context.Context, client *client.Client, db *database.Connection) {
	fmt.Println("Bot is now stopped.")
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
