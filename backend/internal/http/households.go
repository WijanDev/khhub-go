package http

import (
	"net/http"

	"khhub/internal/store"

	"github.com/gin-gonic/gin"
)

type householdRequest struct {
	Name    string `json:"name" binding:"required,max=200"`
	Address string `json:"address" binding:"max=400"`
	Notes   string `json:"notes" binding:"max=1000"`
}

func householdJSON(row store.Household) gin.H {
	return gin.H{
		"id":        row.ID,
		"name":      row.Name,
		"address":   row.Address,
		"notes":     row.Notes,
		"createdAt": row.CreatedAt,
	}
}

func listHouseholds(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := q.ListHouseholds(c.Request.Context())
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudieron cargar las familias")
			return
		}
		out := make([]gin.H, 0, len(rows))
		for _, r := range rows {
			out = append(out, householdJSON(r))
		}
		c.JSON(http.StatusOK, out)
	}
}

func createHousehold(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req householdRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			jsonError(c, http.StatusBadRequest, "datos de familia no válidos")
			return
		}
		row, err := q.CreateHousehold(c.Request.Context(), store.CreateHouseholdParams{
			Name:    trim(req.Name),
			Address: trim(req.Address),
			Notes:   trim(req.Notes),
		})
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo crear la familia")
			return
		}
		c.JSON(http.StatusCreated, householdJSON(row))
	}
}

func updateHousehold(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseUUIDParam(c, "id")
		if !ok {
			return
		}
		var req householdRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			jsonError(c, http.StatusBadRequest, "datos de familia no válidos")
			return
		}
		row, err := q.UpdateHousehold(c.Request.Context(), store.UpdateHouseholdParams{
			ID:      id,
			Name:    trim(req.Name),
			Address: trim(req.Address),
			Notes:   trim(req.Notes),
		})
		if err != nil {
			if isNoRows(err) {
				jsonError(c, http.StatusNotFound, "familia no encontrada")
				return
			}
			jsonError(c, http.StatusInternalServerError, "no se pudo guardar la familia")
			return
		}
		c.JSON(http.StatusOK, householdJSON(row))
	}
}

func deleteHousehold(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseUUIDParam(c, "id")
		if !ok {
			return
		}
		if err := q.DeleteHousehold(c.Request.Context(), id); err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo eliminar la familia")
			return
		}
		c.Status(http.StatusNoContent)
	}
}
