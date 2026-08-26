package http

import (
	"net/http"
	"sync"
	"time"

	"khhub/internal/auth"
	"khhub/internal/config"
	"khhub/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ctxKey string

const userCtxKey ctxKey = "user"

type AuthUser struct {
	ID    uuid.UUID
	Email string
}

func currentUser(c *gin.Context) (AuthUser, bool) {
	v, ok := c.Get(string(userCtxKey))
	if !ok {
		return AuthUser{}, false
	}
	u, ok := v.(AuthUser)
	return u, ok
}

func sessionMiddleware(cfg config.Config, q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(auth.CookieName)
		if err != nil || token == "" {
			c.Next()
			return
		}
		row, err := q.GetSessionByTokenHash(c.Request.Context(), auth.HashToken(token))
		if err != nil {
			c.Next()
			return
		}
		c.Set(string(userCtxKey), AuthUser{ID: row.UserID, Email: row.Email})
		c.Next()
	}
}

func requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := currentUser(c); !ok {
			jsonError(c, http.StatusUnauthorized, "no autenticado")
			c.Abort()
			return
		}
		c.Next()
	}
}

type loginLimiter struct {
	mu      sync.Mutex
	hits    map[string][]time.Time
	window  time.Duration
	maxHits int
}

func newLoginLimiter() *loginLimiter {
	return &loginLimiter{
		hits:    make(map[string][]time.Time),
		window:  15 * time.Minute,
		maxHits: 10,
	}
}

func (l *loginLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-l.window)
	kept := l.hits[ip][:0]
	for _, t := range l.hits[ip] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= l.maxHits {
		l.hits[ip] = kept
		return false
	}
	l.hits[ip] = append(kept, now)
	return true
}
