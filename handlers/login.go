package handlers

import (
	"errors"
	"fmt"
	"html/template"

	"github.com/AndreHeber/go-sqlite-blog/middleware"
	"github.com/AndreHeber/go-sqlite-blog/models"
	"github.com/AndreHeber/go-sqlite-blog/models/users"
)

func ShowLogin(a *middleware.Adapter) error {
	w := a.ResponseWriter
	tmpl, err := template.ParseFiles("static/templates/login.html")
	if err != nil {
		return fmt.Errorf("ShowLogin: %w", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		return fmt.Errorf("ShowLogin: %w", err)
	}
	return nil
}

func TryLogin(a *middleware.Adapter) error {
	r := a.Request

	err := r.ParseForm()
	if err != nil {
		return fmt.Errorf("TryLogin: %w", err)
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	err = login(models.EnvFromAdapter(a), username, password)
	if err != nil {
		if a.ErrorInResponse {
			return fmt.Errorf("TryLogin: %w", err)
		}
		return fmt.Errorf("TryLogin: username or password is invalid")
	}

	return nil
}

func login(env *models.Env, username, password string) error {
	// verify input
	if username == "" || password == "" {
		return errors.New("login: username and password are required")
	}

	// get user from database
	user, err := users.GetUserByUsername(env, username)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	// verify password
	if !verifyPassword(password, user.HashedPassword, user.Salt) {
		return errors.New("login: invalid password")
	}

	return nil
}
