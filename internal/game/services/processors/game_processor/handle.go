package game_processor

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/segmentio/kafka-go"
)

func (p *GameProcessor) HandlePartnershipCreated(ctx context.Context, partnershipID, user1ID, user2ID string) error {
    return p.gameService.HandlePartnershipCreated(ctx, partnershipID, user1ID, user2ID)
}

func (p *GameProcessor) HandlePartnershipTerminated(ctx context.Context, partnershipID string) error {
    return p.gameService.HandlePartnershipTerminated(ctx, partnershipID)
}

func (p *GameProcessor) SendRoomCompleted(ctx context.Context, userID string, roomNumber int32) error {
    event := map[string]interface{}{
        "event_type":  "room_completed",
        "user_id":     userID,
        "room_number": roomNumber,
    }

    data, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %w", err)
    }

    err = p.roomCompletedWriter.WriteMessages(ctx, kafka.Message{
        Value: data,
    })

    if err != nil {
        return fmt.Errorf("failed to write message to kafka: %w", err)
    }

    return nil
}