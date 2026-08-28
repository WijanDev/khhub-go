package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"khhub/internal/config"

	"github.com/gin-gonic/gin"
)

func TestResetSeedRouteAbsentInProductionEvenWithCallback(t *testing.T) {
	t.Parallel()
	h := NewRouter(config.Config{AppEnv: "production"}, nil, func(context.Context) error {
		t.Fatal("production must not call reset")
		return nil
	})
	req := httptest.NewRequest(http.MethodPost, "/dev/reset-seed", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("POST /dev/reset-seed in production: got %d, want 404", rec.Code)
	}
}

func TestResetSeedRouteAbsentWhenCallbackNil(t *testing.T) {
	t.Parallel()
	h := NewRouter(config.Config{AppEnv: "staging"}, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/dev/reset-seed", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("POST /dev/reset-seed with nil callback: got %d, want 404", rec.Code)
	}
}

func TestResetSeedRouteExistsWhenCallbackSet(t *testing.T) {
	t.Parallel()
	for _, env := range []string{"development", "staging"} {
		env := env
		t.Run(env, func(t *testing.T) {
			t.Parallel()
			h := NewRouter(config.Config{AppEnv: env, SessionTTL: testAuthCfg().SessionTTL}, nil, func(context.Context) error {
				t.Fatal("reset must not run without a session")
				return nil
			})
			req := httptest.NewRequest(http.MethodPost, "/dev/reset-seed", nil)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("POST /dev/reset-seed without cookie: got %d, want 401", rec.Code)
			}
		})
	}
}

func TestPostResetSeedNilForbidden(t *testing.T) {
	t.Parallel()
	r := gin.New()
	r.POST("/dev/reset-seed", postResetSeed(nil))
	req := httptest.NewRequest(http.MethodPost, "/dev/reset-seed", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status: got %d want 403 body=%s", rec.Code, rec.Body.String())
	}
	var payload struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("json: %v", err)
	}
	if !strings.Contains(payload.Error, "desarrollo") {
		t.Fatalf("error: %q", payload.Error)
	}
}

func TestPostResetSeedOKWithSession(t *testing.T) {
	t.Parallel()
	q := newMemorySessions()
	called := false
	r := gin.New()
	cfg := testAuthCfg()
	r.Use(sessionMiddleware(cfg, q))
	r.POST("/auth/login", postLogin(cfg, q, newLoginLimiter()))
	authed := r.Group("/")
	authed.Use(requireAuth())
	authed.POST("/dev/reset-seed", postResetSeed(func(context.Context) error {
		called = true
		return nil
	}))

	login := postLoginRequest(r, "admin@example.com", testPassword)
	if login.Code != http.StatusOK {
		t.Fatalf("login: %d %s", login.Code, login.Body.String())
	}
	ck := sessionCookie(login)
	if ck == nil {
		t.Fatal("login did not set khhub_session")
	}

	req := httptest.NewRequest(http.MethodPost, "/dev/reset-seed", nil)
	req.AddCookie(ck)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("reset: got %d body=%s", rec.Code, rec.Body.String())
	}
	if !called {
		t.Fatal("expected reset callback")
	}
	var payload struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if !payload.OK {
		t.Fatalf("payload: %+v", payload)
	}
}
