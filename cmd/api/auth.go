package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/michaelhoman/ShotSeek/internal/store"
)

type RegisterUserPayload struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=72"`
	FirstName string `json:"first_name" validate:"required,max=255"`
	LastName  string `json:"last_name" validate:"required,max=255"`
	Zipcode   string `json:"zip_code" validate:"required,max=12"`
	City      string `json:"city" validate:"required,max=255"`
	State     string `json:"state" validate:"required,max=255"`
}

type userWithToken struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a new user
//	@Description	Registers a new user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	store.User			"User Registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// hash the user password
	user := &store.User{
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Zipcode:   payload.Zipcode,
		City:      payload.City,
		State:     payload.State,
	}

	// has password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// store the user

	// if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
	// 	app.internalServerError(w, r, err)
	// }

	ctx := r.Context()
	plainToken := uuid.New().String()
	// fmt.Println(user.FirstName, user.LastName, " plainToken:", plainToken)
	// store plainToken in the database as hashed token

	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	err := app.store.Users.CreateAndInvite(ctx, user, hashedToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		// case store.ErrDuplicateUser:
		// 	app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// TODO revert the userWithToken after Dev/TESTING - this is insecure and bypasses the email verification
	// userWithToken := &userWithToken{
	// 	User:  user,
	// 	Token: plainToken,
	// }
	// if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
	// 	app.internalServerError(w, r, err)
	// }
	// Example: Sending success message response
	// response := map[string]string{
	// 	"message": "Registration successful! Check your email to verify your account.",
	// }
	confirmationMessage := "Registration successful! Check your email to verify your account."

	fmt.Println(user.FirstName, user.LastName, " plainToken:", plainToken)
	// Return a success message instead of the form
	if err := app.writeMessagePlain(w, http.StatusCreated, confirmationMessage); err != nil {
		app.internalServerError(w, r, err)
	}

	// send email to user with plainToken

}
