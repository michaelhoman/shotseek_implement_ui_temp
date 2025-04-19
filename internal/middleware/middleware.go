package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/michaelhoman/ShotSeek/internal/auth"
	"github.com/michaelhoman/ShotSeek/internal/utils"
)

type contextKey string

const userContextKey contextKey = "user"

// JwtMiddleware validates the JWT and stores the claims in the context
func JwtMiddleware(authHandler *auth.AuthHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Begin JWTMiddleware\n----\n")
			tokenString, err := auth.ExtractJWTToken(r)
			fmt.Printf("\nJWTMiddleware - TokenString: %s\n\n", tokenString)

			ip := authHandler.GetIPAddress(r)
			userAgent := r.UserAgent()
			requestFingerprint := authHandler.GenerateFingerprint(ip, userAgent)

			var claims *auth.Claims // <-- shared claims variable

			if err != nil {
				utils.Logger.Info("JWT invalid or missing, checking refresh token...")

				refreshCookie, err := r.Cookie("refresh_token")
				if err != nil {
					utils.Logger.Warn("No refresh token found. Rejecting.")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				refreshToken := refreshCookie.Value
				userEmail, err := authHandler.ValidateRefreshTokenByHash(r, authHandler.HashToken(refreshToken))
				if err != nil {
					utils.Logger.Warn("Refresh token invalid. Rejecting.")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				// Generate new tokens
				newAuthToken, err := authHandler.GenerateJWTWithFP(userEmail, requestFingerprint)
				if err != nil {
					utils.Logger.Warn("Failed to generate new JWT.")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				newRefreshToken, err := authHandler.GenerateRefreshToken()
				if err != nil {
					utils.Logger.Warn("Failed to generate new refresh token.")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				setAuthCookies(w, newAuthToken, newRefreshToken, authHandler)

				// Validate the newly issued token to extract claims
				claims, err = authHandler.ValidateJWT(r, newAuthToken, requestFingerprint)
				if err != nil {
					utils.Logger.Warn("New JWT validation failed after refresh.")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			} else {
				// Original token was present, so validate it here
				claims, err = authHandler.ValidateJWT(r, tokenString, requestFingerprint)
				if err != nil {
					utils.Logger.Warn("Invalid original JWT.")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}

			fmt.Printf("JWT Fingerprint: %s\n", claims.Fingerprint)
			fmt.Printf("Request Fingerprint: %s\n", requestFingerprint)

			// Store claims in context and continue
			ctx := context.WithValue(r.Context(), userContextKey, claims)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func setAuthCookies(w http.ResponseWriter, authToken, refreshToken string, authHandler *auth.AuthHandler) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    authToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(authHandler.Config.Auth.Token.Exp),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(authHandler.Config.Auth.RefreshToken.Exp),
	})
}
