// discord/handlers.go
package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/tibia-oce/discord-bot/src/logger"
)

// handleBasicCommand responds to the "basic-command" interaction
func handleBasicCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
