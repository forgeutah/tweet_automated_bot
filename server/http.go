package server

import (
	"log"
	"net/http"
	"os"

	"github.com/SoyPete/tweet_automated_bot/internal/botguts"
)

type Tweeter struct {
	AutoBot *botguts.AutoBot
}

// NewTweeter creates a new Tweeter.
func NewTweeterServer(autoBot *botguts.AutoBot) *Tweeter {
	return &Tweeter{
		AutoBot: autoBot,
	}
}

// ServeHTTP runs the HTTP server for gcp scheduler.
func (s *Tweeter) ServeHTTP() error {
	log.Print("starting server...")
	http.HandleFunc("/tweetVideo", s.tweetVideo)
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
		return err

	}
	// will return and close the server?
	return nil
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("I am healthy"))
}

func (s *Tweeter) tweetVideo(w http.ResponseWriter, r *http.Request) {
	// go behind an api
	err := s.AutoBot.TweetYoutubeVideo(r.Context())
	if err != nil {
		log.Printf("error tweeting video: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error tweeting video"))
		return
	}
}
