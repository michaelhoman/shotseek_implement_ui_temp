package main

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) registerPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/auth.html"))
	tmpl.Execute(w, nil)
}

func (app *application) verifyPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/verify.html"))
	tmpl.Execute(w, nil)
}

// Define UI routes in a separate function
func (app *application) uiRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/register", app.registerPageHandler)
	r.Get("/verify", app.verifyPageHandler)
	return r
}
