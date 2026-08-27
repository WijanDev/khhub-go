package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv        string
	HTTPAddr      string
	DatabaseURL   string
	CORSOrigins   []string
	CookieSecure  bool
	SessionTTL    time.Duration
	AdminEmail    string
	AdminPassword string
}

func Load() (Config, error) {
	loadDotEnv()
	cfg := Config{
		AppEnv:        getenv("APP_ENV", "development"),
		HTTPAddr:      getenv("HTTP_ADDR", ":8080"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		CORSOrigins:   splitCSV(os.Getenv("CORS_ORIGINS")),
		CookieSecure:  getenvBool("COOKIE_SECURE", false),
		SessionTTL:    7 * 24 * time.Hour,
		AdminEmail:    getenv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}
	if cfg.DatabaseURL == "" {
		return cfg, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.AdminPassword == "" {
		return cfg, fmt.Errorf("ADMIN_PASSWORD is required")
	}
	if len(cfg.AdminPassword) < 10 && cfg.StrictSecrets() {
		return cfg, fmt.Errorf("ADMIN_PASSWORD must be at least 10 characters in production and staging")
	}
	return cfg, nil
}

func (c Config) Production() bool {
	return c.AppEnv == "production"
}

// StrictSecrets is production or staging (public hosts).
func (c Config) StrictSecrets() bool {
	return c.AppEnv == "production" || c.AppEnv == "staging"
}

// AllowsSeedReset is only local development. Never enable /dev/reset-seed on a public hostname.
func (c Config) AllowsSeedReset() bool {
	return c.AppEnv == "development"
}

// AutoDemoSeed loads fictional publishers when the directory is empty (local + staging).
func (c Config) AutoDemoSeed() bool {
	return c.AppEnv == "development" || c.AppEnv == "staging"
}

func loadDotEnv() {
	for _, path := range []string{".env", "../.env"} {
		if _, err := os.Stat(path); err == nil {
			_ = godotenv.Load(path)
			return
		}
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
