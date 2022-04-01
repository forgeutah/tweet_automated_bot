// Package client is a used to wrapper the twitter client and discord client.
package client

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "NewClient",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// CHECK twitter client
			if client.TweetBot == nil {
				t.Errorf("NewClient() TwitterClient is nil")
			}

			// check discord client
			if client.DiscordBot == nil {
				t.Errorf("NewClient() DiscordClient is nil")
			}
		})
	}
}
