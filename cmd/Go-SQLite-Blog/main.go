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

	adapter := middleware.Init(cfg.ErrorsInResponse)

	slog.Info("Initializing database", "driver", cfg.Database.Driver, "source", cfg.Database.Source)
	db, err := dbService.Init(cfg.Database.Driver, cfg.Database.Source)
	if err != nil {
		slog.Error("main: Error initializing database", "error", err)
	}

	registerService, err := handlers.Init(db, logger)
	if err != nil {
		slog.Error("main: Error initializing register service", "error", err)
	}

	mux := http.NewServeMux()
	showRegister := adapter.HttpToContextHandler(registerService.ShowRegister)
	mux.Handle("GET /register", middleware.LoggingMiddleware(showRegister))

	tryRegister := adapter.HttpToContextHandler(registerService.TryRegister)
	mux.Handle("POST /register", middleware.LoggingMiddleware(tryRegister))

	// register routes
	http.HandleFunc("/", handlers.ArticlesHandler)
	http.HandleFunc("GET /login", handlers.ShowLogin)
	http.HandleFunc("POST /login", handlers.TryLogin)


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
