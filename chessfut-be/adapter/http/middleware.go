package http

import (
	"crypto/subtle"
	"log/slog"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func AdminAuthMiddleware(adminKey string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Admin-Key")
		if key == "" || subtle.ConstantTimeCompare([]byte(key), []byte(adminKey)) != 1 {
			writeErrorMessage(w, http.StatusUnauthorized, "invalid admin key")
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

func writeErrorMessage(w http.ResponseWriter, status int, publicMessage string) {
	writeJSON(w, status, map[string]string{"error": publicMessage})
}

func clientIP(r *http.Request) string {
	fwd := r.Header.Get("X-Forwarded-For")
	if fwd == "" {
		return r.RemoteAddr
	}

	parts := strings.Split(fwd, ",")
	for i := len(parts) - 1; i >= 0; i-- {
		candidate := strings.TrimSpace(parts[i])
		ip := net.ParseIP(candidate)
		if ip != nil && !isPrivateOrInternalIP(ip) {
			return candidate
		}
	}
	return strings.TrimSpace(parts[0])
}

func isPrivateOrInternalIP(ip net.IP) bool {
	return ip.IsPrivate() || ip.IsLoopback()
}
