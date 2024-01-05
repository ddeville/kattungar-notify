package server

import (
	"net/http"
	"slices"
	"strings"
)

func ApiKeyAuth(keys []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var key string
			bearer := r.Header.Get("Authorization")
			if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
				key = bearer[7:]
			}
			if slices.Contains(keys, key) {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}

		})
	}
}
