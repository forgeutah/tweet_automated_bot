package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/SoyPete/tweet_automated_bot/client"
	database "github.com/SoyPete/tweet_automated_bot/db"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())

	fmt.Println("Init is being called.")
	_, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close(ctx)
	defer cancel()

	os.Exit(m.Run())
}

func Test_healthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	healthCheck(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "we are live" {
		t.Errorf("expected \"we are live\" got %v", string(data))
	}
}
