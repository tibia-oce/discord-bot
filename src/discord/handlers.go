// discord/handlers.go
package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tibia-oce/discord-bot/src/logger"
)

func (b *Bot) interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Find the command by name and execute the associated handler
	for _, cmd := range b.getCommands() {
		if i.ApplicationCommandData().Name == cmd.Name {
			cmd.Handler(s, i)
			return
		}
	}
	logger.Warn("Unknown command interaction received.")
}

func (b *Bot) handleBasicCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You executed the basic command!",
		},
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed to respond to basic-command interaction: %v", err))
	}
}

func (b *Bot) handleFollowupsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Followup message will be sent in 5 seconds...",
		},
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed to respond to followups interaction: %v", err))
	}

	time.Sleep(5 * time.Second)
	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "This is a followup message!",
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed to send followup message: %v", err))
	}
}
