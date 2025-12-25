package authService

import (
    "errors"
    "fmt"
    "net/mail"
    "strings"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

func (s *AuthService) validateUsers(users []*models.User) error {
    for _, user := range users {
        if len(user.Username) < s.minNameLen || len(user.Username) > s.maxNameLen {
            return fmt.Errorf("имя пользователя должно быть от %d до %d символов", 
                s.minNameLen, s.maxNameLen)
        }

        if user.ID == 0 && user.Password == "" && user.PasswordHash == "" {
            return errors.New("пароль обязателен для новых пользователей")
        }

        if user.Password != "" && len(user.Password) < 6 {
            return errors.New("пароль должен быть не менее 6 символов")
        }

        if !s.isValidEmail(user.Email) {
            return fmt.Errorf("некорректный email: %s", user.Email)
        }
    }
    return nil
}

func (s *AuthService) validatePartnerships(partnerships []*models.Partnership) error {
    for _, partnership := range partnerships {
        if partnership.Player1ID == partnership.Player2ID {
            return errors.New("игрок 1 и игрок 2 не могут быть одним и тем же пользователем")
        }

        validStatuses := map[string]bool{
            "pending":   true,
            "active":    true,
            "completed": true,
        }

        if !validStatuses[partnership.Status] {
            return fmt.Errorf("некорректный статус: %s", partnership.Status)
        }
    }
    return nil
}

func (s *AuthService) isValidEmail(email string) bool {
    if len(email) < 3 || len(email) > 254 {
        return false
    }

    _, err := mail.ParseAddress(email)
    if err != nil {
        return false
    }

    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }

    if len(parts[1]) == 0 || len(parts[1]) > 253 {
        return false
    }

    return true
}