package models

type Room struct {
	RoomNumber  int32
	Name        string
	Description string
}

type GameSession struct {
	PartnershipID string
	User1ID       string
	User2ID       string
	CurrentRoom   int32
}

type GameState struct {
	CurrentRoom int32
	RoomInfo    *Room
}