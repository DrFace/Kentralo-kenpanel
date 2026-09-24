package website

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// SiteType represents supported runtime execution environments.
type SiteType string

const (
	SiteTypePHP          SiteType = "php"
	SiteTypeStatic       SiteType = "static"
	SiteTypeNodeJS       SiteType = "nodejs"
	SiteTypePython       SiteType = "python"
	SiteTypeReverseProxy SiteType = "reverse_proxy"
	SiteTypeContainer    SiteType = "container"
	SiteTypeWordPress    SiteType = "wordpress"
)

// RedirectRule defines a traffic redirection behavior.
type RedirectRule struct {
	SourcePath  string `json:"source_path"`
	TargetURL   string `json:"target_url"`
	StatusCode  int    `json:"status_code"` // 301, 302, 307, 308
	PreserveURI bool   `json:"preserve_uri"`
}

// WebsiteConfig represents comprehensive site settings (Section 8.1).
type WebsiteConfig struct {
	ID                  string            `json:"id"`
	PrimaryDomain       string            `json:"primary_domain"`
	Aliases             []string          `json:"aliases"`
	Subdomains          []string          `json:"subdomains"`
	WildcardSubdomain   bool              `json:"wildcard_subdomain"`
	Type                SiteType          `json:"type"`
	RuntimeVersion      string            `json:"runtime_version"` // e.g. "8.2", "20"
	DocumentRoot        string            `json:"document_root"`
	PublicDir           string            `json:"public_dir"` // e.g. "public"
	SystemUser          string            `json:"system_user"`
	DedicatedIP         string            `json:"dedicated_ip,omitempty"`
	ForceHTTPS          bool              `json:"force_https"`
	ApexRedirect        string            `json:"apex_redirect"` // "none", "to_www", "to_apex"
	CustomRedirects     []RedirectRule    `json:"custom_redirects"`
	CustomErrorPages    map[int]string    `json:"custom_error_pages"` // 404, 500, 502
	BasicAuthEnabled    bool              `json:"basic_auth_enabled"`
	BasicAuthUser       string            `json:"basic_auth_user,omitempty"`
	BasicAuthHash       string            `json:"-"`
	MaintenanceMode     bool              `json:"maintenance_mode"`
	MaintenanceHTML     string            `json:"maintenance_html,omitempty"`
	EnvVariables        map[string]string `json:"env_variables"`
	Status              string            `json:"status"` // active, suspended, archived
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
}

// LifecycleEngine manages website provisioning, configuration generation and lifecycle transitions.
type LifecycleEngine struct {
	mu    sync.RWMutex
	sites map[string]*WebsiteConfig // key: site ID
}

// NewLifecycleEngine initializes the website manager.
func NewLifecycleEngine() *LifecycleEngine {
	return &LifecycleEngine{
		sites: make(map[string]*WebsiteConfig),
	}
}

// CreateSite initializes a new website with secure isolation and defaults.
func (e *LifecycleEngine) CreateSite(ctx context.Context, domain string, siteType SiteType, runtimeVer string) (*WebsiteConfig, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	domain = strings.ToLower(strings.TrimSpace(domain))
	if domain == "" {
		return nil, fmt.Errorf("primary domain cannot be blank")
	}

	siteID := fmt.Sprintf("site-%s", strings.ReplaceAll(domain, ".", "-"))
	if _, exists := e.sites[siteID]; exists {
		return nil, fmt.Errorf("website %s already exists", domain)
	}

	safeUser := strings.ReplaceAll(domain, ".", "_")
	if len(safeUser) > 16 {
		safeUser = safeUser[:16]
	}

	site := &WebsiteConfig{
		ID:                siteID,
		PrimaryDomain:     domain,
		Aliases:           make([]string, 0),
		Subdomains:        make([]string, 0),
		WildcardSubdomain: false,
		Type:              siteType,
		RuntimeVersion:    runtimeVer,
		DocumentRoot:      fmt.Sprintf("/var/www/%s", domain),
		PublicDir:         "public",
		SystemUser:        safeUser,
		ForceHTTPS:        true,
		ApexRedirect:      "none",
		CustomRedirects:   make([]RedirectRule, 0),
		CustomErrorPages: map[int]string{
			404: "/error/404.html",
			502: "/error/502.html",
		},
		BasicAuthEnabled: false,
		MaintenanceMode:  false,
		EnvVariables:     make(map[string]string),
		Status:           "active",
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	e.sites[siteID] = site
	return site, nil
}

// SuspendSite toggles website suspension.
func (e *LifecycleEngine) SuspendSite(siteID string, reason string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	site, ok := e.sites[siteID]
	if !ok {
		return fmt.Errorf("site %s not found", siteID)
	}
	site.Status = "suspended"
	site.UpdatedAt = time.Now().UTC()
	return nil
}

// SetMaintenanceMode enables or disables maintenance screen.
func (e *LifecycleEngine) SetMaintenanceMode(siteID string, enabled bool, customHTML string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	site, ok := e.sites[siteID]
	if !ok {
		return fmt.Errorf("site %s not found", siteID)
	}
	site.MaintenanceMode = enabled
	if customHTML != "" {
		site.MaintenanceHTML = customHTML
	}
	site.UpdatedAt = time.Now().UTC()
	return nil
}

// CloneSite duplicates an existing site configuration to a new domain name.
func (e *LifecycleEngine) CloneSite(sourceSiteID, targetDomain string) (*WebsiteConfig, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	src, ok := e.sites[sourceSiteID]
	if !ok {
		return nil, fmt.Errorf("source site %s not found", sourceSiteID)
	}

	targetID := fmt.Sprintf("site-%s", strings.ReplaceAll(targetDomain, ".", "-"))
	cloned := *src
	cloned.ID = targetID
	cloned.PrimaryDomain = targetDomain
	cloned.DocumentRoot = fmt.Sprintf("/var/www/%s", targetDomain)
	cloned.CreatedAt = time.Now().UTC()
	cloned.UpdatedAt = time.Now().UTC()

	e.sites[targetID] = &cloned
	return &cloned, nil
}
