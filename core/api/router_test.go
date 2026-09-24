package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPIRouterHealthAndOpenAPI(t *testing.T) {
	router := NewRouter()

	// Test health
	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "healthy") {
		t.Errorf("expected body to contain 'healthy'")
	}

	// Test openapi
	reqSpec := httptest.NewRequest("GET", "/api/v1/openapi.json", nil)
	wSpec := httptest.NewRecorder()
	router.ServeHTTP(wSpec, reqSpec)

	if wSpec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", wSpec.Code)
	}
	if !strings.Contains(wSpec.Body.String(), "3.1.0") {
		t.Errorf("expected OpenAPI 3.1.0 in spec")
	}
}
