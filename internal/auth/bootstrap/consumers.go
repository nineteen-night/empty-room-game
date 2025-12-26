package bootstrap

import (
    "fmt"

    "github.com/nineteen-night/empty-room-game/config"
    authconsumer "github.com/nineteen-night/empty-room-game/internal/auth/consumer/auth_consumer"
    authprocessor "github.com/nineteen-night/empty-room-game/internal/auth/services/processors/auth_processor"
)

func InitRoomCompletedConsumer(cfg *config.Config, authProcessor *authprocessor.AuthProcessor) *authconsumer.AuthConsumer {
    kafkaBrokers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}
    return authconsumer.NewAuthConsumer(authProcessor, kafkaBrokers, cfg.Kafka.RoomCompletedEventsTopic)
}