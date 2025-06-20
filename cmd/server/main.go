package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yelaco/gchess-server/internal/app/server"
	"github.com/yelaco/gchess-server/pkg/util"
)

func main() {
	cfg, err := util.LoadConfig("./configs")
	if err != nil {
		slog.Error("failed to load config:", err)
	}

	// Setup routes
	mux := http.NewServeMux()
	mux.Handle("/game", server.NewGameServeMux(cfg))

	srv := &http.Server{
		Addr:    "0.0.0.0" + cfg.Port,
		Handler: mux,
	}

	// Run server with graceful shutdown
	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("server is running on port: ", cfg.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		slog.Error("could not start server: %w", err)

	case sig := <-stop:
		slog.Info("received signal: %v. shutting down gracefully...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("server shutdown failed:%+v", err)
		}

		slog.Info("server exited")
	}
}
