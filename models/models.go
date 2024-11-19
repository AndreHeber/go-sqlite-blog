package models

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/AndreHeber/go-sqlite-blog/middleware"
)

type Env struct {
	Db *sql.DB
	Ctx context.Context
	Logger *slog.Logger
	LogDbQueries bool
}

func EnvFromAdapter(adapter *middleware.Adapter) *Env {
	return &Env{Db: adapter.Db, Ctx: adapter.Ctx, Logger: adapter.Logger, LogDbQueries: adapter.LogDbQueries}
}