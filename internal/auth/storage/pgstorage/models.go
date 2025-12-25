package pgstorage

type User struct {
    ID           uint64 `db:"id"`
    Username     string `db:"username"`
    Email        string `db:"email"`
    PasswordHash string `db:"password_hash"`
}

type Partnership struct {
    ID        uint64 `db:"id"`
    Player1ID uint64 `db:"player1_id"`
    Player2ID uint64 `db:"player2_id"`
    Status    string `db:"status"`
}

const (
    usersTable = "users"
    partnershipsTable = "partnerships"
    
    userIDColumn = "id"
    usernameColumn = "username"
    emailColumn = "email"
    passwordHashColumn = "password_hash"
    
    partnershipIDColumn = "id"
    player1IDColumn = "player1_id"
    player2IDColumn = "player2_id"
    statusColumn = "status"
)