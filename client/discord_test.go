package client

import (
	"os"
	"testing"
)

func TestClient_RunDiscordBot(t *testing.T) {
	// isIntegrationTest = os.Getenv("INTEGRATION_TEST")
	// if isIntegrationTest != "true" {
	// 	t.Skip("Skipping integration test")
	// }
	testToken := os.Getenv("DISCORD_TOKEN")
	validClient, err := setupDiscord(testToken)
	if err != nil {
		t.Errorf("Error setting up discord client: %s", err)
	}
	invalidClient, err := setupDiscord("")
	if err != nil {
		t.Errorf("Error setting up discord client: %s", err)
	}
	tests := []struct {
		name    string
		c       Client
		clean   func()
		wantErr bool
	}{
		{
			name: "TestClient_RunDiscordBot",
			c: Client{
				DiscordBot: validClient,
			},
			clean: func() {
				validClient.Close()
			},
		},
		{
			name: "TestClient_RunDiscordBot_Invalid",
			c: Client{
				DiscordBot: invalidClient,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.RunDiscordBot()
			if !tt.wantErr && err != nil {
				t.Errorf("Client.RunDiscordBot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		// is this even possible?
		if tt.clean != nil {
			tt.clean()
		}
	}
}

func Test_setupDiscord(t *testing.T) {
	isIntegrationTest := os.Getenv("INTEGRATION_TEST")
	if isIntegrationTest == "true" {
		t.Run("valid discord bot", test_setupDiscord_success)
	}

}

func test_setupDiscord_success(t *testing.T) {
	token := os.Getenv("DISCORD_TOKEN")

	got, err := setupDiscord(token)
	if err != nil {
		t.Errorf("setupDiscord() error = %v, in nil", err)
		return
	}
	t.Cleanup(func() { got.Close() })
}
