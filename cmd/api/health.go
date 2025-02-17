package main

import (
	"net/http"

	//postgres_store "github.com/michaelhoman/ShotSeek/internal/store/postgres"
	"github.com/michaelhoman/ShotSeek/internal/store"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":  "available",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}

	app.store.Posts.Create(r.Context(), &store.Post{})

}
