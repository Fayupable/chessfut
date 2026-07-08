package http

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	violationThreshold = 20
	violationWindow    = 10 * time.Minute
	banDuration        = 24 * time.Hour
)

type visitor struct {
	limiter     *rate.Limiter
	lastSeen    time.Time
	violations  int
	violationAt time.Time
	bannedUntil time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     rate.Limit
	burst    int
}

func NewRateLimiter(r rate.Limit, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     r,
		burst:    burst,
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{limiter: rate.NewLimiter(rl.rate, rl.burst), lastSeen: time.Now()}
		rl.visitors[ip] = v
		return v
	}
	v.lastSeen = time.Now()
	return v
}

func (rl *RateLimiter) isBanned(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		return false
	}
	return time.Now().Before(v.bannedUntil)
}

// recordViolation tracks malformed/rejected requests per IP. An IP that
// racks up too many in a short window (a scanner probing for .env/.php
// paths, for example) gets banned outright for 24 hours instead of just
// throttled — this frees up the shared chess.com rate-limit budget for
// legitimate users instead of making them queue behind scanner noise.
func (rl *RateLimiter) recordViolation(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{limiter: rate.NewLimiter(rl.rate, rl.burst), lastSeen: time.Now()}
		rl.visitors[ip] = v
	}

	if time.Since(v.violationAt) > violationWindow {
		v.violations = 0
	}
	v.violations++
	v.violationAt = time.Now()

	if v.violations >= violationThreshold {
		v.bannedUntil = time.Now().Add(banDuration)
	}
}

// cleanupLoop evicts idle visitors so the map doesn't grow unbounded under
// distributed scraping (many distinct IPs, each hit once). Banned IPs are
// kept around until their ban actually expires, even if idle.
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute && time.Now().After(v.bannedUntil) {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

type statusCapture struct {
	http.ResponseWriter
	status int
}

func (s *statusCapture) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}

func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)

		if rl.isBanned(ip) {
			writeErrorMessage(w, http.StatusForbidden, "too many invalid requests — temporarily banned")
			return
		}

		if !rl.getVisitor(ip).limiter.Allow() {
			writeErrorMessage(w, http.StatusTooManyRequests, "too many requests")
			return
		}

		sc := &statusCapture{ResponseWriter: w, status: http.StatusOK}
		next(sc, r)

		if sc.status == http.StatusBadRequest {
			rl.recordViolation(ip)
		}
	}
}
