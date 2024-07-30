package middleware

import (
	"github.com/shev-dm/TODO-project/internal/hasher"
	"net/http"
	"os"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if len(password) > 0 {
			var jw string
			cookie, err := r.Cookie("token")
			if err == nil {
				jw = cookie.Value
			}
			var valid bool
			signedToken, err := hasher.GenerateToken(password)
			if err != nil {
				return
			}

			if signedToken == jw {
				valid = true
			}
			if !valid {
				http.Error(w, "authentification required", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
