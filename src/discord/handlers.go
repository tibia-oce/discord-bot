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
			Flags:   discordgo.MessageFlagsEphemeral,
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
			Flags:   discordgo.MessageFlagsEphemeral,
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

func handleSelectMenuResponse(s *discordgo.Session, i *discordgo.InteractionCreate, issueChannelID string) {
	selectedValue := i.MessageComponentData().Values[0]
	content := fmt.Sprintf("%s You selected: %s", i.Member.Mention(), selectedValue) // Mention the user in the response

	// Send the response to the specified channel
	_, err := s.ChannelMessageSendComplex(issueChannelID, &discordgo.MessageSend{
		Content: content,
		Flags:   discordgo.MessageFlagsEphemeral, // Ensure only the active user can see the message
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed to send message to channel %s: %v", issueChannelID, err))
	}

	// Clear the dropdown menu in the original interaction message
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    "Issue received! Check the response in the designated channel.",
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{},
		},
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed to respond to select menu selection: %v", err))
	}
}

var selectedOptions = map[string]string{}

func handleExtendedForm(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Please fill out the form and provide an image:",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				// First dropdown for primary selection
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "primary_select_menu",
							Placeholder: "Select a primary option",
							Options: []discordgo.SelectMenuOption{
								{
									Label:       "Server",
									Value:       "server",
									Description: "This is the first primary option",
								},
								{
									Label:       "Client",
									Value:       "client",
									Description: "This is the second primary option",
								},
								{
									Label:       "Map & Assets",
									Value:       "assets",
									Description: "This is the third primary option",
								},
							},
						},
					},
				},
				// Second dropdown for secondary selection
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "secondary_select_menu",
							Placeholder: "Select a secondary option",
							Options: []discordgo.SelectMenuOption{
								{
									Label:       "Option 1",
									Value:       "option1",
									Description: "This is the first secondary option",
								},
								{
									Label:       "Option 2",
									Value:       "option2",
									Description: "This is the second secondary option",
								},
								{
									Label:       "Option 3",
									Value:       "option3",
									Description: "This is the third secondary option",
								},
							},
						},
					},
				},
				// Button to open modal for text input
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Enter Title and Description",
							Style:    discordgo.PrimaryButton,
							CustomID: "open_modal",
						},
					},
				},
			},
		},
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed to send form interaction response: %v", err))
	}
}

// Handler to open a modal for title, description, and image link input
func handleOpenModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "text_input_modal",
			Title:    "Enter Details",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "title_input",
						Label:       "Title",
						Style:       discordgo.TextInputShort,
						Placeholder: "Enter a title",
						Required:    true,
					},
				}},
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "description_input",
						Label:       "Description",
						Style:       discordgo.TextInputParagraph,
						Placeholder: "Enter a description",
						Required:    true,
					},
				}},
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "image_link_input",
						Label:       "Image Link",
						Style:       discordgo.TextInputShort,
						Placeholder: "Enter a URL to the image",
						Required:    false,
					},
				}},
			},
		},
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed to open modal: %v", err))
	}
}

// Handle modal submission
func handleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	title := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	description := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	imageLink := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	primarySelection := selectedOptions["primary_select_menu"]
	secondarySelection := selectedOptions["secondary_select_menu"]

	// Format the message with all user-provided inputs and selections
	content := fmt.Sprintf(
		"**Form Submission Details:**\n\n**Title:** %s\n**Description:** %s\n**Primary Selection:** %s\n**Secondary Selection:** %s",
		title, description, primarySelection, secondarySelection,
	)
	if imageLink != "" {
		content += fmt.Sprintf("\n**Image Link:** %s", imageLink)
	}

	// Respond to the interaction with the final content
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed to respond to modal interaction: %v", err))
	}
}

// Final submit handler that gathers selections and sends them
func handleFormSubmit(s *discordgo.Session, i *discordgo.InteractionCreate, issueChannelID string) {
	primary := selectedOptions["primary_select_menu"]
	secondary := selectedOptions["secondary_select_menu"]
	content := fmt.Sprintf("Form submitted:\nPrimary selection: %s\nSecondary selection: %s", primary, secondary)

	// Send the form data to the specified channel
	_, err := s.ChannelMessageSend(issueChannelID, content)
	if err != nil {
		logger.Error(fmt.Errorf("failed to send form submission to channel %s: %v", issueChannelID, err))
	}

	// Clear form components after submission
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    "Form submitted successfully! Check the designated channel for details.",
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{},
		},
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed to update message after form submit: %v", err))
	}

	// Clear selected options after submission
	selectedOptions = map[string]string{}
}

func handleSelection(s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectedValue := i.MessageComponentData().Values[0]
	selectedOptions[i.MessageComponentData().CustomID] = selectedValue

	// Send an ephemeral confirmation message that the option was recorded
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed to defer response for selection: %v", err))
	}
}
