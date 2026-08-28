package http

import (
	"context"
	"net/http"
	"time"

	"khhub/internal/domain"
	"khhub/internal/store"

	"github.com/gin-gonic/gin"
)

// attendanceWriter is the store surface used by PUT /attendance.
// *store.Queries implements it.
type attendanceWriter interface {
	UpsertAttendance(ctx context.Context, arg store.UpsertAttendanceParams) (store.MeetingAttendance, error)
}

type attendanceRequest struct {
	Date     string `json:"date" binding:"required"`
	Kind     string `json:"kind" binding:"required,oneof=midweek weekend"`
	InPerson int32  `json:"inPerson" binding:"min=0"`
	Online   *int32 `json:"online"`
}

func listAttendance(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		fromStr := c.DefaultQuery("from", time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02"))
		toStr := c.DefaultQuery("to", now.Format("2006-01-02"))
		from, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			jsonError(c, http.StatusBadRequest, "fecha inicial no válida")
			return
		}
		to, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			jsonError(c, http.StatusBadRequest, "fecha final no válida")
			return
		}
		rows, err := q.ListAttendance(c.Request.Context(), store.ListAttendanceParams{
			FromDate: from,
			ToDate:   to,
		})
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo cargar la asistencia")
			return
		}
		out := make([]gin.H, 0, len(rows))
		for _, r := range rows {
			out = append(out, gin.H{
				"id":        r.ID,
				"date":      r.MeetingDate.Format("2006-01-02"),
				"kind":      r.Kind,
				"inPerson":  r.InPerson,
				"online":    r.Online,
				"createdAt": r.CreatedAt,
			})
		}
		c.JSON(http.StatusOK, out)
	}
}

func putAttendance(q attendanceWriter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req attendanceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			jsonError(c, http.StatusBadRequest, "datos de asistencia no válidos")
			return
		}
		if req.Kind != domain.MeetingMidweek && req.Kind != domain.MeetingWeekend {
			jsonError(c, http.StatusBadRequest, "tipo de reunión no válido")
			return
		}
		day, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			jsonError(c, http.StatusBadRequest, "fecha no válida")
			return
		}
		if req.Online != nil && *req.Online < 0 {
			jsonError(c, http.StatusBadRequest, "la asistencia en línea no puede ser negativa")
			return
		}
		row, err := q.UpsertAttendance(c.Request.Context(), store.UpsertAttendanceParams{
			MeetingDate: day,
			Kind:        req.Kind,
			InPerson:    req.InPerson,
			Online:      req.Online,
		})
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo guardar la asistencia")
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":       row.ID,
			"date":     row.MeetingDate.Format("2006-01-02"),
			"kind":     row.Kind,
			"inPerson": row.InPerson,
			"online":   row.Online,
		})
	}
}

func deleteAttendance(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseUUIDParam(c, "id")
		if !ok {
			return
		}
		if err := q.DeleteAttendance(c.Request.Context(), id); err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo eliminar la asistencia")
			return
		}
		c.Status(http.StatusNoContent)
	}
}
