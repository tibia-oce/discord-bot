package main

import (
	"sync"
	"time"

	"github.com/tibia-oce/discord-bot/src/api"
	"github.com/tibia-oce/discord-bot/src/configs"
	"github.com/tibia-oce/discord-bot/src/discord"
	grpc_application "github.com/tibia-oce/discord-bot/src/grpc"
	"github.com/tibia-oce/discord-bot/src/logger"
	"github.com/tibia-oce/discord-bot/src/network"
)

var numberOfServers = 3
var initDelay = 200

func main() {
	logger.Init(configs.GetLogLevel())
	logger.Info("Loading configurations...")

	var wg sync.WaitGroup
	wg.Add(numberOfServers)

	err := configs.Init()
	if err == nil {
		logger.Debug("Environment variables loaded from environment.")
	}

	gConfigs := configs.GetGlobalConfigs()

	go network.StartServer(&wg, gConfigs, &grpc_application.GrpcServer{})
	go network.StartServer(&wg, gConfigs, &api.Api{})

	go func() {
		defer wg.Done()
		bot := &discord.Bot{
			Token:   gConfigs.DiscordBotConfig.Token,
			GuildID: gConfigs.DiscordBotConfig.GuildID,
		}
		if err := bot.Init(); err != nil {
			logger.Error(err)
			return
		}
		defer bot.Close()
	}()

	time.Sleep(time.Duration(initDelay) * time.Millisecond)
	gConfigs.Display()

	wg.Wait()
	logger.Info("Good bye...")
}
