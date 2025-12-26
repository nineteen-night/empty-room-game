package bootstrap

import (
	"fmt"

	"github.com/nineteen-night/empty-room-game/config"
	gameconsumer "github.com/nineteen-night/empty-room-game/internal/game/consumer/game_consumer"
	gameprocessor "github.com/nineteen-night/empty-room-game/internal/game/services/processors/game_processor"
)

func InitPartnershipEventsConsumer(cfg *config.Config, gameProcessor *gameprocessor.GameProcessor) *gameconsumer.GameConsumer {
	kafkaBrokers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}
	return gameconsumer.NewGameConsumer(gameProcessor, kafkaBrokers, cfg.Kafka.PartnershipEventsTopic)
}