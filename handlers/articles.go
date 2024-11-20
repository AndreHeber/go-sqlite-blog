package handlers

import (
	"html/template"
	"net/http"
)

func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/templates/articles.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
