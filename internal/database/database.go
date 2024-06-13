package database

import (
	"fmt"

	"github.com/minhquang053/robinhood-chess/internal/config"
	"github.com/minhquang053/robinhood-chess/internal/logging"

	"github.com/jackc/pgx"
)

func Connect() *pgx.ConnPool {
	cfg := config.Config()
	conn, err := pgx.NewConnPool(cfg)
	if err != nil {
		logging.Fatal(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}
	return conn
}