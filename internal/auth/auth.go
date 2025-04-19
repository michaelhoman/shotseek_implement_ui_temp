package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
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
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=8,max=72"`
	FirstName string  `json:"first_name" validate:"required,max=255"`
	LastName  string  `json:"last_name" validate:"required,max=255"`
	Street    string  `json:"street" validate:"required,max=255"`
	City      string  `json:"city" validate:"required,max=255"`
	State     string  `json:"state" validate:"required,max=255"`
	Zipcode   string  `json:"zip_code" validate:"required,max=12"`
	Country   string  `json:"country" validate:"required,max=255"`
	Latitude  float64 `json:"latitude" validate:"required,max=255"`
	Longitude float64 `json:"longitude" validate:"required,max=255"`
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
	Config     config.Config
	jwtService *JWTService
	JWTAuth    *JWTAuth
}

func NewAuthHandler(store store.Storage, config config.Config, jwtService *JWTService, jwtAuth *JWTAuth) *AuthHandler {
	return &AuthHandler{
		store:      store,
		Config:     config,
		jwtService: jwtService,
		JWTAuth:    jwtAuth,
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

	location := &store.Location{
		Street:    payload.Street,
		City:      payload.City,
		State:     payload.State,
		ZIPCode:   payload.Zipcode,
		Country:   payload.Country,
		Latitude:  payload.Latitude,
		Longitude: payload.Longitude,
	}

	// hash the user password
	user := &store.User{
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
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

	err := a.store.Users.CreateAndInvite(ctx, user, location, hashedToken, a.Config.Mail.Exp)
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
//	@Success		200		{string}	string			"updated Login successful, JWT stored in cookie"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/login [post]
func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginHandler LoginHandler inital call") // Debugging
	tokenStore := a.store.Tokens
	if tokenStore == nil {
		log.Println("Error: tokenStore is nil")
		utils.InternalServerError(w, r, errors.New("internal server error"))
		return
	}

	var payload LoginPayload

	fmt.Println("LoginHandler *1") // Debugging
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	fmt.Println("LoginHandler *2") // Debugging

	if err := utils.Validate.Struct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	fmt.Println("LoginHandler *3") // Debugging

	// Authenticate the user (e.g., check the password against the db)

	user, err := a.store.Users.GetByEmailWithPassword(r.Context(), payload.Email)

	fmt.Println("LoginHandler *4")                            // Debugging
	fmt.Println("LoginHandler payload.Email:", payload.Email) // Debugging

	fmt.Println("LoginHandler user:", user) // Debugging

	if err != nil {
		switch err {
		case store.ErrNotFound:
			utils.UnauthorizedErrorResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	fmt.Println("LoginHandler *5") // Debugging
	// Compare the hashed password
	fmt.Println(payload.Password) // Debugging
	if err := user.Password.Compare(payload.Password); err != nil {
		utils.UnauthorizedErrorResponse(w, r, err)
		return
	}

	ip := a.GetIPAddress(r)
	userAgent := r.UserAgent()
	fingerprint := a.GenerateFingerprint(ip, userAgent) // Optional fingerprint

	fmt.Println("LoginHandler *6") // Debugging

	newRefreshToken, err := a.GenerateRefreshToken()
	fmt.Println("LoginHandler *7 newRefreshToken:", newRefreshToken) // Debugging
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	newRefreshTokenHash := a.HashToken(newRefreshToken)
	fmt.Println("LoginHandler *8 newRefreshTokenHash:", newRefreshTokenHash) // Debugging

	fmt.Println("LoginHandler *9 a.store.Tokens:", a.store.Tokens) // Debugging
	if a.store.Tokens == nil {
		log.Println("Error: a.store.Tokens is nil")
		utils.InternalServerError(w, r, errors.New("internal server error"))
		return
	}

	// Store the refresh token in the databasePrintln("******") // Debugging

	fmt.Println("LoginHandler Calling UpdateRefreshToken") // Debugging
	err = tokenStore.UpdateRefreshToken(r.Context(), user.ID, newRefreshTokenHash, fingerprint, time.Now().Add(a.Config.Auth.RefreshToken.Exp))

	fmt.Println("LoginHandler *10") // Debugging
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
	fmt.Println("LoginHandler **********!!!!!!!!!!!!!!!!*********")  // Debugging
	fmt.Println("LoginHandler Generating JWT for user.ID:", user.ID) // Debugging
	token, err := a.GenerateJWTWithFP(user.ID, fingerprint)          // Pass fingerprint if needed
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
		Expires:  time.Now().Add(a.Config.Auth.Token.Exp), // Cookie expiration time
	})

	// Set the refresh token in a secure, HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		HttpOnly: true,                                           // Ensures it's inaccessible via JavaScript
		Secure:   true,                                           // Only sent over HTTPS
		SameSite: http.SameSiteStrictMode,                        // Adjust as necessary
		Expires:  time.Now().Add(a.Config.Auth.RefreshToken.Exp), // Set expiration based on config
	})

	// Respond to the user (no need to send the token in the body)
	w.Write([]byte("LoginHandler Login successful, JWT stored in cookie"))
}

