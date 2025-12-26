package gameService

import (
	"errors"
)

func (s *GameService) ValidatePartnershipID(partnershipID string) error {
	if partnershipID == "" {
		return errors.New("ID партнёрства не может быть пустым")
	}
	return nil
}

func (s *GameService) ValidateGameSessionCreation(partnershipID, user1ID, user2ID string) error {
	if partnershipID == "" {
		return errors.New("ID партнёрства не может быть пустым")
	}

	if user1ID == "" || user2ID == "" {
		return errors.New("ID пользователей не могут быть пустыми")
	}

	if user1ID == user2ID {
		return errors.New("ID пользователей не могут совпадать")
	}

	return nil
}

func (s *GameService) ValidateRoomNumber(roomNumber int32) error {
	if roomNumber < 1 {
		return errors.New("номер комнаты должен быть положительным")
	}
	return nil
}