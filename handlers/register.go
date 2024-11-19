package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/AndreHeber/go-sqlite-blog/middleware"
	"github.com/AndreHeber/go-sqlite-blog/models"
	"github.com/AndreHeber/go-sqlite-blog/models/users"
)

// ShowRegister renders the register page
func ShowRegister(a *middleware.Adapter) error {
	w := a.ResponseWriter
	tmpl, err := template.ParseFiles("static/templates/register.html")
	if err != nil {
		return fmt.Errorf("ShowRegister: %w", err)
	}
	tmpl.Execute(w, nil)
	return nil
}

// TryRegister handles the registration form submission
func TryRegister(a *middleware.Adapter) error {
	w := a.ResponseWriter
	r := a.Request
	logger := a.Logger

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	logger.Info("Trying to register", "username", username, "password", password, "email", email)

	err := register(models.EnvFromAdapter(a), username, password, email)
	if err != nil {
		return fmt.Errorf("TryRegister: %w", err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

// register saves the user to the database, password is hashed before saving using bcrypt
// also an email is sent to the user with a link to verify their account
func register(env *models.Env, username, password, email string) error {
	// verify input
	if username == "" || password == "" || email == "" {
		return errors.New("register: username, password and email are required")
	}

	encodedHash, encodedSalt := hashPassword(password, 16)

	// save user to database
	user := users.User{
		Username:       username,
		HashedPassword: encodedHash,
		Salt:           encodedSalt,
		Email:          email,
		Verified:       false,
		RoleId:         1,
		CreatedAt:      time.Now(),
		LastLogin:      time.Now(),
	}

	err := users.CreateUser(env, user)
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}

	return nil
}