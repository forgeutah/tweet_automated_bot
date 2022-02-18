package client

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

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
