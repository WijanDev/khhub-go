package domain

import (
	"strings"
	"testing"
	"time"
)

func hoursPtr(v float64) *float64 { return &v }

func TestServiceYear(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		date time.Time
		want int
	}{
		{"day before 2025 year", time.Date(2024, time.August, 31, 0, 0, 0, 0, time.UTC), 2024},
		{"start of 2025 year", time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC), 2025},
		{"mid 2025 year", time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC), 2025},
		{"last day of 2025 year", time.Date(2025, time.August, 31, 0, 0, 0, 0, time.UTC), 2025},
		{"start of 2026 year", time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC), 2026},
		{"calendar 2026 still 2026 year", time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC), 2026},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := ServiceYear(tc.date); got != tc.want {
				t.Fatalf("ServiceYear(%s) = %d, want %d", tc.date.Format("2006-01-02"), got, tc.want)
			}
		})
	}
}

func TestServiceYearBounds(t *testing.T) {
	t.Parallel()
	start, end := ServiceYearBounds(2025)
	wantStart := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)
	wantEnd := time.Date(2025, time.August, 31, 0, 0, 0, 0, time.UTC)
	if !start.Equal(wantStart) || !end.Equal(wantEnd) {
		t.Fatalf("ServiceYearBounds(2025) = %s .. %s", start.Format("2006-01-02"), end.Format("2006-01-02"))
	}
	if ServiceYear(start) != 2025 || ServiceYear(end) != 2025 {
		t.Fatalf("bounds are not both in service year 2025")
	}
}

func TestReportsHours(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		rp, sp  bool
		aux     bool
		shared  bool
		hours   *float64
		wantErr string
	}{
		{name: "publisher without hours"},
		{name: "publisher with hours rejected", hours: hoursPtr(10), wantErr: "hours are only recorded for pioneers"},
		{name: "regular pioneer shared with hours", rp: true, shared: true, hours: hoursPtr(70)},
		{name: "special pioneer shared with hours", sp: true, shared: true, hours: hoursPtr(130)},
		{name: "auxiliary that month with hours", aux: true, shared: true, hours: hoursPtr(30)},
		{name: "pioneer did not share, no hours", rp: true},
		{name: "pioneer shared without hours", rp: true, shared: true, wantErr: "hours are required"},
		{name: "pioneer negative hours", rp: true, shared: true, hours: hoursPtr(-1), wantErr: "hours cannot be negative"},
		{name: "pioneer zero hours allowed", rp: true, shared: true, hours: hoursPtr(0)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ReportsHours(tc.rp, tc.sp, tc.aux, tc.shared, tc.hours)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error %v, want substring %q", err, tc.wantErr)
			}
		})
	}
}
