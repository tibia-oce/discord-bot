// discord/discord.go
package discord

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/tibia-oce/discord-bot/src/logger"
)

type Bot struct {
	Session *discordgo.Session
	Token   string
	GuildID string
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

	b.Session.AddHandler(b.interactionHandler())
	b.Session.AddHandler(b.onReady())
	b.Session.AddHandler(b.onReconnect())
	b.Session.AddHandler(b.onDisconnect())

	// Open WebSocket connection
	err = b.Session.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord connection: %v", err)
	}

	logger.Info("Bot is now running!")

	return b.registerCommands()
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
