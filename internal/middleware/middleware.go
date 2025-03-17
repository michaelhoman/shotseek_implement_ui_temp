package middleware

import (
	"context"
	"fmt"
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
			fmt.Printf("Begin JWTMiddleware\n----\n") // TODO Remove Debugging REMOVE
			// Extract the JWT token using our utility function
			tokenString, err := auth.ExtractJWTToken(r)
			fmt.Printf("\nJWTMiddleware - TokenString: %s\n\n", tokenString) // TODO Remove Debugging REMOVE
			if err != nil {
				fmt.Printf("JWTMiddleware - Exiting - Failed to extract JWT: %s\n", err) // TODO Remove Debugging REMOVE
				utils.Logger.Warn("Failed to extract JWT: ", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ip := authHandler.GetIPAddress(r)
			userAgent := r.UserAgent()
			requestFingerprint := authHandler.GenerateFingerprint(ip, userAgent)

			fmt.Printf("JWTMiddleware - Request Fingerprint: '%s\n", requestFingerprint) // TODO Remove Debugging REMOVE
			// Validate the token
			claims, err := authHandler.ValidateJWT(r, tokenString, requestFingerprint) // Pass fingerprint if needed

			if err != nil {
				fmt.Printf("JWTMiddleware - Exiting 99 - Invalid JWT: %s\n", err) // TODO Remove Debugging REMOVE
				utils.Logger.Warn("Invalid JWT: ", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			fmt.Printf("JWT -------------------------:")                // TODO Remove Debugging REMOVE
			fmt.Printf("JWT Fingerprint: %s\n", claims.Fingerprint)     // TODO Remove Debugging REMOVE
			fmt.Printf("Request Fingerprint: %s\n", requestFingerprint) // TODO Remove Debugging REMOVE

			// Store claims in the context for the next handler
			ctx := context.WithValue(r.Context(), userContextKey, claims)
			r = r.WithContext(ctx)

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
