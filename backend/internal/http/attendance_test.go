package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"khhub/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type memoryAttendance struct {
	byKey map[string]store.MeetingAttendance
}

func newMemoryAttendance() *memoryAttendance {
	return &memoryAttendance{byKey: map[string]store.MeetingAttendance{}}
}

func attendanceKey(day time.Time, kind string) string {
	return day.UTC().Format("2006-01-02") + "|" + kind
}

func (m *memoryAttendance) UpsertAttendance(_ context.Context, arg store.UpsertAttendanceParams) (store.MeetingAttendance, error) {
	key := attendanceKey(arg.MeetingDate, arg.Kind)
	if row, ok := m.byKey[key]; ok {
		row.InPerson = arg.InPerson
		row.Online = arg.Online
		m.byKey[key] = row
		return row, nil
	}
	row := store.MeetingAttendance{
		ID:          uuid.New(),
		MeetingDate: arg.MeetingDate,
		Kind:        arg.Kind,
		InPerson:    arg.InPerson,
		Online:      arg.Online,
		CreatedAt:   time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
	}
	m.byKey[key] = row
	return row, nil
}

func attendanceTestHandler(q attendanceWriter) http.Handler {
	r := gin.New()
	r.PUT("/attendance", putAttendance(q))
	return r
}

func putAttendanceJSON(t *testing.T, h http.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPut, "/attendance", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestPutAttendanceUpsertKeepsOneRow(t *testing.T) {
	t.Parallel()
	q := newMemoryAttendance()
	h := attendanceTestHandler(q)

	first := putAttendanceJSON(t, h, `{"date":"2026-08-02","kind":"weekend","inPerson":40,"online":3}`)
	if first.Code != http.StatusOK {
		t.Fatalf("first put: %d %s", first.Code, first.Body.String())
	}
	var created struct {
		ID       string `json:"id"`
		InPerson int32  `json:"inPerson"`
	}
	if err := json.Unmarshal(first.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.InPerson != 40 {
		t.Fatalf("first payload: %+v", created)
	}

	second := putAttendanceJSON(t, h, `{"date":"2026-08-02","kind":"weekend","inPerson":45}`)
	if second.Code != http.StatusOK {
		t.Fatalf("second put: %d %s", second.Code, second.Body.String())
	}
	var updated struct {
		ID       string `json:"id"`
		Kind     string `json:"kind"`
		InPerson int32  `json:"inPerson"`
	}
	if err := json.Unmarshal(second.Body.Bytes(), &updated); err != nil {
		t.Fatal(err)
	}
	if updated.ID != created.ID {
		t.Fatalf("expected same id, got %s then %s", created.ID, updated.ID)
	}
	if updated.InPerson != 45 || updated.Kind != "weekend" {
		t.Fatalf("updated payload: %+v", updated)
	}
	if len(q.byKey) != 1 {
		t.Fatalf("rows: got %d want 1", len(q.byKey))
	}
}

func TestPutAttendanceRejectsBadInput(t *testing.T) {
	t.Parallel()
	h := attendanceTestHandler(newMemoryAttendance())
	cases := []struct {
		name string
		body string
	}{
		{name: "bad date", body: `{"date":"08-02-2026","kind":"weekend","inPerson":10}`},
		{name: "bad kind", body: `{"date":"2026-08-02","kind":"memorial","inPerson":10}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := putAttendanceJSON(t, h, tc.body)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status: got %d want 400 body=%s", rec.Code, rec.Body.String())
			}
		})
	}
}
