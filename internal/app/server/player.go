package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	Id      string
	Conn    *websocket.Conn
	MatchId string
	Status  Status
	Result  float64

	mu *sync.Mutex
}

type Status uint8

const (
	INIT Status = iota
	CONNECTED
	DISCONNECTED
)
