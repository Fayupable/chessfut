package http

import "net/http"

func AdminAuthMiddleware(adminKey string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Admin-Key")
		if key == "" || key != adminKey {
			writeError(w, http.StatusUnauthorized, "invalid admin key")
			return
		}
		next(w, r)
	}
}
