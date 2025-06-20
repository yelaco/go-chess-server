package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yelaco/gchess-server/pkg/logging"
	"github.com/yelaco/gchess-server/pkg/util"
	"github.com/yelaco/ludofy/internal/domains/dtos"
	"go.uber.org/zap"
)

type Server struct {
	http.Server
	upgrader websocket.Upgrader

	cfg          util.Config
	matches      sync.Map
	totalMatches atomic.Int32
	mu           *sync.Mutex
}

func NewGameServeMux(cfg util.Config) *http.ServeMux {
	srv := &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		mu:  new(sync.Mutex),
		cfg: cfg,
	}

	return srv.setupRoutes()
}

func (s *Server) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		count := s.totalMatches.Load()
		json.NewEncoder(w).Encode(map[string]any{
			"activeMatches": count,
		})
	})

	// Websocket
	http.HandleFunc("/{matchId}", func(w http.ResponseWriter, r *http.Request) {
		playerId := ""
		// if err != nil {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	w.Write([]byte(err.Error()))
		// 	logging.Error("failed to auth: %w", zap.Error(err))
		// 	return
		// }

		conn, err := h.upgrader.Upgrade(w, r, nil)
		if err != nil {
			logging.Error(
				"failed to upgrade connection",
				zap.String("error", err.Error()),
			)
			return
		}
		defer conn.Close()

		matchId := r.PathValue("matchId")
		match, err := h.loadMatch(matchId)
		if err != nil {
			logging.Info("failed to load match", zap.String("error", err.Error()))
			conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(
					websocket.CloseNormalClosure,
					"match failed to load",
				),
				time.Now().Add(5*time.Second),
			)
			return
		}
		match.playerJoin(playerId, conn)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(
					err,
					websocket.CloseNormalClosure,
				) {
					logging.Info(
						"connection closed gracefully",
						zap.String("remote_address", conn.RemoteAddr().String()),
					)
				} else if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseAbnormalClosure,
				) {
					logging.Info(
						"unexpected connection close",
						zap.String("remote_address", conn.RemoteAddr().String()),
						zap.Error(err),
					)
				}
				match.playerDisconnect(playerId)
				break
			}

			err = s.HandleMessage(playerId, match, message)
			if err != nil {
				logging.Error("failed to handle message", zap.Error(err))
				conn.WriteControl(
					websocket.CloseNormalClosure,
					nil,
					time.Now().Add(5*time.Second),
				)
			}
		}
	})

	return mux
}

func (s *Server) HandleMessage(playerId string, match *Match, msg []byte) error {
	if match == nil {
		return fmt.Errorf("match not loaded")
	}
	// err := s.handler.OnHandleMessage(playerId, match.GetHandler(), msg)
	// if err != nil {
	// 	return fmt.Errorf("on handle message: %w", err)
	// }
	return nil
}

func (s *Server) HandleMatchEnd(match *Match) {
	if match == nil {
		return
	}
	// matchRecordReq := MatchRecordRequest{
	// 	MatchId: match.GetId(),
	// 	EndedAt: time.Now(),
	// }

	// if err := s.handler.OnHandleMatchEnd(&matchRecordReq, match.GetHandler()); err != nil {
	// 	logging.Fatal("failed to hanlde match end", zap.Error(err))
	// }

	// payload, err := json.Marshal(matchRecordReq)
	// if err != nil {
	// 	logging.Fatal("failed to marshal match record request", zap.Error(err))
	// }

	s.removeMatch(match.Id)
	logging.Info("match ended", zap.String("match_id", match.Id))
}

func (s *Server) HandleMatchSave(match *Match) {
	// ctx := context.Background()
	// matchStateReq := dtos.MatchStateRequest{
	// 	Id:        utils.GenerateUUID(),
	// 	MatchId:   match.Id,
	// 	Timestamp: time.Now(),
	// }
	// s.handler.OnHandleMatchSave(&matchStateReq, match.GetHandler())

	// Save match
}

func (s *Server) HandleMatchAbort(match *Match) {
	if match == nil {
		return
	}

	matchAbortReq := dtos.MatchAbortRequest{
		MatchId:   match.Id,
		PlayerIds: make([]string, 0, len(match.Players)),
	}

	payload, err := json.Marshal(matchAbortReq)
	if err != nil {
		log.Fatal(err)
	}

	// Abort match

	s.removeMatch(match.Id)
	logging.Info("match aborted", zap.String("match_id", match.Id))
}

func (s *Server) loadMatch(matchId string) (Match, error) {
	ctx := context.Background()

	activeMatch, err := storageClient.GetActiveMatch(ctx, matchId)
	if err != nil {
		return nil, fmt.Errorf("failed to get active match: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	value, loaded := s.matches.Load(matchId)
	if loaded {
		match, ok := value.(Match)
		if ok {
			logging.Info("match loaded")
			return match, nil
		}
		return nil, ErrFailedToLoadMatch
	} else {
		matchStates, _, err := storageClient.FetchMatchStates(
			ctx,
			matchId,
			nil,
			1,
			false,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch match states: %w", err)
		}

		var match Match
		if len(matchStates) > 0 {
			match, err = s.handler.OnMatchResume(activeMatch, matchStates[0])
			if err != nil {
				return nil, fmt.Errorf("failed to resume match: %w", err)
			}
		} else {
			match, err = s.handler.OnMatchCreate(activeMatch)
			if err != nil {
				return nil, fmt.Errorf("failed to create match: %w", err)
			}
		}

		// match.setSaveCallback(s.HandleMatchSave)
		// match.setEndCallback(s.HandleMatchEnd)
		// match.setAbortCallback(s.HandleMatchAbort)
		s.matches.Store(matchId, match)
		s.totalMatches.Add(1)

		go match.start()
		logging.Info("match loaded", zap.String("match_id", matchId))
		return match, nil
	}
}

func (s *Server) removeMatch(matchId string) {
	s.matches.Delete(matchId)
	total := s.totalMatches.Add(-1)
	logging.Info("match removed", zap.Int32("total_matches", total))
}
