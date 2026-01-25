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

type HealthRecordHandler struct {
	healthRepo repository.HealthRecordRepository
	cropRepo   repository.CropRepository
	fieldRepo  repository.FieldRepository
}

func NewHealthRecordHandler(healthRepo repository.HealthRecordRepository, cropRepo repository.CropRepository, fieldRepo repository.FieldRepository) *HealthRecordHandler {
	return &HealthRecordHandler{healthRepo: healthRepo, cropRepo: cropRepo, fieldRepo: fieldRepo}
}

type CreateHealthRecordInput struct {
	CropID          uuid.UUID `json:"crop_id"`
	HealthScore     int       `json:"health_score"`
	ImageURL        string    `json:"image_url"`
	DiseaseDetected *string   `json:"disease_detected,omitempty"`
	Confidence      *float64  `json:"confidence,omitempty"`
	Severity        *string   `json:"severity,omitempty"`
	Notes           *string   `json:"notes,omitempty"`
}

// healthRecordToResource converts a HealthRecord to JSON:API resource
func healthRecordToResource(record *models.HealthRecord) jsonapi.Resource {
	attrs := map[string]interface{}{
		"check_date":   record.CheckDate,
		"health_score": record.HealthScore,
		"image_url":    record.ImageURL,
		"created_at":   record.CreatedAt,
	}
	if record.DiseaseDetected != nil {
		attrs["disease_detected"] = *record.DiseaseDetected
	}
	if record.Confidence != nil {
		attrs["confidence"] = *record.Confidence
	}
	if record.Severity != nil {
		attrs["severity"] = *record.Severity
	}
	if record.Notes != nil {
		attrs["notes"] = *record.Notes
	}

	return jsonapi.Resource{
		ID:         record.ID.String(),
		Type:       "health-records",
		Attributes: attrs,
		Relationships: map[string]jsonapi.Relationship{
			"crop": {
				Data: jsonapi.NewResourceIdentifier("crops", record.CropID.String()),
			},
		},
		Links: &jsonapi.Links{
			Self: "/health-records/" + record.ID.String(),
		},
	}
}

// ListHealthRecords returns health records for a crop
func (h *HealthRecordHandler) ListHealthRecords(w http.ResponseWriter, r *http.Request) {
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

	records, err := h.healthRepo.FindByCropID(cropID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to fetch health records")
		return
	}

	resources := make([]jsonapi.Resource, len(records))
	for i, rec := range records {
		resources[i] = healthRecordToResource(&rec)
	}

	jsonapi.RespondWithResources(w, http.StatusOK, resources)
}

// GetHealthRecord returns a specific health record
func (h *HealthRecordHandler) GetHealthRecord(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	recordID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Bad Request", "invalid record ID")
		return
	}

	record, err := h.healthRepo.FindByID(recordID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusNotFound, "Not Found", "record not found")
		return
	}

	crop, _ := h.cropRepo.FindByID(record.CropID)
	field, _ := h.fieldRepo.FindByID(crop.FieldID)
	if field.UserID != userID {
		jsonapi.RespondWithError(w, http.StatusForbidden, "Forbidden", "access denied")
		return
	}

	jsonapi.RespondWithResource(w, http.StatusOK, healthRecordToResource(record))
}

// CreateHealthRecord creates a new health record
func (h *HealthRecordHandler) CreateHealthRecord(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	var input CreateHealthRecordInput
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

	record := &models.HealthRecord{
		CropID:          input.CropID,
		HealthScore:     input.HealthScore,
		ImageURL:        input.ImageURL,
		DiseaseDetected: input.DiseaseDetected,
		Confidence:      input.Confidence,
		Severity:        input.Severity,
		Notes:           input.Notes,
		CheckDate:       time.Now(),
	}

	if err := h.healthRepo.Create(record); err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to create health record")
		return
	}

	jsonapi.RespondWithResource(w, http.StatusCreated, healthRecordToResource(record))
}

// DeleteHealthRecord deletes a health record
func (h *HealthRecordHandler) DeleteHealthRecord(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	recordID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Bad Request", "invalid record ID")
		return
	}

	record, err := h.healthRepo.FindByID(recordID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusNotFound, "Not Found", "record not found")
		return
	}

	crop, _ := h.cropRepo.FindByID(record.CropID)
	field, _ := h.fieldRepo.FindByID(crop.FieldID)
	if field.UserID != userID {
		jsonapi.RespondWithError(w, http.StatusForbidden, "Forbidden", "access denied")
		return
	}

	if err := h.healthRepo.Delete(recordID); err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to delete record")
		return
	}

	jsonapi.RespondWithMeta(w, http.StatusOK, jsonapi.Meta{"message": "record deleted"})
}
