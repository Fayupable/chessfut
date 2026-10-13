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

var cloudflareRanges = []string{
	"173.245.48.0/20", "103.21.244.0/22", "103.22.200.0/22", "103.31.4.0/22",
	"141.101.64.0/18", "108.162.192.0/18", "190.93.240.0/20", "188.114.96.0/20",
	"197.234.240.0/22", "198.41.128.0/17", "162.158.0.0/15", "104.16.0.0/13",
	"104.24.0.0/14", "172.64.0.0/13", "131.0.72.0/22",
	"2400:cb00::/32", "2606:4700::/32", "2803:f800::/32", "2405:b500::/32",
	"2405:8100::/32", "2a06:98c0::/29", "2c0f:f248::/32",
}

var trustedProxyNets = parseCIDRs(cloudflareRanges)

func parseCIDRs(values []string) []*net.IPNet {
	nets := make([]*net.IPNet, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if !strings.Contains(v, "/") {
			if ip := net.ParseIP(v); ip != nil {
				if ip.To4() != nil {
					v += "/32"
				} else {
					v += "/128"
				}
			}
		}
		if _, block, err := net.ParseCIDR(v); err == nil {
			nets = append(nets, block)
		}
	}
	return nets
}

// SetTrustedProxies registers extra proxy addresses (for example the
// server's own public IP, which appears in the forwarded chain when the
// frontend calls the API through Cloudflare) on top of Cloudflare's ranges.
// It must be called once at startup, before the server begins serving.
func SetTrustedProxies(extra []string) {
	trustedProxyNets = append(parseCIDRs(cloudflareRanges), parseCIDRs(extra)...)
}

func isTrustedHop(ip net.IP) bool {
	if ip.IsPrivate() || ip.IsLoopback() {
		return true
	}
	for _, block := range trustedProxyNets {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// clientIP walks X-Forwarded-For from the right, skipping hops we control
// (private/Docker addresses, Cloudflare edges, configured proxies), and
// returns the first address that is not one of ours. Each proxy appends the
// peer it saw to the end, so the right side is the trustworthy side.
func clientIP(r *http.Request) string {
	fwd := r.Header.Get("X-Forwarded-For")
	if fwd == "" {
		return remoteHost(r.RemoteAddr)
	}

	parts := strings.Split(fwd, ",")
	for i := len(parts) - 1; i >= 0; i-- {
		candidate := strings.TrimSpace(parts[i])
		ip := net.ParseIP(candidate)
		if ip != nil && !isTrustedHop(ip) {
			return candidate
		}
	}
	return strings.TrimSpace(parts[0])
}

func remoteHost(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}
