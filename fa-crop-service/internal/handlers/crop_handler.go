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

type CropHandler struct {
	cropRepo  repository.CropRepository
	fieldRepo repository.FieldRepository
}

func NewCropHandler(cropRepo repository.CropRepository, fieldRepo repository.FieldRepository) *CropHandler {
	return &CropHandler{cropRepo: cropRepo, fieldRepo: fieldRepo}
}

type CreateCropInput struct {
	FieldID         uuid.UUID         `json:"field_id"`
	CropType        string            `json:"crop_type"`
	Variety         *string           `json:"variety,omitempty"`
	PlantingDate    string            `json:"planting_date"`
	ExpectedHarvest *string           `json:"expected_harvest,omitempty"`
	Status          models.CropStatus `json:"status,omitempty"`
}

type UpdateCropInput struct {
	Variety         *string            `json:"variety,omitempty"`
	Status          *models.CropStatus `json:"status,omitempty"`
	ExpectedHarvest *string            `json:"expected_harvest,omitempty"`
	ActualHarvest   *string            `json:"actual_harvest,omitempty"`
}

// cropToResource converts a Crop model to JSON:API resource
func cropToResource(crop *models.Crop) jsonapi.Resource {
	attrs := map[string]interface{}{
		"crop_type":     crop.CropType,
		"planting_date": crop.PlantingDate.Format("2006-01-02"),
		"status":        crop.Status,
		"created_at":    crop.CreatedAt,
		"updated_at":    crop.UpdatedAt,
	}
	if crop.Variety != nil {
		attrs["variety"] = *crop.Variety
	}
	if crop.ExpectedHarvest != nil {
		attrs["expected_harvest"] = crop.ExpectedHarvest.Format("2006-01-02")
	}
	if crop.ActualHarvest != nil {
		attrs["actual_harvest"] = crop.ActualHarvest.Format("2006-01-02")
	}

	return jsonapi.Resource{
		ID:         crop.ID.String(),
		Type:       "crops",
		Attributes: attrs,
		Relationships: map[string]jsonapi.Relationship{
			"field": {
				Data: jsonapi.NewResourceIdentifier("fields", crop.FieldID.String()),
			},
		},
		Links: &jsonapi.Links{
			Self: "/crops/" + crop.ID.String(),
		},
	}
}

// ListCrops returns all crops for the authenticated user
func (h *CropHandler) ListCrops(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		jsonapi.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", "missing user ID")
		return
	}

	crops, err := h.cropRepo.FindByUserID(userID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to fetch crops")
		return
	}

	resources := make([]jsonapi.Resource, len(crops))
	for i, c := range crops {
		resources[i] = cropToResource(&c)
	}

	jsonapi.RespondWithResources(w, http.StatusOK, resources)
}

// GetCrop returns a specific crop
func (h *CropHandler) GetCrop(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	cropID, err := uuid.Parse(chi.URLParam(r, "id"))
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

	jsonapi.RespondWithResource(w, http.StatusOK, cropToResource(crop))
}

// CreateCrop creates a new crop
func (h *CropHandler) CreateCrop(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		jsonapi.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", "missing user ID")
		return
	}

	var input CreateCropInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Bad Request", "invalid request body")
		return
	}

	field, err := h.fieldRepo.FindByID(input.FieldID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusNotFound, "Not Found", "field not found")
		return
	}
	if field.UserID != userID {
		jsonapi.RespondWithError(w, http.StatusForbidden, "Forbidden", "access denied")
		return
	}

	plantingDate, err := time.Parse("2006-01-02", input.PlantingDate)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Validation Error", "invalid planting_date format (use YYYY-MM-DD)")
		return
	}

	crop := &models.Crop{
		FieldID:      input.FieldID,
		CropType:     input.CropType,
		Variety:      input.Variety,
		PlantingDate: plantingDate,
		Status:       models.StatusPlanted,
	}

	if input.Status != "" {
		crop.Status = input.Status
	}

	if input.ExpectedHarvest != nil {
		expectedHarvest, _ := time.Parse("2006-01-02", *input.ExpectedHarvest)
		crop.ExpectedHarvest = &expectedHarvest
	}

	if err := h.cropRepo.Create(crop); err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to create crop")
		return
	}

	jsonapi.RespondWithResource(w, http.StatusCreated, cropToResource(crop))
}

// UpdateCrop updates a crop
func (h *CropHandler) UpdateCrop(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	cropID, err := uuid.Parse(chi.URLParam(r, "id"))
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

	var input UpdateCropInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Bad Request", "invalid request body")
		return
	}

	if input.Variety != nil {
		crop.Variety = input.Variety
	}
	if input.Status != nil {
		crop.Status = *input.Status
	}
	if input.ExpectedHarvest != nil {
		t, _ := time.Parse("2006-01-02", *input.ExpectedHarvest)
		crop.ExpectedHarvest = &t
	}
	if input.ActualHarvest != nil {
		t, _ := time.Parse("2006-01-02", *input.ActualHarvest)
		crop.ActualHarvest = &t
	}

	if err := h.cropRepo.Update(crop); err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to update crop")
		return
	}

	jsonapi.RespondWithResource(w, http.StatusOK, cropToResource(crop))
}

// DeleteCrop deletes a crop
func (h *CropHandler) DeleteCrop(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	cropID, err := uuid.Parse(chi.URLParam(r, "id"))
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

	if err := h.cropRepo.Delete(cropID); err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to delete crop")
		return
	}

	jsonapi.RespondWithMeta(w, http.StatusOK, jsonapi.Meta{"message": "crop deleted"})
}

// ListCropsByField returns crops for a specific field
func (h *CropHandler) ListCropsByField(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	fieldID, err := uuid.Parse(chi.URLParam(r, "fieldId"))
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Bad Request", "invalid field ID")
		return
	}

	field, err := h.fieldRepo.FindByID(fieldID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusNotFound, "Not Found", "field not found")
		return
	}
	if field.UserID != userID {
		jsonapi.RespondWithError(w, http.StatusForbidden, "Forbidden", "access denied")
		return
	}

	crops, err := h.cropRepo.FindByFieldID(fieldID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to fetch crops")
		return
	}

	resources := make([]jsonapi.Resource, len(crops))
	for i, c := range crops {
		resources[i] = cropToResource(&c)
	}

	jsonapi.RespondWithResources(w, http.StatusOK, resources)
}
