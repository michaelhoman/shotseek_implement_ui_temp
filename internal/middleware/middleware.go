package middleware

import (
	"context"
	"net/http"

	"github.com/michaelhoman/ShotSeek/internal/auth"
	"github.com/michaelhoman/ShotSeek/internal/utils"
)

type contextKey string

const userContextKey contextKey = "user"

// JwtMiddleware validates the JWT and stores the claims in the context
func JwtMiddleware(authHandler *auth.AuthHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the JWT token using our utility function
			tokenString, err := auth.ExtractJWTToken(r)
			if err != nil {
				utils.Logger.Warn("Failed to extract JWT: ", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Validate the token
			claims, err := authHandler.ValidateJWT(r, tokenString, "") // Pass fingerprint if needed
			if err != nil {
				utils.Logger.Warn("Invalid JWT: ", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Store claims in the context for the next handler
			ctx := context.WithValue(r.Context(), userContextKey, claims)
			r = r.WithContext(ctx)

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
