package database

import (
	"database/sql"
	"github.com/tibia-oce/discord-bot/src/configs"
	"github.com/tibia-oce/discord-bot/src/logger"
)

const (
	DefaultMaxDbOpenConns = 100
)

func PullConnection(gConfigs configs.GlobalConfigs) *sql.DB {
	DB, err := sql.Open("mysql", gConfigs.DBConfigs.GetConnectionString())
	if err != nil {
		logger.Panic(err)
	}

	DB.SetMaxOpenConns(DefaultMaxDbOpenConns)

	return DB
}
