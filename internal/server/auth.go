package server

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/ddeville/kattungar-notify/internal/store"
)

func ApiKeyAuth(keys []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := getAuthToken(r)
			if slices.Contains(keys, key) {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}

		})
	}
}

type key int

const DeviceAuthContextKey key = iota

func DeviceAuth(store *store.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := getAuthToken(r)
			device, err := store.GetDevice(key)
			if err == nil && device != nil {
				ctx := context.WithValue(r.Context(), DeviceAuthContextKey, device)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}

		})
	}
}

func getAuthToken(r *http.Request) string {
	var key string
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToLower(bearer[0:6]) == "bearer" {
		key = bearer[7:]
	}
	return key
}
