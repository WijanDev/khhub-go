package http

import (
	"net/http"
	"time"

	"khhub/internal/auth"
	"khhub/internal/config"
	"khhub/internal/store"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

func setSessionCookie(c *gin.Context, cfg config.Config, token string, ttl time.Duration) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(auth.CookieName, token, int(ttl.Seconds()), "/", "", cfg.CookieSecure, true)
}

func clearSessionCookie(c *gin.Context, cfg config.Config) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(auth.CookieName, "", -1, "/", "", cfg.CookieSecure, true)
}

func postLogin(cfg config.Config, q *store.Queries, limiter *loginLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.allow(c.ClientIP()) {
			jsonError(c, http.StatusTooManyRequests, "demasiados intentos; espera unos minutos")
			return
		}
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			jsonError(c, http.StatusBadRequest, "correo o contraseña no válidos")
			return
		}
		user, err := q.GetUserByEmail(c.Request.Context(), trim(req.Email))
		if err != nil || !auth.CheckPassword(user.PasswordHash, req.Password) {
			jsonError(c, http.StatusUnauthorized, "correo o contraseña incorrectos")
			return
		}
		plain, hash, err := auth.NewSessionToken()
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo iniciar sesión")
			return
		}
		expires := time.Now().Add(cfg.SessionTTL)
		if _, err := q.CreateSession(c.Request.Context(), store.CreateSessionParams{
			UserID:    user.ID,
			TokenHash: hash,
			ExpiresAt: expires,
		}); err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo iniciar sesión")
			return
		}
		_ = q.DeleteExpiredSessions(c.Request.Context())
		setSessionCookie(c, cfg, plain, cfg.SessionTTL)
		c.JSON(http.StatusOK, gin.H{"id": user.ID, "email": user.Email})
	}
}

func postLogout(cfg config.Config, q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie(auth.CookieName); err == nil && token != "" {
			_ = q.DeleteSessionByTokenHash(c.Request.Context(), auth.HashToken(token))
		}
		clearSessionCookie(c, cfg)
		c.Status(http.StatusNoContent)
	}
}

func getMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := currentUser(c)
		if !ok {
			jsonError(c, http.StatusUnauthorized, "no autenticado")
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": u.ID, "email": u.Email})
	}
}

func postChangePassword(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := currentUser(c)
		if !ok {
			jsonError(c, http.StatusUnauthorized, "no autenticado")
			return
		}
		var req changePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			jsonError(c, http.StatusBadRequest, "la nueva contraseña debe tener al menos 8 caracteres")
			return
		}
		user, err := q.GetUserByID(c.Request.Context(), u.ID)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo cambiar la contraseña")
			return
		}
		if !auth.CheckPassword(user.PasswordHash, req.CurrentPassword) {
			jsonError(c, http.StatusUnauthorized, "la contraseña actual no es correcta")
			return
		}
		hash, err := auth.HashPassword(req.NewPassword)
		if err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}
		if err := q.UpdateUserPassword(c.Request.Context(), store.UpdateUserPasswordParams{
			ID:           u.ID,
			PasswordHash: hash,
		}); err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo cambiar la contraseña")
			return
		}
		c.Status(http.StatusNoContent)
	}
}
