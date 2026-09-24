package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DrFace/Kentralo-kenpanel/pkg/engines/doctor"
)

// APIHandler aggregates REST API routes for KenPanel.
type APIHandler struct {
	doctorEngine *doctor.WebStackDoctor
}

func NewAPIHandler() *APIHandler {
	return &APIHandler{
		doctorEngine: doctor.NewWebStackDoctor(),
	}
}

// RegisterRoutes attaches the v1 REST API routes to a standard HTTP ServeMux.
func (h *APIHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/health", h.HealthCheck)
	mux.HandleFunc("GET /api/v1/doctor/diagnose", h.DiagnoseDomain)
}

func (h *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"service": "kenpanel-api",
	})
}

func (h *APIHandler) DiagnoseDomain(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, `{"error":"missing domain parameter"}`, http.StatusBadRequest)
		return
	}

	report, err := h.doctorEngine.Diagnose(r.Context(), domain, "", "")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(report)
}
