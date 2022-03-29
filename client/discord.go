package client

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// TODO: where should this go?
const guildID = "922613112119631913"
const tweetBotRole = "939282540991225897"

// TODO: update this function to accepts a variable number of arguments or generic discordgo.New() connection.
func setupDiscord(token string) (*discordgo.Session, error) {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		err = fmt.Errorf("error creating Discord session: %w", err)
		return nil, err
	}

	return dg, nil
}

// RunDiscord runs the discord client. This should be called in a goroutine in the main function. It starts by settings up the
// websocket connection to the discord server. A message is get to the general channel of the server stating that the bot has
// connected. Then it listens for messages from the discord server and handles them accordingly. Message handlers are defined
// as slashcommands in the client/slashcommands.go file.
func (c *Client) RunDiscordBot() error {
	// Open a websocket connection to Discord and begin listening.
	err := c.DiscordBot.Open()
	if err != nil {
		return fmt.Errorf("error opening connection %w", err)
	}

	// TODO: how to test without hittings this channel?

	// ChannelMessageSend returns a response messge and and error
	// _, err = c.DiscordBot.ChannelMessageSend("922613112585207833", "ForgeFoundation Twitter Bot is now online!")
	// if err != nil {
	// 	return fmt.Errorf("error sending startup message %w", err)
	// }
	err = c.configureSlashCommands()
	if err != nil {
		return fmt.Errorf("error configuring slash commands %w", err)
	}

	return nil
}
