package jsonapi

import (
	"encoding/json"
	"net/http"
)

// Document represents a JSON:API top-level document
type Document struct {
	Data     interface{}  `json:"data,omitempty"`
	Errors   []Error      `json:"errors,omitempty"`
	Meta     *Meta        `json:"meta,omitempty"`
	Links    *Links       `json:"links,omitempty"`
	Included []Resource   `json:"included,omitempty"`
	JSONAPI  *VersionInfo `json:"jsonapi,omitempty"`
}

// Resource represents a JSON:API resource object
type Resource struct {
	ID            string                  `json:"id"`
	Type          string                  `json:"type"`
	Attributes    map[string]interface{}  `json:"attributes,omitempty"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
	Links         *Links                  `json:"links,omitempty"`
	Meta          map[string]interface{}  `json:"meta,omitempty"`
}

// Relationship represents a JSON:API relationship object
type Relationship struct {
	Data  interface{}            `json:"data,omitempty"`
	Links *Links                 `json:"links,omitempty"`
	Meta  map[string]interface{} `json:"meta,omitempty"`
}

// ResourceIdentifier represents a JSON:API resource identifier
type ResourceIdentifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// Error represents a JSON:API error object
type Error struct {
	ID     string                 `json:"id,omitempty"`
	Status string                 `json:"status"`
	Code   string                 `json:"code,omitempty"`
	Title  string                 `json:"title"`
	Detail string                 `json:"detail,omitempty"`
	Source *ErrorSource           `json:"source,omitempty"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

// ErrorSource represents the source of a JSON:API error
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
	Header    string `json:"header,omitempty"`
}

// Links represents a JSON:API links object
type Links struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
	First   string `json:"first,omitempty"`
	Last    string `json:"last,omitempty"`
	Prev    string `json:"prev,omitempty"`
	Next    string `json:"next,omitempty"`
}

// Meta represents a JSON:API meta object
type Meta map[string]interface{}

// VersionInfo represents the jsonapi member
type VersionInfo struct {
	Version string `json:"version"`
}

// Response helpers

// RespondWithData sends a JSON:API document with data
func RespondWithData(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Document{
		Data:    data,
		JSONAPI: &VersionInfo{Version: "1.1"},
	})
}

// RespondWithResource sends a single resource
func RespondWithResource(w http.ResponseWriter, status int, resource Resource) {
	RespondWithData(w, status, resource)
}

// RespondWithResources sends a collection of resources
func RespondWithResources(w http.ResponseWriter, status int, resources []Resource) {
	RespondWithData(w, status, resources)
}

// RespondWithError sends a JSON:API error document
func RespondWithError(w http.ResponseWriter, status int, title, detail string) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Document{
		Errors: []Error{{
			Status: http.StatusText(status),
			Title:  title,
			Detail: detail,
		}},
		JSONAPI: &VersionInfo{Version: "1.1"},
	})
}

// RespondWithErrors sends multiple JSON:API errors
func RespondWithErrors(w http.ResponseWriter, status int, errors []Error) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Document{
		Errors:  errors,
		JSONAPI: &VersionInfo{Version: "1.1"},
	})
}

// RespondWithMeta sends a JSON:API document with only meta
func RespondWithMeta(w http.ResponseWriter, status int, meta Meta) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Document{
		Meta:    &meta,
		JSONAPI: &VersionInfo{Version: "1.1"},
	})
}

// NewResource creates a new resource with the given type and ID
func NewResource(resourceType, id string, attributes map[string]interface{}) Resource {
	return Resource{
		ID:         id,
		Type:       resourceType,
		Attributes: attributes,
	}
}

// NewResourceIdentifier creates a resource identifier
func NewResourceIdentifier(resourceType, id string) ResourceIdentifier {
	return ResourceIdentifier{
		Type: resourceType,
		ID:   id,
	}
}
