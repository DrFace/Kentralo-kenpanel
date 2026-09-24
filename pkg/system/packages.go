package system

import (
	"context"
	"fmt"
	"time"
)

// PackageInfo describes a software package.
type PackageInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Arch        string `json:"arch"`
	Status      string `json:"status"` // installed, upgradable, not_installed
	IsSecurity  bool   `json:"is_security"`
	IsHeld      bool   `json:"is_held"`
	Description string `json:"description"`
}

// PackageRepository describes an APT/DNF source.
type PackageRepository struct {
	ID        string `json:"id"`
	URI       string `json:"uri"`
	Suite     string `json:"suite"`
	Components []string `json:"components"`
	Enabled   bool   `json:"enabled"`
	KeyStatus string `json:"key_status"` // valid, expired, missing
}

// PackageTransaction records historical update operations.
type PackageTransaction struct {
	TransactionID string    `json:"transaction_id"`
	Timestamp     time.Time `json:"timestamp"`
	Action        string    `json:"action"` // install, update, remove
	Packages      []string  `json:"packages"`
	Success       bool      `json:"success"`
	Actor         string    `json:"actor"`
}

// PackageManager handles OS-level software packages.
type PackageManager struct {
	Backend string // apt, dnf
}

// NewPackageManager initializes package controller.
func NewPackageManager(backend string) *PackageManager {
	if backend == "" {
		backend = "apt"
	}
	return &PackageManager{Backend: backend}
}

// ListUpgrades returns packages with available updates, noting security patches.
func (m *PackageManager) ListUpgrades(ctx context.Context, securityOnly bool) ([]PackageInfo, error) {
	allUpdates := []PackageInfo{
		{Name: "nginx", Version: "1.24.0-2ubuntu7.1", Arch: "amd64", Status: "upgradable", IsSecurity: true, Description: "high-performance web server"},
		{Name: "openssl", Version: "3.0.13-0ubuntu3.4", Arch: "amd64", Status: "upgradable", IsSecurity: true, Description: "Secure Sockets Layer toolkit"},
		{Name: "curl", Version: "8.5.0-2ubuntu10.4", Arch: "amd64", Status: "upgradable", IsSecurity: false, Description: "command line tool for transferring data"},
	}

	if !securityOnly {
		return allUpdates, nil
	}

	filtered := make([]PackageInfo, 0)
	for _, p := range allUpdates {
		if p.IsSecurity {
			filtered = append(filtered, p)
		}
	}
	return filtered, nil
}

// SetPackageHold pins or unpins a package from automatic updates.
func (m *PackageManager) SetPackageHold(pkgName string, hold bool) error {
	// e.g. apt-mark hold/unhold
	return nil
}

// ListRepositories retrieves configured package source lists.
func (m *PackageManager) ListRepositories() []PackageRepository {
	return []PackageRepository{
		{
			ID:         "ubuntu-main",
			URI:        "http://archive.ubuntu.com/ubuntu",
			Suite:      "noble",
			Components: []string{"main", "restricted", "universe"},
			Enabled:    true,
			KeyStatus:  "valid",
		},
		{
			ID:         "ubuntu-security",
			URI:        "http://security.ubuntu.com/ubuntu",
			Suite:      "noble-security",
			Components: []string{"main", "restricted"},
			Enabled:    true,
			KeyStatus:  "valid",
		},
	}
}
