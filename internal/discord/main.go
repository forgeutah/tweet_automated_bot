package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const guildID = "922613112119631913"

func main() {

	token := os.Getenv("DISCORD_TOKEN")
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.AddHandler(sendTweet)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// Cleanly close down the Discord session.
	defer dg.Close()

	//make command
	cmd := discordgo.ApplicationCommand{
		Name:        "tweet_gw",
		Description: "Send a tweet in the gowest channel",
	}

	// message we are online
	(dg.ChannelMessageSend("922613112585207833", "the Forge has it's eyes on you!"))
	_, err = dg.ApplicationCommandCreate(dg.State.User.ID, guildID, &cmd)
	if err != nil {
		fmt.Print(fmt.Errorf("cannot create '%v' command: %w", cmd.Name, err))
		panic(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func sendTweet(s *discordgo.Session, it *discordgo.InteractionCreate) {
	if it.Type == discordgo.InteractionApplicationCommand {
		if it.Message.Content == "!tweet_gw" {
			s.ChannelMessageSend(it.ChannelID, "Tweet sent!")
		}
	}
}
