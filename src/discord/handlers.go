// discord/handlers.go
package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/tibia-oce/discord-bot/src/github"
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
			Content: "Please select the repository and issue type:",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				// Dropdown for repository selection
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "repository_select_menu",
							Placeholder: "Select a repository",
							Options: []discordgo.SelectMenuOption{
								{
									Label:       "Server",
									Value:       "server",
									Description: "Issues related to the server",
								},
								{
									Label:       "Client",
									Value:       "client",
									Description: "Issues related to the client",
								},
								{
									Label:       "Map & Assets",
									Value:       "assets",
									Description: "Issues related to map and assets",
								},
							},
						},
					},
				},
				// Dropdown for issue type selection
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "issue_type_select_menu",
							Placeholder: "Select an issue type",
							Options: []discordgo.SelectMenuOption{
								{
									Label:       "Bug Report",
									Value:       "bug",
									Description: "Report a bug",
								},
								{
									Label:       "Feature Request",
									Value:       "feature",
									Description: "Request a new feature",
								},
								{
									Label:       "Other",
									Value:       "other",
									Description: "Other types of issues",
								},
							},
						},
					},
				},
				// Button to open modal for issue details
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Enter Issue Details",
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
			CustomID: "issue_details_modal",
			Title:    "Issue Details",
			Components: []discordgo.MessageComponent{
				// Title input
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "issue_title_input",
						Label:       "Issue Title",
						Style:       discordgo.TextInputShort,
						Placeholder: "Brief title of the issue",
						Required:    true,
					},
				}},
				// Expected behavior input
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "expected_behavior_input",
						Label:       "Expected Behavior",
						Style:       discordgo.TextInputParagraph,
						Placeholder: "Describe the expected behavior",
						Required:    true,
					},
				}},
				// Current behavior input
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "current_behavior_input",
						Label:       "Current Behavior",
						Style:       discordgo.TextInputParagraph,
						Placeholder: "Describe the current behavior",
						Required:    true,
					},
				}},
				// Optional image link input
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "image_link_input",
						Label:       "Image Link (optional)",
						Style:       discordgo.TextInputShort,
						Placeholder: "URL to an image or screenshot",
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

func handleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate, ghClient *github.GitHubClient) {
	data := i.ModalSubmitData()
	title := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	description := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	imageLink := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	repository := selectedOptions["repository_select_menu"]
	issueType := selectedOptions["issue_type_select_menu"]

	// Consolidated content to display in Discord and log in GitHub placeholder
	content := fmt.Sprintf(
		"**Form Submission Details:**\n\n**Title:** %s\n**Description:** %s\n**Repository:** %s\n**Issue Type:** %s",
		title, description, repository, issueType,
	)
	if imageLink != "" {
		content += fmt.Sprintf("\n**Image Link:** %s", imageLink)
	}

	ghClient.CreateIssue(repository, issueType, title, description, imageLink)

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
	primary := selectedOptions["repository_select_menu"]
	secondary := selectedOptions["issue_type_select_menu"]
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
