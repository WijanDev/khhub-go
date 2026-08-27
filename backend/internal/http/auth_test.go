package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"khhub/internal/auth"
	"khhub/internal/config"
	"khhub/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const testPassword = "correct-horse"

var (
	testUserID       = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	testPasswordHash string
)

func init() {
	gin.SetMode(gin.TestMode)
	h, err := auth.HashPassword(testPassword)
	if err != nil {
		panic(err)
	}
	testPasswordHash = h
}

type memorySessions struct {
	user   store.User
	byHash map[string]store.GetSessionByTokenHashRow
}

func newMemorySessions() *memorySessions {
	return &memorySessions{
		user: store.User{
			ID:           testUserID,
			Email:        "admin@example.com",
			PasswordHash: testPasswordHash,
		},
		byHash: map[string]store.GetSessionByTokenHashRow{},
	}
}

func (m *memorySessions) GetUserByEmail(_ context.Context, email string) (store.User, error) {
	if !strings.EqualFold(strings.TrimSpace(email), m.user.Email) {
		return store.User{}, pgx.ErrNoRows
	}
	return m.user, nil
}

func (m *memorySessions) CreateSession(_ context.Context, arg store.CreateSessionParams) (store.Session, error) {
	id := uuid.New()
	m.byHash[arg.TokenHash] = store.GetSessionByTokenHashRow{
		ID:        id,
		UserID:    arg.UserID,
		TokenHash: arg.TokenHash,
		ExpiresAt: arg.ExpiresAt,
		Email:     m.user.Email,
	}
	return store.Session{ID: id, UserID: arg.UserID, TokenHash: arg.TokenHash, ExpiresAt: arg.ExpiresAt}, nil
}

func (m *memorySessions) DeleteExpiredSessions(context.Context) error { return nil }

func (m *memorySessions) DeleteSessionByTokenHash(_ context.Context, tokenHash string) error {
	delete(m.byHash, tokenHash)
	return nil
}

func (m *memorySessions) GetSessionByTokenHash(_ context.Context, tokenHash string) (store.GetSessionByTokenHashRow, error) {
	row, ok := m.byHash[tokenHash]
	if !ok || !row.ExpiresAt.After(time.Now()) {
		return store.GetSessionByTokenHashRow{}, pgx.ErrNoRows
	}
	return row, nil
}

func testAuthCfg() config.Config {
	return config.Config{AppEnv: "development", SessionTTL: time.Hour, CookieSecure: false}
}

func authTestHandler(q sessionQuerier) http.Handler {
	return authTestHandlerWithLimiter(q, newLoginLimiter())
}

