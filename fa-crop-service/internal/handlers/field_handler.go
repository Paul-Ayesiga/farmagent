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

type FieldHandler struct {
	fieldRepo repository.FieldRepository
}

func NewFieldHandler(fieldRepo repository.FieldRepository) *FieldHandler {
	return &FieldHandler{fieldRepo: fieldRepo}
}

type CreateFieldInput struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	SizeAcres float64 `json:"size_acres"`
	SoilType  *string `json:"soil_type,omitempty"`
}

type UpdateFieldInput struct {
	Name      *string  `json:"name,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	SizeAcres *float64 `json:"size_acres,omitempty"`
	SoilType  *string  `json:"soil_type,omitempty"`
}

// fieldToResource converts a Field model to JSON:API resource
func fieldToResource(field *models.Field) jsonapi.Resource {
	attrs := map[string]interface{}{
		"name":       field.Name,
		"latitude":   field.Latitude,
		"longitude":  field.Longitude,
		"size_acres": field.SizeAcres,
		"created_at": field.CreatedAt,
		"updated_at": field.UpdatedAt,
	}
	if field.SoilType != nil {
		attrs["soil_type"] = *field.SoilType
	}

	return jsonapi.Resource{
		ID:         field.ID.String(),
		Type:       "fields",
		Attributes: attrs,
		Links: &jsonapi.Links{
			Self: "/fields/" + field.ID.String(),
		},
	}
}

// ListFields returns all fields for the authenticated user
func (h *FieldHandler) ListFields(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		jsonapi.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", "missing user ID")
		return
	}

	fields, err := h.fieldRepo.FindByUserID(userID)
	if err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to fetch fields")
		return
	}

	resources := make([]jsonapi.Resource, len(fields))
	for i, f := range fields {
		resources[i] = fieldToResource(&f)
	}

	jsonapi.RespondWithResources(w, http.StatusOK, resources)
}

// GetField returns a specific field
func (h *FieldHandler) GetField(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	fieldID, err := uuid.Parse(chi.URLParam(r, "id"))
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

	jsonapi.RespondWithResource(w, http.StatusOK, fieldToResource(field))
}

// CreateField creates a new field
func (h *FieldHandler) CreateField(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		jsonapi.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", "missing user ID")
		return
	}

	var input CreateFieldInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Bad Request", "invalid request body")
		return
	}

	if input.Name == "" {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Validation Error", "name is required")
		return
	}

	field := &models.Field{
		UserID:    userID,
		Name:      input.Name,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		SizeAcres: input.SizeAcres,
		SoilType:  input.SoilType,
	}

	if err := h.fieldRepo.Create(field); err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to create field")
		return
	}

	jsonapi.RespondWithResource(w, http.StatusCreated, fieldToResource(field))
}

// UpdateField updates a field
func (h *FieldHandler) UpdateField(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	fieldID, err := uuid.Parse(chi.URLParam(r, "id"))
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

	var input UpdateFieldInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonapi.RespondWithError(w, http.StatusBadRequest, "Bad Request", "invalid request body")
		return
	}

	if input.Name != nil {
		field.Name = *input.Name
	}
	if input.Latitude != nil {
		field.Latitude = *input.Latitude
	}
	if input.Longitude != nil {
		field.Longitude = *input.Longitude
	}
	if input.SizeAcres != nil {
		field.SizeAcres = *input.SizeAcres
	}
	if input.SoilType != nil {
		field.SoilType = input.SoilType
	}
	field.UpdatedAt = time.Now()

	if err := h.fieldRepo.Update(field); err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to update field")
		return
	}

	jsonapi.RespondWithResource(w, http.StatusOK, fieldToResource(field))
}

// DeleteField deletes a field
func (h *FieldHandler) DeleteField(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	fieldID, err := uuid.Parse(chi.URLParam(r, "id"))
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

	if err := h.fieldRepo.Delete(fieldID); err != nil {
		jsonapi.RespondWithError(w, http.StatusInternalServerError, "Internal Error", "failed to delete field")
		return
	}

	jsonapi.RespondWithMeta(w, http.StatusOK, jsonapi.Meta{"message": "field deleted"})
}
