package main

import (
	"fmt"
	"os"

	"github.com/nineteen-night/empty-room-game/config"
	"github.com/nineteen-night/empty-room-game/internal/game/bootstrap"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	gamestorage := bootstrap.InitPGStorage(cfg)
	defer gamestorage.Close()

	gameservice := bootstrap.InitGameService(gamestorage, cfg)
	gameprocessor := bootstrap.InitGameProcessor(gameservice, cfg)
	gameservice.SetEventSender(gameprocessor)

	partnershipconsumer := bootstrap.InitPartnershipEventsConsumer(cfg, gameprocessor)
	gameapi := bootstrap.InitGameServiceAPI(gameservice)

	bootstrap.AppRun(*gameapi, partnershipconsumer)
}