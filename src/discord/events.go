// discord/events.go
package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tibia-oce/discord-bot/src/logger"
)

func (b *Bot) interactionHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	}
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
		logger.Warn(fmt.Sprintf("Bot disconnected. Reconnecting... Reason: %v", d))
		for {
			err := s.Open()
			if err == nil {
				logger.Info("Reconnected successfully")
				break
			}
			logger.Error(fmt.Errorf("reconnect failed: %v", err))
			time.Sleep(5 * time.Second)
		}
	}
}
