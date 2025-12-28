package auth

import (
	"context"
	"net/http"
)

// Keys for storing values inside context
type contextKey string

const (
	ContextUserIDKey   contextKey = "user_id"
	ContextUsernameKey contextKey = "username"
)

// RequireJWT protects routes and extracts claims into context
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Extract token from Authorization header
		// authHeader := r.Header.Get("Authorization")
		// if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		// 	http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		// 	return
		// }
		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Missing auth cookie", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		// Validate JWT
		claims, err := ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), ContextUserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, ContextUsernameKey, claims.Username)

		// Continue request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID returns authenticated user ID from context
func GetUserID(r *http.Request) int {
	val := r.Context().Value(ContextUserIDKey)
	if id, ok := val.(int); ok {
		return id
	}
	return 0
}

// GetUsername returns authenticated username from context
func GetUsername(r *http.Request) string {
	val := r.Context().Value(ContextUsernameKey)
	if name, ok := val.(string); ok {
		return name
	}
	return ""
}
