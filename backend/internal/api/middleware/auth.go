package middleware

import (
	"context"
	"github.com/pythonistD/Guess-The-Flag/internal/service/user"
	"net/http"
	"strings"
)

func JWTMiddleware(tokenManager user.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/auth") {
				next.ServeHTTP(w, r)
				return
			}
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			userId, err := tokenManager.ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, "Error while validating token", http.StatusUnauthorized)
			}
			ctx := context.WithValue(r.Context(), "userId", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
