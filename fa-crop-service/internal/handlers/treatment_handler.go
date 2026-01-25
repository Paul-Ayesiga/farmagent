package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/farmagent/fa-crop-service/internal/models"
	"github.com/farmagent/fa-crop-service/internal/repository"
	"github.com/farmagent/fa-crop-service/pkg/jsonapi"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TreatmentHandler struct {
	treatmentRepo repository.TreatmentRepository
	cropRepo      repository.CropRepository
	fieldRepo     repository.FieldRepository
}

func NewTreatmentHandler(treatmentRepo repository.TreatmentRepository, cropRepo repository.CropRepository, fieldRepo repository.FieldRepository) *TreatmentHandler {
	return &TreatmentHandler{treatmentRepo: treatmentRepo, cropRepo: cropRepo, fieldRepo: fieldRepo}
}

type CreateTreatmentInput struct {
	CropID          uuid.UUID  `json:"crop_id"`
	HealthRecordID  *uuid.UUID `json:"health_record_id,omitempty"`
	DiseaseName     string     `json:"disease_name"`
	TreatmentName   string     `json:"treatment_name"`
	TreatmentType   string     `json:"treatment_type"`
	ApplicationDate string     `json:"application_date"`
	Cost            *float64   `json:"cost,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
}

type UpdateTreatmentInput struct {
	Effectiveness *int    `json:"effectiveness,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}

// treatmentToResource converts a Treatment to JSON:API resource
func treatmentToResource(treatment *models.Treatment) jsonapi.Resource {
	attrs := map[string]interface{}{
		"disease_name":     treatment.DiseaseName,
		"treatment_name":   treatment.TreatmentName,
		"treatment_type":   treatment.TreatmentType,
		"application_date": treatment.ApplicationDate.Format("2006-01-02"),
		"created_at":       treatment.CreatedAt,
	}
	if treatment.Cost != nil {
		attrs["cost"] = *treatment.Cost
	}
	if treatment.Effectiveness != nil {
		attrs["effectiveness"] = *treatment.Effectiveness
	}
	if treatment.Notes != nil {
		attrs["notes"] = *treatment.Notes
	}

	relationships := map[string]jsonapi.Relationship{
		"crop": {
			Data: jsonapi.NewResourceIdentifier("crops", treatment.CropID.String()),
		},
	}
	if treatment.HealthRecordID != nil {
		relationships["health-record"] = jsonapi.Relationship{
			Data: jsonapi.NewResourceIdentifier("health-records", treatment.HealthRecordID.String()),
		}
	}

	return jsonapi.Resource{
		ID:            treatment.ID.String(),
		Type:          "treatments",
		Attributes:    attrs,
		Relationships: relationships,
		Links: &jsonapi.Links{
			Self: "/treatments/" + treatment.ID.String(),
		},
	}
}

// ListTreatments returns treatments for a crop
func (h *TreatmentHandler) ListTreatments(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	cropID, err := uuid.Parse(chi.URLParam(r, "cropId"))
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Bad Request", "invalid crop ID")
		return
	}

	crop, err := h.cropRepo.FindByID(cropID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusNotFound, "Not Found", "crop not found")
		return
	}
	field, err := h.fieldRepo.FindByID(crop.FieldID)
	if err != nil || field.UserID != userID {
		jsonapi.RespondWithError(w, http.StatusForbidden, "Forbidden", "access denied")
		return
	}

	treatments, err := h.treatmentRepo.FindByCropID(cropID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to fetch treatments")
		return
	}

	resources := make([]jsonapi.Resource, len(treatments))
	for i, t := range treatments {
		resources[i] = treatmentToResource(&t)
	}

	jsonapi.RespondWithResources(w, http.StatusOK, resources)
}

// GetTreatment returns a specific treatment
func (h *TreatmentHandler) GetTreatment(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	treatmentID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Bad Request", "invalid treatment ID")
		return
	}

	treatment, err := h.treatmentRepo.FindByID(treatmentID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusNotFound, "Not Found", "treatment not found")
		return
	}

	crop, _ := h.cropRepo.FindByID(treatment.CropID)
	field, _ := h.fieldRepo.FindByID(crop.FieldID)
	if field.UserID != userID {
		jsonapi.RespondWithError(w, http.StatusForbidden, "Forbidden", "access denied")
		return
	}

	jsonapi.RespondWithResource(w, http.StatusOK, treatmentToResource(treatment))
}

// CreateTreatment creates a new treatment
func (h *TreatmentHandler) CreateTreatment(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	var input CreateTreatmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Bad Request", "invalid request body")
		return
	}

	crop, err := h.cropRepo.FindByID(input.CropID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusNotFound, "Not Found", "crop not found")
		return
	}
	field, err := h.fieldRepo.FindByID(crop.FieldID)
	if err != nil || field.UserID != userID {
		jsonapi.RespondWithError(w, http.StatusForbidden, "Forbidden", "access denied")
		return
	}

	appDate, err := time.Parse("2006-01-02", input.ApplicationDate)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Validation Error", "invalid application_date format")
		return
	}

	treatment := &models.Treatment{
		CropID:          input.CropID,
		HealthRecordID:  input.HealthRecordID,
		DiseaseName:     input.DiseaseName,
		TreatmentName:   input.TreatmentName,
		TreatmentType:   input.TreatmentType,
		ApplicationDate: appDate,
		Cost:            input.Cost,
		Notes:           input.Notes,
	}

	if err := h.treatmentRepo.Create(treatment); err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to create treatment")
		return
	}

	jsonapi.RespondWithResource(w, http.StatusCreated, treatmentToResource(treatment))
}

// UpdateTreatment updates a treatment (mainly for effectiveness rating)
func (h *TreatmentHandler) UpdateTreatment(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	treatmentID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Bad Request", "invalid treatment ID")
		return
	}

	treatment, err := h.treatmentRepo.FindByID(treatmentID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusNotFound, "Not Found", "treatment not found")
		return
	}

	crop, _ := h.cropRepo.FindByID(treatment.CropID)
	field, _ := h.fieldRepo.FindByID(crop.FieldID)
	if field.UserID != userID {
		jsonapi.RespondWithError(w, http.StatusForbidden, "Forbidden", "access denied")
		return
	}

	var input UpdateTreatmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Bad Request", "invalid request body")
		return
	}

	if input.Effectiveness != nil {
		treatment.Effectiveness = input.Effectiveness
	}
	if input.Notes != nil {
		treatment.Notes = input.Notes
	}

	if err := h.treatmentRepo.Update(treatment); err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to update treatment")
		return
	}

	jsonapi.RespondWithResource(w, http.StatusOK, treatmentToResource(treatment))
}

// DeleteTreatment deletes a treatment
func (h *TreatmentHandler) DeleteTreatment(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	treatmentID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Bad Request", "invalid treatment ID")
		return
	}

	treatment, err := h.treatmentRepo.FindByID(treatmentID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusNotFound, "Not Found", "treatment not found")
		return
	}

	crop, _ := h.cropRepo.FindByID(treatment.CropID)
	field, _ := h.fieldRepo.FindByID(crop.FieldID)
	if field.UserID != userID {
		jsonapi.RespondWithError(w, http.StatusForbidden, "Forbidden", "access denied")
		return
	}

	if err := h.treatmentRepo.Delete(treatmentID); err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to delete treatment")
		return
	}

	jsonapi.RespondWithMeta(w, http.StatusOK, jsonapi.Meta{"message": "treatment deleted"})
}
