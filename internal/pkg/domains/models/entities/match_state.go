package entities

import "time"

type MatchState struct {
	ID           string
	MatchID      string
	PlayerStates []PlayerState
	GameState    string
	Move         Move
	Ply          int
	Timestamp    time.Time
}

type PlayerState struct {
	Clock  string
	Status string
}

type Move struct {
	PlayerID string
	Uci      string
}
