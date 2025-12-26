package authService

import (
    "context"
    "errors"
    "fmt"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

func (s *AuthService) CreatePartnership(ctx context.Context, currentUserID, partnerUsername string) (*models.Partnership, error) {
    if len(partnerUsername) < s.minNameLen || len(partnerUsername) > s.maxNameLen {
        return nil, fmt.Errorf("имя партнёра должно быть от %d до %d символов",
            s.minNameLen, s.maxNameLen)
    }

    partner, err := s.authStorage.GetUserByUsername(ctx, partnerUsername)
    if err != nil {
        return nil, fmt.Errorf("ошибка получения партнёра: %w", err)
    }
    if partner == nil {
        return nil, errors.New("партнёр не найден")
    }

    if err := s.ValidatePartnershipCreation(currentUserID, partner.ID); err != nil {
        return nil, err
    }

    existingPartnership, err := s.authStorage.GetPartnershipBetweenUsers(ctx, currentUserID, partner.ID)
    if err != nil {
        return nil, fmt.Errorf("ошибка проверки существующего партнёрства: %w", err)
    }
    
    if existingPartnership != nil {
        return nil, errors.New("партнёрство уже существует")
    }

    partnershipID, err := s.authStorage.CreatePartnership(ctx, currentUserID, partner.ID)
    if err != nil {
        return nil, fmt.Errorf("ошибка создания партнёрства: %w", err)
    }

    if s.eventSender != nil {
        if err := s.eventSender.SendPartnershipCreated(ctx, partnershipID, currentUserID, partner.ID); err != nil {
            fmt.Printf("не удалось отправить событие: %v\n", err)
        }
    }
    
    return &models.Partnership{
        ID:       partnershipID,
        User1ID:  currentUserID,
        User2ID:  partner.ID,
    }, nil
}

func (s *AuthService) TerminatePartnership(ctx context.Context, partnershipID string) error {
    if err := s.ValidatePartnershipTermination(partnershipID); err != nil {
        return err
    }

    if err := s.authStorage.TerminatePartnership(ctx, partnershipID); err != nil {
        return fmt.Errorf("ошибка расторжения партнёрства: %w", err)
    }

    if s.eventSender != nil {
        if err := s.eventSender.SendPartnershipTerminated(ctx, partnershipID); err != nil {
            fmt.Printf("Предупреждение: не удалось отправить событие: %v\n", err)
        }
    }
    
    return nil
}