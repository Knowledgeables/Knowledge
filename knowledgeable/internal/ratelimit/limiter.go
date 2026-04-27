package ratelimit

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// IPLimiter is a sliding-window rate limiter keyed by IP address.
type IPLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	max      int
	window   time.Duration
}

func NewIPLimiter(max int, window time.Duration) *IPLimiter {
	l := &IPLimiter{
		requests: make(map[string][]time.Time),
		max:      max,
		window:   window,
	}
	go l.cleanup()
	return l
}

func (l *IPLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)

	prev := l.requests[ip]
	valid := prev[:0]
	for _, t := range prev {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= l.max {
		l.requests[ip] = valid
		return false
	}

	l.requests[ip] = append(valid, now)
	return true
}

// cleanup removes stale entries once per window to prevent unbounded memory growth.
func (l *IPLimiter) cleanup() {
	ticker := time.NewTicker(l.window)
	defer ticker.Stop()
	for range ticker.C {
		l.mu.Lock()
		cutoff := time.Now().Add(-l.window)
		for ip, times := range l.requests {
			valid := times[:0]
			for _, t := range times {
				if t.After(cutoff) {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(l.requests, ip)
			} else {
				l.requests[ip] = valid
			}
		}
		l.mu.Unlock()
	}
}

// RealIP returns the client IP, preferring X-Real-IP set by nginx over RemoteAddr.
func RealIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
