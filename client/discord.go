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

func setupDiscord(token string) *discordgo.Session {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return nil
	}

	dg.AddHandler(messageCreate)

	dg.AddHandler(sendTweet)

	dg.AddHandler(sendNewman)

	return dg
}

func (c *Client) RunDiscordBot() {
	// Open a websocket connection to Discord and begin listening.
	err := c.DiscordBot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// Cleanly close down the Discord session.
	defer c.DiscordBot.Close()

	// TODO: this needs to link to twittttterss?
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}

func (c *Client) configureSlashCommands() {

	//make command
	cmd := discordgo.ApplicationCommand{
		Name:        "tweet_gw",
		Description: "Send a tweet in the gowest channel",
		// Options:     cmdMap,
	}
	cmdNewman := discordgo.ApplicationCommand{
		Name:        "newman",
		Description: "nyan cat gif share",
		// Options:     cmdMap,
	}

	// message we are online
	c.DiscordBot.ChannelMessageSend("922613112585207833", "the Forge has it's eyes on you!")
	_, err := c.DiscordBot.ApplicationCommandCreate(c.DiscordBot.State.User.ID, guildID, &cmd)
	if err != nil {
		fmt.Print(fmt.Errorf("cannot create '%v' command: %w", cmd.Name, err))
		panic(err)
	}

	_, err = c.DiscordBot.ApplicationCommandCreate(c.DiscordBot.State.User.ID, guildID, &cmdNewman)
	if err != nil {
		fmt.Print(fmt.Errorf("cannot create '%v' command: %w", cmd.Name, err))
		panic(err)
	}

}
