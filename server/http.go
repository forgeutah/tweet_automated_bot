package server

import (
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
