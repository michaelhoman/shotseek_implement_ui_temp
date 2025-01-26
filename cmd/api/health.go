package main

import (
	"net/http"

	"github.com/michaelhoman/ShotSeek/internal/store/postgres"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))

	app.store.Posts.Create(r.Context(), &postgres.Post{})

}
