package handlers

import (
	"context"
	"database/sql"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/AndreHeber/go-sqlite-blog/middleware"
	"github.com/AndreHeber/go-sqlite-blog/models/users"
	"golang.org/x/crypto/bcrypt"
)

type RegisterService struct {
	db *sql.DB
	logger *slog.Logger
}

func Init(db *sql.DB, logger *slog.Logger) (*RegisterService, error) {
	err := users.CreateSchema(db)
	if err != nil {
		return nil, err
	}
	return &RegisterService{db: db, logger: logger}, nil
}

// register saves the user to the database, password is hashed before saving using bcrypt
// also an email is sent to the user with a link to verify their account
func register(db *sql.DB, username, password, email string) error {
	// verify input
	if username == "" || password == "" || email == "" {
		return errors.New("register: username, password and email are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// save user to database
	user := users.User{
		Username:       username,
		HashedPassword: string(hashedPassword),
		Email:          email,
		Verified:       false,
		RoleId:         1,
		CreatedAt:      time.Now(),
		LastLogin:      time.Now(),
	}

	err = users.CreateUser(db, user)
	if err != nil {
		return err
	}

	return nil
}

// ShowRegister renders the register page
func (rs *RegisterService) ShowRegister(ctx context.Context) error {
	w := middleware.GetResponseWriter(ctx)
	tmpl, err := template.ParseFiles("static/templates/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	tmpl.Execute(w, nil)
	return nil
}

// TryRegister handles the registration form submission
func (rs *RegisterService) TryRegister(ctx context.Context) error {
	w := middleware.GetResponseWriter(ctx)	
	r := middleware.GetRequest(ctx)

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	slog.Info("Trying to register", "username", username, "password", password, "email", email)

	err := register(rs.db, username, password, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
