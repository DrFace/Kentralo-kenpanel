package service

import (
	"context"
)

// WebServerAdapter controls HTTP/HTTPS reverse proxies and web server daemons (Nginx, Apache, Caddy).
type WebServerAdapter interface {
	Name() string
	TestConfig(ctx context.Context) (bool, string, error)
	DeployVHost(ctx context.Context, domain string, configContent []byte) error
	RemoveVHost(ctx context.Context, domain string) error
	Reload(ctx context.Context) error
}

// DatabaseAdapter controls database engines (MariaDB, MySQL, PostgreSQL).
type DatabaseAdapter interface {
	Engine() string
	CreateDatabase(ctx context.Context, dbName, collation string) error
	DropDatabase(ctx context.Context, dbName string) error
	CreateUser(ctx context.Context, username, password, host string) error
	GrantPermissions(ctx context.Context, username, dbName, host string, privileges []string) error
	DumpDatabase(ctx context.Context, dbName, targetPath string) error
	RestoreDatabase(ctx context.Context, dbName, sourcePath string) error
}

// DNSAdapter manages DNS zones and resource records.
type DNSAdapter interface {
	ProviderName() string
	CreateZone(ctx context.Context, zoneName string) error
	SetRecord(ctx context.Context, zoneName, recordName, recordType, content string, ttl int) error
	DeleteRecord(ctx context.Context, zoneName, recordName, recordType string) error
}

// FirewallAdapter manages host port exposure and packet filtering rules (UFW, firewalld, nftables).
type FirewallAdapter interface {
	Name() string
	ListOpenPorts(ctx context.Context) ([]int, error)
	AllowPort(ctx context.Context, port int, proto string, comment string) error
	DenyPort(ctx context.Context, port int, proto string) error
	ApplyRules(ctx context.Context) error
}
