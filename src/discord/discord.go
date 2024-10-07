package discord

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tibia-oce/discord-bot/src/logger"
)

type Bot struct {
	Session *discordgo.Session
	Token   string
	GuildID string
}

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

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You executed the basic command!",
				},
			})
		},
		"followups": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Respond with a message, then follow up with another message
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Followup message will be sent in 5 seconds...",
				},
			})

			time.Sleep(5 * time.Second)
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "This is a followup message!",
			})
		},
	}
)

func (b *Bot) Init() error {
	if b.Token == "" {
		return fmt.Errorf("no Discord token provided")
	}

	session, err := discordgo.New("Bot " + b.Token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %v", err)
	}

	b.Session = session
	b.Session.ShouldReconnectOnError = true

	// Add handler for interactions (slash commands)
	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// Add handler for the bot connection to Discord
	b.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Info(fmt.Sprintf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator))
	})

	// Reconnect handler to log reconnection attempts
	b.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Resumed) {
		logger.Info("Bot reconnected to Discord")
	})

	// Ensure the bot tries to reconnect when disconnected
	b.Session.AddHandler(func(s *discordgo.Session, d *discordgo.Disconnect) {
		logger.Warn(fmt.Sprintf("Bot disconnected. Reconnecting... Reason: %v", d))
		for {
			err := s.Open()
			if err == nil {
				logger.Info("Reconnected successfully")
				break
			}
			logger.Error(fmt.Errorf("reconnect failed: %v", err))
			time.Sleep(5 * time.Second) // Backoff before retrying
		}
	})

	// Open WebSocket connection
	err = b.Session.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord connection: %v", err)
	}

	logger.Info("Bot is now running!")

	// Register commands
	for _, cmd := range commands {
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.GuildID, cmd)
		if err != nil {
			return fmt.Errorf("cannot create command '%v': %v", cmd.Name, err)
		}
	}

	return nil
}

// Close gracefully closes the Discord bot session
func (b *Bot) Close() error {
	if b.Session != nil {
		return b.Session.Close()
	}
	return nil
}

// WaitForShutdown waits for a system interrupt (Ctrl+C) and shuts down the bot
func (b *Bot) WaitForShutdown() {
	// Wait here until CTRL-C or other term signal is received.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	logger.Info("Gracefully shutting down.")
	b.Close()
}
