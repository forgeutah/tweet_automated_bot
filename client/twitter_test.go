package client

import (
	"os"
	"testing"
)

func TestClient_SendTweet(t *testing.T) {
	isIntegrationTest := os.Getenv("INTEGRATION_TEST")
	if isIntegrationTest != "true" {
		t.Skip("Skipping integration test")
	}
	c, _ := NewClient()
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		bot     *Client
		args    args
		wantErr bool
	}{
		{
			name: "Send tweet but no credentials",
			bot:  c,
			args: args{
				message: "Test tweet",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.bot.SendTweet(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Client.SendTweet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
