package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func ConnectDB(driver string, source string) (*sql.DB, error) {
	return sql.Open(driver, source)
}

func CloseDB(db *sql.DB) error {
	return db.Close()
}

func CreateTables(logger *slog.Logger, db *sql.DB, initFile string) error {
	cleanPath := filepath.Clean(initFile)
	if strings.Contains(cleanPath, "..") {
		logger.Error("CreateTables: Invalid file path", "path", cleanPath)
		return fmt.Errorf("invalid file path")
	}

	// read tables.sql
	content, err := os.ReadFile(initFile)
	if err != nil {
		logger.Error("CreateTables: Error reading init file", "error", err)
		return err
	}

	// execute the content
	_, err = db.Exec(string(content))
	if err != nil {
		logger.Error("CreateTables: Error executing init file", "error", err)
		return err
	}

	return nil
}

func Init(logger *slog.Logger, driver string, source string, initFile string) (*sql.DB, error) {
	db, err := ConnectDB(driver, source)
	if err != nil {
		return nil, err
	}

	err = CreateTables(logger, db, initFile)
	if err != nil {
		return nil, err
	}

	return db, nil
}
