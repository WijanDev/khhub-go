package http

import (
	"net/http"
	"time"

	"khhub/internal/domain"
	"khhub/internal/store"

	"github.com/gin-gonic/gin"
)

func getDashboard(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		year, month, ok := parseYearMonth(c)
		if !ok {
			return
		}
		ctx := c.Request.Context()

		pubs, err := q.ListPublishers(ctx)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo cargar el panel")
			return
		}
		reports, err := q.ListReportsForMonth(ctx, store.ListReportsForMonthParams{
			Year:  int16(year),
			Month: int16(month),
		})
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo cargar el panel")
			return
		}

		start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, -1)
		attendance, err := q.ListAttendance(ctx, store.ListAttendanceParams{FromDate: start, ToDate: end})
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo cargar el panel")
			return
		}

		cong, err := q.GetCongregation(ctx)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo cargar el panel")
			return
		}

		var regular, irregular, inactive, students int
		for _, p := range pubs {
			switch {
			case p.SpiritualStatus == domain.SpiritualStudent:
				students++
			case p.ActivityStatus == domain.ActivityRegular:
				regular++
			case p.ActivityStatus == domain.ActivityIrregular:
				irregular++
			default:
				inactive++
			}
		}

		var reported, shared, studies int
		var hours float64
		missing := make([]gin.H, 0)
		for _, r := range reports {
			if r.ID == nil {
				missing = append(missing, gin.H{
					"publisherId": r.PublisherID,
					"firstName":   r.FirstName,
					"lastName":    r.LastName,
				})
				continue
			}
			reported++
			if boolOr(r.SharedInMinistry) {
				shared++
			}
			studies += int(intOr(r.BibleStudies))
			if r.Hours != nil {
				hours += *r.Hours
			}
		}

		shouldReport := len(reports)
		participation := 0.0
		if shouldReport > 0 {
			participation = float64(shared) / float64(shouldReport) * 100
		}

		var midweekSum, weekendSum, midweekN, weekendN int
		var midweekOnline, weekendOnline int
		for _, a := range attendance {
			if a.Kind == domain.MeetingMidweek {
				midweekSum += int(a.InPerson)
				midweekN++
				if a.Online != nil {
					midweekOnline += int(*a.Online)
				}
			} else {
				weekendSum += int(a.InPerson)
				weekendN++
				if a.Online != nil {
					weekendOnline += int(*a.Online)
				}
			}
		}

		avg := func(sum, n int) *float64 {
			if n == 0 {
				return nil
			}
			v := float64(sum) / float64(n)
			return &v
		}

		c.JSON(http.StatusOK, gin.H{
			"congregation": gin.H{"name": cong.Name, "number": cong.Number},
			"year":         year,
			"month":        month,
			"serviceYear":  domain.ServiceYear(start),
			"publishers": gin.H{
				"total":      len(pubs),
				"students":   students,
				"active":     regular + irregular,
				"regular":    regular,
				"irregular":  irregular,
				"inactive":   inactive,
			},
			"reports": gin.H{
				"shouldReport":   shouldReport,
				"reported":       reported,
				"missingCount":   len(missing),
				"missing":        missing,
				"shared":         shared,
				"participation":  participation,
				"bibleStudies":   studies,
				"pioneerHours":   hours,
			},
			"attendance": gin.H{
				"midweekMeetings":   midweekN,
				"weekendMeetings":   weekendN,
				"midweekAvg":        avg(midweekSum, midweekN),
				"weekendAvg":        avg(weekendSum, weekendN),
				"midweekOnlineAvg":  avg(midweekOnline, midweekN),
				"weekendOnlineAvg":  avg(weekendOnline, weekendN),
			},
		})
	}
}