func authTestHandlerWithLimiter(q sessionQuerier, limiter *loginLimiter) http.Handler {
	cfg := testAuthCfg()
	r := gin.New()
	r.Use(sessionMiddleware(cfg, q))
	r.POST("/auth/login", postLogin(cfg, q, limiter))
	r.POST("/auth/logout", postLogout(cfg, q))
	authed := r.Group("/")
	authed.Use(requireAuth())
	authed.GET("/auth/me", getMe())
	authed.GET("/publishers", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func postLoginRequest(h http.Handler, email, password string) *httptest.ResponseRecorder {
	body := `{"email":"` + email + `","password":"` + password + `"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "192.0.2.10:1234"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func sessionCookie(rec *httptest.ResponseRecorder) *http.Cookie {
	for _, c := range rec.Result().Cookies() {
		if c.Name == auth.CookieName {
			return c
		}
	}
	return nil
}

func TestAuthHandlers(t *testing.T) {
	t.Parallel()

	type want struct {
		status      int
		errorSubstr string
		cookie      bool
		email       string
	}

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		cookie string
		want   want
	}{
		{
			name:   "login success sets khhub_session",
			method: http.MethodPost,
			path:   "/auth/login",
			body:   `{"email":"admin@example.com","password":"correct-horse"}`,
			want:   want{status: http.StatusOK, cookie: true, email: "admin@example.com"},
		},
		{
			name:   "login bad password",
			method: http.MethodPost,
			path:   "/auth/login",
			body:   `{"email":"admin@example.com","password":"wrong-password"}`,
			want:   want{status: http.StatusUnauthorized, errorSubstr: "incorrectos"},
		},
		{
			name:   "login unknown email",
			method: http.MethodPost,
			path:   "/auth/login",
			body:   `{"email":"nobody@example.com","password":"correct-horse"}`,
			want:   want{status: http.StatusUnauthorized, errorSubstr: "incorrectos"},
		},
		{
			name:   "login missing fields",
			method: http.MethodPost,
			path:   "/auth/login",
			body:   `{"email":"admin@example.com"}`,
			want:   want{status: http.StatusBadRequest, errorSubstr: "no válidos"},
		},
		{
			name:   "me without cookie",
			method: http.MethodGet,
			path:   "/auth/me",
			want:   want{status: http.StatusUnauthorized, errorSubstr: "no autenticado"},
		},
		{
			name:   "publishers without cookie",
			method: http.MethodGet,
			path:   "/publishers",
			want:   want{status: http.StatusUnauthorized, errorSubstr: "no autenticado"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := authTestHandler(newMemorySessions())
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			if tc.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tc.want.status {
				t.Fatalf("status: got %d want %d body=%s", rec.Code, tc.want.status, rec.Body.String())
			}

			ck := sessionCookie(rec)
			if tc.want.cookie {
				if ck == nil || ck.Value == "" {
					t.Fatal("expected Set-Cookie khhub_session")
				}
				if !ck.HttpOnly {
					t.Fatal("khhub_session must be HttpOnly")
				}
				if ck.Path != "/" {
					t.Fatalf("cookie path: got %q", ck.Path)
				}
			} else if ck != nil && ck.Value != "" && ck.MaxAge >= 0 {
				t.Fatalf("did not expect a session cookie, got %q", ck.Value)
			}

			if tc.want.errorSubstr != "" {
				var payload struct {
					Error string `json:"error"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
					t.Fatalf("json: %v body=%s", err, rec.Body.String())
				}
				if !strings.Contains(payload.Error, tc.want.errorSubstr) {
					t.Fatalf("error %q does not contain %q", payload.Error, tc.want.errorSubstr)
				}
			}
			if tc.want.email != "" {
				var payload struct {
					Email string `json:"email"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
					t.Fatalf("json: %v", err)
				}
				if payload.Email != tc.want.email {
					t.Fatalf("email: got %q", payload.Email)
				}
			}
		})
	}
}

func TestAuthMeWithSessionCookie(t *testing.T) {
	t.Parallel()
	q := newMemorySessions()
	h := authTestHandler(q)

	login := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(
		`{"email":"admin@example.com","password":"correct-horse"}`,
	))
	login.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	h.ServeHTTP(loginRec, login)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login: %d %s", loginRec.Code, loginRec.Body.String())
	}
	ck := sessionCookie(loginRec)
	if ck == nil {
		t.Fatal("login did not set khhub_session")
	}

	me := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	me.AddCookie(ck)
	meRec := httptest.NewRecorder()
	h.ServeHTTP(meRec, me)
	if meRec.Code != http.StatusOK {
		t.Fatalf("GET /auth/me: %d %s", meRec.Code, meRec.Body.String())
	}
	var payload struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(meRec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.ID != testUserID.String() || payload.Email != "admin@example.com" {
		t.Fatalf("me payload: %+v", payload)
	}
}

func TestProtectedRouteWithoutCookieOnFullRouter(t *testing.T) {
	t.Parallel()
	h := NewRouter(testAuthCfg(), nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/congregation", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("GET /congregation: got %d", rec.Code)
	}
	var payload struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Error != "no autenticado" {
		t.Fatalf("error: %q", payload.Error)
	}
}

func TestLoginNthFailureReturns429(t *testing.T) {
	t.Parallel()
	limiter := newLoginLimiterWith(3, time.Hour)
	h := authTestHandlerWithLimiter(newMemorySessions(), limiter)
	for i := 0; i < 3; i++ {
		rec := postLoginRequest(h, "admin@example.com", "wrong-password")
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("fail %d: got %d %s", i+1, rec.Code, rec.Body.String())
		}
	}
	rec := postLoginRequest(h, "admin@example.com", "wrong-password")
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("N+1: got %d %s", rec.Code, rec.Body.String())
	}
	var payload struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(payload.Error, "demasiados") {
		t.Fatalf("error: %q", payload.Error)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After")
	}
}

func TestLoginSuccessClearsRateLimit(t *testing.T) {
	t.Parallel()
	limiter := newLoginLimiterWith(3, time.Hour)
	h := authTestHandlerWithLimiter(newMemorySessions(), limiter)
	for i := 0; i < 2; i++ {
		rec := postLoginRequest(h, "admin@example.com", "wrong-password")
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("fail %d: got %d", i+1, rec.Code)
		}
	}
	ok := postLoginRequest(h, "admin@example.com", testPassword)
	if ok.Code != http.StatusOK {
		t.Fatalf("success: got %d %s", ok.Code, ok.Body.String())
	}
	for i := 0; i < 3; i++ {
		rec := postLoginRequest(h, "admin@example.com", "wrong-password")
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("after clear fail %d: got %d %s", i+1, rec.Code, rec.Body.String())
		}
	}
}
