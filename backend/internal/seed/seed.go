package seed

import (
	"context"
	"fmt"
	"strings"
	"time"

	"khhub/internal/domain"
	"khhub/internal/service"
	"khhub/internal/store"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// demoDataTables are wiped by Reset. users and sessions must never appear here.
var demoDataTables = []string{
	"field_service_reports",
	"meeting_attendance",
	"publishers",
	"households",
}

func truncateDemoSQL() string {
	return "TRUNCATE " + strings.Join(demoDataTables, ", ") + " RESTART IDENTITY CASCADE"
}

// Reset wipes congregation data (not users/sessions) and loads the demo sample.
func Reset(ctx context.Context, pool *pgxpool.Pool, q *store.Queries) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin seed tx: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := q.WithTx(tx)

	if _, err := tx.Exec(ctx, truncateDemoSQL()); err != nil {
		return fmt.Errorf("truncate demo tables: %w", err)
	}

	if _, err := qtx.UpdateCongregation(ctx, store.UpdateCongregationParams{
		Name:       "Demo Norte",
		Number:     "99999",
		MidweekDay: 4,
		WeekendDay: 0,
	}); err != nil {
		return fmt.Errorf("congregation: %w", err)
	}

	garcia, err := qtx.CreateHousehold(ctx, store.CreateHouseholdParams{
		Name: "García", Address: "Calle Mayor 12", Notes: "Grupo 1",
	})
	if err != nil {
		return err
	}
	lopez, err := qtx.CreateHousehold(ctx, store.CreateHouseholdParams{
		Name: "López", Address: "Avenida del Parque 4", Notes: "Grupo 1",
	})
	if err != nil {
		return err
	}
	martin, err := qtx.CreateHousehold(ctx, store.CreateHouseholdParams{
		Name: "Martín", Address: "Plaza Nueva 8", Notes: "Grupo 2",
	})
	if err != nil {
		return err
	}

	gid, lid, mid := garcia.ID, lopez.ID, martin.ID
	pubs := []pubSpec{
		{"Carlos", "García", "male", "publisher", "600111001", &gid, d(1998, 6, 12), d(1996, 3, 1), true, false, false, false, "regular"},
		{"Elena", "García", "female", "publisher", "600111002", &gid, d(2001, 9, 8), d(1999, 5, 1), false, false, false, false, "regular"},
		{"Lucía", "García", "female", "student", "", &gid, pgtype.Date{}, pgtype.Date{}, false, false, false, false, "none"},
		{"Miguel", "López", "male", "publisher", "600222001", &lid, d(2010, 4, 17), d(2008, 1, 1), false, true, false, false, "irregular"},
		{"Ana", "López", "female", "publisher", "600222002", &lid, d(2012, 11, 3), d(2010, 6, 1), false, false, true, false, "pioneer"},
		{"Pablo", "Martín", "male", "publisher", "600333001", &mid, d(2005, 2, 20), d(2004, 9, 1), true, false, false, false, "regular"},
		{"Sofía", "Martín", "female", "unbaptized_publisher", "600333002", &mid, pgtype.Date{}, d(2024, 10, 1), false, false, false, false, "new"},
		{"Javier", "Ruiz", "male", "publisher", "600444001", nil, d(1995, 7, 1), d(1993, 1, 1), false, false, false, false, "inactive"},
		{"Marta", "Vega", "female", "publisher", "600555001", nil, d(2018, 3, 15), d(2016, 9, 1), false, false, false, false, "aux"},
		{"Diego", "Nieto", "male", "publisher", "600666001", nil, d(2008, 8, 22), d(2007, 4, 1), false, true, false, false, "regular"},
	}

	asOf := domain.MonthFromTime(time.Now())
	ids := make([]uuid.UUID, 0, len(pubs))
	for _, p := range pubs {
		row, err := qtx.CreatePublisher(ctx, store.CreatePublisherParams{
			HouseholdID:          p.household,
			FirstName:            p.first,
			LastName:             p.last,
			Gender:               p.gender,
			Phone:                p.phone,
			Email:                "",
			BaptismDate:          p.baptism,
			StartedPreachingDate: p.started,
			SpiritualStatus:      p.spiritual,
			IsElder:              p.elder,
			IsMinisterialServant: p.ms,
			IsRegularPioneer:     p.rp,
			IsSpecialPioneer:     p.sp,
		})
		if err != nil {
			return fmt.Errorf("publisher %s: %w", p.last, err)
		}
		ids = append(ids, row.ID)
		if err := writeReports(ctx, qtx, row.ID, p, asOf); err != nil {
			return err
		}
	}

	if err := writeAttendance(ctx, qtx, asOf); err != nil {
		return err
	}

	for _, id := range ids {
		if err := recomputeActivity(ctx, qtx, id, asOf); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit seed: %w", err)
	}
	return nil
}

