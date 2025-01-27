package main

import (
	"net/http"

	"github.com/michaelhoman/ShotSeek/internal/store/postgres"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":  "available",
		"env":     app.config.env,
		"version": version,
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}

	app.store.Posts.Create(r.Context(), &postgres.Post{})

}
