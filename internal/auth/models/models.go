package models

type User struct {
    ID           uint64
    Username     string
    Email        string
    Password     string
    PasswordHash string
}

type Partnership struct {
    ID        uint64
    Player1ID uint64
    Player2ID uint64
    Status    string
}