func writeReports(ctx context.Context, q *store.Queries, id uuid.UUID, p pubSpec, asOf domain.Month) error {
	if p.pattern == "none" {
		return nil
	}
	months := domain.LastNMonths(asOf, 7)
	for i, m := range months {
		shared, studies, hrs, aux, late := reportFor(p.pattern, i, len(months), p.rp)
		if p.pattern == "inactive" {
			continue
		}
		if p.pattern == "new" && i < 4 {
			continue
		}
		if _, err := q.UpsertReport(ctx, store.UpsertReportParams{
			PublisherID:      id,
			Year:             int16(m.Year),
			Month:            int16(m.Month),
			SharedInMinistry: shared,
			BibleStudies:     studies,
			Hours:            hrs,
			AuxiliaryPioneer: aux,
			Late:             late,
			Remarks:          "",
		}); err != nil {
			return fmt.Errorf("report %s %s: %w", p.last, m.Key(), err)
		}
	}
	return nil
}

func reportFor(pattern string, i, n int, regularPioneer bool) (shared bool, studies int32, hours *float64, aux, late bool) {
	switch pattern {
	case "regular":
		return true, int32(i % 3), nil, false, false
	case "pioneer":
		h := 52.0 + float64(i)
		return true, 2, &h, false, false
	case "irregular":
		shared = i != 2 && i != 5
		return shared, 1, nil, false, i == 4
	case "aux":
		aux = i == n-1
		shared = true
		if aux {
			h := 30.0
			hours = &h
		}
		return shared, 1, hours, aux, false
	case "new":
		return true, 1, nil, false, false
	default:
		if regularPioneer {
			h := 50.0
			return true, 1, &h, false, false
		}
		return true, 0, nil, false, false
	}
}

func writeAttendance(ctx context.Context, q *store.Queries, asOf domain.Month) error {
	start := time.Date(asOf.Year, asOf.Month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
	for week := 0; week < 6; week++ {
		weekend := start.AddDate(0, 0, week*7)
		for weekend.Weekday() != time.Sunday {
			weekend = weekend.AddDate(0, 0, 1)
		}
		midweek := weekend.AddDate(0, 0, 4)
		online := int32(3 + week%2)
		if _, err := q.UpsertAttendance(ctx, store.UpsertAttendanceParams{
			MeetingDate: weekend,
			Kind:        domain.MeetingWeekend,
			InPerson:    int32(42 + week),
			Online:      &online,
		}); err != nil {
			return err
		}
		if _, err := q.UpsertAttendance(ctx, store.UpsertAttendanceParams{
			MeetingDate: midweek,
			Kind:        domain.MeetingMidweek,
			InPerson:    int32(35 + week),
			Online:      &online,
		}); err != nil {
			return err
		}
	}
	return nil
}

func recomputeActivity(ctx context.Context, q *store.Queries, publisherID uuid.UUID, asOf domain.Month) error {
	window := domain.LastNMonths(asOf, 6)
	from, to := window[0], window[len(window)-1]
	pub, err := q.GetPublisherForReport(ctx, publisherID)
	if err != nil {
		return err
	}
	rows, err := q.ListSharesForPublisher(ctx, store.ListSharesForPublisherParams{
		PublisherID: publisherID,
		FromYear:    int16(from.Year),
		FromMonth:   int16(from.Month),
		ToYear:      int16(to.Year),
		ToMonth:     int16(to.Month),
	})
	if err != nil {
		return err
	}
	shares := make([]domain.MonthShare, 0, len(rows))
	for _, r := range rows {
		shares = append(shares, domain.MonthShare{
			Year: int(r.Year), Month: time.Month(r.Month), Shared: r.SharedInMinistry,
		})
	}
	var started *time.Time
	if pub.StartedPreachingDate.Valid {
		t := pub.StartedPreachingDate.Time
		started = &t
	}
	return q.UpdatePublisherActivity(ctx, store.UpdatePublisherActivityParams{
		ID:             publisherID,
		ActivityStatus: service.ActivityStatus(shares, started, asOf),
	})
}

func d(year, month, day int) pgtype.Date {
	return pgtype.Date{Time: time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), Valid: true}
}

type pubSpec struct {
	first, last, gender, spiritual, phone string
	household                             *uuid.UUID
	baptism, started                      pgtype.Date
	elder, ms, rp, sp                     bool
	pattern                               string
}
