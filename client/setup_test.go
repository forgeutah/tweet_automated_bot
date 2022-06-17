// Package client is a used to wrapper the twitter client and discord client.
package client

import (
	"testing"
)

func TestSetupTwitterClient(t *testing.T) {

	tests := []struct {
		name              string
		jsonFileName      string
		wantNumberClients int
		wantErr           bool
	}{
		{
			name:              "pass:Setup_One_client",
			jsonFileName:      "test_twitter_single_valid.json",
			wantNumberClients: 1,
		},
		{
			name:              "pass:Setup_Many_clients",
			jsonFileName:      "test_twitter_many_valid.json",
			wantNumberClients: 2,
		},
		{
			name:         "fail:cannot_read_json",
			jsonFileName: "test_twitter_single_invalid.json",
			wantErr:      true,
		},
		{
			name:    "fail:no_file_found",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotClients, err := SetupTwitterClients(tt.jsonFileName)
			if !tt.wantErr && err != nil {
				t.Fatal(err)
			}
			if len(gotClients) != tt.wantNumberClients {
				t.Fatalf("expected length: %d\n actual length: %d\n", tt.wantNumberClients, len(gotClients))
			}
		})
	}
}
