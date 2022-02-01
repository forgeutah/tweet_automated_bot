package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	tweetBotCommand := &discordgo.ApplicationCommand{
		Name:        "tweet",
		Description: "tweet test to the forge foundation twiiter account.",
		Type:        discordgo.ChatApplicationCommand,
	}
	_, err = dg.ApplicationCommandCreate(dg.State.User.ID, "", tweetBotCommand)
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	dg.AddHandler(handleTweets)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// this function will handle tweets sent via discord slash commands.
func handleTweets(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "tweet" {
		fmt.Println("tweet command received")
	}
}
