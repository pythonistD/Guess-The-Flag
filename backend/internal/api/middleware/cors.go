package middleware

import (
	"log"
	"net/http"
)

func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Логируем все входящие запросы
			log.Printf("🌐 CORS: %s %s from %s", r.Method, r.URL.Path, r.Header.Get("Origin"))

			// Всегда устанавливаем CORS заголовки первыми
			origin := r.Header.Get("Origin")
			if origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				log.Printf("✅ CORS: Allowed origin %s", origin)
			} else if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				log.Printf("⚠️  CORS: Allowing origin %s", origin)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				log.Printf("🔓 CORS: No origin header, using wildcard")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Обрабатываем preflight запросы
			if r.Method == "OPTIONS" {
				log.Printf("✈️  CORS: Handling preflight request for %s", r.URL.Path)
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
