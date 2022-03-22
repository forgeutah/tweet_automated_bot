package client

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func (c *Client) configureSlashCommands() error {
	log.Println("Setup Slash Commands")
	c.DiscordBot.AddHandler(c.messageCreate)

	//make command
	cmd := discordgo.ApplicationCommand{
		Name:        "tweet_gw",
		Description: "Send a tweet in the gowest channel",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "tweet-body", // all names is discord api must be lowercase and without space.
				Description: "content of the tweet",
				Required:    true,
			},
		},
	}
	cmdNewman := discordgo.ApplicationCommand{
		Name:        "newman",
		Description: "nyan cat gif share",
	}

	// add commands to bot
	_, err := c.DiscordBot.ApplicationCommandCreate(c.DiscordBot.State.User.ID, guildID, &cmd)
	if err != nil {
		return fmt.Errorf("cannot create '%v' command: %w", cmd.Name, err)

	}
	_, err = c.DiscordBot.ApplicationCommandCreate(c.DiscordBot.State.User.ID, guildID, &cmdNewman)
	if err != nil {
		return fmt.Errorf("cannot create '%v' command: %w", cmd.Name, err)
	}
	return nil
}

func (c *Client) messageCreate(s *discordgo.Session, it *discordgo.InteractionCreate) {
	switch it.ApplicationCommandData().Name {
	case "tweet_gw":
		c.sendTweet(s, it)
	case "newman":
		sendNewman(s, it)
	default:
	}

}

func (c *Client) sendTweet(s *discordgo.Session, it *discordgo.InteractionCreate) {
	discordResponse := "tweet Sent"
	embedImage := &discordgo.MessageEmbed{
		Title: "Sent",
		Image: &discordgo.MessageEmbedImage{
			URL: "https://media.giphy.com/media/Qs79cNS60bhY9UC1dP/giphy.gif",
		},
	}

	if !haveValidRoles(it.Member.Roles) {
		discordResponse = "you are not authorized to use this command"
		embedImage.Title = "Not Authorized"
		embedImage.Image.URL = "https://media.giphy.com/media/aU8vURhmTjX4He06SF/giphy.gif"
	} else {

		//get tweet message
		tweetBody := it.ApplicationCommandData().Options[0].StringValue()

		err := c.SendTweet(tweetBody)
		if err != nil {
			discordResponse = "error sending tweet"
			embedImage.Title = "Error"
			embedImage.Image.URL = "https://media.giphy.com/media/3oxHQn7gZW2wBGafxm/giphy.gif"
		}
	}

	// make response
	s.InteractionRespond(it.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: discordResponse,
			Embeds:  []*discordgo.MessageEmbed{embedImage},
		},
	})

}

// check if user have valid roles
func haveValidRoles(roles []string) bool {
	for _, role := range roles {
		// tweet_bot role
		if role == tweetBotRole {
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
