package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/SoyPete/tweet_automated_bot/client"
	database "github.com/SoyPete/tweet_automated_bot/db"
)

func main() {
	ctx := context.Background()

	// Connect to Supabase (migrations run automatically)
	db, err := database.ConnectSupabase()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("error closing database: %v", closeErr)
		}
	}()

	// Verify database health
	if healthErr := db.HealthCheck(); healthErr != nil {
		log.Fatal("database health check failed:", healthErr)
	}
	log.Println("Database health check passed")

	// TODO: Phase 2+ will implement the new multi-platform bot architecture
	// Legacy Twitter-only bot code has been removed

	http.HandleFunc("/health", healthCheck)

	//handle for ctrl+c
	client, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}
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

func shutDown(_ context.Context, _ *client.Client, db *database.SupabaseConnection) {
	log.Println("Bot is now stopped.")
	if err := db.Close(); err != nil {
		log.Printf("error closing database: %v", err)
	}
	os.Exit(0)
}

// healthCheck is a http handler for health check to make sure the server is up.
func healthCheck(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("we are live"))
	if err != nil {
		log.Println(err)
	}
}
