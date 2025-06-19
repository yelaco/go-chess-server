package main

import (
	"log/slog"

	"github.com/yelaco/gchess-server/internal/api"
	"github.com/yelaco/gchess-server/internal/database"
	"github.com/yelaco/gchess-server/pkg/agent"
	"github.com/yelaco/gchess-server/pkg/logging"
	"github.com/yelaco/gchess-server/pkg/util"
	"go.uber.org/zap"
)

func main() {
	agent := agent.NewAgent()
	database.InitDB()
	defer database.CloseDB()

	cfg, err := util.LoadConfig("./configs")
	if err != nil {
		slog.Error("failed to load config:", err)
	}

	go func() {
		if err := agent.StartGameServer(cfg.WsPort); err != nil {
			logging.Fatal("game server failed to start", zap.Error(err))
		}
	}()

	go func() {
		if err := api.StartRESTServer(cfg.HttpPort); err != nil {
			logging.Fatal("rest server failed to start", zap.Error(err))
		}
	}()

	select {}
}
