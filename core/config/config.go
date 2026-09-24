package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var (
	Version   = "2.3.0"
	GitCommit = "dev"
	BuildDate = "2026-09-22"
)

// Config represents the complete KenPanel control plane configuration.
type Config struct {
	ListenAddr        string `json:"listen_addr"`
	DataDir           string `json:"data_dir"`
	LogDir            string `json:"log_dir"`
	RunDir            string `json:"run_dir"`
	DatabaseURL       string `json:"database_url"`
	JWTSecret         string `json:"jwt_secret"`
	AllowRegistration bool   `json:"allow_registration"`
	TelemetryEnabled  bool   `json:"telemetry_enabled"`
	TLSCertFile       string `json:"tls_cert_file"`
	TLSKeyFile        string `json:"tls_key_file"`

	mu sync.RWMutex
}

// DefaultConfig returns safe, production-oriented defaults.
func DefaultConfig() *Config {
	return &Config{
		ListenAddr:        "127.0.0.1:8443",
		DataDir:           "/var/lib/kenpanel",
		LogDir:            "/var/log/kenpanel",
		RunDir:            "/run/kenpanel",
		DatabaseURL:       "sqlite:///var/lib/kenpanel/kenpanel.db",
		JWTSecret:         "",
		AllowRegistration: false,
		TelemetryEnabled:  false, // Non-negotiable: zero forced telemetry
		TLSCertFile:       "/etc/kenpanel/certs/server.crt",
		TLSKeyFile:        "/etc/kenpanel/certs/server.key",
	}
}

// Load reads and parses configuration from a JSON/conf file.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return cfg, nil
}
