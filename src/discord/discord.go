// discord/discord.go
package discord

import (
	"fmt"
	"math"
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
	b.Session.AddHandler(b.interactionHandler)

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
func (b *Bot) interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commandName := i.ApplicationCommandData().Name
	userID := i.Interaction.Member.User.Username
	channelID := i.ChannelID
	logger.Info(fmt.Sprintf("Command triggered: %s, User ID: %s, Channel ID: %s", commandName, userID, channelID))
	if handler, ok := commandHandlers[commandName]; ok {
		handler(s, i)
	} else {
		logger.Warn("Unknown command interaction received.")
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
		// Add more commands here
	}
}

// maps command names to their handler functions
var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"basic-command": handleBasicCommand,
}
