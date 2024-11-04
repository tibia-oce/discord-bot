// discord/handlers.go
package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/tibia-oce/discord-bot/src/logger"
)

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

func handleButtonPrompt(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Would you like to proceed?",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Yes",
							Style:    discordgo.SuccessButton,
							CustomID: "prompt_yes",
						},
						discordgo.Button{
							Label:    "No",
							Style:    discordgo.DangerButton,
							CustomID: "prompt_no",
						},
					},
				},
			},
		},
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed to respond to button prompt: %v", err))
	}
}

func handleYesResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "You chose Yes!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed to respond to yes button: %v", err))
	}
}

func handleNoResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "You chose No.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed to respond to no button: %v", err))
	}
}
