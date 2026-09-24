package os

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/DrFace/Kentralo-kenpanel/pkg/protocol"
)

// DebianAdapter implements the OSAdapter interface for Debian and Ubuntu distributions.
type DebianAdapter struct{}

func NewDebianAdapter() *DebianAdapter {
	return &DebianAdapter{}
}

func (a *DebianAdapter) Name() string {
	return "debian"
}

func (a *DebianAdapter) DiscoverCapabilities(ctx context.Context) (*protocol.NodeCapabilities, error) {
	hostname, _ := os.Hostname()

	caps := &protocol.NodeCapabilities{
		Hostname:         hostname,
		Kernel:           runtime.GOOS,
		OS:               "debian",
		OSVersion:        "12",
		Arch:             runtime.GOARCH,
		TotalMemoryMB:    4096,
		CPUCores:         runtime.NumCPU(),
		PackageManagers:  []string{"apt"},
		ServiceManagers:  []string{"systemd"},
		WebServers:       []string{"nginx", "apache2"},
		Runtimes:         map[string][]string{"php": {"8.1", "8.2", "8.3", "8.4"}, "node": {"20", "22"}},
		Databases:        []string{"mariadb", "postgres"},
		StorageStacks:    []string{"ext4", "xfs"},
		ContainerEngines: []string{"docker"},
		Firewalls:        []string{"ufw", "nftables"},
		DiscoveredAt:     time.Now().UTC(),
	}

	return caps, nil
}

func (a *DebianAdapter) InstallPackage(ctx context.Context, pkgName string) error {
	// Sanitize package name to avoid injection
	if strings.ContainsAny(pkgName, ";|&$`") {
		return fmt.Errorf("invalid characters in package name: %s", pkgName)
	}

	cmd := exec.CommandContext(ctx, "apt-get", "install", "-y", "--no-install-recommends", pkgName)
	cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("apt-get install failed for %s: %s (%w)", pkgName, string(output), err)
	}
	return nil
}

func (a *DebianAdapter) RemovePackage(ctx context.Context, pkgName string) error {
	if strings.ContainsAny(pkgName, ";|&$`") {
		return fmt.Errorf("invalid characters in package name: %s", pkgName)
	}

	cmd := exec.CommandContext(ctx, "apt-get", "remove", "-y", pkgName)
	cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("apt-get remove failed for %s: %s (%w)", pkgName, string(output), err)
	}
	return nil
}

func (a *DebianAdapter) StartService(ctx context.Context, serviceName string) error {
	return exec.CommandContext(ctx, "systemctl", "start", serviceName).Run()
}

func (a *DebianAdapter) StopService(ctx context.Context, serviceName string) error {
	return exec.CommandContext(ctx, "systemctl", "stop", serviceName).Run()
}

func (a *DebianAdapter) ReloadService(ctx context.Context, serviceName string) error {
	return exec.CommandContext(ctx, "systemctl", "reload", serviceName).Run()
}

func (a *DebianAdapter) GetServiceStatus(ctx context.Context, serviceName string) (string, error) {
	cmd := exec.CommandContext(ctx, "systemctl", "is-active", serviceName)
	output, err := cmd.Output()
	status := strings.TrimSpace(string(output))
	if err != nil && status == "" {
		return "inactive", nil
	}
	return status, nil
}

func (a *DebianAdapter) WriteFileAtomic(ctx context.Context, path string, content []byte, mode uint32, owner, group string) error {
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, "kenpanel-atomic-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	if _, err := tmpFile.Write(content); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed writing to temp file: %w", err)
	}
	if err := tmpFile.Chmod(os.FileMode(mode)); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to chmod temp file: %w", err)
	}
	tmpFile.Close()

	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("failed atomic rename to %s: %w", path, err)
	}
	return nil
}

func (a *DebianAdapter) ReadFileSafe(ctx context.Context, path string, maxBytes int64) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, maxBytes)
	n, err := f.Read(buf)
	if err != nil && n == 0 {
		return nil, err
	}
	return buf[:n], nil
}

func (a *DebianAdapter) EnsureDirectory(ctx context.Context, path string, mode uint32, owner, group string) error {
	return os.MkdirAll(path, os.FileMode(mode))
}
