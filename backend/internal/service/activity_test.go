package service

import (
	"testing"
	"time"

	"khhub/internal/domain"
)

func TestActivityStatus(t *testing.T) {
	asOf := domain.Month{Year: 2026, Month: time.March}
	// window: Oct 2025 .. Mar 2026
	allShared := []domain.MonthShare{
		{Year: 2025, Month: time.October, Shared: true},
		{Year: 2025, Month: time.November, Shared: true},
		{Year: 2025, Month: time.December, Shared: true},
		{Year: 2026, Month: time.January, Shared: true},
		{Year: 2026, Month: time.February, Shared: true},
		{Year: 2026, Month: time.March, Shared: true},
	}
	if got := ActivityStatus(allShared, nil, asOf); got != domain.ActivityRegular {
		t.Fatalf("all shared: got %s", got)
	}

	missed := append([]domain.MonthShare{}, allShared...)
	missed[2].Shared = false
	if got := ActivityStatus(missed, nil, asOf); got != domain.ActivityIrregular {
		t.Fatalf("missed one: got %s", got)
	}

	if got := ActivityStatus(nil, nil, asOf); got != domain.ActivityInactive {
		t.Fatalf("none: got %s", got)
	}

	started := time.Date(2026, time.January, 10, 0, 0, 0, 0, time.UTC)
	newPub := []domain.MonthShare{
		{Year: 2026, Month: time.January, Shared: true},
		{Year: 2026, Month: time.February, Shared: true},
		{Year: 2026, Month: time.March, Shared: true},
	}
	if got := ActivityStatus(newPub, &started, asOf); got != domain.ActivityRegular {
		t.Fatalf("new publisher: got %s", got)
	}
}
