package middleware

import (
	"github.com/iamvkosarev/learning-cards/internal/config"
	"net/http"
	"strconv"
)

func CorsWithOptions(next http.Handler, options config.CorsOptions) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			allowed := false
			for _, allowedOrigin := range options.AllowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}

			if allowed || len(options.AllowedOrigins) == 0 {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(options.MaxAge))

				if options.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}
