package http

import (
	"net/http"
	"time"

	"khhub/internal/domain"
	"khhub/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type publisherRequest struct {
	HouseholdID          *uuid.UUID `json:"householdId"`
	FirstName            string     `json:"firstName" binding:"required,max=100"`
	LastName             string     `json:"lastName" binding:"required,max=100"`
	Gender               string     `json:"gender" binding:"required,oneof=male female"`
	Phone                string     `json:"phone" binding:"max=50"`
	Email                string     `json:"email" binding:"omitempty,email,max=200"`
	BaptismDate          *string    `json:"baptismDate"`
	StartedPreachingDate *string    `json:"startedPreachingDate"`
	SpiritualStatus      string     `json:"spiritualStatus" binding:"required,oneof=student unbaptized_publisher publisher"`
	IsElder              bool       `json:"isElder"`
	IsMinisterialServant bool       `json:"isMinisterialServant"`
	IsRegularPioneer     bool       `json:"isRegularPioneer"`
	IsSpecialPioneer     bool       `json:"isSpecialPioneer"`
}

func parseDatePtr(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func datePtrJSON(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.Format("2006-01-02")
}

func publisherJSON(
	id uuid.UUID,
	householdID *uuid.UUID,
	householdName *string,
	firstName, lastName, gender, phone, email string,
	baptism, started *time.Time,
	spiritual string,
	elder, ms, rp, sp bool,
	activity string,
	created time.Time,
) gin.H {
	return gin.H{
		"id":                   id,
		"householdId":          householdID,
		"householdName":        householdName,
		"firstName":            firstName,
		"lastName":             lastName,
		"gender":               gender,
		"phone":                phone,
		"email":                email,
		"baptismDate":          datePtrJSON(baptism),
		"startedPreachingDate": datePtrJSON(started),
		"spiritualStatus":      spiritual,
		"isElder":              elder,
		"isMinisterialServant": ms,
		"isRegularPioneer":     rp,
		"isSpecialPioneer":     sp,
		"activityStatus":       activity,
		"isActive":             activity == domain.ActivityRegular || activity == domain.ActivityIrregular,
		"createdAt":            created,
	}
}

func listPublishers(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := q.ListPublishers(c.Request.Context())
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudieron cargar los publicadores")
			return
		}
		out := make([]gin.H, 0, len(rows))
		for _, r := range rows {
			out = append(out, publisherJSON(
				r.ID, r.HouseholdID, r.HouseholdName, r.FirstName, r.LastName, r.Gender, r.Phone, r.Email,
				ptrFromDate(r.BaptismDate), ptrFromDate(r.StartedPreachingDate), r.SpiritualStatus,
				r.IsElder, r.IsMinisterialServant, r.IsRegularPioneer, r.IsSpecialPioneer,
				r.ActivityStatus, r.CreatedAt,
			))
		}
		c.JSON(http.StatusOK, out)
	}
}

func getPublisher(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseUUIDParam(c, "id")
		if !ok {
			return
		}
		r, err := q.GetPublisher(c.Request.Context(), id)
		if err != nil {
			if isNoRows(err) {
				jsonError(c, http.StatusNotFound, "publicador no encontrado")
				return
			}
			jsonError(c, http.StatusInternalServerError, "no se pudo cargar el publicador")
			return
		}
		c.JSON(http.StatusOK, publisherJSON(
			r.ID, r.HouseholdID, r.HouseholdName, r.FirstName, r.LastName, r.Gender, r.Phone, r.Email,
			ptrFromDate(r.BaptismDate), ptrFromDate(r.StartedPreachingDate), r.SpiritualStatus,
			r.IsElder, r.IsMinisterialServant, r.IsRegularPioneer, r.IsSpecialPioneer,
			r.ActivityStatus, r.CreatedAt,
		))
	}
}

