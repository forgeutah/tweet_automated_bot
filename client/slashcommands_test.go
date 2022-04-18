package client

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
)

func TestClient_configureSlashCommands(t *testing.T) {
	type fields struct {
		TweetBot   *twitter.Client
		DiscordBot *discordgo.Session
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				TweetBot:   tt.fields.TweetBot,
				DiscordBot: tt.fields.DiscordBot,
			}
			if err := c.configureSlashCommands(); (err != nil) != tt.wantErr {
				t.Errorf("Client.configureSlashCommands() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_messageCreate(t *testing.T) {
	type fields struct {
		TweetBot   *twitter.Client
		DiscordBot *discordgo.Session
	}
	type args struct {
		s  *discordgo.Session
		it *discordgo.InteractionCreate
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				TweetBot:   tt.fields.TweetBot,
				DiscordBot: tt.fields.DiscordBot,
			}
			c.messageCreate(tt.args.s, tt.args.it)
		})
	}
}

func TestClient_sendTweet(t *testing.T) {
	type fields struct {
		TweetBot   *twitter.Client
		DiscordBot *discordgo.Session
	}
	type args struct {
		s  *discordgo.Session
		it *discordgo.InteractionCreate
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				TweetBot:   tt.fields.TweetBot,
				DiscordBot: tt.fields.DiscordBot,
			}
			c.sendTweet(tt.args.s, tt.args.it)
		})
	}
}

func Test_haveValidRoles(t *testing.T) {
	tests := []struct {
		name  string
		roles []string
		want  bool
	}{
		{
			name: "Valid roles",
			roles: []string{
				"939282540991225897",
			},
			want: true,
		},
		{
			name: "Invalid roles",
			roles: []string{
				"123432523593580293",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := haveValidRoles(tt.roles); got != tt.want {
				t.Errorf("haveValidRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sendNewman(t *testing.T) {
	type args struct {
		s  *discordgo.Session
		it *discordgo.InteractionCreate
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendNewman(tt.args.s, tt.args.it)
		})
	}
}
