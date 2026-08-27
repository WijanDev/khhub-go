package http

import (
	"strings"
	"sync"
	"time"
)

const (
	loginFailMax    = 10
	loginFailWindow = 15 * time.Minute
)

type loginLimiter struct {
	mu      sync.Mutex
	hits    map[string][]time.Time
	window  time.Duration
	maxHits int
}

func newLoginLimiter() *loginLimiter {
	return newLoginLimiterWith(loginFailMax, loginFailWindow)
}

func newLoginLimiterWith(maxHits int, window time.Duration) *loginLimiter {
	return &loginLimiter{
		hits:    make(map[string][]time.Time),
		window:  window,
		maxHits: maxHits,
	}
}

func limiterKeys(ip, email string) []string {
	return []string{
		"ip:" + ip,
		"email:" + strings.ToLower(strings.TrimSpace(email)),
	}
}

func (l *loginLimiter) pruneLocked(key string, now time.Time) []time.Time {
	cutoff := now.Add(-l.window)
	var kept []time.Time
	for _, t := range l.hits[key] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) == 0 {
		delete(l.hits, key)
		return nil
	}
	l.hits[key] = kept
	return kept
}

func (l *loginLimiter) blocked(ip, email string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	for _, key := range limiterKeys(ip, email) {
		if len(l.pruneLocked(key, now)) >= l.maxHits {
			return true
		}
	}
	return false
}

func (l *loginLimiter) fail(ip, email string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	for _, key := range limiterKeys(ip, email) {
		kept := l.pruneLocked(key, now)
		l.hits[key] = append(kept, now)
	}
}

func (l *loginLimiter) clear(ip, email string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, key := range limiterKeys(ip, email) {
		delete(l.hits, key)
	}
}
