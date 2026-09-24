package resource

import (
	"fmt"
	"sync"
	"time"
)

// LifecycleState defines the resource lifecycle.
type LifecycleState string

const (
	LifecyclePending       LifecycleState = "pending"
	LifecycleProvisioning  LifecycleState = "provisioning"
	LifecycleActive        LifecycleState = "active"
	LifecycleDegraded      LifecycleState = "degraded"
	LifecycleDecommissioned LifecycleState = "decommissioned"
	LifecycleDeleted       LifecycleState = "deleted"
)

// UniversalResource defines the canonical resource model (Specification Section 5.4).
type UniversalResource struct {
	ID               string                 `json:"id"`                 // Stable UUID
	Name             string                 `json:"name"`               // Human-readable identifier
	Type             string                 `json:"type"`               // site, database, mailbox, cron, container
	Provider         string                 `json:"provider"`           // local, aws, hetzner, digitalocean
	NodeID           string                 `json:"node_id"`            // Server/Node where resource lives
	OrgID            string                 `json:"org_id"`             // Organization hierarchy
	ProjectID        string                 `json:"project_id"`         // Project
	Environment      string                 `json:"environment"`        // production, staging, development
	OwnerID          string                 `json:"owner_id"`           // User actor
	Tags             []string               `json:"tags"`
	Labels           map[string]string      `json:"labels"`
	Annotations      map[string]string      `json:"annotations"`
	DesiredState     map[string]interface{} `json:"desired_state"`
	ObservedState    map[string]interface{} `json:"observed_state"`
	Health           string                 `json:"health"`             // healthy, warning, critical, unknown
	Lifecycle        LifecycleState         `json:"lifecycle"`
	Dependencies     []string               `json:"dependencies"`       // IDs of prerequisite resources
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	LastReconciledAt time.Time              `json:"last_reconciled_at"`
	DeletedAt        *time.Time             `json:"deleted_at,omitempty"` // Soft-delete metadata
	Actor            string                 `json:"actor"`              // Last modifying user/job
}

// ResourceManager provides in-memory and database-backed management of resources.
type ResourceManager struct {
	mu        sync.RWMutex
	resources map[string]*UniversalResource
}

// NewResourceManager initializes resource tracking.
func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		resources: make(map[string]*UniversalResource),
	}
}

// Register adds or updates a resource in the global hierarchy.
func (m *ResourceManager) Register(r *UniversalResource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if r.ID == "" {
		return fmt.Errorf("resource ID cannot be empty")
	}
	if r.OrgID == "" || r.ProjectID == "" {
		return fmt.Errorf("resource must belong to an organization and project hierarchy")
	}

	r.UpdatedAt = time.Now().UTC()
	if r.CreatedAt.IsZero() {
		r.CreatedAt = r.UpdatedAt
	}
	if r.Lifecycle == "" {
		r.Lifecycle = LifecycleActive
	}

	m.resources[r.ID] = r
	return nil
}

// Get retrieves a resource by UUID.
func (m *ResourceManager) Get(id string) (*UniversalResource, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res, ok := m.resources[id]
	return res, ok
}

// SoftDelete marks a resource as deleted without purging immediately.
func (m *ResourceManager) SoftDelete(id string, actor string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	res, ok := m.resources[id]
	if !ok {
		return fmt.Errorf("resource %s not found", id)
	}

	now := time.Now().UTC()
	res.DeletedAt = &now
	res.Lifecycle = LifecycleDeleted
	res.Actor = actor
	return nil
}

// Search queries resources by organization, project, type or health.
func (m *ResourceManager) Search(orgID, projectID, resType, health string) []*UniversalResource {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make([]*UniversalResource, 0)
	for _, r := range m.resources {
		if r.DeletedAt != nil {
			continue // ignore soft-deleted items by default
		}
		if orgID != "" && r.OrgID != orgID {
			continue
		}
		if projectID != "" && r.ProjectID != projectID {
			continue
		}
		if resType != "" && r.Type != resType {
			continue
		}
		if health != "" && r.Health != health {
			continue
		}
		results = append(results, r)
	}
	return results
}
