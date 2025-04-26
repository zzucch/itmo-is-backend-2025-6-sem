package server

import (
	"net/http"
)

func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=1")
		next.ServeHTTP(w, r)
	})
}
