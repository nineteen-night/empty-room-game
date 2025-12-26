package models

type User struct {
    ID              string
    Username        string
    Email           string
    Password        string
    PasswordHash    string
    MaxRoomReached  int32
}

type Partnership struct {
    ID       string
    User1ID  string
    User2ID  string
}