package auth_processor

import (
    "context"

    "github.com/segmentio/kafka-go"
)

type authService interface {
    HandleRoomCompleted(ctx context.Context, userID string, roomNumber int32) error
}

type EventSender interface {
    SendPartnershipCreated(ctx context.Context, partnershipID, user1ID, user2ID string) error
    SendPartnershipTerminated(ctx context.Context, partnershipID string) error
}

type AuthProcessor struct {
    authService       authService
    partnershipWriter *kafka.Writer
}

func NewAuthProcessor(authService authService, kafkaBrokers []string) *AuthProcessor {
    partnershipWriter := &kafka.Writer{
        Addr:     kafka.TCP(kafkaBrokers...),
        Topic:    "partnership_events",
        Balancer: &kafka.LeastBytes{},
    }
    
    return &AuthProcessor{
        authService:       authService,
        partnershipWriter: partnershipWriter,
    }
}