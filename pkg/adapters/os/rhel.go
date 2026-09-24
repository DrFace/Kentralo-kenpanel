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

	"github.com/Kentralo/kenpanel/pkg/protocol"
)

// RHELAdapter implements the OSAdapter interface for RHEL, AlmaLinux, and Rocky Linux.
type RHELAdapter struct{}

func NewRHELAdapter() *RHELAdapter {
	return &RHELAdapter{}
}

func (a *RHELAdapter) Name() string {
	return "rhel"
}

func (a *RHELAdapter) DiscoverCapabilities(ctx context.Context) (*protocol.NodeCapabilities, error) {
	hostname, _ := os.Hostname()

	caps := &protocol.NodeCapabilities{
		Hostname:         hostname,
		Kernel:           runtime.GOOS,
		OS:               "almalinux",
		OSVersion:        "9",
		Arch:             runtime.GOARCH,
		TotalMemoryMB:    4096,
		CPUCores:         runtime.NumCPU(),
		PackageManagers:  []string{"dnf", "rpm"},
		ServiceManagers:  []string{"systemd"},
		WebServers:       []string{"nginx", "httpd"},
		Runtimes:         map[string][]string{"php": {"8.1", "8.2", "8.3"}, "node": {"20"}},
		Databases:        []string{"mariadb", "postgres"},
		StorageStacks:    []string{"xfs", "ext4"},
		ContainerEngines: []string{"podman", "docker"},
		Firewalls:        []string{"firewalld", "nftables"},
		DiscoveredAt:     time.Now().UTC(),
	}

	return caps, nil
}

func (a *RHELAdapter) InstallPackage(ctx context.Context, pkgName string) error {
	if strings.ContainsAny(pkgName, ";|&$`") {
		return fmt.Errorf("invalid characters in package name: %s", pkgName)
	}

	cmd := exec.CommandContext(ctx, "dnf", "install", "-y", pkgName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("dnf install failed for %s: %s (%w)", pkgName, string(output), err)
	}
	return nil
}

func (a *RHELAdapter) RemovePackage(ctx context.Context, pkgName string) error {
	if strings.ContainsAny(pkgName, ";|&$`") {
		return fmt.Errorf("invalid characters in package name: %s", pkgName)
	}

	cmd := exec.CommandContext(ctx, "dnf", "remove", "-y", pkgName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("dnf remove failed for %s: %s (%w)", pkgName, string(output), err)
	}
	return nil
}

func (a *RHELAdapter) StartService(ctx context.Context, serviceName string) error {
	return exec.CommandContext(ctx, "systemctl", "start", serviceName).Run()
}

func (a *RHELAdapter) StopService(ctx context.Context, serviceName string) error {
	return exec.CommandContext(ctx, "systemctl", "stop", serviceName).Run()
}

func (a *RHELAdapter) ReloadService(ctx context.Context, serviceName string) error {
	return exec.CommandContext(ctx, "systemctl", "reload", serviceName).Run()
}

func (a *RHELAdapter) GetServiceStatus(ctx context.Context, serviceName string) (string, error) {
	cmd := exec.CommandContext(ctx, "systemctl", "is-active", serviceName)
	output, err := cmd.Output()
	status := strings.TrimSpace(string(output))
	if err != nil && status == "" {
		return "inactive", nil
	}
	return status, nil
}

func (a *RHELAdapter) WriteFileAtomic(ctx context.Context, path string, content []byte, mode uint32, owner, group string) error {
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, "kenpanel-rhel-atomic-*.tmp")
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

func (a *RHELAdapter) ReadFileSafe(ctx context.Context, path string, maxBytes int64) ([]byte, error) {
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

func (a *RHELAdapter) EnsureDirectory(ctx context.Context, path string, mode uint32, owner, group string) error {
	return os.MkdirAll(path, os.FileMode(mode))
}
