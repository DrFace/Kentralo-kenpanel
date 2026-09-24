package registry

import (
	_ "embed"
	"os"
	"testing"
)

func TestRequirementRegistry(t *testing.T) {
	data, err := os.ReadFile("data/requirements.json")
	if err != nil {
		t.Fatalf("failed to read requirements.json: %v", err)
	}

	reg := NewRegistry()
	if err := reg.LoadFromJSON(data); err != nil {
		t.Fatalf("failed to parse requirements: %v", err)
	}

	all := reg.All()
	if len(all) == 0 {
		t.Fatalf("expected registered requirements, found 0")
	}

	req, found := reg.Get("KP-REQ-0001")
	if !found {
		t.Errorf("expected to find KP-REQ-0001")
	}
	if req.Subsystem != "webstack" {
		t.Errorf("expected subsystem webstack, got %s", req.Subsystem)
	}

	transplants := reg.FilterBySubsystem("transplant")
	if len(transplants) == 0 {
		t.Errorf("expected transplant requirements")
	}
}
