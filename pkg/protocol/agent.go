package protocol

import (
	"time"
)

// AgentRole defines the privilege mode of the agent.
type AgentRole string

const (
	AgentRolePrivilegedDaemon AgentRole = "privileged_daemon" // runs as root/SYSTEM
)

// NodeCapabilities represents the self-discovered hardware and software features of a node.
type NodeCapabilities struct {
	Hostname       string            `json:"hostname"`
	Kernel         string            `json:"kernel"`
	OS             string            `json:"os"`
	OSVersion      string            `json:"os_version"`
	Arch           string            `json:"arch"`
	TotalMemoryMB  uint64            `json:"total_memory_mb"`
	CPUCores       int               `json:"cpu_cores"`
	PackageManagers []string         `json:"package_managers"` // "apt", "dnf", "apk", "pacman"
	ServiceManagers []string         `json:"service_managers"` // "systemd", "openrc", "windows_service"
	WebServers     []string          `json:"web_servers"`      // "nginx", "apache", "caddy"
	Runtimes       map[string][]string `json:"runtimes"`       // "php": ["8.1", "8.2", "8.3"], "node": ["20", "22"]
	Databases      []string          `json:"databases"`        // "mariadb", "postgres", "redis"
	StorageStacks  []string          `json:"storage_stacks"`   // "ext4", "xfs", "btrfs", "zfs", "lvm"
	ContainerEngines []string        `json:"container_engines"`// "docker", "podman", "containerd"
	Firewalls      []string          `json:"firewalls"`        // "ufw", "firewalld", "nftables"
	DiscoveredAt   time.Time         `json:"discovered_at"`
}

// RPCRequest is a strictly typed command sent from control plane to agent.
type RPCRequest struct {
	ID        string                 `json:"id"`
	Action    string                 `json:"action"` // e.g. "webserver:reload", "service:status", "cert:deploy"
	Subsystem string                 `json:"subsystem"`
	Params    map[string]interface{} `json:"params"`
	Timestamp time.Time              `json:"timestamp"`
}

// RPCResponse is the result returned by the agent.
type RPCResponse struct {
	RequestID string                 `json:"request_id"`
	Success   bool                   `json:"success"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Error     string                 `json:"error,omitempty"`
	DurationMS int64                 `json:"duration_ms"`
}
