package bootstrap

import (
    "fmt"
    "github.com/nineteen-night/empty-room-game/config"
    gameConsumer "github.com/nineteen-night/empty-room-game/internal/game/consumer/game_consumer"
    gameProcessor "github.com/nineteen-night/empty-room-game/internal/game/services/processors/game_processor"
)

func InitGameSessionsConsumer(cfg *config.Config, gameProcessor *gameProcessor.GameProcessor) *gameConsumer.GameConsumer {
    kafkaBrokers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}
    return gameConsumer.NewGameConsumer(gameProcessor, kafkaBrokers, cfg.Kafka.GameSessionsUpsertTopic)
}

func InitGameStatesConsumer(cfg *config.Config, gameProcessor *gameProcessor.GameProcessor) *gameConsumer.GameConsumer {
    kafkaBrokers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}
    return gameConsumer.NewGameConsumer(gameProcessor, kafkaBrokers, cfg.Kafka.GameStatesUpsertTopic)
}