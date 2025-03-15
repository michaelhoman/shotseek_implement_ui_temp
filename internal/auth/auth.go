package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/michaelhoman/ShotSeek/internal/config"
	"github.com/michaelhoman/ShotSeek/internal/store"
	"github.com/michaelhoman/ShotSeek/internal/utils"
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

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type Claims struct {
	Fingerprint          string `json:"fp"` // Fingerprint (optional)
	jwt.RegisteredClaims        // Contains standard claims like exp, iss, aud, iat, etc.
}

type AuthHandler struct {
	store      store.Storage
	config     config.Config
	jwtService *JWTService
}

func NewAuthHandler(store store.Storage, config config.Config, jwtService *JWTService) *AuthHandler {
	return &AuthHandler{
		store:      store,
		config:     config,
		jwtService: jwtService,
	}
}

// type userWithToken struct {
// 	*store.User
// 	Token string `json:"token"`
// }

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
func (a *AuthHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
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
		utils.InternalServerError(w, r, err)
		return
	}
	fmt.Println("user.Password", user.Password)
	// store the user

	// if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
	// 	utils.InternalServerError(w, r, err)
	// }

	ctx := r.Context()
	plainToken := uuid.New().String()
	// fmt.Println(user.FirstName, user.LastName, " plainToken:", plainToken)
	// store plainToken in the database as hashed token

	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	err := a.store.Users.CreateAndInvite(ctx, user, hashedToken, a.config.Mail.Exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			utils.BadRequestResponse(w, r, err)
		// case store.ErrDuplicateUser:
		// 	utils.BadRequestResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	// TODO revert the userWithToken after Dev/TESTING - this is insecure and bypasses the email verification
	// userWithToken := &userWithToken{
	// 	User:  user,
	// 	Token: plainToken,
	// }
	// if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
	// 	utils.InternalServerError(w, r, err)
	// }
	// Example: Sending success message response
	// response := map[string]string{
	// 	"message": "Registration successful! Check your email to verify your account.",
	// }
	confirmationMessage := "Registration successful! Check your email to verify your account."

	fmt.Println(user.FirstName, user.LastName, " plainToken:", plainToken)
	// Return a success message instead of the form
	if err := utils.WriteMessagePlain(w, http.StatusCreated, confirmationMessage); err != nil {
		utils.InternalServerError(w, r, err)
	}

	// send email to user with plainToken

}

// var jwtSigningKey = []byte(os.Getenv("JWT_SIGNING_KEY")) // Replace with a secure key
// Change this to your actual issuer

// Login handler that sets JWT in an HTTP-Only cookie

// loginHandler godoc
//
//	@Summary		Login a user
//	@Description	Login a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		LoginPayload	true	"User credentials"
//	@Success		200		{string}	string			"Login successful, JWT stored in cookie"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/login [post]
func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginHandler inital call") // Debugging
	var payload LoginPayload

	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	// Authenticate the user (e.g., check the password against the db)

	user, err := a.store.Users.GetByEmail(r.Context(), payload.Email)

	fmt.Println("user:", user) // Debugging

	if err != nil {
		switch err {
		case store.ErrNotFound:
			utils.UnauthorizedErrorResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	fmt.Println(payload.Password) // Debugging
	if err := user.Password.Compare(payload.Password); err != nil {
		utils.UnauthorizedErrorResponse(w, r, err)
		return
	}

	ip := getIPAddress(r)
	userAgent := r.UserAgent()
	fingerprint := generateFingerprint(ip, userAgent) // Optional fingerprint

	token, err := a.generateJWT(payload.Email, fingerprint)
	if err != nil {
		fmt.Println("Error generating JWT:", err) // TODO Remove Debugging
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Set the JWT in an HTTP-only, Secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,                                    // Ensures it's inaccessible via JavaScript
		Secure:   true,                                    // Only sent over HTTPS
		SameSite: http.SameSiteStrictMode,                 // Prevent CSRF
		Expires:  time.Now().Add(a.config.Auth.Token.Exp), // Cookie expiration time
	})

	// Respond to the user (no need to send the token in the body)
	w.Write([]byte("Login successful, JWT stored in cookie"))
}

// ExtractJWTToken extracts the JWT token from the Authorization header
func ExtractJWTToken(r *http.Request) (string, error) {
	// Extract JWT from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	// Token format validation (Bearer <token>)
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	return tokenParts[1], nil
}
