// discord/commands.go
package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "basic-command",
			Description: "Basic command",
		},
		{
			Name:        "followups",
			Description: "Followup messages",
		},
	}
)

func (b *Bot) registerCommands() error {
	for _, cmd := range commands {
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.GuildID, cmd)
		if err != nil {
			return fmt.Errorf("cannot create command '%v': %v", cmd.Name, err)
		}
	}
	return nil
}
