package main

import (
	"context"
	"fmt"
	"log"
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

	db.Ping()
	defer db.Close(ctx)
	defer cancel()

	os.Exit(m.Run())
}
