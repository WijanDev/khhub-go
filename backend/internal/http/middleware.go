package http

import (
	"context"
	"net/http"

	"khhub/internal/auth"
	"khhub/internal/config"
	"khhub/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// sessionQuerier is the store surface used by login, logout, and the session cookie.
// *store.Queries implements it.
type sessionQuerier interface {
	GetUserByEmail(ctx context.Context, email string) (store.User, error)
	CreateSession(ctx context.Context, arg store.CreateSessionParams) (store.Session, error)
	DeleteExpiredSessions(ctx context.Context) error
	DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (store.GetSessionByTokenHashRow, error)
}

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

func sessionMiddleware(cfg config.Config, q sessionQuerier) gin.HandlerFunc {
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
