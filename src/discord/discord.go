// discord/discord.go
package discord

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tibia-oce/discord-bot/src/github"
	"github.com/tibia-oce/discord-bot/src/logger"
)

type Bot struct {
	Session        *discordgo.Session
	Token          string
	GuildID        string
	AppID          string
	IssueChannelID string
}

type Command struct {
	Name        string
	Description string
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

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
	b.Session.AddHandler(b.onReady())
	b.Session.AddHandler(b.onReconnect())
	b.Session.AddHandler(b.onDisconnect())

	ghClient := github.NewGitHubClient()
	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		b.interactionHandler(s, i, ghClient)
	})

	// Open WebSocket connection
	if err := b.Session.Open(); err != nil {
		return fmt.Errorf("error opening Discord connection: %v", err)
	}

	if err := b.startLatencyMonitor(); err != nil {
		return fmt.Errorf("error starting latency monitor: %v", err)
	}

	if err := b.registerCommands(); err != nil {
		return fmt.Errorf("error registering commands: %v", err)
	}

	logger.Info("Bot is now running!")

	return nil
}

func (b *Bot) onReady() func(s *discordgo.Session, r *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Info(fmt.Sprintf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator))
	}
}

func (b *Bot) onReconnect() func(s *discordgo.Session, r *discordgo.Resumed) {
	return func(s *discordgo.Session, r *discordgo.Resumed) {
		logger.Info("Bot reconnected to Discord")
	}
}

func (b *Bot) onDisconnect() func(s *discordgo.Session, d *discordgo.Disconnect) {
	return func(s *discordgo.Session, d *discordgo.Disconnect) {
		logger.Warn("Bot disconnected. Attempting to reconnect...")
		retryReconnection(s)
	}
}

func (b *Bot) Close() error {
	if b.Session != nil {
		return b.Session.Close()
	}
	return nil
}

func (b *Bot) WaitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	logger.Info("Gracefully shutting down.")
	b.Close()
}

// Check latency every 30 seconds
// https://github.com/bwmarrin/discordgo/issues/908
func (b *Bot) startLatencyMonitor() error {
	if b.Session == nil {
		return fmt.Errorf("session is not initialized")
	}

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			latency := b.Session.HeartbeatLatency()
			logger.Info(fmt.Sprintf("Current heartbeat latency: %v", latency))
		}
	}()
	return nil
}

func retryReconnection(session *discordgo.Session) {
	const maxRetries = 5
	const baseDelay = 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		err := session.Open()
		if err == nil {
			logger.Info("Reconnected successfully")
			return
		}

		logger.Error(fmt.Errorf("reconnect attempt %d failed: %v; retrying in %v seconds", i+1, err, int(math.Pow(2, float64(i)))))
		time.Sleep(time.Duration(math.Pow(2, float64(i))) * baseDelay)
	}

	logger.Info("All reconnect attempts failed")
}

// routes incoming interactions to the appropriate handler function
func (b *Bot) interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate, ghClient *github.GitHubClient) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		commandName := i.ApplicationCommandData().Name
		if handler, ok := commandHandlers[commandName]; ok {
			handler(s, i)
		} else {
			logger.Warn(fmt.Sprintf("Unknown slash command interaction received: %s", commandName))
		}

	case discordgo.InteractionMessageComponent:
		switch i.MessageComponentData().CustomID {
		case "prompt_yes":
			handleYesResponse(s, i)
		case "prompt_no":
			handleNoResponse(s, i)
		case "repository_select_menu":
			handleSelection(s, i)
		case "issue_type_select_menu":
			handleSelection(s, i)
		case "select_menu":
			handleSelectMenuResponse(s, i, b.IssueChannelID)
		case "open_modal":
			handleOpenModal(s, i)
		case "form_submit", "text_input_submit":
			handleFormSubmit(s, i, b.IssueChannelID)
		default:
			logger.Warn(fmt.Sprintf("Unknown component interaction received: %s", i.MessageComponentData().CustomID))
		}

	case discordgo.InteractionModalSubmit:
		if i.ModalSubmitData().CustomID == "issue_details_modal" {
			handleModalSubmit(s, i, ghClient)
		} else {
			logger.Warn(fmt.Sprintf("Unknown modal submission received: %s", i.ModalSubmitData().CustomID))
		}
	}
}

// registers all commands with Discord
func (b *Bot) registerCommands() error {
	for _, cmd := range getCommands() {
		if _, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.GuildID, &discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
		}); err != nil {
			return fmt.Errorf("failed to register command %s: %v", cmd.Name, err)
		}
	}
	return nil
}

// returns a list of available commands
func getCommands() []Command {
	return []Command{
		{
			Name:        "basic-command",
			Description: "Basic command",
			Handler:     handleBasicCommand,
		},
		{
			Name:        "button-prompt",
			Description: "Displays a simple Yes/No prompt",
			Handler:     handleButtonPrompt,
		},
		{
			Name:        "select-menu",
			Description: "Displays a select menu with three choices",
			Handler:     handleExtendedForm,
		},
	}
}

// maps command names to their handler functions
var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"basic-command": handleBasicCommand,
	"button-prompt": handleButtonPrompt,
	"select-menu":   handleExtendedForm,
}
