package authService

import (
    "errors"
    "fmt"
    "net/mail"
    "strings"

)

// Валидация регистрации
func (s *AuthService) ValidateRegistration(username, email, password string) error {
    if len(username) < s.minNameLen || len(username) > s.maxNameLen {
        return fmt.Errorf("имя пользователя должно быть от %d до %d символов", 
            s.minNameLen, s.maxNameLen)
    }

    if len(password) < 6 {
        return errors.New("пароль должен быть не менее 6 символов")
    }

    if !s.IsValidEmail(email) {
        return fmt.Errorf("некорректный email: %s", email)
    }

    return nil
}

// Валидация входа
func (s *AuthService) ValidateLogin(username, password string) error {
    if len(username) < s.minNameLen || len(username) > s.maxNameLen {
        return fmt.Errorf("имя пользователя должно быть от %d до %d символов",
            s.minNameLen, s.maxNameLen)
    }

    if len(password) < 1 {
        return errors.New("пароль не может быть пустым")
    }

    return nil
}

// Валидация создания партнёрства
func (s *AuthService) ValidatePartnershipCreation(user1ID, user2ID string) error {
    if user1ID == "" || user2ID == "" {
        return errors.New("ID пользователей не могут быть пустыми")
    }

    if user1ID == user2ID {
        return errors.New("нельзя создать партнёрство с самим собой")
    }

    return nil
}

// Валидация расторжения партнёрства
func (s *AuthService) ValidatePartnershipTermination(partnershipID string) error {
    if partnershipID == "" {
        return errors.New("ID партнёрства не может быть пустым")
    }
    return nil
}

// Проверка email
func (s *AuthService) IsValidEmail(email string) bool {
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

	// Проверяем длину local part (до @)
	if len(parts[0]) == 0 || len(parts[0]) > 64 {
		return false
	}

	// Проверяем длину domain part (после @)
	if len(parts[1]) == 0 || len(parts[1]) > 253 {
		return false
	}

	return true
}

// ValidateRoomCompletion - валидация завершения комнаты
func (s *AuthService) ValidateRoomCompletion(userID string, roomNumber int32) error {
    if userID == "" {
        return errors.New("ID пользователя не может быть пустым")
    }

    if roomNumber < 1 {
        return errors.New("номер комнаты должен быть положительным")
    }

    return nil
}