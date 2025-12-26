package game_processor

import (
    "context"

    "github.com/nineteen-night/empty-room-game/internal/game/models"
    "github.com/segmentio/kafka-go"
)

type gameService interface {
    GetGameState(ctx context.Context, partnershipID string) (*models.GameState, error)
    MoveToNextRoom(ctx context.Context, partnershipID string) (*models.GameState, error)

    HandlePartnershipCreated(ctx context.Context, partnershipID, user1ID, user2ID string) error
    HandlePartnershipTerminated(ctx context.Context, partnershipID string) error
}

type GameProcessor struct {
    gameService         gameService  //
    roomCompletedWriter *kafka.Writer
}

func NewGameProcessor(gameService gameService, kafkaBrokers []string, topicName string) *GameProcessor {
    roomCompletedWriter := &kafka.Writer{
        Addr:     kafka.TCP(kafkaBrokers...),
        Topic:    topicName,
        Balancer: &kafka.LeastBytes{},
    }

    return &GameProcessor{
        gameService:         gameService,
        roomCompletedWriter: roomCompletedWriter,
    }
}