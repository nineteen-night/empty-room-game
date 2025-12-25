package models

type GameSession struct {
    ID              uint64
    PartnershipID   uint64
    CurrentRoom     string
    Status          string 
    CurrentPlayerID uint64
}

type GameState struct {
    ID            uint64
    GameSessionID uint64
    Inventory     string 
    PuzzlesSolved string 
    CurrentRoomID uint64
}