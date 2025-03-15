package ui

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Define UI route handlers as standalone functions
func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/auth.html"))
	tmpl.Execute(w, nil)
}

func VerifyPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/verify.html"))
	tmpl.Execute(w, nil)
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

// Register UI routes on an existing router
func RegisterUIRoutes(r *chi.Mux) {
	r.Get("/register", RegisterPageHandler)
	r.Get("/verify", VerifyPageHandler)
	r.Get("/login", LoginPageHandler)
}
