package authService

import (
	"context"
	"errors"
	"testing"

	"github.com/nineteen-night/empty-room-game/internal/auth/models"
	"github.com/nineteen-night/empty-room-game/internal/auth/services/authService/mocks"
	"github.com/stretchr/testify/suite"
	"gotest.tools/v3/assert"
)

type PartnershipServiceSuite struct {
	suite.Suite
	ctx          context.Context
	authStorage  *mocks.AuthStorage
	eventSender  *mocks.EventSender
	authService  *AuthService
}

func (s *PartnershipServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.authStorage = mocks.NewAuthStorage(s.T())
	s.eventSender = mocks.NewEventSender(s.T())
	s.authService = NewAuthService(s.ctx, s.authStorage, 3, 50)
	s.authService.SetEventSender(s.eventSender)
}

//успешное создание партнёрства
func (s *PartnershipServiceSuite) TestCreatePartnership_Success() {
	currentUserID := "user-1"
	partnerUsername := "partneruser"
	partnerID := "user-2"
	partnershipID := "partnership-123"

	partner := &models.User{
		ID:             partnerID,
		Username:       partnerUsername,
		Email:          "partner@example.com",
		PasswordHash:   "hash",
		MaxRoomReached: 0,
	}

	s.authStorage.EXPECT().GetUserByUsername(s.ctx, partnerUsername).Return(partner, nil)
	s.authStorage.EXPECT().GetPartnershipBetweenUsers(s.ctx, currentUserID, partnerID).Return(nil, nil)
	s.authStorage.EXPECT().CreatePartnership(s.ctx, currentUserID, partnerID).Return(partnershipID, nil)
	s.eventSender.EXPECT().SendPartnershipCreated(s.ctx, partnershipID, currentUserID, partnerID).Return(nil)

	partnership, err := s.authService.CreatePartnership(s.ctx, currentUserID, partnerUsername)

	assert.NilError(s.T(), err)
	assert.Equal(s.T(), partnershipID, partnership.ID)
	assert.Equal(s.T(), currentUserID, partnership.User1ID)
	assert.Equal(s.T(), partnerID, partnership.User2ID)
}

//несуществующий партнёр
func (s *PartnershipServiceSuite) TestCreatePartnership_PartnerNotFound() {
	currentUserID := "user-1"
	partnerUsername := "nonexistent"

	s.authStorage.EXPECT().GetUserByUsername(s.ctx, partnerUsername).Return(nil, nil)

	partnership, err := s.authService.CreatePartnership(s.ctx, currentUserID, partnerUsername)

	assert.Error(s.T(), err, "партнёр не найден")
	assert.Assert(s.T(), partnership == nil)
}

//создание партнёрства с самим собой
func (s *PartnershipServiceSuite) TestCreatePartnership_WithSelf() {
	currentUserID := "user-1"
	partnerUsername := "currentuser"

	partner := &models.User{
		ID:             currentUserID,
		Username:       partnerUsername,
		Email:          "current@example.com",
		PasswordHash:   "hash",
		MaxRoomReached: 0,
	}

	s.authStorage.EXPECT().GetUserByUsername(s.ctx, partnerUsername).Return(partner, nil)

	partnership, err := s.authService.CreatePartnership(s.ctx, currentUserID, partnerUsername)

	assert.Error(s.T(), err, "нельзя создать партнёрство с самим собой")
	assert.Assert(s.T(), partnership == nil)
}

//существующее партнёрство
func (s *PartnershipServiceSuite) TestCreatePartnership_AlreadyExists() {
	currentUserID := "user-1"
	partnerUsername := "partneruser"
	partnerID := "user-2"

	partner := &models.User{
		ID:             partnerID,
		Username:       partnerUsername,
		Email:          "partner@example.com",
		PasswordHash:   "hash",
		MaxRoomReached: 0,
	}

	existingPartnership := &models.Partnership{
		ID:      "existing-partnership",
		User1ID: currentUserID,
		User2ID: partnerID,
	}

	s.authStorage.EXPECT().GetUserByUsername(s.ctx, partnerUsername).Return(partner, nil)
	s.authStorage.EXPECT().GetPartnershipBetweenUsers(s.ctx, currentUserID, partnerID).Return(existingPartnership, nil)

	partnership, err := s.authService.CreatePartnership(s.ctx, currentUserID, partnerUsername)

	assert.Error(s.T(), err, "партнёрство уже существует")
	assert.Assert(s.T(), partnership == nil)
}

