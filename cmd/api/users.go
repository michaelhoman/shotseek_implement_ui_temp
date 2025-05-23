package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	// store "github.com/michaelhoman/ShotSeek/internal/store/postgres"

	"github.com/michaelhoman/ShotSeek/internal/store"
	"github.com/michaelhoman/ShotSeek/internal/utils"
)

type userKey string
type locationKey string

const userCtx userKey = "user"
const locationCtx locationKey = "location"

type CreateUserPayload struct {
	Email string `json:"email" validate:"required,email"`
	// Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Zipcode   string `json:"zip_code" validate:"required"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
}

// CreateUser godoc
//
//	@Summary		Creates a user
//	@Description	Creates a new user from payload
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserPayload	true	"User payload"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users [post]
// func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

// 	var payload CreateUserPayload
// 	if err := readJSON(w, r, &payload); err != nil {
// 		utils.BadRequestResponse(w, r, err)
// 		return
// 	}

// 	if err := Validate.Struct(payload); err != nil {
// 		utils.BadRequestResponse(w, r, err)
// 		return
// 	}

// 	user := store.User{
// 		Email: payload.Email,
// 		// Password:  payload.Password,
// 		FirstName: payload.FirstName,
// 		LastName:  payload.LastName,
// 		Zipcode:   payload.Zipcode,
// 		City:      payload.City,
// 		State:     payload.State,
// 	}

// 	ctx := r.Context()

// 	usersStore := app.store.Users
// 	if err := usersStore.Create(ctx, &user); err != nil {
// 		utils.InternalServerError(w, r, err)
// 		return
// 	}

// 	// w.WriteHeader(http.StatusCreated)

// 	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
// 		utils.InternalServerError(w, r, err)
// 		return
// 	}
// }

type UpdateUserPayload struct {
	Email *string `json:"email" validate:"omitempty,email"`
	// Password  *string `json:"password" validate:"omitempty,min=6"`
	FirstName *string `json:"first_name" validate:"omitempty"`
	LastName  *string `json:"last_name" validate:"omitempty"`
	Zipcode   *string `json:"zip_code" validate:"omitempty"`
	City      *string `json:"city" validate:"omitempty"`
	State     *string `json:"state" validate:"omitempty"`
}

// GetUser godoc
//
//	@Summary		Fetches a user by ID string/uuid
//	@Description	Fetches a user by ID string/uuid
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := utils.JsonResponse(w, http.StatusOK, user); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}

// GetCurrentUser godoc
//
//	@Summary		Fetches the current user
//	@Description	Fetches the current user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/ [get]
func (app *application) getCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if err := utils.JsonResponse(w, http.StatusOK, user); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}

// UpdateUser godoc
//
//	@Summary		Updates a user
//	@Description	Updates a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"User ID"
//	@Param			payload	body		UpdateUserPayload	true	"User payload"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [patch]
func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	location := getLocationFromCtx(r)
	var payload UpdateUserPayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if payload.Email != nil {
		user.Email = *payload.Email
	}
	// if payload.Password != nil {
	// user.Password = *payload.Password
	// }

	if payload.FirstName != nil {
		user.FirstName = *payload.FirstName
	}

	if payload.LastName != nil {
		user.LastName = *payload.LastName
	}
	if payload.Zipcode != nil {
		location.ZIPCode = *payload.Zipcode
	}
	if payload.City != nil {
		location.City = *payload.City
	}
	if payload.State != nil {
		location.State = *payload.State
	}

	ctx := r.Context()

	usersStore := app.store.Users
	if err := usersStore.Update(ctx, user, location); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser godoc
//
//	@Summary		Deletes a user
//	@Description	Deletes a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		204	{object}	nil
//	@Failure		400	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [delete]
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	ctx := r.Context()

	usersStore := app.store.Users
	if err := usersStore.Delete(ctx, user.ID); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ActivateUser godoc
//
//	@Summary		Activates a user
//	@Description	Activates a user by invitation token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/authentication/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token") // Get token from path

	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			utils.BadRequestResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	if err := utils.JsonResponse(w, http.StatusNoContent, "User activated"); err != nil {
		utils.InternalServerError(w, r, err)
	}
}

func (app *application) getUserLocationHandler(w http.ResponseWriter, r *http.Request) {
	location := getLocationFromCtx(r)

	if err := utils.JsonResponse(w, http.StatusOK, location); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		id, err := convertToUUID(idParam)
		if err != nil {
			fmt.Println("Error parsing UUID:", err)
		} else {
			fmt.Println("Parsed UUID:", id)
		}

		if err != nil {
			utils.InternalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				utils.NotFoundResponse(w, r, err)
			default:
				utils.InternalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) currentUserContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Extract token from Authorization header
		tokenStr, err := utils.ExtractBearerToken(r)
		if err != nil {
			utils.BadRequestResponse(w, r, fmt.Errorf("missing or malformed Authorization header"))
			return
		}

		// Step 2: Decode and validate token
		claims, err := utils.DecodeJWT(tokenStr, app.jwtAuth.PublicKey)
		if err != nil {
			utils.BadRequestResponse(w, r, fmt.Errorf("invalid token: %w", err))
			return
		}

		// Step 3: Extract user ID from `sub` claim
		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			utils.BadRequestResponse(w, r, errors.New("missing or invalid subject claim"))
			return
		}

		userID, err := uuid.Parse(sub)
		if err != nil {
			utils.BadRequestResponse(w, r, errors.New("invalid user ID format in token"))
			return
		}

		// Step 4: Fetch user from DB
		user, err := app.store.Users.GetByID(r.Context(), userID)
		if err != nil {
			utils.InternalServerError(w, r, fmt.Errorf("failed to load user: %w", err))
			return
		}

		// Step 5: Add user to context and continue
		ctx := context.WithValue(r.Context(), userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

func convertToUUID(idParam string) (uuid.UUID, error) {
	id, err := uuid.Parse(idParam)
	if err != nil {
		return uuid.Nil, err // Return nil UUID on error
	}
	return id, nil
}

// In your handler file, e.g., location.go or debug.go
