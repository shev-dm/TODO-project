package middleware

import (
	"net/http"

	"github.com/shev-dm/TODO-project/internal/hasher"
)

func Authentication(password string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
					http.Error(w, "Authentication required", http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
