package http

import (
	"context"
	"net/http"
	"time"

	"khhub/internal/config"
	"khhub/internal/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, q *store.Queries, resetSeed func(ctx context.Context) error) http.Handler {
	if cfg.Production() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	// Traefik sits on the Docker network; trust private hops so ClientIP is the browser.
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"})
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(secureHeaders())

	if len(cfg.CORSOrigins) > 0 {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     cfg.CORSOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	limiter := newLoginLimiter()
	r.Use(sessionMiddleware(cfg, q))

	r.POST("/auth/login", postLogin(cfg, q, limiter))
	r.POST("/auth/logout", postLogout(cfg, q))

	authed := r.Group("/")
	authed.Use(requireAuth())
	{
		authed.GET("/auth/me", getMe())
		authed.POST("/auth/change-password", postChangePassword(q))

		authed.GET("/congregation", getCongregation(q, resetSeed != nil))
		authed.PUT("/congregation", putCongregation(q, resetSeed != nil))
		if resetSeed != nil {
			authed.POST("/dev/reset-seed", postResetSeed(resetSeed))
		}

		authed.GET("/households", listHouseholds(q))
		authed.POST("/households", createHousehold(q))
		authed.PUT("/households/:id", updateHousehold(q))
		authed.DELETE("/households/:id", deleteHousehold(q))

		authed.GET("/publishers", listPublishers(q))
		authed.POST("/publishers", createPublisher(q))
		authed.GET("/publishers/:id", getPublisher(q))
		authed.PUT("/publishers/:id", updatePublisher(q))
		authed.DELETE("/publishers/:id", deletePublisher(q))

		authed.GET("/reports", listReports(q))
		authed.PUT("/reports", putReports(q))

		authed.GET("/attendance", listAttendance(q))
		authed.PUT("/attendance", putAttendance(q))
		authed.DELETE("/attendance/:id", deleteAttendance(q))

		authed.GET("/dashboard", getDashboard(q))
	}

	return r
}

func secureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "same-origin")
		c.Header("Cache-Control", "no-store")
		c.Next()
	}
}
