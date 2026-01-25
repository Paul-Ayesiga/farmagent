package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/farmagent/fa-crop-service/internal/handlers"
	"github.com/farmagent/fa-crop-service/internal/models"
	"github.com/farmagent/fa-crop-service/pkg/jsonapi"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ===== Mock Repositories =====

type MockFieldRepository struct {
	mock.Mock
}

func (m *MockFieldRepository) Create(field *models.Field) error {
	args := m.Called(field)
	return args.Error(0)
}

func (m *MockFieldRepository) FindByID(id uuid.UUID) (*models.Field, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Field), args.Error(1)
}

func (m *MockFieldRepository) FindByUserID(userID string) ([]models.Field, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Field), args.Error(1)
}

func (m *MockFieldRepository) Update(field *models.Field) error {
	args := m.Called(field)
	return args.Error(0)
}

func (m *MockFieldRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockCropRepository struct {
	mock.Mock
}

func (m *MockCropRepository) Create(crop *models.Crop) error {
	args := m.Called(crop)
	return args.Error(0)
}

func (m *MockCropRepository) FindByID(id uuid.UUID) (*models.Crop, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Crop), args.Error(1)
}

func (m *MockCropRepository) FindByFieldID(fieldID uuid.UUID) ([]models.Crop, error) {
	args := m.Called(fieldID)
	return args.Get(0).([]models.Crop), args.Error(1)
}

func (m *MockCropRepository) FindByUserID(userID string) ([]models.Crop, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Crop), args.Error(1)
}

func (m *MockCropRepository) Update(crop *models.Crop) error {
	args := m.Called(crop)
	return args.Error(0)
}

func (m *MockCropRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// ===== JSON:API Response Helpers =====

type JSONAPIDocument struct {
	Data    interface{}     `json:"data"`
	Errors  []jsonapi.Error `json:"errors,omitempty"`
	Meta    jsonapi.Meta    `json:"meta,omitempty"`
	JSONAPI *struct {
		Version string `json:"version"`
	} `json:"jsonapi,omitempty"`
}

// ===== Field Tests =====

func TestListFields_Success_JSONAPI(t *testing.T) {
	mockFieldRepo := new(MockFieldRepository)
	handler := handlers.NewFieldHandler(mockFieldRepo)

	userID := "user-123"
	fields := []models.Field{
		{ID: uuid.New(), UserID: userID, Name: "Field 1", SizeAcres: 2.5},
		{ID: uuid.New(), UserID: userID, Name: "Field 2", SizeAcres: 3.0},
	}

	mockFieldRepo.On("FindByUserID", userID).Return(fields, nil)

	r := chi.NewRouter()
	r.Get("/fields", handler.ListFields)

	req := httptest.NewRequest(http.MethodGet, "/fields", nil)
	req.Header.Set("X-User-ID", userID)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/vnd.api+json", rec.Header().Get("Content-Type"))

	var doc JSONAPIDocument
	json.Unmarshal(rec.Body.Bytes(), &doc)
	assert.NotNil(t, doc.Data)
	assert.Equal(t, "1.1", doc.JSONAPI.Version)
	mockFieldRepo.AssertExpectations(t)
}

func TestListFields_MissingUserID_JSONAPI(t *testing.T) {
	mockFieldRepo := new(MockFieldRepository)
	handler := handlers.NewFieldHandler(mockFieldRepo)

	r := chi.NewRouter()
	r.Get("/fields", handler.ListFields)

	req := httptest.NewRequest(http.MethodGet, "/fields", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "application/vnd.api+json", rec.Header().Get("Content-Type"))

	var doc JSONAPIDocument
	json.Unmarshal(rec.Body.Bytes(), &doc)
	assert.NotEmpty(t, doc.Errors)
	assert.Equal(t, "Unauthorized", doc.Errors[0].Title)
}

func TestCreateField_Success_JSONAPI(t *testing.T) {
	mockFieldRepo := new(MockFieldRepository)
	handler := handlers.NewFieldHandler(mockFieldRepo)

	userID := "user-123"
	mockFieldRepo.On("Create", mock.AnythingOfType("*models.Field")).Return(nil)

	r := chi.NewRouter()
	r.Post("/fields", handler.CreateField)

	body := map[string]interface{}{
		"name":       "Test Field",
		"latitude":   0.3476,
		"longitude":  32.5825,
		"size_acres": 2.5,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/fields", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/vnd.api+json", rec.Header().Get("Content-Type"))

	var doc map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &doc)
	data := doc["data"].(map[string]interface{})
	assert.Equal(t, "fields", data["type"])
	assert.NotEmpty(t, data["id"])
	assert.NotNil(t, data["attributes"])
	mockFieldRepo.AssertExpectations(t)
}

