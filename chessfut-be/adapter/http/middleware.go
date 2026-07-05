package http

import (
	"log/slog"
	"net/http"
	"runtime"
	"strings"
	"time"
)

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

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var before runtime.MemStats
		runtime.ReadMemStats(&before)

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rec, r)

		var after runtime.MemStats
		runtime.ReadMemStats(&after)

		slog.Info("http_request",
			"method", r.Method,
			"path", r.URL.Path,
			"username", r.PathValue("username"),
			"ip", clientIP(r),
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"mem_alloc_delta_bytes", int64(after.Alloc)-int64(before.Alloc),
		)
	}
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.Split(fwd, ",")[0]
	}
	return r.RemoteAddr
}
