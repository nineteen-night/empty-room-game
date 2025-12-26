package pgstorage

type User struct {
    ID              string `db:"id"`
    Username        string `db:"username"`
    Email           string `db:"email"`
    PasswordHash    string `db:"password_hash"`
    MaxRoomReached  int32  `db:"max_room_reached"`
}

type Partnership struct {
    ID       string `db:"id"`
    User1ID  string `db:"user1_id"`
    User2ID  string `db:"user2_id"`
    Status   string `db:"status"`
}

const (
    usersTable       = "users"
    partnershipsTable = "partnerships"
    
    userIDColumn         = "id"
    usernameColumn       = "username"
    emailColumn          = "email"
    passwordHashColumn   = "password_hash"
    maxRoomReachedColumn = "max_room_reached"
    
    partnershipIDColumn = "id"
    user1IDColumn       = "user1_id"
    user2IDColumn       = "user2_id"
    statusColumn        = "status"
)