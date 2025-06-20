package entities

import "time"

type Match struct {
	ID        string
	Players   []Player
	GameMode  string
	StartedAt *time.Time
	CreatedAt time.Time
}

type Player struct {
	ID string
}
