package main

import (
	"net/http"

	store "github.com/michaelhoman/ShotSeek/internal/store/postgres"
)

type CreateUserPayload struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Zipcode   string `json:"zip_code" validate:"required"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreateUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := store.User{
		Email:     payload.Email,
		Password:  payload.Password,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Zipcode:   payload.Zipcode,
		City:      payload.City,
		State:     payload.State,
	}

	ctx := r.Context()

	usersStore := app.store.Users
	if err := usersStore.Create(ctx, &user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
