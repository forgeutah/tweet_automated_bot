package client

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// TODO: where should this go?
const guildID = "922613112119631913"
const tweetBotRole = "939282540991225897"

func setupDiscord(token string) (*discordgo.Session, error) {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		err = fmt.Errorf("error creating Discord session: %w", err)
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
	log.Println("sending message to discord bot")

	// ChannelMessageSend returns a response messge and and error
	_, err = c.DiscordBot.ChannelMessageSend("922613112585207833", "ForgeFoundation Twitter Bot is now online!")
	if err != nil {
		fmt.Println("error sending startup message,", err)
		return
	}
	err = c.configureSlashCommands()
	if err != nil {
		fmt.Println("error configuring slash commands,", err)
		return
	}
}
