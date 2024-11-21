package models

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/AndreHeber/go-sqlite-blog/middleware"
)

type Env struct {
	DB           *sql.DB
	Ctx          context.Context
	Logger       *slog.Logger
	LogDBQueries bool
}

func EnvFromAdapter(adapter *middleware.Adapter) *Env {
	return &Env{DB: adapter.DB, Ctx: adapter.Ctx, Logger: adapter.Logger, LogDBQueries: adapter.LogDBQueries}
}
