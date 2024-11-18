package users

import (
	"database/sql"
	_ "embed"
	"time"
)

type User struct {
	Id             uint64
	Username       string
	HashedPassword string
	Email          string
	Verified       bool
	RoleId         uint64
	CreatedAt      time.Time
	LastLogin      time.Time
}

//go:embed schema.sql
var table string

func CreateSchema(db *sql.DB) error {
	_, err := db.Exec(table)
	return err
}

//go:embed insert.sql
var insert string

func CreateUser(db *sql.DB, user User) error {
	_, err := db.Exec(insert, user.Username, user.HashedPassword, user.Email, user.Verified, user.RoleId, user.CreatedAt, user.LastLogin)
	return err
}
