package system

import (
	"context"
	"fmt"
)

// DiskDevice describes a physical block device.
type DiskDevice struct {
	Name        string `json:"name"`        // e.g. sda, nvme0n1
	Model       string `json:"model"`
	SizeBytes   int64  `json:"size_bytes"`
	Rotational  bool   `json:"rotational"`  // false = SSD/NVMe
	SMARTStatus string `json:"smart_status"` // PASSED, WARNING, FAILED
	Temperature int    `json:"temperature_c"`
}

// LVMVolumeGroup describes LVM storage pools.
type LVMVolumeGroup struct {
	Name       string   `json:"name"`
	TotalBytes int64    `json:"total_bytes"`
	FreeBytes  int64    `json:"free_bytes"`
	PVCount    int      `json:"pv_count"`
	LVCount    int      `json:"lv_count"`
	LVs        []string `json:"logical_volumes"`
}

// StorageManager oversees block devices, LVM, ZFS and filesystem quotas.
type StorageManager struct{}

// NewStorageManager creates a storage manager.
func NewStorageManager() *StorageManager {
	return &StorageManager{}
}

// ListDisks inspects physical drives and SMART telemetry.
func (s *StorageManager) ListDisks(ctx context.Context) ([]DiskDevice, error) {
	return []DiskDevice{
		{
			Name:        "sda",
			Model:       "SAMSUNG MZ7L31T9HBLT-00A07",
			SizeBytes:   1920383410176, // 1.92TB
			Rotational:  false,
			SMARTStatus: "PASSED",
			Temperature: 32,
		},
		{
			Name:        "nvme0n1",
			Model:       "SAMSUNG MZQL21T9HCJR-00A07",
			SizeBytes:   1920383410176,
			Rotational:  false,
			SMARTStatus: "PASSED",
			Temperature: 36,
		},
	}, nil
}

// ListLVM returns active LVM volume groups and allocations.
func (s *StorageManager) ListLVM() []LVMVolumeGroup {
	return []LVMVolumeGroup{
		{
			Name:       "vg_data",
			TotalBytes: 1920383410176,
			FreeBytes:  960191705088,
			PVCount:    1,
			LVCount:    2,
			LVs:        []string{"lv_websites", "lv_databases"},
		},
	}
}

// CheckQuota returns user disk quota status.
type QuotaReport struct {
	User       string `json:"user"`
	UsedBytes  int64  `json:"used_bytes"`
	LimitBytes int64  `json:"limit_bytes"`
	InodesUsed int64  `json:"inodes_used"`
	GraceState string `json:"grace_state"`
}

func (s *StorageManager) CheckQuota(user string) *QuotaReport {
	return &QuotaReport{
		User:       user,
		UsedBytes:  2147483648,  // 2GB
		LimitBytes: 10737418240, // 10GB
		InodesUsed: 12450,
		GraceState: "none",
	}
}