func TestCreateField_MissingName_JSONAPI(t *testing.T) {
	mockFieldRepo := new(MockFieldRepository)
	handler := handlers.NewFieldHandler(mockFieldRepo)

	r := chi.NewRouter()
	r.Post("/fields", handler.CreateField)

	body := map[string]interface{}{
		"latitude":   0.3476,
		"longitude":  32.5825,
		"size_acres": 2.5,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/fields", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var doc JSONAPIDocument
	json.Unmarshal(rec.Body.Bytes(), &doc)
	assert.NotEmpty(t, doc.Errors)
	assert.Equal(t, "Validation Error", doc.Errors[0].Title)
}

func TestGetField_NotOwner_JSONAPI(t *testing.T) {
	mockFieldRepo := new(MockFieldRepository)
	handler := handlers.NewFieldHandler(mockFieldRepo)

	fieldID := uuid.New()
	field := &models.Field{
		ID:     fieldID,
		UserID: "other-user",
		Name:   "Someone's Field",
	}

	mockFieldRepo.On("FindByID", fieldID).Return(field, nil)

	r := chi.NewRouter()
	r.Get("/fields/{id}", handler.GetField)

	req := httptest.NewRequest(http.MethodGet, "/fields/"+fieldID.String(), nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)

	var doc JSONAPIDocument
	json.Unmarshal(rec.Body.Bytes(), &doc)
	assert.NotEmpty(t, doc.Errors)
	assert.Equal(t, "Forbidden", doc.Errors[0].Title)
}

func TestDeleteField_Success_JSONAPI(t *testing.T) {
	mockFieldRepo := new(MockFieldRepository)
	handler := handlers.NewFieldHandler(mockFieldRepo)

	userID := "user-123"
	fieldID := uuid.New()
	field := &models.Field{
		ID:     fieldID,
		UserID: userID,
		Name:   "My Field",
	}

	mockFieldRepo.On("FindByID", fieldID).Return(field, nil)
	mockFieldRepo.On("Delete", fieldID).Return(nil)

	r := chi.NewRouter()
	r.Delete("/fields/{id}", handler.DeleteField)

	req := httptest.NewRequest(http.MethodDelete, "/fields/"+fieldID.String(), nil)
	req.Header.Set("X-User-ID", userID)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var doc map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &doc)
	assert.NotNil(t, doc["meta"])
	mockFieldRepo.AssertExpectations(t)
}

// ===== Crop Tests =====

func TestListCrops_Success_JSONAPI(t *testing.T) {
	mockCropRepo := new(MockCropRepository)
	mockFieldRepo := new(MockFieldRepository)
	handler := handlers.NewCropHandler(mockCropRepo, mockFieldRepo)

	userID := "user-123"
	crops := []models.Crop{
		{ID: uuid.New(), CropType: "maize", Status: models.StatusPlanted, PlantingDate: time.Now()},
		{ID: uuid.New(), CropType: "beans", Status: models.StatusGrowing, PlantingDate: time.Now()},
	}

	mockCropRepo.On("FindByUserID", userID).Return(crops, nil)

	r := chi.NewRouter()
	r.Get("/crops", handler.ListCrops)

	req := httptest.NewRequest(http.MethodGet, "/crops", nil)
	req.Header.Set("X-User-ID", userID)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/vnd.api+json", rec.Header().Get("Content-Type"))
}

func TestCreateCrop_Success_JSONAPI(t *testing.T) {
	mockCropRepo := new(MockCropRepository)
	mockFieldRepo := new(MockFieldRepository)
	handler := handlers.NewCropHandler(mockCropRepo, mockFieldRepo)

	userID := "user-123"
	fieldID := uuid.New()
	field := &models.Field{ID: fieldID, UserID: userID}

	mockFieldRepo.On("FindByID", fieldID).Return(field, nil)
	mockCropRepo.On("Create", mock.AnythingOfType("*models.Crop")).Return(nil)

	r := chi.NewRouter()
	r.Post("/crops", handler.CreateCrop)

	body := map[string]interface{}{
		"field_id":      fieldID.String(),
		"crop_type":     "maize",
		"planting_date": time.Now().Format("2006-01-02"),
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/crops", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/vnd.api+json", rec.Header().Get("Content-Type"))

	var doc map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &doc)
	data := doc["data"].(map[string]interface{})
	assert.Equal(t, "crops", data["type"])
	assert.NotNil(t, data["relationships"])
}

func TestCreateCrop_FieldNotOwned_JSONAPI(t *testing.T) {
	mockCropRepo := new(MockCropRepository)
	mockFieldRepo := new(MockFieldRepository)
	handler := handlers.NewCropHandler(mockCropRepo, mockFieldRepo)

	fieldID := uuid.New()
	field := &models.Field{ID: fieldID, UserID: "other-user"}

	mockFieldRepo.On("FindByID", fieldID).Return(field, nil)

	r := chi.NewRouter()
	r.Post("/crops", handler.CreateCrop)

	body := map[string]interface{}{
		"field_id":      fieldID.String(),
		"crop_type":     "maize",
		"planting_date": "2026-01-15",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/crops", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)

	var doc JSONAPIDocument
	json.Unmarshal(rec.Body.Bytes(), &doc)
	assert.NotEmpty(t, doc.Errors)
}

func TestCreateCrop_InvalidDate_JSONAPI(t *testing.T) {
	mockCropRepo := new(MockCropRepository)
	mockFieldRepo := new(MockFieldRepository)
	handler := handlers.NewCropHandler(mockCropRepo, mockFieldRepo)

	userID := "user-123"
	fieldID := uuid.New()
	field := &models.Field{ID: fieldID, UserID: userID}

	mockFieldRepo.On("FindByID", fieldID).Return(field, nil)

	r := chi.NewRouter()
	r.Post("/crops", handler.CreateCrop)

	body := map[string]interface{}{
		"field_id":      fieldID.String(),
		"crop_type":     "maize",
		"planting_date": "invalid-date",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/crops", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var doc JSONAPIDocument
	json.Unmarshal(rec.Body.Bytes(), &doc)
	assert.Equal(t, "Validation Error", doc.Errors[0].Title)
}
