package main

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/AndreHeber/go-sqlite-blog/config"
	dbService "github.com/AndreHeber/go-sqlite-blog/db"
	"github.com/AndreHeber/go-sqlite-blog/middleware"
)

func TestAPI(t *testing.T) {

	// reading env variables
	cfg := config.Config{
		LogLevel: &slog.LevelVar{},
		Database: config.DatabaseConfig{
			Driver: "sqlite3",
			Source: "./test.db",
		},
		ErrorsInResponse: false,
		IPRateLimit:      10,
		BurstRateLimit:   10,
	}
	cfg.LogLevel.Set(slog.LevelDebug)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	slog.SetDefault(logger)

	_ = os.Remove(cfg.Database.Source)

	db, err := dbService.Init(cfg.Database.Driver, cfg.Database.Source, "../../db/tables.sql")
	if err != nil {
		slog.Error("main: Error initializing database", "error", err)
		t.Fatalf("Error initializing database: %v", err)
	}

	adapter := middleware.Init(logger, db, cfg.ErrorsInResponse, cfg.Database.LogQueries, cfg.IPRateLimit, cfg.BurstRateLimit)
	mux := setupRouter(adapter)

	// Create test server
	server := httptest.NewServer(mux)
	defer server.Close()

	defer os.Remove(cfg.Database.Source)

	t.Run("Test /register endpoint", func(t *testing.T) {
		// make a POST request to /register with Form data
		requestBody := bytes.NewBufferString(`username=testuser&password=testpassword&email=test@test.com`)
		request, err := http.NewRequest("POST", server.URL+"/register", requestBody)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// make the request
		response, err := server.Client().Do(request)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}

		// check the response status code
		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
		}
	})
}
