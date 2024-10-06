package network

import (
	"fmt"
	"github.com/tibia-oce/discord-bot/src/configs"
	"github.com/tibia-oce/discord-bot/src/logger"
	"sync"
)

type ServerInterface interface {
	Initialize(gConfigs configs.GlobalConfigs) error
	Run(globalConfigs configs.GlobalConfigs) error
	GetName() string
}

func StartServer(
	wg *sync.WaitGroup,
	gConfigs configs.GlobalConfigs,
	server ServerInterface,
) {
	if err := server.Initialize(gConfigs); err == nil {
		logger.Info(fmt.Sprintf("Starting %s server...", server.GetName()))
		logger.Error(server.Run(gConfigs))
	}

	wg.Done()
	logger.Warn(fmt.Sprintf("Server %s is gone...", server.GetName()))
}
