package models

import (
	"time"
)

// Organization represents a multi-tenant isolation boundary.
type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User represents a KenPanel operator, customer, or API service account.
type User struct {
	ID           string    `json:"id"`
	OrgID        string    `json:"org_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	MFAEnabled   bool      `json:"mfa_enabled"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Node represents an enrolled server node running kenpanel-agent.
type Node struct {
	ID           string            `json:"id"`
	OrgID        string            `json:"org_id"`
	Hostname     string            `json:"hostname"`
	PublicIP     string            `json:"public_ip"`
	PrivateIP    string            `json:"private_ip"`
	OS           string            `json:"os"`
	OSVersion    string            `json:"os_version"`
	Arch         string            `json:"arch"`
	Status       string            `json:"status"` // "online", "offline", "unreachable", "maintenance"
	Capabilities map[string]string `json:"capabilities"`
	LastSeenAt   time.Time         `json:"last_seen_at"`
	CreatedAt    time.Time         `json:"created_at"`
}

// Website represents a hosted web property.
type Website struct {
	ID            string    `json:"id"`
	NodeID        string    `json:"node_id"`
	PrimaryDomain string    `json:"primary_domain"`
	Aliases       []string  `json:"aliases"`
	DocumentRoot  string    `json:"document_root"`
	WebStackMode  string    `json:"webstack_mode"` // "nginx_apache", "nginx_phpfpm", "caddy_standalone"
	RuntimeType   string    `json:"runtime_type"`  // "php", "node", "python", "static", "docker"
	RuntimeVer    string    `json:"runtime_version"`
	SSLStatus     string    `json:"ssl_status"` // "active", "expired", "pending", "none"
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Database represents a managed database instance (MySQL/MariaDB/PostgreSQL).
type Database struct {
	ID        string    `json:"id"`
	NodeID    string    `json:"node_id"`
	Engine    string    `json:"engine"` // "mariadb", "mysql", "postgres"
	Name      string    `json:"name"`
	Collation string    `json:"collation"`
	SizeBytes int64     `json:"size_bytes"`
	CreatedAt time.Time `json:"created_at"`
}

// Mailbox represents an email account under a managed mail domain.
type Mailbox struct {
	ID          string    `json:"id"`
	NodeID      string    `json:"node_id"`
	Email       string    `json:"email"`
	QuotaBytes  int64     `json:"quota_bytes"`
	UsedBytes   int64     `json:"used_bytes"`
	IsForwarder bool      `json:"is_forwarder"`
	ForwardTo   []string  `json:"forward_to,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// Checkpoint represents a ChangeGuard pre-modification snapshot.
type Checkpoint struct {
	ID          string            `json:"id"`
	NodeID      string            `json:"node_id"`
	Subsystem   string            `json:"subsystem"`
	Description string            `json:"description"`
	FileHashes  map[string]string `json:"file_hashes"`
	BackupPath  string            `json:"backup_path"`
	CreatedAt   time.Time         `json:"created_at"`
}
