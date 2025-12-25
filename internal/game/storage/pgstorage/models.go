package pgstorage

type GameSession struct {
    ID              uint64 `db:"id"`
    PartnershipID   uint64 `db:"partnership_id"`
    CurrentRoom     string `db:"current_room"`
    Status          string `db:"status"`
    CurrentPlayerID uint64 `db:"current_player_id"`
}

type GameState struct {
    ID            uint64 `db:"id"`
    GameSessionID uint64 `db:"game_session_id"`
    Inventory     string `db:"inventory"`
    PuzzlesSolved string `db:"puzzles_solved"`
    CurrentRoomID uint64 `db:"current_room_id"`
}

const (
    gameSessionsTable = "game_sessions"
    gameStatesTable   = "game_states"
    
    gameSessionIDColumn   = "id"
    partnershipIDColumn   = "partnership_id"
    currentRoomColumn     = "current_room"
    statusColumn          = "status"
    currentPlayerIDColumn = "current_player_id"
    
    gameStateIDColumn     = "id"
    gameSessionIDColumn2  = "game_session_id"
    inventoryColumn       = "inventory"
    puzzlesSolvedColumn   = "puzzles_solved"
    currentRoomIDColumn   = "current_room_id"
)