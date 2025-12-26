package main

import (
    "fmt"
    "os"

    "github.com/nineteen-night/empty-room-game/config"
    "github.com/nineteen-night/empty-room-game/internal/auth/bootstrap"
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

    authstorage := bootstrap.InitPGStorage(cfg)
    defer authstorage.Close()
    
    authservice := bootstrap.InitAuthService(authstorage, cfg)
    authprocessor := bootstrap.InitAuthProcessor(authservice, cfg)
    authservice.SetEventSender(authprocessor)
    
    roomcompletedconsumer := bootstrap.InitRoomCompletedConsumer(cfg, authprocessor)
    authapi := bootstrap.InitAuthServiceAPI(authservice)

    bootstrap.AppRun(*authapi, roomcompletedconsumer)
}