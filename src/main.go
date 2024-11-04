package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	wg := &sync.WaitGroup{} // Using a pointer to sync.WaitGroup to maintain the reference
	wg.Add(numberOfServers)

	err := configs.Init()
	if err == nil {
		logger.Debug("Environment variables loaded from environment.")
	}

	gConfigs := configs.GetGlobalConfigs()

	go network.StartServer(wg, gConfigs, &grpc_application.GrpcServer{})
	go network.StartServer(wg, gConfigs, &api.Api{})

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

	setupSignalHandling(wg)
}

func setupSignalHandling(wg *sync.WaitGroup) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shutdown(ctx, wg)
}

func shutdown(ctx context.Context, wg *sync.WaitGroup) {
	done := make(chan bool, 1)
	go func() {
		wg.Wait() // Ensure all operations are completed
		logger.Info("All servers gracefully shut down.")
		done <- true
	}()

	select {
	case <-done:
		logger.Info("Shutdown completed")
	case <-ctx.Done():
		logger.Warn("Shutdown timed out")
	}
}
