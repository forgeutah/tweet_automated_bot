package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

/*
Todo:
* add command to single gowest channel
* limit who can call the command using discord roles
* add a command to forge utah twiter
* add roles to discord for gowest organizer and meetup organizer
* connect to twitter client
*/

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

	// channelOption := discordgo.ApplicationCommandOption{
	// 	Type:        discordgo.ApplicationCommandOptionChannel,
	// 	Name:        "channel-option",
	// 	Description: "limit channel?",
	// 	// Channel type mask
	// 	ChannelTypes: []discordgo.ChannelType{
	// 		discordgo.ChannelTypeGuildText,
	// 		discordgo.ChannelTypeGuildVoice,
	// 	},
	// 	Required: true,
	// }

	// roleOption := discordgo.ApplicationCommandOption{
	// 	Type:        discordgo.ApplicationCommandOptionRole,
	// 	Name:        "role-option",
	// 	Description: "Role option",
	// 	Required:    true,
	// }

	// var cmdMap []*discordgo.ApplicationCommandOption
	// cmdMap = append(cmdMap, &channelOption)
	// cmdMap = append(cmdMap, &roleOption)
	//make command
	cmd := discordgo.ApplicationCommand{
		Name:        "tweet_gw",
		Description: "Send a tweet in the gowest channel",
		// Options:     cmdMap,
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
	// check user role

	// role id 939282540991225897
	fmt.Println(it.Member.Roles, it.Member.User.Username)

	if haveValidRoles(it.Member.Roles) {
		// send tweet
		s.ChannelMessageSend("939270685468008520", "tweet sent")
	} else {
		s.ChannelMessageSend("939270685468008520", "you are not authorized to use this command")
	}
	s.InteractionRespond(it.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "tweet send to twitter client",
		},
	})
}

func haveValidRoles(roles []string) bool {
	for _, role := range roles {
		if role == "939282540991225897" {
			return true
		}
	}
	return false
}
