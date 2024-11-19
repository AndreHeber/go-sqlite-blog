package db

import (
	"database/sql"
	"os"

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

func CreateTables(db *sql.DB, initFile string) error {
	// read tables.sql
	content, err := os.ReadFile(initFile)
	if err != nil {
		return err
	}

	// execute the content
	_, err = db.Exec(string(content))
	if err != nil {
		return err
	}

	return nil
}

func Init(driver string, source string, initFile string) (*sql.DB, error) {
	db, err := ConnectDB(driver, source)
	if err != nil {
		return nil, err
	}

	err = CreateTables(db, initFile)
	if err != nil {
		return nil, err
	}

	return db, nil
}