// LogoutHandler godoc
//
//	@Summary		Logout a user
//	@Description	Logout a user
//	@Tags			users
//	@Produce		json
//	@Success		200	{string}	string	"Logout successful"
//	@Failure		500	{object}	error
//	@Router			/authentication/logout [post]
func (a *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the auth_token cookie by setting MaxAge to -1 (expires immediately)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",                  // Ensure it applies to the entire domain
		HttpOnly: true,                 // Maintain security
		Secure:   true,                 // Use Secure for HTTPS
		SameSite: http.SameSiteLaxMode, // Adjust as needed
		MaxAge:   -1,                   // Expires immediately
		Expires:  time.Unix(0, 0),      // Alternative expiration method
	})

	// Clear the refresh_token cookie by setting MaxAge to -1 (expires immediately)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",                  // Ensure it applies to the entire domain
		HttpOnly: true,                 // Maintain security
		Secure:   true,                 // Use Secure for HTTPS
		SameSite: http.SameSiteLaxMode, // Adjust as needed
		MaxAge:   -1,                   // Expires immediately
		Expires:  time.Unix(0, 0),      // Alternative expiration method
	})

	// Optionally, send a response confirming logout
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}

// ExtractJWTToken extracts the JWT token from the Authorization header
// func ExtractJWTToken(r *http.Request) (string, error) {
// 	// Extract JWT from Authorization header
// 	authHeader := r.Header.Get("Authorization")
// 	if authHeader == "" {
// 		return "", errors.New("missing Authorization header")
// 	}

// 	// Token format validation (Bearer <token>)
// 	tokenParts := strings.Split(authHeader, " ")
// 	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
// 		return "", errors.New("invalid Authorization header format")
// 	}

//		return tokenParts[1], nil
//	}
//2// func ExtractJWTToken(r *http.Request) (string, error) {
// 	// 1. Try to get the token from the Authorization header
// 	authHeader := r.Header.Get("Authorization")
// 	if authHeader != "" {
// 		tokenParts := strings.Split(authHeader, " ")
// 		if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
// 			return tokenParts[1], nil
// 		}
// 		return "", errors.New("invalid Authorization header format")
// 	}

// 	// 2. If no Authorization header, check the cookie
// 	cookie, err := r.Cookie("auth_token")
// 	if err == nil {
// 		return cookie.Value, nil
// 	}

//		return "", errors.New("no JWT found in Authorization header or cookie")
//	}
func ExtractJWTToken(r *http.Request) (string, error) {
	// Check the "auth_token" cookie first
	cookie, err := r.Cookie("auth_token")
	fmt.Println("cookie:", cookie) // Debugging
	if err == nil {
		return cookie.Value, nil
	}

	// Fall back to Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing token")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	fmt.Println("tokenParts[1]:", tokenParts[1]) // Debugging
	return tokenParts[1], nil
}

// RefreshHandler will handle the refresh logic for the auth token

