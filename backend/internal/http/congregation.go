package http

import (
	"context"
	"net/http"

	"khhub/internal/store"

	"github.com/gin-gonic/gin"
)

type congregationRequest struct {
	Name       string `json:"name" binding:"required,max=200"`
	Number     string `json:"number" binding:"max=50"`
	MidweekDay int16  `json:"midweekDay" binding:"min=0,max=6"`
	WeekendDay int16  `json:"weekendDay" binding:"min=0,max=6"`
}

func congregationJSON(row store.Congregation, seedResetEnabled bool) gin.H {
	return gin.H{
		"name":              row.Name,
		"number":            row.Number,
		"midweekDay":        row.MidweekDay,
		"weekendDay":        row.WeekendDay,
		"updatedAt":         row.UpdatedAt,
		"seedResetEnabled":  seedResetEnabled,
	}
}

func getCongregation(q *store.Queries, seedResetEnabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		row, err := q.GetCongregation(c.Request.Context())
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo cargar la congregación")
			return
		}
		c.JSON(http.StatusOK, congregationJSON(row, seedResetEnabled))
	}
}

func putCongregation(q *store.Queries, seedResetEnabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req congregationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			jsonError(c, http.StatusBadRequest, "datos de congregación no válidos")
			return
		}
		row, err := q.UpdateCongregation(c.Request.Context(), store.UpdateCongregationParams{
			Name:       trim(req.Name),
			Number:     trim(req.Number),
			MidweekDay: req.MidweekDay,
			WeekendDay: req.WeekendDay,
		})
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo guardar la congregación")
			return
		}
		c.JSON(http.StatusOK, congregationJSON(row, seedResetEnabled))
	}
}

func postResetSeed(reset func(ctx context.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if reset == nil {
			jsonError(c, http.StatusForbidden, "el restablecimiento de datos solo está disponible en desarrollo")
			return
		}
		if err := reset(c.Request.Context()); err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudieron restablecer los datos de demostración")
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}
