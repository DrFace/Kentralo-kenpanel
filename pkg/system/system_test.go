package system

import (
	"context"
	"testing"
)

func TestInventoryCollector(t *testing.T) {
	collector := NewInventoryCollector()
	inv, err := collector.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error collecting inventory: %v", err)
	}

	if inv.Hostname == "" {
		t.Errorf("expected hostname to be set")
	}
	if inv.CPUCores < 1 {
		t.Errorf("expected at least 1 core, got %d", inv.CPUCores)
	}
}

func TestPackageManagerUpgrades(t *testing.T) {
	pm := NewPackageManager("apt")
	upgrades, err := pm.ListUpgrades(context.Background(), false)
	if err != nil {
		t.Fatalf("unexpected error listing upgrades: %v", err)
	}
	if len(upgrades) == 0 {
		t.Fatalf("expected upgrades")
	}

	secOnly, err := pm.ListUpgrades(context.Background(), true)
	if err != nil {
		t.Fatalf("unexpected error listing security upgrades: %v", err)
	}
	for _, p := range secOnly {
		if !p.IsSecurity {
			t.Errorf("expected only security upgrades")
		}
	}
}

func TestProcessManagerTop(t *testing.T) {
	pm := NewProcessManager()
	procs, err := pm.ListProcesses(context.Background())
	if err != nil {
		t.Fatalf("unexpected error listing processes: %v", err)
	}

	top := pm.TopConsumers(procs, "cpu", 2)
	if len(top) != 2 {
		t.Errorf("expected 2 top processes, got %d", len(top))
	}
}