//слишком короткое имя партнёра
func (s *PartnershipServiceSuite) TestCreatePartnership_UsernameTooShort() {
	currentUserID := "user-1"
	partnerUsername := "ab"

	partnership, err := s.authService.CreatePartnership(s.ctx, currentUserID, partnerUsername)

	assert.ErrorContains(s.T(), err, "имя партнёра должно быть от 3 до 50 символов")
	assert.Assert(s.T(), partnership == nil)
}

//ошибка отправки события
func (s *PartnershipServiceSuite) TestCreatePartnership_EventSendError() {
	currentUserID := "user-1"
	partnerUsername := "partneruser"
	partnerID := "user-2"
	partnershipID := "partnership-123"
	eventErr := errors.New("kafka error")

	partner := &models.User{
		ID:             partnerID,
		Username:       partnerUsername,
		Email:          "partner@example.com",
		PasswordHash:   "hash",
		MaxRoomReached: 0,
	}

	s.authStorage.EXPECT().GetUserByUsername(s.ctx, partnerUsername).Return(partner, nil)
	s.authStorage.EXPECT().GetPartnershipBetweenUsers(s.ctx, currentUserID, partnerID).Return(nil, nil)
	s.authStorage.EXPECT().CreatePartnership(s.ctx, currentUserID, partnerID).Return(partnershipID, nil)
	s.eventSender.EXPECT().SendPartnershipCreated(s.ctx, partnershipID, currentUserID, partnerID).Return(eventErr)

	partnership, err := s.authService.CreatePartnership(s.ctx, currentUserID, partnerUsername)

	assert.NilError(s.T(), err)
	assert.Equal(s.T(), partnershipID, partnership.ID)
}

//успешное расторжение партнёрства
func (s *PartnershipServiceSuite) TestTerminatePartnership_Success() {
	partnershipID := "partnership-123"

	s.authStorage.EXPECT().TerminatePartnership(s.ctx, partnershipID).Return(nil)
	s.eventSender.EXPECT().SendPartnershipTerminated(s.ctx, partnershipID).Return(nil)

	err := s.authService.TerminatePartnership(s.ctx, partnershipID)

	assert.NilError(s.T(), err)
}

//тестирует несуществующее партнёрство
func (s *PartnershipServiceSuite) TestTerminatePartnership_NotFound() {
	partnershipID := "non-existent"
	notFoundErr := errors.New("partnership not found")

	s.authStorage.EXPECT().TerminatePartnership(s.ctx, partnershipID).Return(notFoundErr)

	err := s.authService.TerminatePartnership(s.ctx, partnershipID)

	assert.ErrorContains(s.T(), err, "ошибка расторжения партнёрства")
	assert.ErrorContains(s.T(), err, notFoundErr.Error())
}

//тестирует пустой ID
func (s *PartnershipServiceSuite) TestTerminatePartnership_EmptyID() {
	partnershipID := ""

	err := s.authService.TerminatePartnership(s.ctx, partnershipID)

	assert.Error(s.T(), err, "ID партнёрства не может быть пустым")
}

//тестирует ошибку отправки события
func (s *PartnershipServiceSuite) TestTerminatePartnership_EventSendError() {
	partnershipID := "partnership-123"
	eventErr := errors.New("kafka error")

	s.authStorage.EXPECT().TerminatePartnership(s.ctx, partnershipID).Return(nil)
	s.eventSender.EXPECT().SendPartnershipTerminated(s.ctx, partnershipID).Return(eventErr)

	err := s.authService.TerminatePartnership(s.ctx, partnershipID)

	assert.NilError(s.T(), err)
}

func TestPartnershipServiceSuite(t *testing.T) {
	suite.Run(t, new(PartnershipServiceSuite))
}