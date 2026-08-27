package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"khhub/internal/config"
)

func TestHealthHasNoAPIPrefix(t *testing.T) {
	h := NewRouter(config.Config{AppEnv: "development"}, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /health: got %d", rec.Code)
	}
	old := httptest.NewRecorder()
	h.ServeHTTP(old, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if old.Code == http.StatusOK {
		t.Fatal("GET /api/health must not remain mounted")
	}
}
