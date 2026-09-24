package transplant

import (
	"time"
)

// ManifestVersion specifies the serialization schema version.
const ManifestVersion = "v1.0"

// MigrationManifest is the universal, normalized migration artifact.
type MigrationManifest struct {
	Version      string                 `json:"version"`
	ExportedAt   time.Time              `json:"exported_at"`
	SourcePanel  string                 `json:"source_panel"` // "cpanel", "plesk", "cyberpanel", "aapanel", "kenpanel"
	SourceHost   string                 `json:"source_host"`
	Accounts     []AccountManifest      `json:"accounts"`
	Databases    []DatabaseManifest     `json:"databases"`
	MailDomains  []MailDomainManifest   `json:"mail_domains"`
	CronJobs     []CronJobManifest      `json:"cron_jobs"`
	LossWarnings []string               `json:"semantic_loss_warnings"`
	Supplemental map[string]interface{} `json:"supplemental"`
}

// AccountManifest models a migrated tenant or website owner account.
type AccountManifest struct {
	Username   string            `json:"username"`
	HomeDir    string            `json:"home_dir"`
	Websites   []WebsiteManifest `json:"websites"`
	QuotaBytes int64             `json:"quota_bytes"`
}

// WebsiteManifest models a migrated domain/virtual host.
type WebsiteManifest struct {
	PrimaryDomain string   `json:"primary_domain"`
	Aliases       []string `json:"aliases"`
	DocRoot       string   `json:"doc_root"`
	PHPVersion    string   `json:"php_version,omitempty"`
	SSLCertPem    string   `json:"ssl_cert_pem,omitempty"`
	SSLKeyPem     string   `json:"ssl_key_pem,omitempty"`
}

// DatabaseManifest models an exported database with credentials and privileges.
type DatabaseManifest struct {
	Engine    string   `json:"engine"` // "mariadb", "mysql", "postgres"
	Name      string   `json:"name"`
	DumpFile  string   `json:"dump_file"` // Relative path in transfer bundle
	Users     []string `json:"users"`
	Collation string   `json:"collation"`
}

// MailDomainManifest models mailboxes and routing.
type MailDomainManifest struct {
	Domain    string            `json:"domain"`
	Mailboxes []MailboxManifest `json:"mailboxes"`
}

// MailboxManifest models an individual mail account.
type MailboxManifest struct {
	Address        string   `json:"address"`
	PasswordFormat string   `json:"password_format"` // "dovecot_sha512", "md5", "crypt"
	PasswordHash   string   `json:"password_hash"`
	QuotaBytes     int64    `json:"quota_bytes"`
	Forwarders     []string `json:"forwarders,omitempty"`
}

// CronJobManifest models scheduled cron jobs.
type CronJobManifest struct {
	Schedule string `json:"schedule"`
	Command  string `json:"command"`
	User     string `json:"user"`
}
