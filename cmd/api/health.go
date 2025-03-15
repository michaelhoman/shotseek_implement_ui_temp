package main

import (
	"net/http"

	//postgres_store "github.com/michaelhoman/ShotSeek/internal/store/postgres"
	"github.com/michaelhoman/ShotSeek/internal/store"
	"github.com/michaelhoman/ShotSeek/internal/utils"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":  "available",
		"env":     app.config.Env,
		"version": version,
	}

	if err := utils.JsonResponse(w, http.StatusOK, data); err != nil {
		utils.InternalServerError(w, r, err)
	}

	app.store.Posts.Create(r.Context(), &store.Post{})

}