func createPublisher(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, baptism, started, ok := bindPublisher(c)
		if !ok {
			return
		}
		r, err := q.CreatePublisher(c.Request.Context(), store.CreatePublisherParams{
			HouseholdID:          req.HouseholdID,
			FirstName:            trim(req.FirstName),
			LastName:             trim(req.LastName),
			Gender:               req.Gender,
			Phone:                trim(req.Phone),
			Email:                trim(req.Email),
			BaptismDate:          dateFromPtr(baptism),
			StartedPreachingDate: dateFromPtr(started),
			SpiritualStatus:      req.SpiritualStatus,
			IsElder:              req.IsElder,
			IsMinisterialServant: req.IsMinisterialServant,
			IsRegularPioneer:     req.IsRegularPioneer,
			IsSpecialPioneer:     req.IsSpecialPioneer,
		})
		if err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo crear el publicador")
			return
		}
		c.JSON(http.StatusCreated, publisherJSON(
			r.ID, r.HouseholdID, nil, r.FirstName, r.LastName, r.Gender, r.Phone, r.Email,
			ptrFromDate(r.BaptismDate), ptrFromDate(r.StartedPreachingDate), r.SpiritualStatus,
			r.IsElder, r.IsMinisterialServant, r.IsRegularPioneer, r.IsSpecialPioneer,
			r.ActivityStatus, r.CreatedAt,
		))
	}
}

func updatePublisher(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseUUIDParam(c, "id")
		if !ok {
			return
		}
		req, baptism, started, ok := bindPublisher(c)
		if !ok {
			return
		}
		r, err := q.UpdatePublisher(c.Request.Context(), store.UpdatePublisherParams{
			ID:                   id,
			HouseholdID:          req.HouseholdID,
			FirstName:            trim(req.FirstName),
			LastName:             trim(req.LastName),
			Gender:               req.Gender,
			Phone:                trim(req.Phone),
			Email:                trim(req.Email),
			BaptismDate:          dateFromPtr(baptism),
			StartedPreachingDate: dateFromPtr(started),
			SpiritualStatus:      req.SpiritualStatus,
			IsElder:              req.IsElder,
			IsMinisterialServant: req.IsMinisterialServant,
			IsRegularPioneer:     req.IsRegularPioneer,
			IsSpecialPioneer:     req.IsSpecialPioneer,
		})
		if err != nil {
			if isNoRows(err) {
				jsonError(c, http.StatusNotFound, "publicador no encontrado")
				return
			}
			jsonError(c, http.StatusInternalServerError, "no se pudo guardar el publicador")
			return
		}
		c.JSON(http.StatusOK, publisherJSON(
			r.ID, r.HouseholdID, nil, r.FirstName, r.LastName, r.Gender, r.Phone, r.Email,
			ptrFromDate(r.BaptismDate), ptrFromDate(r.StartedPreachingDate), r.SpiritualStatus,
			r.IsElder, r.IsMinisterialServant, r.IsRegularPioneer, r.IsSpecialPioneer,
			r.ActivityStatus, r.CreatedAt,
		))
	}
}

func deletePublisher(q *store.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseUUIDParam(c, "id")
		if !ok {
			return
		}
		if err := q.DeletePublisher(c.Request.Context(), id); err != nil {
			jsonError(c, http.StatusInternalServerError, "no se pudo eliminar el publicador")
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func bindPublisher(c *gin.Context) (publisherRequest, *time.Time, *time.Time, bool) {
	var req publisherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, "datos de publicador no válidos")
		return req, nil, nil, false
	}
	if req.IsElder && req.IsMinisterialServant {
		jsonError(c, http.StatusBadRequest, "un hermano no puede ser anciano y siervo ministerial a la vez")
		return req, nil, nil, false
	}
	baptism, err := parseDatePtr(req.BaptismDate)
	if err != nil {
		jsonError(c, http.StatusBadRequest, "fecha de bautismo no válida")
		return req, nil, nil, false
	}
	started, err := parseDatePtr(req.StartedPreachingDate)
	if err != nil {
		jsonError(c, http.StatusBadRequest, "fecha de inicio en la predicación no válida")
		return req, nil, nil, false
	}
	return req, baptism, started, true
}
