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
	db, err := dbService.Init(cfg.Database.Driver, cfg.Database.Source)
	if err != nil {
		slog.Error("main: Error initializing database", "error", err)
	}

	adapter := middleware.Init(logger, db, cfg.ErrorsInResponse, cfg.Database.LogQueries, cfg.IpRateLimit, cfg.BurstRateLimit)

	mux := http.NewServeMux()

	mux.Handle("GET /health", adapter.HttpToContextHandler(handlers.Health))
	mux.Handle("GET /time-consuming", adapter.HttpToContextHandler(handlers.TimeConsumingHandler))

	mux.Handle("GET /register", adapter.HttpToContextHandler(handlers.ShowRegister))
	mux.Handle("POST /register", adapter.HttpToContextHandler(handlers.TryRegister))

	mux.Handle("GET /login", adapter.HttpToContextHandler(handlers.ShowLogin))
	mux.Handle("POST /login", adapter.HttpToContextHandler(handlers.TryLogin))

	// register routes
	http.HandleFunc("/", handlers.ArticlesHandler)


	http.HandleFunc("/settings", handlers.SettingsHandler)
	http.HandleFunc("/articles", handlers.ArticlesHandler)
	// http.HandleFunc("/categories", handlers.CategoriesHandler)
	// http.HandleFunc("/tags", handlers.TagsHandler)
	// http.HandleFunc("/media", handlers.MediaHandler)
	// http.HandleFunc("/pages", handlers.PagesHandler)

	// start server
	slog.Info("Starting server", "port", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), mux)
	if err != nil {
		slog.Error("main: Error starting server", "error", err)
	}
}
