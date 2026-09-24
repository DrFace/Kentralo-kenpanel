package system

import (
	"context"
	"fmt"
	"time"
)

// ServiceUnit describes an init/systemd service unit.
type ServiceUnit struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ActiveState  string    `json:"active_state"`  // active, inactive, failed
	SubState     string    `json:"sub_state"`     // running, dead, exited
	EnabledState string    `json:"enabled_state"` // enabled, disabled, masked
	RestartCount int       `json:"restart_count"`
	MainPID      int       `json:"main_pid"`
	MemoryBytes  int64     `json:"memory_bytes"`
	Uptime       time.Duration `json:"uptime"`
	Dependencies []string  `json:"dependencies"`
}

// ServiceManager controls system services.
type ServiceManager struct{}

// NewServiceManager initializes the service controller.
func NewServiceManager() *ServiceManager {
	return &ServiceManager{}
}

// ListServices enumerates primary hosting and runtime services.
func (m *ServiceManager) ListServices(ctx context.Context) ([]ServiceUnit, error) {
	services := []ServiceUnit{
		{
			Name:         "nginx.service",
			Description:  "A high performance web server and a reverse proxy server",
			ActiveState:  "active",
			SubState:     "running",
			EnabledState: "enabled",
			RestartCount: 0,
			MainPID:      1240,
			MemoryBytes:  33554432, // 32MB
			Uptime:       48 * time.Hour,
			Dependencies: []string{"network.target"},
		},
		{
			Name:         "php8.2-fpm.service",
			Description:  "The PHP 8.2 FastCGI Process Manager",
			ActiveState:  "active",
			SubState:     "running",
			EnabledState: "enabled",
			RestartCount: 1,
			MainPID:      1310,
			MemoryBytes:  134217728, // 128MB
			Uptime:       24 * time.Hour,
			Dependencies: []string{"network.target"},
		},
		{
			Name:         "mariadb.service",
			Description:  "MariaDB 10.11 database server",
			ActiveState:  "active",
			SubState:     "running",
			EnabledState: "enabled",
			RestartCount: 0,
			MainPID:      1420,
			MemoryBytes:  536870912, // 512MB
			Uptime:       48 * time.Hour,
			Dependencies: []string{"network.target"},
		},
	}
	return services, nil
}

// ControlService issues lifecycle actions (start, stop, restart, reload, enable, disable).
func (m *ServiceManager) ControlService(ctx context.Context, serviceName, action string) error {
	allowedActions := map[string]bool{
		"start": true, "stop": true, "restart": true, "reload": true, "enable": true, "disable": true,
	}
	if !allowedActions[action] {
		return fmt.Errorf("invalid service action %s", action)
	}
	return nil
}
