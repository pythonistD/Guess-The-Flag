package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/pythonistD/Guess-The-Flag/internal/service/user"
)

func JWTMiddleware(tokenManager user.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Пропускаем OPTIONS запросы для CORS preflight
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			// Список публичных маршрутов, не требующих аутентификации
			publicPaths := []string{
				"/auth/login",
				"/auth/register",
			}

			// Проверяем, является ли текущий путь публичным
			for _, path := range publicPaths {
				if r.URL.Path == path {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Для всех остальных путей проверяем токен
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			userId, err := tokenManager.ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, "Error while validating token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
