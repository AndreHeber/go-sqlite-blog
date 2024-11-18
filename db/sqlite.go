package db

import (
	"database/sql"

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

// func CreateTables(db *sql.DB) error {
// 	// read tables.sql
// 	content, err := os.ReadFile("./db/tables.sql")
// 	if err != nil {
// 		return err
// 	}

// 	// execute the content
// 	_, err = db.Exec(string(content))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func Init(driver string, source string) (*sql.DB, error) {
	return ConnectDB(driver, source)
}
