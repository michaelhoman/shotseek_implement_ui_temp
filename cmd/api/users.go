package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	store "github.com/michaelhoman/ShotSeek/internal/store/postgres"
)

type userKey string

const userCtx postKey = "user"

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

	// w.WriteHeader(http.StatusCreated)

	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type UpdateUserPayload struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	Password  *string `json:"password" validate:"omitempty,min=6"`
	FirstName *string `json:"first_name" validate:"omitempty"`
	LastName  *string `json:"last_name" validate:"omitempty"`
	Zipcode   *string `json:"zip_code" validate:"omitempty"`
	City      *string `json:"city" validate:"omitempty"`
	State     *string `json:"state" validate:"omitempty"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	// ctx := r.Context()

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	var payload UpdateUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Email != nil {
		user.Email = *payload.Email
	}
	if payload.Password != nil {
		user.Password = *payload.Password
	}

	if payload.FirstName != nil {
		user.FirstName = *payload.FirstName
	}

	if payload.LastName != nil {
		user.LastName = *payload.LastName
	}
	if payload.Zipcode != nil {
		user.Zipcode = *payload.Zipcode
	}
	if payload.City != nil {
		user.City = *payload.City
	}
	if payload.State != nil {
		user.State = *payload.State
	}

	ctx := r.Context()

	usersStore := app.store.Users
	if err := usersStore.Update(ctx, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(idParam, 10, 64)

		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	post, _ := r.Context().Value(userCtx).(*store.User)
	return post
}
