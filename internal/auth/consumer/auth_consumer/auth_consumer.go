package auth_consumer

import (
    "context"
    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

type authProcessor interface {
    HandleUsers(ctx context.Context, user *models.User) error
    HandlePartnerships(ctx context.Context, partnership *models.Partnership) error
}

type AuthConsumer struct {
    authProcessor authProcessor
    kafkaBroker   []string
    topicName     string
}

func NewAuthConsumer(authProcessor authProcessor, kafkaBroker []string, topicName string) *AuthConsumer {
    return &AuthConsumer{
        authProcessor: authProcessor,
        kafkaBroker:   kafkaBroker,
        topicName:     topicName,
    }
}