package client

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func configureSlashCommands(dg *discordgo.Session) error {

	dg.AddHandler(messageCreate)

	dg.AddHandler(sendTweet)

	dg.AddHandler(sendNewman)
	//make command
	cmd := discordgo.ApplicationCommand{
		Name:        "tweet_gw",
		Description: "Send a tweet in the gowest channel",
		// Options:     cmdMap,
	}
	// cmdNewman := discordgo.ApplicationCommand{
	// 	Name:        "newman",
	// 	Description: "nyan cat gif share",
	// 	// Options:     cmdMap,
	// }
	fmt.Println(dg.State.User.ID)

	// message we are online
	_, err := dg.ApplicationCommandCreate(dg.State.User.ID, guildID, &cmd)
	if err != nil {
		return fmt.Errorf("cannot create '%v' command: %w", cmd.Name, err)

	}

	// _, err = c.DiscordBot.ApplicationCommandCreate(c.DiscordBot.State.User.ID, guildID, &cmdNewman)
	// if err != nil {
	// 	return fmt.Errorf("cannot create '%v' command: %w", cmd.Name, err)
	// }
	return nil
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
	})
}

// lock out jacoboco
func haveValidRoles(roles []string) bool {
	for _, role := range roles {
		// tweet_bot role
		if role == "939282540991225897" {
			return true
		}
	}
	return false

}

// send nyan cat gifs
func sendNewman(s *discordgo.Session, it *discordgo.InteractionCreate) {
	embedImage := &discordgo.MessageEmbed{
		Title: "Newman",
		Image: &discordgo.MessageEmbedImage{
			URL: "https://media.giphy.com/media/sIIhZliB2McAo/giphy.gif",
		},
	}
	s.InteractionRespond(it.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey there! Congratulations, you called for a newman!",
			Embeds:  []*discordgo.MessageEmbed{embedImage},
		},
	})
}
