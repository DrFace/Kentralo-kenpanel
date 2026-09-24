package transplant

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestServerClonePlan(t *testing.T) {
	engine := NewServerCloneEngine(ServerCloneConfig{
		SourceHost:          "old-host.example.com",
		TargetHost:          "new-host.example.com",
		NewHostname:         "srv02.example.com",
		RegenerateMachineID: true,
		RenewSSHHostKeys:    true,
		IPRemapTable: map[string]string{
			"192.168.1.50": "10.0.0.100",
		},
		PreserveSourceState: true,
		DryRun:              true,
	})

	plan, err := engine.GeneratePlan(context.Background())
	if err != nil {
		t.Fatalf("unexpected error generating plan: %v", err)
	}

	if !plan.CanExecute {
		t.Fatalf("expected plan to be executable")
	}

	if len(plan.Actions) == 0 {
		t.Fatalf("expected plan actions to be generated")
	}

	// Verify hostname substitution
	if plan.Substitutions["old-host.example.com"] != "srv02.example.com" {
		t.Errorf("expected hostname substitution, got %s", plan.Substitutions["old-host.example.com"])
	}
}

func TestCPanelDiscoveryNonExistent(t *testing.T) {
	parser := NewCPanelDiscoveryParser(filepath.Join(os.TempDir(), "nonexistent-cpmove.tar.gz"))
	_, err := parser.Parse(context.Background())
	if err == nil {
		t.Errorf("expected error for non-existent archive")
	}
}
