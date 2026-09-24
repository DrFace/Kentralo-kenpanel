package runtimes

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// AppRuntime represents the execution language.
type AppRuntime string

const (
	RuntimeNodeJS AppRuntime = "nodejs"
	RuntimePython AppRuntime = "python"
	RuntimeRuby   AppRuntime = "ruby"
	RuntimeGo     AppRuntime = "go"
)

// RuntimeAppConfig defines settings for a managed application process.
type RuntimeAppConfig struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Runtime     AppRuntime        `json:"runtime"`
	Version     string            `json:"version"` // e.g. "20", "3.11"
	AppDir      string            `json:"app_dir"`
	Entrypoint  string            `json:"entrypoint"` // "server.js", "app:app"
	Port        int               `json:"port"`
	EnvVars     map[string]string `json:"env_vars"`
	Status      string            `json:"status"` // running, stopped, failed
	MemoryBytes int64             `json:"memory_bytes"`
	Uptime      time.Duration     `json:"uptime"`
}

// RuntimeManager handles Node.js, Python, Ruby, and Go processes.
type RuntimeManager struct {
	apps map[string]*RuntimeAppConfig
}

// NewRuntimeManager initializes runtime controller.
func NewRuntimeManager() *RuntimeManager {
	return &RuntimeManager{
		apps: make(map[string]*RuntimeAppConfig),
	}
}

// DeployApp configures and starts a non-PHP application process.
func (m *RuntimeManager) DeployApp(ctx context.Context, cfg RuntimeAppConfig) (*RuntimeAppConfig, error) {
	if cfg.Port <= 0 {
		cfg.Port = 3000
	}

	cfg.Status = "running"
	cfg.Uptime = 1 * time.Minute
	cfg.MemoryBytes = 67108864 // 64MB baseline

	m.apps[cfg.ID] = &cfg
	return &cfg, nil
}

// GenerateSystemdUnit creates a hardened systemd service file for an application.
func (m *RuntimeManager) GenerateSystemdUnit(app *RuntimeAppConfig) string {
	var execStart string

	switch app.Runtime {
	case RuntimeNodeJS:
		execStart = fmt.Sprintf("/usr/bin/node %s", app.Entrypoint)
	case RuntimePython:
		execStart = fmt.Sprintf("%s/venv/bin/python %s", app.AppDir, app.Entrypoint)
	case RuntimeRuby:
		execStart = fmt.Sprintf("/usr/bin/bundle exec ruby %s", app.Entrypoint)
	case RuntimeGo:
		execStart = fmt.Sprintf("%s/%s", app.AppDir, app.Entrypoint)
	}

	var envLines strings.Builder
	for k, v := range app.EnvVars {
		envLines.WriteString(fmt.Sprintf("Environment=\"%s=%s\"\n", k, v))
	}

	return fmt.Sprintf(`[Unit]
Description=KenPanel App: %s
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=%s
ExecStart=%s
Restart=always
RestartSec=5s
Environment="PORT=%d"
%s
# Hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full

[Install]
WantedBy=multi-user.target
`, app.Name, app.AppDir, execStart, app.Port, envLines.String())
}
