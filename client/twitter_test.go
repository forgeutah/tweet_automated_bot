package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
)

// testServer returns an http Client, ServeMux, and Server. The client proxies
// requests to the server and handlers can be registered on the mux to handle
// requests. The caller must close the test server.
func testServer() (*http.Client, *http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	transport := &RewriteTransport{&http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}}
	client := &http.Client{Transport: transport}
	return client, mux, server
}

// RewriteTransport rewrites https requests to http to avoid TLS cert issues
// during testing.
type RewriteTransport struct {
	Transport http.RoundTripper
}

// RoundTrip rewrites the request scheme to http and calls through to the
// composed RoundTripper or if it is nil, to the http.DefaultTransport.
func (t *RewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	if t.Transport == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return t.Transport.RoundTrip(req)
}

func testClientSetup(t *testing.T) *Client {
	httpClient, mux, server := testServer()
	t.Cleanup(server.Close)
	mux.HandleFunc("/1.1/statuses/update.json", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			t.Errorf("Expected method to be POST, got %s", r.Method)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading body: %v", err)
		}

		t.Cleanup(func() {
			r.Body.Close()
		})

		if string(body) != "status=Test+tweet" {
			t.Errorf("Expected body to be 'status=Test tweet', got %s", string(body))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id": 581980947630845953, "text": "very informative tweet"}`)
	})

	client := twitter.NewClient(httpClient)
	return &Client{
		TweetBot: client,
	}
}
func testClientFailSetup(t *testing.T) *Client {
	httpClient, mux, server := testServer()
	t.Cleanup(server.Close)
	mux.HandleFunc("/1.1/statuses/update.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected method to be POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, `{"errors":[{"code": 187,"message":"Status is a duplicate"}]}`)
	})

	client := twitter.NewClient(httpClient)
	return &Client{
		TweetBot: client,
	}
}
func testClient500Setup(t *testing.T) *Client {
	httpClient, mux, server := testServer()
	t.Cleanup(server.Close)
	mux.HandleFunc("/1.1/statuses/update.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected method to be POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := twitter.NewClient(httpClient)
	return &Client{
		TweetBot: client,
	}
}

func TestClient_SendTweet(t *testing.T) {
	tests := []struct {
		name    string
		bot     *Client
		message string
		wantErr bool
	}{
		{
			name:    "Send tweet",
			bot:     testClientSetup(t),
			message: "Test tweet",
			wantErr: false,
		},
		{
			name:    "Send tweet fail",
			bot:     testClientFailSetup(t),
			wantErr: true,
		},
		{
			name:    "Error on status code",
			bot:     testClient500Setup(t),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.bot.SendTweet(tt.message); (err != nil) != tt.wantErr {
				t.Errorf("Client.SendTweet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
