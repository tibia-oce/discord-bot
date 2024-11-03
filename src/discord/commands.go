// discord/commands.go
package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func (b *Bot) registerCommands() error {
	commands := b.getCommands()

	for _, cmd := range commands {
		if _, err := b.registerCommand(cmd.Name, cmd.Description); err != nil {
			return fmt.Errorf("failed to register command %s: %v", cmd.Name, err)
		}
	}

	return nil
}

func (b *Bot) registerCommand(name, description string) (*discordgo.ApplicationCommand, error) {
	cmd := &discordgo.ApplicationCommand{
		Name:        name,
		Description: description,
	}
	return b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.GuildID, cmd)
}

func (b *Bot) getCommands() []Command {
	return []Command{
		{
			Name:        "basic-command",
			Description: "Basic command",
			Handler:     b.handleBasicCommand,
		},
		{
			Name:        "followups",
			Description: "Followup messages",
			Handler:     b.handleFollowupsCommand,
		},
	}
}
