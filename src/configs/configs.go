package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/tibia-oce/discord-bot/src/logger"
)

type Config interface {
	Format() string
}

type GlobalConfigs struct {
	DBConfigs        DBConfigs
	ServerConfigs    ServerConfigs
	DiscordBotConfig DiscordBotConfigs
}

func Init() error {
	return godotenv.Load(".env")
}

func (c *GlobalConfigs) Display() {
	logger.Info(c.DBConfigs.format())
	logger.Info(c.ServerConfigs.Format())
	logger.Info(c.DiscordBotConfig.Format())
}

func GetGlobalConfigs() GlobalConfigs {
	return GlobalConfigs{
		DBConfigs:        GetDBConfigs(),
		ServerConfigs:    GetServerConfigs(),
		DiscordBotConfig: GetDiscordBotConfigs(),
	}
}

func GetEnvInt(key string, defaultValue ...int) int {
	value := os.Getenv(key)
	if len(value) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		logger.Error(err)
		return 0
	}

	return intValue
}

func GetEnvStr(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if len(value) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return value
}
