package domain

import (
	"fmt"
	"time"
)

const (
	GenderMale   = "male"
	GenderFemale = "female"

	SpiritualStudent             = "student"
	SpiritualUnbaptizedPublisher = "unbaptized_publisher"
	SpiritualPublisher           = "publisher"

	ActivityRegular   = "regular"
	ActivityIrregular = "irregular"
	ActivityInactive  = "inactive"

	MeetingMidweek = "midweek"
	MeetingWeekend = "weekend"
)

type Month struct {
	Year  int
	Month time.Month
}

func (m Month) Key() string {
	return fmt.Sprintf("%04d-%02d", m.Year, int(m.Month))
}

func (m Month) AddMonths(n int) Month {
	t := time.Date(m.Year, m.Month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, n, 0)
	return Month{Year: t.Year(), Month: t.Month()}
}

func MonthFromTime(t time.Time) Month {
	return Month{Year: t.Year(), Month: t.Month()}
}

// ServiceYear is the JW service year that contains t.
// The 2025 service year ran 1 Sep 2024 – 31 Aug 2025.
func ServiceYear(t time.Time) int {
	if t.Month() >= time.September {
		return t.Year() + 1
	}
	return t.Year()
}

func ServiceYearBounds(year int) (start, end time.Time) {
	start = time.Date(year-1, time.September, 1, 0, 0, 0, 0, time.UTC)
	end = time.Date(year, time.August, 31, 0, 0, 0, 0, time.UTC)
	return start, end
}

func LastNMonths(asOf Month, n int) []Month {
	out := make([]Month, n)
	for i := 0; i < n; i++ {
		out[n-1-i] = asOf.AddMonths(-i)
	}
	return out
}

type MonthShare struct {
	Year   int
	Month  time.Month
	Shared bool
}

func ReportsHours(isRegularPioneer, isSpecialPioneer, auxiliary, shared bool, hours *float64) error {
	hourReporter := isRegularPioneer || isSpecialPioneer || auxiliary
	if !hourReporter {
		if hours != nil {
			return fmt.Errorf("hours are only recorded for pioneers")
		}
		return nil
	}
	if hours != nil && *hours < 0 {
		return fmt.Errorf("hours cannot be negative")
	}
	if shared && hours == nil {
		return fmt.Errorf("hours are required when a pioneer shared in the ministry")
	}
	return nil
}

func MustReport(spiritualStatus string) bool {
	return spiritualStatus == SpiritualPublisher || spiritualStatus == SpiritualUnbaptizedPublisher
}

func IsHourReporter(isRegularPioneer, isSpecialPioneer, auxiliary bool) bool {
	return isRegularPioneer || isSpecialPioneer || auxiliary
}
