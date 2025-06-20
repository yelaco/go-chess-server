package server

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/yelaco/gchess-server/pkg/logging"
	"github.com/yelaco/gchess-server/pkg/util"
	"go.uber.org/zap"
)

type Match struct {
	Id      string
	Players map[string]Player
	moveCh  chan Move

	endCallback   func(Match)
	saveCallback  func(Match)
	abortCallback func(Match)

	ended bool
	mu    *sync.Mutex
}

type playerStatusResponse struct {
	Type     string `json:"type"`
	PlayerId string `json:"playerId"`
	Status   string `json:"status"`
}

type errorResponse struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}

func NewMatch(id string, players map[string]Player) Match {
	return Match{
		Id:      id,
		Players: players,
		moveCh:  make(chan Move),
		mu:      new(sync.Mutex),
	}
}

func (m *Match) start() {
	for move := range m.moveCh {
		_, exist := m.Players[move.PlayerId]
		if !exist {
			// player.WriteJson(errorResponse{
			// 	Type:  "error",
			// 	Error: ErrStatusInvalidPlayerId,
			// })
			continue
		}
		// m.handler.HandleMove(player, move)
	}
}

func (m *Match) Abort() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ended {
		return
	}
	m.ended = true
	if !util.IsClosed(m.moveCh) {
		close(m.moveCh)
	}
	// m.DisconnectPlayers("match aborted", time.Now().Add(5*time.Second))
	// m.handler.OnMatchAbort()
	// m.abortCallback(m)
}

func (m *Match) Save() {
	// m.handler.OnMatchSave()
	// m.saveCallback(m)
}

func (m *Match) End() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ended {
		return
	}
	m.ended = true
	if !util.IsClosed(m.moveCh) {
		close(m.moveCh)
	}
	// m.handler.OnMatchEnd()
	// m.DisconnectPlayers("match ended", time.Now().Add(5*time.Second))
	// m.endCallback(m)
}

func (m *Match) ProcessMove(move Move) {
	m.moveCh <- move
}

func (m *Match) GetPlayerWithId(id string) (Player, bool) {
	player, exist := m.Players[id]
	return player, exist
}

func (m *Match) IsEnded() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.ended
}

func (m *Match) playerJoin(playerId string, conn *websocket.Conn) {
	if m == nil {
		return
	}

	_, exist := m.GetPlayerWithId(playerId)
	if !exist {
		logging.Fatal("invalid player id", zap.String("player_id", playerId))
		return
	}

	// init, err := m.handler.OnPlayerJoin(player)
	// if err != nil {
	// 	logging.Error("on player join", zap.Error(err))
	// }
	// if init {
	// 	err := storageClient.UpdateActiveMatch(
	// 		context.Background(),
	// 		m.GetId(),
	// 		storage.ActiveMatchUpdateOptions{
	// 			StartedAt: aws.Time(time.Now()),
	// 		},
	// 	)
	// 	if err != nil {
	// 		logging.Error("failed to update match", zap.Error(err))
	// 	}
	// }

	// player.setConn(conn)
	// m.handler.OnPlayerSync(player)

	// m.notifyAboutPlayerStatus(playerStatusResponse{
	// 	Type:     "playerStatus",
	// 	PlayerId: playerId,
	// 	Status:   player.GetStatus(),
	// })
}

func (m *Match) playerDisconnect(playerId string) {
	if m == nil {
		return
	}

	_, exist := m.GetPlayerWithId(playerId)
	if !exist {
		logging.Fatal("invalid player id", zap.String("player_id", playerId))
		return
	}
	// player.setConn(nil)

	// m.handler.OnPlayerLeave(player)

	// m.notifyAboutPlayerStatus(playerStatusResponse{
	// 	Type:     "playerStatus",
	// 	PlayerId: playerId,
	// 	Status:   player.Status,
	// })
}

// func (m *Match) notifyAboutPlayerStatus(resp playerStatusResponse) {
// 	for _, player := range m.Players {
// 		if player.GetId() == resp.PlayerId {
// 			continue
// 		}
// 		err := player.WriteJson(resp)
// 		if err != nil {
// 			logging.Error(
// 				"couldn't notify player: ",
// 				zap.String("player_id", player.GetId()),
// 			)
// 		}
// 	}
// }
//
// func (m *Match) DisconnectPlayers(msg string, deadline time.Time) {
// 	for _, player := range m.Players {
// 		player.WriteControl(
// 			websocket.CloseMessage,
// 			websocket.FormatCloseMessage(
// 				websocket.CloseNormalClosure,
// 				msg,
// 			),
// 			deadline,
// 		)
// 	}
// }
