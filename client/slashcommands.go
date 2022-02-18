package client

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func configureSlashCommands(dg *discordgo.Session) error {

	dg.AddHandler(messageCreate)

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
	_, err := dg.ApplicationCommandCreate(dg.State.User.ID, guildID, &cmd)
	if err != nil {
		return fmt.Errorf("cannot create '%v' command: %w", cmd.Name, err)

	}

	_, err = dg.ApplicationCommandCreate(dg.State.User.ID, guildID, &cmdNewman)
	if err != nil {
		return fmt.Errorf("cannot create '%v' command: %w", cmd.Name, err)
	}
	return nil
}

func messageCreate(s *discordgo.Session, it *discordgo.InteractionCreate) {
	switch it.ApplicationCommandData().Name {
	case "tweet_gw":
		sendTweet(s, it)
	case "newman":
		sendNewman(s, it)
	default:
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
