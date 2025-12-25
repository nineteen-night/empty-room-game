package gameService

import (
    "encoding/json" 
    "errors"
    "fmt"
    "strings"

    "github.com/nineteen-night/empty-room-game/internal/game/models"
)

func (s *GameService) validateGameSessions(sessions []*models.GameSession) error {
    for _, session := range sessions {
        if session.PartnershipID == 0 {
            return errors.New("partnership_id обязателен")
        }

        validStatuses := map[string]bool{
            "active":    true,
            "paused":    true,
            "completed": true,
        }

        if !validStatuses[session.Status] {
            return fmt.Errorf("некорректный статус: %s", session.Status)
        }

        if strings.TrimSpace(session.CurrentRoom) == "" {
            return errors.New("current_room не может быть пустым")
        }
    }
    return nil
}

func (s *GameService) validateGameStates(states []*models.GameState) error {
    for _, state := range states {
        if state.GameSessionID == 0 {
            return errors.New("game_session_id обязателен")
        }

        if state.Inventory != "" && !isValidJSON(state.Inventory) {
            return errors.New("inventory должен быть валидным JSON")
        }
        if state.PuzzlesSolved != "" && !isValidJSON(state.PuzzlesSolved) {
            return errors.New("puzzles_solved должен быть валидным JSON")
        }
    }
    return nil
}

func isValidJSON(str string) bool {
    var js interface{}
    return json.Unmarshal([]byte(str), &js) == nil
}