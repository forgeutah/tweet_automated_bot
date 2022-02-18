package client

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// TODO: where should this go?
const guildID = "922613112119631913"

func setupDiscord(token string) (*discordgo.Session, error) {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		err = fmt.Errorf("error creating Discord session: %w", err)
		return nil, err
	}

	RunDiscordBot(dg)

	err = configureSlashCommands(dg)
	if err != nil {
		err = fmt.Errorf("error configuring slash commands: %w", err)
		return nil, err
	}
	return dg, nil
}

func RunDiscordBot(dg *discordgo.Session) {
	// Open a websocket connection to Discord and begin listening.
	err := dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	dg.ChannelMessageSend("922613112585207833", "the Forge has it's eyes on you!")

	// Cleanly close down the Discord session.
	defer dg.Close()

	// TODO: this needs to link to twittttterss?
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
