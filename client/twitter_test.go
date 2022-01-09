package client

import (
	"testing"
)

func TestClient_SendTweet(t *testing.T) {
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
			bot:  NewClient(),
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
