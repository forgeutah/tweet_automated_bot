// Package client is a used to wrapper the twitter client and discord client.
package client

import (
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		want    *Client
		wantErr bool
	}{
		{
			name: "NewClient",
			want: &Client{
				TweetBot:   &twitter.Client{},
				DiscordBot: &discordgo.Session{},
				ShutDown:   make(chan os.Signal, 1),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewClient() = %v, want %v", *got, *tt.want)
			}
		})
	}
}
