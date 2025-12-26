package bootstrap

import (
	"fmt"

	"github.com/nineteen-night/empty-room-game/config"
	gameprocessor "github.com/nineteen-night/empty-room-game/internal/game/services/processors/game_processor"
	gameservice "github.com/nineteen-night/empty-room-game/internal/game/services/gameService"
)

func InitGameProcessor(gameService *gameservice.GameService, cfg *config.Config) *gameprocessor.GameProcessor {
	kafkaBrokers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}
	return gameprocessor.NewGameProcessor(gameService, kafkaBrokers, cfg.Kafka.RoomCompletedEventsTopic)
}