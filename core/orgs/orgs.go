package orgs

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Organization represents a top-level tenant.
type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	OwnerID   string    `json:"owner_id"`
	Status    string    `json:"status"` // active, suspended, deleted
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Project belongs to an organization and contains environments.
type Project struct {
	ID           string        `json:"id"`
	OrgID        string        `json:"org_id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Environments []Environment `json:"environments"`
	CreatedAt    time.Time     `json:"created_at"`
}

// Environment represents an deployment stage (prod, staging, dev).
type Environment struct {
	Name        string `json:"name"`        // production, staging, development
	Description string `json:"description"`
	IsProtected bool   `json:"is_protected"` // e.g. requires dual-approval for high-risk actions
}

// Team represents a group of users within an organization.
type Team struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	Name      string    `json:"name"`
	MemberIDs []string  `json:"member_ids"`
	CreatedAt time.Time `json:"created_at"`
}

// OrganizationManager handles tenancy operations.
type OrganizationManager struct {
	mu       sync.RWMutex
	orgs     map[string]*Organization
	projects map[string]*Project
	teams    map[string]*Team
}

// NewOrganizationManager initializes the tenancy manager.
func NewOrganizationManager() *OrganizationManager {
	return &OrganizationManager{
		orgs:     make(map[string]*Organization),
		projects: make(map[string]*Project),
		teams:    make(map[string]*Team),
	}
}

// CreateOrganization creates a new tenant organization.
func (m *OrganizationManager) CreateOrganization(name, slug, ownerID string) (*Organization, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	org := &Organization{
		ID:        fmt.Sprintf("org-%d", time.Now().UnixNano()),
		Name:      name,
		Slug:      slug,
		OwnerID:   ownerID,
		Status:    "active",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	m.orgs[org.ID] = org
	return org, nil
}

// CreateProject establishes a project with default prod/staging/dev environments.
func (m *OrganizationManager) CreateProject(orgID, name, description string) (*Project, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.orgs[orgID]; !ok {
		return nil, errors.New("organization not found")
	}

	proj := &Project{
		ID:          fmt.Sprintf("proj-%d", time.Now().UnixNano()),
		OrgID:       orgID,
		Name:        name,
		Description: description,
		Environments: []Environment{
			{Name: "production", Description: "Live production traffic", IsProtected: true},
			{Name: "staging", Description: "Pre-release testing mirror", IsProtected: false},
			{Name: "development", Description: "Local development builds", IsProtected: false},
		},
		CreatedAt: time.Now().UTC(),
	}

	m.projects[proj.ID] = proj
	return proj, nil
}

// SuspendOrganization suspends an organization and its resources.
func (m *OrganizationManager) SuspendOrganization(orgID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	org, ok := m.orgs[orgID]
	if !ok {
		return errors.New("organization not found")
	}
	org.Status = "suspended"
	org.UpdatedAt = time.Now().UTC()
	return nil
}
