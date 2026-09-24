package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// APIResponse provides standard JSON wrapper.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Time    time.Time   `json:"time"`
}

// Router dispatches REST API requests for KenPanel control plane.
type Router struct {
	mux *http.ServeMux
}

// NewRouter sets up all versioned endpoints.
func NewRouter() *Router {
	r := &Router{mux: http.NewServeMux()}
	r.registerRoutes()
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	r.mux.ServeHTTP(w, req)
}

func (r *Router) registerRoutes() {
	r.mux.HandleFunc("/api/v1/health", r.handleHealth)
	r.mux.HandleFunc("/api/v1/system/inventory", r.handleInventory)
	r.mux.HandleFunc("/api/v1/sites", r.handleSites)
	r.mux.HandleFunc("/api/v1/nodes", r.handleNodes)
	r.mux.HandleFunc("/api/v1/openapi.json", r.handleOpenAPI)
}

func (r *Router) handleHealth(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"status":  "healthy",
			"version": "2.3.0",
		},
		Time: time.Now().UTC(),
	})
}

func (r *Router) handleInventory(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"os":   "Ubuntu 24.04 LTS",
			"arch": "x86_64",
			"cpus": 4,
			"ram":  8192,
		},
		Time: time.Now().UTC(),
	})
}

func (r *Router) handleSites(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data: []map[string]string{
			{"domain": "example.com", "type": "php", "status": "active"},
		},
		Time: time.Now().UTC(),
	})
}

func (r *Router) handleNodes(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data: []map[string]string{
			{"node_id": "node-01", "role": "controller", "status": "healthy"},
		},
		Time: time.Now().UTC(),
	})
}

func (r *Router) handleOpenAPI(w http.ResponseWriter, req *http.Request) {
	schema := GetOpenAPISpec()
	w.Write([]byte(schema))
}

// GetOpenAPISpec returns the OpenAPI 3.1 JSON definition.
func GetOpenAPISpec() string {
	return `{
  "openapi": "3.1.0",
  "info": {
    "title": "KenPanel Control Plane REST API",
    "version": "2.3.0",
    "description": "Unified hosting, application, container and fleet operations API."
  },
  "paths": {
    "/api/v1/health": {
      "get": {
        "summary": "Health status probe",
        "responses": { "200": { "description": "System operational" } }
      }
    },
    "/api/v1/system/inventory": {
      "get": {
        "summary": "Hardware and OS inventory",
        "responses": { "200": { "description": "Inventory profile" } }
      }
    },
    "/api/v1/sites": {
      "get": {
        "summary": "List managed websites",
        "responses": { "200": { "description": "List of websites" } }
      }
    },
    "/api/v1/nodes": {
      "get": {
        "summary": "List fleet worker nodes",
        "responses": { "200": { "description": "List of nodes" } }
      }
    }
  }
}`
}
