package auth_consumer

import (
    "context"
)

type authProcessor interface {
    HandleRoomCompleted(ctx context.Context, userID string, roomNumber int32) error
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