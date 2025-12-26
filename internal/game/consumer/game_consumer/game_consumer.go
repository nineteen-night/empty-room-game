package game_consumer

import (
	"context"
)

type gameProcessor interface {
	HandlePartnershipCreated(ctx context.Context, partnershipID, user1ID, user2ID string) error
	HandlePartnershipTerminated(ctx context.Context, partnershipID string) error
}

type GameConsumer struct {
	gameProcessor gameProcessor
	kafkaBrokers  []string
	topicName     string
}

func NewGameConsumer(gameProcessor gameProcessor, kafkaBrokers []string, topicName string) *GameConsumer {
	return &GameConsumer{
		gameProcessor: gameProcessor,
		kafkaBrokers:  kafkaBrokers,
		topicName:     topicName,
	}
}