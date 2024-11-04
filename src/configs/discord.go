package configs

import (
	"fmt"
)

const (
	EnvDiscordTokenKey = "DISCORD_TOKEN"
	EnvDiscordGuildID  = "DISCORD_SERVER_ID"
	EnvDiscordAppID    = "DISCORD_APP_ID"
)

type DiscordBotConfigs struct {
	Token   string
	GuildID string
	AppID   string
	Config
}

func (botConfigs *DiscordBotConfigs) Format() string {
	return fmt.Sprintf(
		"Discord Bot Token: %s",
		botConfigs.Token[:10],
	)
}

func GetDiscordBotConfigs() DiscordBotConfigs {
	return DiscordBotConfigs{
		Token:   GetEnvStr(EnvDiscordTokenKey, ""),
		GuildID: GetEnvStr(EnvDiscordGuildID, ""),
		AppID:   GetEnvStr(EnvDiscordAppID, ""),
	}
}
