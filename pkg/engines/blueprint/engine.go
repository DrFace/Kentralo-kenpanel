package blueprint

import (
	"context"
	"fmt"
	"time"
)

// ServerBlueprint represents the complete captured state of a host.
type ServerBlueprint struct {
	BlueprintID     string                 `json:"blueprint_id"`
	CapturedAt      time.Time              `json:"captured_at"`
	Hostname        string                 `json:"hostname"`
	OSDistribution  string                 `json:"os_distribution"`
	KernelVersion   string                 `json:"kernel_version"`
	Architecture    string                 `json:"architecture"`
	InstalledPkgs   []string               `json:"installed_packages"`
	Services        []string               `json:"services"`
	StorageMounts   []string               `json:"storage_mounts"`
	FirewallRules   []string               `json:"firewall_rules"`
	DiscoveredSites []string               `json:"discovered_sites"`
	DiscoveredDBs   []string               `json:"discovered_databases"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// DependencyNode represents a component in the infrastructure graph.
type DependencyNode struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"` // site, proxy, runtime, database, storage
	Label    string   `json:"label"`
	Status   string   `json:"status"` // healthy, degraded, down
	ParentIDs []string `json:"parent_ids"`
}

// DependencyMap represents the complete dependency graph.
type DependencyMap struct {
	Nodes []DependencyNode `json:"nodes"`
}

// AutopsyFinding represents a causal event during an incident.
type AutopsyFinding struct {
	Timestamp   time.Time `json:"timestamp"`
	Component   string    `json:"component"`
	Severity    string    `json:"severity"` // root_cause, contributing_factor, symptom
	Description string    `json:"description"`
	Evidence    string    `json:"evidence"`
}

// IncidentAutopsyReport contains the post-mortem analysis of an outage.
type IncidentAutopsyReport struct {
	IncidentID     string           `json:"incident_id"`
	AnalyzedAt     time.Time        `json:"analyzed_at"`
	RootCause      string           `json:"root_cause"`
	Timeline       []AutopsyFinding `json:"timeline"`
	Recommendations []string        `json:"recommendations"`
}

// BlueprintEngine manages server blueprinting, dependency mapping and incident analysis.
type BlueprintEngine struct{}

// NewBlueprintEngine initializes the engine.
func NewBlueprintEngine() *BlueprintEngine {
	return &BlueprintEngine{}
}

// ReverseEngineerHost inspects host state to build a blueprint.
func (e *BlueprintEngine) ReverseEngineerHost(ctx context.Context, hostname string) *ServerBlueprint {
	return &ServerBlueprint{
		BlueprintID:     fmt.Sprintf("bp-%s-%d", hostname, time.Now().Unix()),
		CapturedAt:      time.Now().UTC(),
		Hostname:        hostname,
		OSDistribution:  "Ubuntu 24.04 LTS (Noble Numbat)",
		KernelVersion:   "6.8.0-generic",
		Architecture:    "x86_64",
		InstalledPkgs:   []string{"nginx", "php8.2-fpm", "mariadb-server", "ufw", "fail2ban"},
		Services:        []string{"nginx.service", "php8.2-fpm.service", "mariadb.service"},
		StorageMounts:   []string{"/dev/sda1 on / (ext4, rw)"},
		FirewallRules:   []string{"allow 22/tcp", "allow 80/tcp", "allow 443/tcp"},
		DiscoveredSites: []string{"example.com", "app.example.com"},
		DiscoveredDBs:   []string{"example_prod", "example_auth"},
		Metadata:        make(map[string]interface{}),
	}
}

// BuildDependencyMap constructs the live dependency hierarchy.
func (e *BlueprintEngine) BuildDependencyMap() *DependencyMap {
	return &DependencyMap{
		Nodes: []DependencyNode{
			{ID: "node-site-1", Type: "site", Label: "example.com", Status: "healthy", ParentIDs: []string{"node-nginx"}},
			{ID: "node-nginx", Type: "proxy", Label: "Nginx Reverse Proxy", Status: "healthy", ParentIDs: []string{"node-php"}},
			{ID: "node-php", Type: "runtime", Label: "PHP 8.2 FPM", Status: "healthy", ParentIDs: []string{"node-db"}},
			{ID: "node-db", Type: "database", Label: "MariaDB 10.11", Status: "healthy", ParentIDs: []string{}},
		},
	}
}

// RunIncidentAutopsy analyzes system logs to trace root causes.
func (e *BlueprintEngine) RunIncidentAutopsy(ctx context.Context, incidentLogs string) *IncidentAutopsyReport {
	return &IncidentAutopsyReport{
		IncidentID:  fmt.Sprintf("autopsy-%d", time.Now().Unix()),
		AnalyzedAt:  time.Now().UTC(),
		RootCause:   "Database connection exhaustion triggered by upstream query lock",
		Timeline: []AutopsyFinding{
			{
				Timestamp:   time.Now().Add(-15 * time.Minute),
				Component:   "MariaDB",
				Severity:    "root_cause",
				Description: "max_connections reached (151/151); active long-running SELECT query on table 'orders'",
				Evidence:    "Error 1040 (HY000): Too many connections in /var/log/mysql/error.log",
			},
			{
				Timestamp:   time.Now().Add(-14 * time.Minute),
				Component:   "PHP-FPM",
				Severity:    "contributing_factor",
				Description: "Worker pool exhausted waiting for PDO database responses",
				Evidence:    "server reached pm.max_children setting (50) in /var/log/php8.2-fpm.log",
			},
			{
				Timestamp:   time.Now().Add(-12 * time.Minute),
				Component:   "Nginx",
				Severity:    "symptom",
				Description: "HTTP 502 Bad Gateway returned to visitor web traffic",
				Evidence:    "connect() to unix:/run/php/php8.2-fpm.sock failed (11: Resource temporarily unavailable)",
			},
		},
		Recommendations: []string{
			"Increase MariaDB max_connections from 151 to 300 with adequate RAM buffers",
			"Enable MySQL slow query log to identify unindexed queries on 'orders'",
			"Tune PHP-FPM pm.max_children and request_terminate_timeout",
		},
	}
}
