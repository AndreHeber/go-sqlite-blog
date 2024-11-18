package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

func ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func TryLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	slog.Info("Trying to login", "username", username, "password", password)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
