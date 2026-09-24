package os

import (
	"context"

	"github.com/Kentralo/kenpanel/pkg/protocol"
)

// OSAdapter defines the uniform contract across all supported operating systems.
type OSAdapter interface {
	// Name returns the canonical OS family identifier (e.g., "debian", "rhel", "alpine", "windows").
	Name() string

	// DiscoverCapabilities inspects the host and returns detailed system capabilities.
	DiscoverCapabilities(ctx context.Context) (*protocol.NodeCapabilities, error)

	// ManagePackage installs, updates, or removes OS packages safely.
	InstallPackage(ctx context.Context, pkgName string) error
	RemovePackage(ctx context.Context, pkgName string) error

	// ServiceControl starts, stops, reloads, or checks service status.
	StartService(ctx context.Context, serviceName string) error
	StopService(ctx context.Context, serviceName string) error
	ReloadService(ctx context.Context, serviceName string) error
	GetServiceStatus(ctx context.Context, serviceName string) (string, error)

	// FileOperations performs privileged filesystem actions with strict path escaping protection.
	WriteFileAtomic(ctx context.Context, path string, content []byte, mode uint32, owner, group string) error
	ReadFileSafe(ctx context.Context, path string, maxBytes int64) ([]byte, error)
	EnsureDirectory(ctx context.Context, path string, mode uint32, owner, group string) error
}
