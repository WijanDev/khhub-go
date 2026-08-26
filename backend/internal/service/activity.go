package service

import (
	"time"

	"khhub/internal/domain"
)

// ActivityStatus derives regular / irregular / inactive from the last six
// calendar months through asOf. Months before startedPreaching are ignored.
// Regular and irregular are both "active" (shared at least once).
func ActivityStatus(shares []domain.MonthShare, startedPreaching *time.Time, asOf domain.Month) string {
	window := domain.LastNMonths(asOf, 6)
	sharedByKey := make(map[string]bool, len(shares))
	for _, s := range shares {
		sharedByKey[domain.Month{Year: s.Year, Month: s.Month}.Key()] = s.Shared
	}

	considered := 0
	sharedCount := 0
	for _, m := range window {
		if startedPreaching != nil {
			start := domain.MonthFromTime(*startedPreaching)
			if m.Year < start.Year || (m.Year == start.Year && m.Month < start.Month) {
				continue
			}
		}
		considered++
		if sharedByKey[m.Key()] {
			sharedCount++
		}
	}

	if considered == 0 || sharedCount == 0 {
		return domain.ActivityInactive
	}
	if sharedCount == considered {
		return domain.ActivityRegular
	}
	return domain.ActivityIrregular
}

func IsActive(status string) bool {
	return status == domain.ActivityRegular || status == domain.ActivityIrregular
}
