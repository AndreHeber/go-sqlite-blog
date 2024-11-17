package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() (*sql.DB, error) {
	return sql.Open("sqlite3", "./blog.db")
}

func CloseDB(db *sql.DB) error {
	return db.Close()
}

func CreateTables(db *sql.DB) error {
	// read tables.sql
	content, err := os.ReadFile("tables.sql")
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

// check if db exists, if not create it
func CheckDBExists(db *sql.DB) (bool, error) {
	_, err := os.Stat("./blog.db")
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}