// RefreshHandler godoc
//
//	@Summary		Refresh the JWT token via valid Refresh token
//	@Description	Refresh the JWT token
//	@Tags			users
//	@Produce		json
//	@Success		200	{string}	string	"JWT refreshed successfully"
//	@Failure		401	{object}	error
//	@Failure		500	{object}	error
//	@Router			/authentication/refresh [post]
func (a *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RefreshHandler initial call") // Debugging

	// Step 1: Get the refresh_token from the cookies
	cookie, err := r.Cookie("refresh_token")
	// fmt.Println("cookie:", cookie) // Debugging
	fmt.Println("cookie:", cookie) // Debugging
	if err != nil {
		fmt.Println("Error getting refresh token cookie:", err) // Debugging
		utils.UnauthorizedErrorResponse(w, r, errors.New("no refresh token"))
		return
	}

	// Step 2: Validate the refresh token
	refreshToken := cookie.Value
	refreshTokenHash := a.HashToken(refreshToken)

	fmt.Println("refreshTokenHash:", refreshTokenHash) // Debugging

	// Validate the refresh token hash
	fmt.Println("ValidateRefreshTokenByHash Called")
	userID, err := a.ValidateRefreshTokenByHash(r, refreshTokenHash) // Validate the refresh token logic

	if err != nil {
		utils.UnauthorizedErrorResponse(w, r, err)
		return
	}

	fmt.Println("userID:", userID) // Debugging
	// Step 3: Generate a new JWT (auth_token)
	// You might want to pass a fingerprint here if you’re using one
	ip := a.GetIPAddress(r)
	userAgent := r.UserAgent()
	fingerprint := a.GenerateFingerprint(ip, userAgent)           // Optional fingerprint
	fmt.Println("fingerprint:", fingerprint)                      // Debugging
	newAuthToken, err := a.GenerateJWTWithFP(userID, fingerprint) // Pass fingerprint if needed
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	// Step 4: Set the new JWT in the auth_token cookie
	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "auth_token",
	// 	Value:    newAuthToken,
	// 	Path:     "/",
	// 	HttpOnly: true,                                    // Prevents JS access to the cookie
	// 	Secure:   true,                                    // Ensures cookie is sent only over HTTPS
	// 	SameSite: http.SameSiteStrictMode,                 // Prevents CSRF
	// 	Expires:  time.Now().Add(a.Config.Auth.Token.Exp), // Set expiration time for the JWT cookie
	// })

	// Step 5: Respond to the user with a success message
	tokenStore := a.store.Tokens

	fmt.Println("*6") // Debugging

	newRefreshToken, err := a.GenerateRefreshToken()
	fmt.Println("* 7 newRefreshToken:", newRefreshToken) // Debugging
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	newRefreshTokenHash := a.HashToken(newRefreshToken)
	fmt.Println("*8 newRefreshTokenHash:", newRefreshTokenHash) // Debugging

	fmt.Println("*9 a.store.Tokens:", a.store.Tokens) // Debugging
	if a.store.Tokens == nil {
		log.Println("Error: a.store.Tokens is nil")
		utils.InternalServerError(w, r, errors.New("internal server error"))
		return
	}

	// Store the refresh token in the database
	fmt.Println("Calling UpdateRefreshToken") // Debugging
	err = tokenStore.UpdateRefreshToken(r.Context(), userID, newRefreshTokenHash, fingerprint, time.Now().Add(a.Config.Auth.RefreshToken.Exp))

	fmt.Println("*10") // Debugging
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	// Remove the old refresh token from cookies

	fmt.Println("do i get here------------") // Debugging
	// Set the refresh token in a secure, HTTP-only cookie

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		HttpOnly: true,                                           // Ensures it's inaccessible via JavaScript
		Secure:   true,                                           // Only sent over HTTPS
		SameSite: http.SameSiteStrictMode,                        // Adjust as necessary
		Expires:  time.Now().Add(a.Config.Auth.RefreshToken.Exp), // Set expiration based on config
	})
	// Set the JWT in an HTTP-only, Secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    newAuthToken,
		Path:     "/",
		HttpOnly: true,                                    // Ensures it's inaccessible via JavaScript
		Secure:   true,                                    // Only sent over HTTPS
		SameSite: http.SameSiteStrictMode,                 // Prevent CSRF
		Expires:  time.Now().Add(a.Config.Auth.Token.Exp), // Cookie expiration time
	})

	// Set the refresh token in a secure, HTTP-only cookie
	fmt.Println("*11 ")
	w.Write([]byte("JWT refreshed successfully"))
}

// ValidateRefreshTokenByHash checks if the refresh token is valid
func (a *AuthHandler) ValidateRefreshTokenByHash(r *http.Request, refreshToken string) (uuid.UUID, error) {
	// Step 1: Check if the refresh token exists in the database

	tokenRecord, err := a.store.Tokens.GetByRefreshTokenHash(r.Context(), refreshToken) // Use r.Context() here
	if err != nil {
		return uuid.Nil, errors.New("invalid or expired refresh token")
	}

	// Step 2: Ensure the token has not expired
	if tokenRecord.ExpiresAt.Before(time.Now()) {
		return uuid.Nil, errors.New("refresh token has expired")
	}

	// Step 3: Return the email associated with the token
	return uuid.Parse(tokenRecord.UserID)
}
