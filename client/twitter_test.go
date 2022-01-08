package client

import (
	"reflect"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SendTweet(t *testing.T) {
	type fields struct {
		TweetBot *twitter.Client
	}
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Send tweet",
			fields: fields{
				TweetBot: &twitter.Client{},
			},
			args: args{
				message: "Test tweet",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				TweetBot: tt.fields.TweetBot,
			}
			if err := c.SendTweet(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Client.SendTweet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
