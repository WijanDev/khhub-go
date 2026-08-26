package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"khhub/internal/domain"
	"khhub/internal/service"
	"khhub/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type reportUpsert struct {
	PublisherID       uuid.UUID `json:"publisherId" binding:"required"`
	SharedInMinistry  bool      `json:"sharedInMinistry"`
	BibleStudies      int32     `json:"bibleStudies" binding:"min=0"`
	Hours             *float64  `json:"hours"`
	AuxiliaryPioneer  bool      `json:"auxiliaryPioneer"`
	Late              bool      `json:"late"`
	Remarks           string    `json:"remarks" binding:"max=500"`
}

type reportsBatchRequest struct {
	Year    int            `json:"year" binding:"required,min=2000,max=2100"`
	Month   int            `json:"month" binding:"required,min=1,max=12"`
	Reports []reportUpsert `json:"reports" binding:"required,dive"`
}

func parseYearMonth(c *gin.Context) (int, int, bool) {
	now := time.Now()
	year, err := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))
	if err != nil || year < 2000 || year > 2100 {
		jsonError(c, http.StatusBadRequest, "año no válido")
		return 0, 0, false
	}
	month, err := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	if err != nil || month < 1 || month > 12 {
		jsonError(c, http.StatusBadRequest, "mes no válido")
		return 0, 0, false
	}
	return year, month, true
}

func listReports(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		year, month, ok := parseYearMonth(c)
		if !ok {
			return
		}
		rows, err := q.ListReportsForMonth(c.Request.Context(), store.ListReportsForMonthParams{
			Year:  int16(year),
			Month: int16(month),
		})
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudieron cargar los informes")
			return
		}
		out := make([]gin.H, 0, len(rows))
		missing := make([]gin.H, 0)
		for _, r := range rows {
			hasReport := r.ID != nil
			hourReporter := r.IsRegularPioneer || r.IsSpecialPioneer || (r.AuxiliaryPioneer != nil && *r.AuxiliaryPioneer)
			item := gin.H{
				"publisherId":        r.PublisherID,
				"firstName":          r.FirstName,
				"lastName":           r.LastName,
				"spiritualStatus":    r.SpiritualStatus,
				"isRegularPioneer":   r.IsRegularPioneer,
				"isSpecialPioneer":   r.IsSpecialPioneer,
				"hourReporter":       hourReporter,
				"hasReport":          hasReport,
				"id":                 r.ID,
				"sharedInMinistry":   boolOr(r.SharedInMinistry),
				"bibleStudies":       intOr(r.BibleStudies),
				"hours":              r.Hours,
				"auxiliaryPioneer":   boolOr(r.AuxiliaryPioneer),
				"late":               boolOr(r.Late),
				"remarks":            stringOr(r.Remarks),
			}
			out = append(out, item)
			if !hasReport {
				missing = append(missing, gin.H{
					"publisherId": r.PublisherID,
					"firstName":   r.FirstName,
					"lastName":    r.LastName,
				})
			}
		}
		asOf := domain.Month{Year: year, Month: time.Month(month)}
		c.JSON(http.StatusOK, gin.H{
			"year":        year,
			"month":       month,
			"serviceYear": domain.ServiceYear(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)),
			"reports":     out,
			"missing":     missing,
			"asOf":        asOf.Key(),
		})
	}
}

func putReports(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req reportsBatchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			jsonError(c, http.StatusBadRequest, "datos de informes no válidos")
			return
		}
		ctx := c.Request.Context()
		saved := make([]store.FieldServiceReport, 0, len(req.Reports))
		seen := make(map[uuid.UUID]struct{})
		for _, r := range req.Reports {
			if _, dup := seen[r.PublisherID]; dup {
				jsonError(c, http.StatusBadRequest, "publicador duplicado en el lote")
				return
			}
			seen[r.PublisherID] = struct{}{}
			pub, err := q.GetPublisherForReport(ctx, r.PublisherID)
			if err != nil {
				if isNoRows(err) {
					jsonError(c, http.StatusBadRequest, "publicador no encontrado")
					return
				}
				jsonError(c, http.StatusInternalServerError, "no se pudieron guardar los informes")
				return
			}
			if !domain.MustReport(pub.SpiritualStatus) {
				jsonError(c, http.StatusBadRequest, "los estudiantes no envían informe de predicación")
				return
			}
			if err := domain.ReportsHours(pub.IsRegularPioneer, pub.IsSpecialPioneer, r.AuxiliaryPioneer, r.SharedInMinistry, r.Hours); err != nil {
				jsonError(c, http.StatusBadRequest, err.Error())
				return
			}
			row, err := q.UpsertReport(ctx, store.UpsertReportParams{
				PublisherID:      r.PublisherID,
				Year:             int16(req.Year),
				Month:            int16(req.Month),
				SharedInMinistry: r.SharedInMinistry,
				BibleStudies:     r.BibleStudies,
				Hours:            r.Hours,
				AuxiliaryPioneer: r.AuxiliaryPioneer,
				Late:             r.Late,
				Remarks:          trim(r.Remarks),
			})
			if err != nil {
				jsonError(c, http.StatusInternalServerError, "no se pudieron guardar los informes")
				return
			}
			saved = append(saved, row)
			if err := recomputeActivity(c.Request.Context(), q, r.PublisherID, domain.Month{Year: req.Year, Month: time.Month(req.Month)}); err != nil {
				jsonError(c, http.StatusInternalServerError, "no se pudo actualizar el estado de actividad")
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"saved": len(saved)})
	}
}

func recomputeActivity(ctx context.Context, q *store.Queries, publisherID uuid.UUID, asOf domain.Month) error {
	window := domain.LastNMonths(asOf, 6)
	from := window[0]
	to := window[len(window)-1]
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
			Year:   int(r.Year),
			Month:  time.Month(r.Month),
			Shared: r.SharedInMinistry,
		})
	}
	status := service.ActivityStatus(shares, ptrFromDate(pub.StartedPreachingDate), asOf)
	return q.UpdatePublisherActivity(ctx, store.UpdatePublisherActivityParams{
		ID:             publisherID,
		ActivityStatus: status,
	})
}

func boolOr(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

func intOr(v *int32) int32 {
	if v == nil {
		return 0
	}
	return *v
}

func stringOr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
