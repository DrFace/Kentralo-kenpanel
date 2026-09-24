package benchmark

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestBenchmarkSafeExecution(t *testing.T) {
	engine := NewEngine(false, DefaultSafetyLimits)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := engine.Run(ctx, ProfileCPU, "")
	if err != nil {
		t.Fatalf("unexpected error running CPU benchmark: %v", err)
	}

	if res.Score <= 0 {
		t.Errorf("expected positive score, got %f", res.Score)
	}
	if res.Context.KenPanelVersion != "2.3.0" {
		t.Errorf("unexpected KenPanel version: %s", res.Context.KenPanelVersion)
	}
}

func TestBenchmarkProductionSafetyGuard(t *testing.T) {
	engine := NewEngine(true, DefaultSafetyLimits) // production mode
	ctx := context.Background()

	// 1. Without ack token - MUST FAIL
	_, err := engine.Run(ctx, ProfileCPU, "")
	if err == nil {
		t.Fatal("expected error when running benchmark on production without ack token")
	}

	// 2. With invalid token - MUST FAIL
	_, err = engine.Run(ctx, ProfileCPU, "INVALID-TOKEN")
	if err == nil {
		t.Fatal("expected error when running benchmark on production with invalid ack token")
	}

	// 3. With valid token - MUST SUCCEED
	res, err := engine.Run(ctx, ProfileCPU, "ACK-PROD-STRESS-12345678-ABCD")
	if err != nil {
		t.Fatalf("expected success with valid token, got: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestBenchmarkCompareAndExport(t *testing.T) {
	engine := NewEngine(false, DefaultSafetyLimits)
	ctx := context.Background()

	res1, _ := engine.Run(ctx, ProfileMemory, "")
	res2, _ := engine.Run(ctx, ProfileMemory, "")

	cmp, err := engine.Compare(res1.ID, res2.ID)
	if err != nil {
		t.Fatalf("unexpected error comparing runs: %v", err)
	}
	if cmp.Profile != ProfileMemory {
		t.Errorf("expected profile memory, got %s", cmp.Profile)
	}

	// Test JSON export
	jsonBytes, err := engine.ExportJSON()
	if err != nil || len(jsonBytes) == 0 {
		t.Fatalf("failed JSON export: %v", err)
	}

	// Test CSV export
	csvStr, err := engine.ExportCSV()
	if err != nil || !strings.Contains(csvStr, "memory") {
		t.Fatalf("failed CSV export: %v", err)
	}
}
