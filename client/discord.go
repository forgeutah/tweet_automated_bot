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

	if err != nil {
		err = fmt.Errorf("error configuring slash commands: %w", err)
		return nil, err
	}
	return dg, nil
}

func (c *Client) RunDiscordBot() {
	// Open a websocket connection to Discord and begin listening.
	err := c.DiscordBot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// sent to general channel
	c.DiscordBot.ChannelMessageSend("922613112585207833", "ForgeFoundation Twitter Bot is now online!")
	err = c.configureSlashCommands()
	if err != nil {
		fmt.Println("error configuring slash commands,", err)
		return
	}

	// Cleanly close down the Discord session.
	defer c.DiscordBot.Close()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
