package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"khhub/internal/domain"
	"khhub/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type memoryReports struct {
	pubs    map[uuid.UUID]store.GetPublisherForReportRow
	upserts int
}

func (m *memoryReports) GetPublisherForReport(_ context.Context, id uuid.UUID) (store.GetPublisherForReportRow, error) {
	row, ok := m.pubs[id]
	if !ok {
		return store.GetPublisherForReportRow{}, pgx.ErrNoRows
	}
	return row, nil
}

func (m *memoryReports) UpsertReport(_ context.Context, arg store.UpsertReportParams) (store.FieldServiceReport, error) {
	m.upserts++
	return store.FieldServiceReport{ID: uuid.New(), PublisherID: arg.PublisherID, Year: arg.Year, Month: arg.Month}, nil
}

func (m *memoryReports) ListSharesForPublisher(context.Context, store.ListSharesForPublisherParams) ([]store.ListSharesForPublisherRow, error) {
	return nil, nil
}

func (m *memoryReports) UpdatePublisherActivity(context.Context, store.UpdatePublisherActivityParams) error {
	return nil
}

func reportsTestHandler(q reportStore) http.Handler {
	r := gin.New()
	r.PUT("/reports", putReports(q))
	return r
}

func hoursPtr(v float64) *float64 { return &v }

func TestPutReportsHourRules(t *testing.T) {
	t.Parallel()

	publisherID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	tests := []struct {
		name       string
		pub        store.GetPublisherForReportRow
		hours      *float64
		aux        bool
		shared     bool
		wantStatus int
		wantSubstr string
		wantSaved  bool
	}{
		{
			name: "regular pioneer with hours",
			pub: store.GetPublisherForReportRow{
				ID: publisherID, SpiritualStatus: domain.SpiritualPublisher, IsRegularPioneer: true,
			},
			hours: hoursPtr(70), shared: true,
			wantStatus: http.StatusOK, wantSaved: true,
		},
		{
			name: "special pioneer with hours",
			pub: store.GetPublisherForReportRow{
				ID: publisherID, SpiritualStatus: domain.SpiritualPublisher, IsSpecialPioneer: true,
			},
			hours: hoursPtr(130), shared: true,
			wantStatus: http.StatusOK, wantSaved: true,
		},
		{
			name: "auxiliary that month with hours",
			pub: store.GetPublisherForReportRow{
				ID: publisherID, SpiritualStatus: domain.SpiritualPublisher,
			},
			hours: hoursPtr(30), aux: true, shared: true,
			wantStatus: http.StatusOK, wantSaved: true,
		},
		{
			name: "publisher without hours",
			pub: store.GetPublisherForReportRow{
				ID: publisherID, SpiritualStatus: domain.SpiritualPublisher,
			},
			shared:     true,
			wantStatus: http.StatusOK, wantSaved: true,
		},
		{
			name: "publisher with hours rejected",
			pub: store.GetPublisherForReportRow{
				ID: publisherID, SpiritualStatus: domain.SpiritualPublisher,
			},
			hours: hoursPtr(10), shared: true,
			wantStatus: http.StatusBadRequest, wantSubstr: "hours are only recorded for pioneers",
		},
		{
			name: "pioneer shared without hours rejected",
			pub: store.GetPublisherForReportRow{
				ID: publisherID, SpiritualStatus: domain.SpiritualPublisher, IsRegularPioneer: true,
			},
			shared:     true,
			wantStatus: http.StatusBadRequest, wantSubstr: "hours are required",
		},
		{
			name: "pioneer negative hours rejected",
			pub: store.GetPublisherForReportRow{
				ID: publisherID, SpiritualStatus: domain.SpiritualPublisher, IsRegularPioneer: true,
			},
			hours: hoursPtr(-1), shared: true,
			wantStatus: http.StatusBadRequest, wantSubstr: "hours cannot be negative",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			q := &memoryReports{pubs: map[uuid.UUID]store.GetPublisherForReportRow{publisherID: tc.pub}}
			body, err := json.Marshal(reportsBatchRequest{
				Year:  2026,
				Month: 3,
				Reports: []reportUpsert{{
					PublisherID:      publisherID,
					SharedInMinistry: tc.shared,
					Hours:            tc.hours,
					AuxiliaryPioneer: tc.aux,
				}},
			})
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPut, "/reports", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			reportsTestHandler(q).ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status: got %d want %d body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if tc.wantSubstr != "" && !strings.Contains(rec.Body.String(), tc.wantSubstr) {
				t.Fatalf("body %q does not contain %q", rec.Body.String(), tc.wantSubstr)
			}
			if tc.wantSaved && q.upserts != 1 {
				t.Fatalf("upserts: got %d want 1", q.upserts)
			}
			if !tc.wantSaved && q.upserts != 0 {
				t.Fatalf("upserts: got %d want 0", q.upserts)
			}
		})
	}
}
