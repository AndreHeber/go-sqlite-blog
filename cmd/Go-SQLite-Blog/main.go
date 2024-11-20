package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/AndreHeber/go-sqlite-blog/config"
	dbService "github.com/AndreHeber/go-sqlite-blog/db"
	"github.com/AndreHeber/go-sqlite-blog/handlers"
	"github.com/AndreHeber/go-sqlite-blog/middleware"
)

func main() {
	// reading env variables
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("main: Error loading config", "error", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	slog.SetDefault(logger)

	if cfg.Database.Reset {
		slog.Info("Resetting database")
		// delete database file
		os.Remove(cfg.Database.Source)
	}

	slog.Info("Initializing database", "driver", cfg.Database.Driver, "source", cfg.Database.Source)
	db, err := dbService.Init(cfg.Database.Driver, cfg.Database.Source, "./db/tables.sql")
	if err != nil {
		slog.Error("main: Error initializing database", "error", err)
		os.Exit(1)
	}

	adapter := middleware.Init(logger, db, cfg.ErrorsInResponse, cfg.Database.LogQueries, cfg.IPRateLimit, cfg.BurstRateLimit)
	mux := setupRouter(adapter)

	// start server
	slog.Info("Starting server", "port", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), mux)
	if err != nil {
		slog.Error("main: Error starting server", "error", err)
	}
}

func setupRouter(adapter *middleware.Adapter) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /health", adapter.HTTPToContextHandler(handlers.Health))
	mux.Handle("GET /time-consuming", adapter.HTTPToContextHandler(handlers.TimeConsumingHandler))

	mux.Handle("GET /register", adapter.HTTPToContextHandler(handlers.ShowRegister))
	mux.Handle("POST /register", adapter.HTTPToContextHandler(handlers.TryRegister))

	mux.Handle("GET /login", adapter.HTTPToContextHandler(handlers.ShowLogin))
	mux.Handle("POST /login", adapter.HTTPToContextHandler(handlers.TryLogin))

	return mux
}
