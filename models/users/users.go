package users

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"github.com/AndreHeber/go-sqlite-blog/models"
)

type User struct {
	ID             uint64
	Username       string
	HashedPassword string
	Salt           string
	Email          string
	Verified       bool
	RoleID         uint64
	CreatedAt      time.Time
	LastLogin      time.Time
}

//go:embed insert.sql
var insert string

func CreateUser(env *models.Env, user User) error {
	_, err := env.Db.ExecContext(env.Ctx, insert, user.Username, user.HashedPassword, user.Salt, user.Email, user.Verified, user.RoleID, user.CreatedAt, user.LastLogin)
	if err != nil {
		env.Logger.Error("models: CreateUser", "error", err, "sql", insert, "user", user)
		return fmt.Errorf("CreateUser: %w", err)
	}

	if env.LogDbQueries {
		env.Logger.Info("models: CreateUser", "sql", insert, "user", user)
	}

	return nil
}

//go:embed select_where_username.sql
var selectWhereUsername string

func GetUserByUsername(env *models.Env, username string) (User, error) {
	var user User
	err := env.Db.QueryRowContext(env.Ctx, selectWhereUsername, username).Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Salt, &user.Email, &user.Verified, &user.RoleID, &user.CreatedAt, &user.LastLogin)
	if err == sql.ErrNoRows {
		return User{}, fmt.Errorf("GetUserByUsername: user not found")
	}
	if err != nil {
		env.Logger.Error("models: GetUserByUsername", "error", err, "sql", selectWhereUsername, "username", username)
		return User{}, fmt.Errorf("GetUserByUsername: %w", err)
	}
	if env.LogDbQueries {
		env.Logger.Info("models: GetUserByUsername", "sql", selectWhereUsername, "username", username)
	}

	return user, nil
}
